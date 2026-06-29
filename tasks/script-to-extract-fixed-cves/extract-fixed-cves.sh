#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
CONFIG_FILE="$SCRIPT_DIR/config.yaml"
REPORT="$SCRIPT_DIR/report.md"
WORK_DIR=$(mktemp -d)

YAML_OUTPUT=true
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

# ---- Pre-flight: validate all configured project directories ----

validation_errors=0
seen_invalid_projects=""
for i in "${!ENTRIES[@]}"; do
  entry="${ENTRIES[$i]}"
  project="${entry%%:*}"
  branch="${entry##*:}"
  [ -n "$FILTER" ] && [[ "$project" != *"$FILTER"* ]] && continue
  pdir="$REPO_ROOT/projects/$project"
  if [ ! -e "$pdir/.git" ]; then
    echo "[ERROR] [$project @ $branch] Project directory is not a git repo: $pdir" >&2
    # Print the fix hint once per unique project to avoid repetition
    if [[ "$seen_invalid_projects" != *"|${project}|"* ]]; then
      echo "        Fix: git submodule update --init projects/$project" >&2
      seen_invalid_projects="${seen_invalid_projects}|${project}|"
    fi
    validation_errors=$((validation_errors + 1))
  fi
done
if [ "$validation_errors" -gt 0 ]; then
  echo "[ERROR] $validation_errors entry/entries failed pre-flight validation — aborting." >&2
  exit 1
fi

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
  # Try to find the remote by scanning already-fetched tracking branches.
  # grep exits 1 when there is no match, so guard with || true to avoid a
  # silent set -e exit before the caller can print a useful error message.
  local remote
  remote=$(git -C "$pdir" branch -r | grep -E "/${branch}$" | head -1 | sed 's/^[[:space:]]*//' | cut -d'/' -f1 || true)
  if [ -z "$remote" ]; then
    # Branch not fetched yet — fall back to the repo's first configured remote.
    remote=$(git -C "$pdir" remote 2>/dev/null | head -1 || true)
  fi
  echo "$remote"
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

