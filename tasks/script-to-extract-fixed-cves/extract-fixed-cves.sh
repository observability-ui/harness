#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
CONFIG_FILE="$SCRIPT_DIR/config.yaml"
REPORT="$SCRIPT_DIR/report.md"
WORK_DIR=$(mktemp -d)

YAML_OUTPUT=false
FILTER=""
for arg in "$@"; do
  case "$arg" in
    --yaml) YAML_OUTPUT=true ;;
    *) FILTER="$arg" ;;
  esac
done

mkdir -p "$WORK_DIR/sections" "$WORK_DIR/summaries" "$WORK_DIR/refs" "$WORK_DIR/yaml"

HAS_GOVULNCHECK=true
if ! command -v govulncheck &>/dev/null; then
  echo "[WARN] govulncheck not found — Go vulnerability checks will be skipped" >&2
  HAS_GOVULNCHECK=false
fi

if ! command -v jq &>/dev/null; then
  echo "[ERROR] jq is required. Install: brew install jq" >&2
  exit 1
fi

if ! command -v npm &>/dev/null; then
  echo "[ERROR] npm is required." >&2
  exit 1
fi

if ! command -v yq &>/dev/null; then
  echo "[ERROR] yq is required to parse config.yaml. Install: brew install yq" >&2
  exit 1
fi

HAS_GH=true
if ! command -v gh &>/dev/null; then
  echo "[WARN] gh CLI not found — CVE mapping will be skipped" >&2
  HAS_GH=false
fi

# ---- Load config ----

if [ ! -f "$CONFIG_FILE" ]; then
  echo "[ERROR] Config file not found: $CONFIG_FILE" >&2
  exit 1
fi

DEFAULT_ROLLBACK=$(yq -r '.default_rollback // 3' "$CONFIG_FILE")

ENTRIES=()
ENTRY_COMPONENTS=()
ENTRY_ROLLBACKS=()

entry_count=$(yq '.entries | length' "$CONFIG_FILE")
for (( i=0; i<entry_count; i++ )); do
  project=$(yq -r ".entries[$i].project" "$CONFIG_FILE")
  branch=$(yq -r ".entries[$i].branch" "$CONFIG_FILE")
  component=$(yq -r ".entries[$i].component" "$CONFIG_FILE")
  rollback=$(yq -r ".entries[$i].rollback // $DEFAULT_ROLLBACK" "$CONFIG_FILE")
  ENTRIES+=("${project}:${branch}")
  ENTRY_COMPONENTS+=("$component")
  ENTRY_ROLLBACKS+=("$rollback")
done

echo "Loaded ${#ENTRIES[@]} entries from $CONFIG_FILE (default rollback: $DEFAULT_ROLLBACK)" >&2

# ---- Cleanup ----

cleanup() {
  for refs_file in "$WORK_DIR"/refs/*.txt; do
    [ -f "$refs_file" ] || continue
    while IFS='|' read -r pdir ref; do
      echo "[cleanup] Restoring $pdir" >&2
      git -C "$pdir" checkout -- . 2>/dev/null || true
      git -C "$pdir" checkout "$ref" 2>/dev/null || true
    done < "$refs_file"
  done
  rm -rf "$WORK_DIR"
}
trap cleanup EXIT

# ---- Helper functions ----

resolve_remote() {
  local pdir="$1" branch="$2"
  git -C "$pdir" branch -r | grep -E "/${branch}$" | head -1 | sed 's/^[[:space:]]*//' | cut -d'/' -f1
}

detect_frontend_dir() {
  local pdir="$1"
  if [ -f "$pdir/web/package.json" ]; then
    echo "$pdir/web"
  elif [ -f "$pdir/package.json" ]; then
    echo "$pdir"
  else
    echo ""
  fi
}

extract_npm_urls() {
  local json_file="$1"
  [ -s "$json_file" ] || return 0
  jq -r '
    .vulnerabilities // {} | to_entries[] | .value.via[]? |
    select(type == "object") | .url // empty
  ' "$json_file" 2>/dev/null | sort -u || true
}