# Returns the lowercase severity tier (critical|high|medium|low|unknown) for a Go OSV entry.
# Reads from database_specific.severity first; falls back to parsing the CVSS base score
# embedded in the vector string via awk when that field is absent.
extract_go_severity() {
  local json_file="$1" vid="$2"
  [ -s "$json_file" ] || { echo "unknown"; return; }

  local tier
  tier=$(jq -r --arg id "$vid" '
    select(.osv != null and .osv.id == $id) |
    # Prefer an explicit tier stored by the Go vuln DB
    if .osv.database_specific.severity != null then
      .osv.database_specific.severity | ascii_downcase
    # Fall back to CVSS vector — emit the numeric base score so awk can bucket it
    elif (.osv.severity // [] | map(select(.type == "CVSS_V3" or .type == "CVSS_V3_1")) | length) > 0 then
      (.osv.severity[] | select(.type == "CVSS_V3" or .type == "CVSS_V3_1") | .score)
    else
      "unknown"
    end
  ' "$json_file" 2>/dev/null | head -1)

  case "$tier" in
    critical|high|medium|low|unknown) echo "$tier" ;;
    # Received a raw CVSS vector — extract the base score digit(s) after "CVSS:x.x/"
    # The Go vuln DB stores pre-computed base scores in database_specific; this path
    # handles the rare case where only the vector string is available.
    CVSS:*)
      local base_score
      base_score=$(echo "$tier" | awk -F'/' '{
        for (i=1; i<=NF; i++) {
          if ($i ~ /^BS:/) { gsub("BS:", "", $i); print $i; exit }
        }
        print "0"
      }')
      # If we could not extract BS: field, default to unknown (include it to be safe)
      if [[ "$base_score" == "0" ]]; then
        echo "unknown"
      else
        awk -v s="$base_score" 'BEGIN {
          if      (s+0 >= 9.0) print "critical"
          else if (s+0 >= 7.0) print "high"
          else if (s+0 >= 4.0) print "medium"
          else                  print "low"
        }'
      fi
      ;;
    *) echo "unknown" ;;
  esac
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
    echo "$label [ERROR] Not found or not a git repo: $pdir" >&2
    exit 1
  fi

  if ! git -C "$pdir" diff --quiet 2>/dev/null || ! git -C "$pdir" diff --cached --quiet 2>/dev/null; then
    echo "$label [ERROR] Uncommitted changes in $pdir — refusing to modify working tree" >&2
    exit 1
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
    echo "$label [ERROR] Cannot find remote for branch $branch in $pdir" >&2
    exit 1
  fi

  echo "$label Fetching $remote/$branch..." >&2
  if ! git -C "$pdir" fetch "$remote" "$branch" 2>/dev/null; then
    echo "$label [ERROR] Failed to fetch $remote/$branch" >&2
    exit 1
  fi

  if ! git -C "$pdir" rev-parse "$remote/$branch~$rollback" &>/dev/null; then
    echo "$label [ERROR] Branch $branch has fewer than $rollback commits (need $rollback)" >&2
    exit 1
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

  # --- Filter NPM by severity (high/critical only) ---
  # Intermediate files are pre-populated here so the section-writing block
  # and the YAML fragment can simply cat/read them without re-looping.
  local npm_table_file="$etmp/npm_table.md"
  local npm_high_count=0
  touch "$npm_table_file"
  if [ -n "$frontend_dir" ] && [ "$npm_count" -gt 0 ]; then
    while IFS= read -r url; do
      local detail pkg sev title cve
      detail=$(extract_npm_detail "$npm_old" "$url")
      pkg=$(echo "$detail" | cut -f1)
      sev=$(echo "$detail" | cut -f2)
      title=$(echo "$detail" | cut -f3)
      # Skip anything below high
      [[ "$sev" == "high" || "$sev" == "critical" ]] || continue
      cve=$(fetch_cve_for_ghsa "$url")
      echo "${url}|${cve}" >> "$etmp/cve_cache.txt"
      printf '| %s | %s | %s | %s | %s |\n' "$url" "$cve" "$pkg" "$sev" "$title" \
        >> "$npm_table_file"
      npm_high_count=$((npm_high_count + 1))
    done < "$fixed_npm"
  fi

  # --- Filter Go by severity (high/critical only) ---
  local go_table_file="$etmp/go_table.md"
  local go_aliases_file="$etmp/go_aliases.txt"
  local go_high_count=0
  touch "$go_table_file" "$go_aliases_file"
  if [ -f "$pdir/go.mod" ] && [ "$HAS_GOVULNCHECK" = true ] && [ "$go_count" -gt 0 ]; then
    while IFS= read -r vid; do
      local detail aliases module summary sev_tier display_sev
      detail=$(extract_go_detail "$go_old" "$vid")
      aliases=$(echo "$detail" | cut -f2)
      module=$(echo "$detail" | cut -f3)
      summary=$(echo "$detail" | cut -f4)
      sev_tier=$(extract_go_severity "$go_old" "$vid")
      # Skip medium and low; keep high, critical, and unknown (unknown = no data, include to be safe)
      [[ "$sev_tier" == "medium" || "$sev_tier" == "low" ]] && continue
      display_sev="${sev_tier^^}"
      [[ "$display_sev" == "UNKNOWN" ]] && display_sev="N/A"
      printf '| %s | %s | %s | %s | %s |\n' "$vid" "$aliases" "$module" "$display_sev" "$summary" \
        >> "$go_table_file"
      echo "$aliases" >> "$go_aliases_file"
      go_high_count=$((go_high_count + 1))
    done < "$fixed_go"
  fi

  echo "$label $branch: fixed $npm_high_count npm (high/critical), $go_high_count go (high/critical)" >&2

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
    elif [ "$npm_high_count" -gt 0 ]; then
      echo "| Advisory | CVE | Package | Severity | Title |"
      echo "| -------- | --- | ------- | -------- | ----- |"
      cat "$npm_table_file"
    else
      echo "No high/critical NPM vulnerabilities were fixed."
    fi
    echo ""

    echo "### Go Vulnerabilities Fixed"
    echo ""
    if [ ! -f "$pdir/go.mod" ]; then
      echo "N/A — no go.mod in this project."
    elif [ "$HAS_GOVULNCHECK" = false ]; then
      echo "Skipped — govulncheck not installed."
    elif [ "$go_high_count" -gt 0 ]; then
      echo "| ID | CVE/Aliases | Module | Severity | Summary |"
      echo "| -- | ----------- | ------ | -------- | ------- |"
      cat "$go_table_file"
    else
      echo "No high/critical Go vulnerabilities were fixed."
    fi
    echo ""
  } > "$section_file"

  echo "$project|$branch|$npm_high_count|$go_high_count|ok" > "$summary_line"

  # --- YAML fragment ---
  if [ "$YAML_OUTPUT" = true ]; then
    local yaml_file="$WORK_DIR/yaml/${padded_index}.yaml"
    {
      # NPM CVEs — already filtered to high/critical via cve_cache.txt
      if [ -f "$etmp/cve_cache.txt" ]; then
        while IFS='|' read -r _url cve; do
          if [ "$cve" != "N/A" ] && [ -n "$cve" ]; then
            echo "  - key: $cve"
            echo "    component: $component"
          fi
        done < "$etmp/cve_cache.txt"
      fi
      # Go CVEs — already filtered to high/critical via go_aliases_file
      if [ -s "$go_aliases_file" ]; then
        while IFS= read -r aliases; do
          IFS=', ' read -ra cve_list <<< "$aliases"
          for cve in "${cve_list[@]}"; do
            cve=$(echo "$cve" | tr -d ' ')
            if [[ "$cve" == CVE-* ]]; then
              echo "  - key: $cve"
              echo "    component: $component"
            fi
          done
        done < "$go_aliases_file"
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
declare -A pid_project  # maps PID → project name for error reporting
for project in "${unique_projects[@]}"; do
  if [ -n "$FILTER" ] && [[ "$project" != *"$FILTER"* ]]; then
    continue
  fi
  process_project "$project" &
  pid=$!
  pids+=("$pid")
  pid_project[$pid]="$project"
done

# Wait for background jobs — kill all remaining and abort on the first failure
for pid in "${pids[@]+"${pids[@]}"}"; do
  if ! wait "$pid"; then
    echo "[ERROR] Project group '${pid_project[$pid]}' failed — killing remaining jobs and aborting." >&2
    for remaining in "${pids[@]+"${pids[@]}"}"; do
      kill "$remaining" 2>/dev/null || true
    done
    exit 1
  fi
done

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