extract_npm_detail() {
  local json_file="$1" url="$2"
  [ -s "$json_file" ] || return 0
  jq -r --arg url "$url" '
    .vulnerabilities // {} | to_entries[] | .value.via[]? |
    select(type == "object" and .url == $url) |
    [.name // "N/A", .severity // "N/A", .title // "N/A"] | @tsv
  ' "$json_file" 2>/dev/null | head -1 || true
}

fetch_cve_for_ghsa() {
  local url="$1"
  if [ "$HAS_GH" = false ]; then echo "N/A"; return; fi
  local ghsa_id="${url##*/}"
  gh api "/advisories/$ghsa_id" 2>/dev/null | jq -r '.cve_id // "N/A"' 2>/dev/null || echo "N/A"
}

extract_go_ids() {
  local json_file="$1"
  [ -s "$json_file" ] || return 0
  jq -r 'select(.finding != null) | .finding.osv' "$json_file" 2>/dev/null | sort -u || true
}

extract_go_detail() {
  local json_file="$1" vid="$2"
  [ -s "$json_file" ] || return 0
  jq -r --arg id "$vid" '
    select(.osv != null and .osv.id == $id) | .osv |
    [
      .id,
      (.aliases // [] | join(", ")),
      (if .affected then [.affected[].package.name] | unique | join(", ") else "N/A" end),
      (.summary // "N/A" | gsub("\n"; " "))
    ] | @tsv
  ' "$json_file" 2>/dev/null | head -1 || true
}

# ---- Per-entry processing ----

process_entry() {
  local entry_index="$1" project="$2" branch="$3" local_dir="$4" rollback="$5" component="$6"
  local pdir="$REPO_ROOT/projects/$local_dir"
  local padded_index
  padded_index=$(printf "%04d" "$entry_index")
  local label="[$local_dir]"

  local section_file="$WORK_DIR/sections/${padded_index}.md"
  local summary_line="$WORK_DIR/summaries/${padded_index}.txt"
  local refs_file="$WORK_DIR/refs/${local_dir}.txt"

  echo "$label Processing $branch..." >&2

  if [ ! -e "$pdir/.git" ]; then
    echo "$label [SKIP] Not found or not a git repo" >&2
    echo "$project|$branch|0|0|skipped" > "$summary_line"
    touch "$section_file"
    return
  fi

  if ! git -C "$pdir" diff --quiet 2>/dev/null || ! git -C "$pdir" diff --cached --quiet 2>/dev/null; then
    echo "$label [SKIP] Uncommitted changes — refusing to modify working tree" >&2
    echo "$project|$branch|0|0|dirty" > "$summary_line"
    touch "$section_file"
    return
  fi

  # Save original ref (one file per project, sequential within project so no race)
  if [ ! -s "$refs_file" ]; then
    local ref
    ref=$(git -C "$pdir" symbolic-ref --short HEAD 2>/dev/null || git -C "$pdir" rev-parse HEAD)
    echo "${pdir}|${ref}" > "$refs_file"
  fi

  local remote
  remote=$(resolve_remote "$pdir" "$branch")
  if [ -z "$remote" ]; then
    echo "$label [SKIP] Cannot find remote for branch $branch" >&2
    echo "$project|$branch|0|0|no-remote" > "$summary_line"
    touch "$section_file"
    return
  fi

  echo "$label Fetching $remote/$branch..." >&2
  if ! git -C "$pdir" fetch "$remote" "$branch" 2>/dev/null; then
    echo "$label [SKIP] Failed to fetch $remote/$branch" >&2
    echo "$project|$branch|0|0|fetch-failed" > "$summary_line"
    touch "$section_file"
    return
  fi

  if ! git -C "$pdir" rev-parse "$remote/$branch~$rollback" &>/dev/null; then
    echo "$label [SKIP] Branch $branch has fewer than $rollback commits" >&2
    echo "$project|$branch|0|0|too-few-commits" > "$summary_line"
    touch "$section_file"
    return
  fi

  local etmp="$WORK_DIR/${local_dir}_${branch}"
  mkdir -p "$etmp"

  # --- Capture commit log for the range ---
  local commits_file="$etmp/commits.txt"
  git -C "$pdir" log --oneline "$remote/$branch~$rollback..$remote/$branch" > "$commits_file" 2>/dev/null || true

  # --- Checkout HEAD (tip of branch) ---
  git -C "$pdir" checkout "$remote/$branch" --detach 2>/dev/null
  local frontend_dir
  frontend_dir=$(detect_frontend_dir "$pdir")

  local npm_head="$etmp/npm_head.json"
  touch "$npm_head"
  if [ -n "$frontend_dir" ]; then
    echo "$label [npm] Auditing $branch at HEAD..." >&2
    (cd "$frontend_dir" && npm audit --json 2>/dev/null) > "$npm_head" || true
  fi

  local go_head="$etmp/go_head.json"
  touch "$go_head"
  if [ -f "$pdir/go.mod" ] && [ "$HAS_GOVULNCHECK" = true ]; then
    echo "$label [go] govulncheck $branch at HEAD..." >&2
    (cd "$pdir" && govulncheck -json ./... 2>/dev/null) > "$go_head" || true
  fi

  # --- Checkout HEAD~N ---
  git -C "$pdir" checkout -- . 2>/dev/null || true
  git -C "$pdir" checkout "$remote/$branch~$rollback" --detach 2>/dev/null
  local frontend_dir_old
  frontend_dir_old=$(detect_frontend_dir "$pdir")

  local npm_old="$etmp/npm_old.json"
  touch "$npm_old"
  if [ -n "$frontend_dir_old" ]; then
    echo "$label [npm] Auditing $branch at HEAD~$rollback..." >&2
    (cd "$frontend_dir_old" && npm install >/dev/null 2>&1 && npm audit --json 2>/dev/null) > "$npm_old" || true
  fi

  local go_old="$etmp/go_old.json"
  touch "$go_old"
  if [ -f "$pdir/go.mod" ] && [ "$HAS_GOVULNCHECK" = true ]; then
    echo "$label [go] govulncheck $branch at HEAD~$rollback..." >&2
    (cd "$pdir" && govulncheck -json ./... 2>/dev/null) > "$go_old" || true
  fi

  # --- Restore ---
  git -C "$pdir" checkout -- . 2>/dev/null || true
  local original_ref
  original_ref=$(grep "^${pdir}|" "$refs_file" 2>/dev/null | head -1 | cut -d'|' -f2)
  git -C "$pdir" checkout "$original_ref" 2>/dev/null || true

  # --- Diff NPM ---
  local urls_head="$etmp/urls_head.txt" urls_old="$etmp/urls_old.txt" fixed_npm="$etmp/fixed_npm.txt"
  extract_npm_urls "$npm_head" > "$urls_head"
  extract_npm_urls "$npm_old" > "$urls_old"
  comm -23 "$urls_old" "$urls_head" > "$fixed_npm" 2>/dev/null || true
  local npm_count
  npm_count=$(wc -l < "$fixed_npm" | tr -d ' ')

  # --- Diff Go ---
  local ids_head="$etmp/ids_head.txt" ids_old="$etmp/ids_old.txt" fixed_go="$etmp/fixed_go.txt"
  extract_go_ids "$go_head" > "$ids_head"
  extract_go_ids "$go_old" > "$ids_old"
  comm -23 "$ids_old" "$ids_head" > "$fixed_go" 2>/dev/null || true
  local go_count
  go_count=$(wc -l < "$fixed_go" | tr -d ' ')

  echo "$label $branch: fixed $npm_count npm, $go_count go" >&2

  # --- Write report section ---
  {
    echo "## $project ($branch)"
    echo ""

    echo "### Commits analyzed"
    echo ""
    if [ -s "$commits_file" ]; then
      while IFS= read -r commit_line; do
        echo "- \`${commit_line}\`"
      done < "$commits_file"
    else
      echo "No commits in range."
    fi
    echo ""

    echo "### NPM Vulnerabilities Fixed"
    echo ""
    if [ -z "$frontend_dir" ]; then
      echo "N/A — no frontend in this project."
    elif [ "$npm_count" -gt 0 ]; then
      echo "| Advisory | CVE | Package | Severity | Title |"
      echo "| -------- | --- | ------- | -------- | ----- |"
      while IFS= read -r url; do
        local detail pkg sev title cve
        detail=$(extract_npm_detail "$npm_old" "$url")
        pkg=$(echo "$detail" | cut -f1)
        sev=$(echo "$detail" | cut -f2)
        title=$(echo "$detail" | cut -f3)
        cve=$(fetch_cve_for_ghsa "$url")
        echo "${url}|${cve}" >> "$etmp/cve_cache.txt"
        echo "| $url | $cve | $pkg | $sev | $title |"
      done < "$fixed_npm"
    else
      echo "No NPM vulnerabilities were fixed."
    fi
    echo ""

    echo "### Go Vulnerabilities Fixed"
    echo ""
    if [ ! -f "$pdir/go.mod" ]; then
      echo "N/A — no go.mod in this project."
    elif [ "$HAS_GOVULNCHECK" = false ]; then
      echo "Skipped — govulncheck not installed."
    elif [ "$go_count" -gt 0 ]; then
      echo "| ID | CVE/Aliases | Module | Summary |"
      echo "| -- | ----------- | ------ | ------- |"
      while IFS= read -r vid; do
        local detail aliases module summary
        detail=$(extract_go_detail "$go_old" "$vid")
        aliases=$(echo "$detail" | cut -f2)
        module=$(echo "$detail" | cut -f3)
        summary=$(echo "$detail" | cut -f4)
        echo "| $vid | $aliases | $module | $summary |"
      done < "$fixed_go"
    else
      echo "No Go vulnerabilities were fixed."
    fi
    echo ""
  } > "$section_file"

  echo "$project|$branch|$npm_count|$go_count|ok" > "$summary_line"

  # --- YAML fragment ---
  if [ "$YAML_OUTPUT" = true ]; then
    local yaml_file="$WORK_DIR/yaml/${padded_index}.yaml"
    {
      # NPM CVEs from cache
      if [ -f "$etmp/cve_cache.txt" ]; then
        while IFS='|' read -r _url cve; do
          if [ "$cve" != "N/A" ] && [ -n "$cve" ]; then
            echo "  - key: $cve"
            echo "    component: $component"
          fi
        done < "$etmp/cve_cache.txt"
      fi
      # Go CVEs from aliases
      if [ -s "$fixed_go" ]; then
        while IFS= read -r vid; do
          local detail aliases
          detail=$(extract_go_detail "$go_old" "$vid")
          aliases=$(echo "$detail" | cut -f2)
          IFS=', ' read -ra cve_list <<< "$aliases"
          for cve in "${cve_list[@]}"; do
            cve=$(echo "$cve" | tr -d ' ')
            if [[ "$cve" == CVE-* ]]; then
              echo "  - key: $cve"
              echo "    component: $component"
            fi
          done
        done < "$fixed_go"
      fi
    } > "$yaml_file"
  fi
}

# ---- Per-project processing (sequential within project) ----

process_project() {
  local target_project="$1"
  for i in "${!ENTRIES[@]}"; do
    local entry="${ENTRIES[$i]}"
    local project="${entry%%:*}"
    local branch="${entry##*:}"

    if [ "$project" != "$target_project" ]; then
      continue
    fi

    if [ -n "$FILTER" ] && [[ "$project" != *"$FILTER"* ]]; then
      continue
    fi

    process_entry "$i" "$project" "$branch" "$project" "${ENTRY_ROLLBACKS[$i]}" "${ENTRY_COMPONENTS[$i]}"
  done
}

# ---- Main: parallel dispatch ----

# Build unique project list (preserving order of first appearance)
unique_projects=()
seen_projects=""
for entry in "${ENTRIES[@]}"; do
  project="${entry%%:*}"
  if [[ "$seen_projects" != *"|${project}|"* ]]; then
    unique_projects+=("$project")
    seen_projects="${seen_projects}|${project}|"
  fi
done

echo "Starting parallel processing of ${#unique_projects[@]} projects..." >&2

pids=()
for project in "${unique_projects[@]}"; do
  if [ -n "$FILTER" ] && [[ "$project" != *"$FILTER"* ]]; then
    continue
  fi
  process_project "$project" &
  pids+=($!)
done

# Wait for all background jobs
failed=0
for pid in "${pids[@]+"${pids[@]}"}"; do
  if ! wait "$pid"; then
    failed=$((failed + 1))
  fi
done

if [ "$failed" -gt 0 ]; then
  echo "[WARN] $failed project group(s) had errors" >&2
fi

# ---- Assemble report ----

{
  echo "# Fixed CVEs Report"
  echo ""
  echo "Generated: $(date +%Y-%m-%d)"
  echo ""
} > "$REPORT"

# Concatenate sections in original entry order
for section in "$WORK_DIR"/sections/*.md; do
  [ -f "$section" ] || continue
  [ -s "$section" ] || continue
  cat "$section" >> "$REPORT"
done

# Summary table
{
  echo "## Summary"
  echo ""
  echo "| Project | Branch | NPM Fixed | Go Fixed | Total | Status |"
  echo "| ------- | ------ | --------- | -------- | ----- | ------ |"
  for summary in "$WORK_DIR"/summaries/*.txt; do
    [ -f "$summary" ] || continue
    while IFS='|' read -r p b n g s; do
      total=$((n + g))
      echo "| $p | $b | $n | $g | $total | $s |"
    done < "$summary"
  done
  echo ""
} >> "$REPORT"

# ---- Assemble YAML (if requested) ----

if [ "$YAML_OUTPUT" = true ]; then
  YAML_REPORT="$SCRIPT_DIR/report.yaml"
  echo "cves:" > "$YAML_REPORT"
  for yf in "$WORK_DIR"/yaml/*.yaml; do
    [ -f "$yf" ] || continue
    [ -s "$yf" ] || continue
    cat "$yf" >> "$YAML_REPORT"
  done
  echo "YAML saved to $YAML_REPORT" >&2
fi

echo "" >&2
echo "Report saved to $REPORT" >&2
