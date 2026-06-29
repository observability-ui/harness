# Plan: Script to Extract Fixed CVEs

## Problem

After dependency update commits are cherry-picked to release branches, we need a way to identify which CVEs were fixed by those updates. Currently
this requires manually running `npm audit` and `govulncheck` before and after the commits, then comparing the results — tedious and error-prone across
14 project/branch combinations. The script automates this: for each project+branch pair, it diffs the vulnerability state at HEAD vs HEAD~3 and
produces a markdown report listing the CVEs that were resolved.

## Current State

| Component                    | File / Location                          | Current Behavior                                                                                        |
| ---------------------------- | ---------------------------------------- | ------------------------------------------------------------------------------------------------------- |
| perses-operator              | `projects/perses-operator/`              | Go-only project. Branch `release-coo-1.5` on `rhobs` remote. No frontend.                               |
| monitoring-plugin            | `projects/monitoring-plugin/`            | Go backend + frontend in `web/`. Release branches `release-coo-ocp-4.{15,19,22}` on `upstream` remote.  |
| logging-view-plugin          | `projects/logging-view-plugin/`          | Go backend + frontend in `web/`. Release branches `release-coo-ocp-4.{12,15,22}` on `upstream` remote.  |
| distributed-tracing-plugin   | `projects/distributed-tracing-plugin/`   | Go backend + frontend in `web/`. Release branches `release-coo-ocp-4.{12,15,19,22}` on `origin` remote. |
| troubleshooting-panel-plugin | `projects/troubleshooting-panel-plugin/` | Go backend + frontend in `web/`. Release branches `release-coo-ocp-4.{19,22}` on `origin` remote.       |
| Script location              | `tasks/script-to-extract-fixed-cves/`    | Only `spec.md` exists. No script yet.                                                                   |

### Remote/Branch Mapping

The script must resolve the correct remote for each branch since projects use different remote names:

| Project                      | Remote     | Branch pattern        |
| ---------------------------- | ---------- | --------------------- |
| perses-operator              | `rhobs`    | `release-coo-1.5`     |
| monitoring-plugin            | `upstream` | `release-coo-ocp-4.*` |
| logging-view-plugin          | `upstream` | `release-coo-ocp-4.*` |
| distributed-tracing-plugin   | `origin`   | `release-coo-ocp-4.*` |
| troubleshooting-panel-plugin | `origin`   | `release-coo-ocp-4.*` |

### Frontend Location

All UI plugins have their `package.json` in `web/`, not the project root. The perses-operator has no frontend at all. The script must auto-detect the
frontend directory by checking for `web/package.json` first, then falling back to root `package.json`.

## Changes

### Phase 1: Core Script — Git Operations and Project Loop

**Dependency:** None **Parallel with:** None

#### Files Modified

| File                                                       | Change                      |
| ---------------------------------------------------------- | --------------------------- |
| `tasks/script-to-extract-fixed-cves/extract-fixed-cves.sh` | New file — main bash script |

#### Details

Create a bash script that:

1. **Defines the project/branch configuration as an array:**

```bash
ENTRIES=(
  "perses-operator:release-coo-1.5"
  "monitoring-plugin:release-coo-ocp-4.15"
  "monitoring-plugin:release-coo-ocp-4.19"
  "monitoring-plugin:release-coo-ocp-4.22"
  "logging-view-plugin:release-coo-ocp-4.12"
  "logging-view-plugin:release-coo-ocp-4.15"
  "logging-view-plugin:release-coo-ocp-4.22"
  "distributed-tracing-console-plugin:release-coo-ocp-4.12"
  "distributed-tracing-console-plugin:release-coo-ocp-4.15"
  "distributed-tracing-console-plugin:release-coo-ocp-4.19"
  "distributed-tracing-console-plugin:release-coo-ocp-4.22"
  "troubleshooting-panel-console-plugin:release-coo-ocp-4.19"
  "troubleshooting-panel-console-plugin:release-coo-ocp-4.22"
)
```

Note: The spec uses repository names (e.g., `distributed-tracing-console-plugin`) but the local submodule directories use shorter names (e.g.,
`distributed-tracing-plugin`). The script must map repo names → local directory names:

```bash
get_local_dir() {
  case "$1" in
    distributed-tracing-console-plugin) echo "distributed-tracing-plugin" ;;
    troubleshooting-panel-console-plugin) echo "troubleshooting-panel-plugin" ;;
    *) echo "$1" ;;
  esac
}
```

2. **For each entry, resolves the remote that tracks the branch:**

```bash
resolve_remote() {
  local project_dir="$1" branch="$2"
  git -C "$project_dir" branch -r | grep -E "/${branch}$" | head -1 | cut -d'/' -f1 | tr -d ' '
}
```

3. **Fetches the remote, checks out the branch at detached HEAD, and records the original ref for restoration:**

```bash
original_ref=$(git -C "$project_dir" rev-parse HEAD)
git -C "$project_dir" fetch "$remote" "$branch"
git -C "$project_dir" checkout "$remote/$branch" --detach
```

4. **After running audits (Phase 2), restores the original state:**

```bash
git -C "$project_dir" checkout "$original_ref" --detach
git -C "$project_dir" checkout -  # or restore to original branch
```

#### Phase 1 Verification

- Run `bash -n tasks/script-to-extract-fixed-cves/extract-fixed-cves.sh` to verify syntax
- Verify the remote resolution logic works: `git -C ./projects/monitoring-plugin branch -r | grep -E "/release-coo-ocp-4.15$"`

---

### Phase 2: NPM Audit Diff Logic

**Dependency:** Phase 1 **Parallel with:** Phase 3 (different concern, same file)

#### Files Modified

| File                                                       | Change                       |
| ---------------------------------------------------------- | ---------------------------- |
| `tasks/script-to-extract-fixed-cves/extract-fixed-cves.sh` | Add npm audit diff functions |

#### Details

Add functions to detect the frontend directory, run `npm audit` at two points, and diff the results.

##### Frontend Directory Detection

```bash
detect_frontend_dir() {
  local project_dir="$1"
  if [ -f "$project_dir/web/package.json" ]; then
    echo "$project_dir/web"
  elif [ -f "$project_dir/package.json" ]; then
    echo "$project_dir"
  else
    echo ""
  fi
}
```

##### NPM Audit Diff

```bash
run_npm_audit() {
  local frontend_dir="$1" output_file="$2"
  (cd "$frontend_dir" && npm install --ignore-scripts 2>/dev/null && npm audit --json 2>/dev/null) > "$output_file"
}

extract_npm_cves() {
  local audit_json="$1"
  # npm audit --json outputs vulnerabilities keyed by package name
  # Each vulnerability has a "via" array containing objects with "url" fields (advisory URLs containing CVE refs)
  jq -r '.vulnerabilities | to_entries[] | .value.via[]? | 
    if type == "object" then .url // empty else empty end' "$audit_json" | 
    sort -u
}
```

The diff logic:

1. At HEAD: run `npm audit --json` → `audit_head.json`
2. At HEAD~3: run `npm install --ignore-scripts` then `npm audit --json` → `audit_head3.json`
3. Extract advisory URLs from both
4. Use `comm -23` to find URLs in HEAD~3 but not in HEAD (= fixed)

For each fixed advisory URL, also extract the CVE ID, severity, package name, and vulnerability title from the HEAD~3 audit JSON for the report.

#### Phase 2 Verification

- Test against monitoring-plugin on `release-coo-ocp-4.22` to verify npm audit output parsing
- Verify `jq` correctly extracts advisory URLs from npm audit JSON

---

### Phase 3: govulncheck Diff Logic

**Dependency:** Phase 1 **Parallel with:** Phase 2 (different concern, same file)

#### Files Modified

| File                                                       | Change                         |
| ---------------------------------------------------------- | ------------------------------ |
| `tasks/script-to-extract-fixed-cves/extract-fixed-cves.sh` | Add govulncheck diff functions |

#### Details

Add functions to run `govulncheck` at two points and diff the results. This applies to all projects since they all have `go.mod`.

##### govulncheck Diff

```bash
run_govulncheck() {
  local project_dir="$1" output_file="$2"
  (cd "$project_dir" && govulncheck -json ./... 2>/dev/null) > "$output_file"
}

extract_go_vulns() {
  local vuln_json="$1"
  # govulncheck -json outputs one JSON object per line with "finding" entries
  # Each finding has an "osv" field with the vulnerability ID (e.g., GO-2024-3321)
  jq -r 'select(.finding != null) | .finding.osv' "$vuln_json" | sort -u
}
```

The diff logic mirrors NPM:

1. At HEAD: run `govulncheck -json ./...` → `govulncheck_head.json`
2. At HEAD~3: run `govulncheck -json ./...` → `govulncheck_head3.json`
3. Extract vulnerability IDs from both
4. Use `comm -23` to find IDs in HEAD~3 but not in HEAD (= fixed)

For each fixed vulnerability, also extract the aliases (CVE IDs), affected module, and summary from the HEAD~3 output for the report. The
`govulncheck -json` output includes `osv` entries with full details:

```bash
extract_go_vuln_details() {
  local vuln_json="$1" vuln_id="$2"
  jq -r --arg id "$vuln_id" '
    select(.osv != null and .osv.id == $id) |
    .osv | {id, aliases: (.aliases // []), summary, affected: [.affected[]?.package.name]}
  ' "$vuln_json"
}
```

#### Phase 3 Verification

- Test against perses-operator on `release-coo-1.5` to verify govulncheck output parsing
- Verify `jq` correctly extracts OSV IDs from govulncheck JSON

---

### Phase 4: Markdown Report Generation

**Dependency:** Phase 2, Phase 3 **Parallel with:** None

#### Files Modified

| File                                                       | Change                         |
| ---------------------------------------------------------- | ------------------------------ |
| `tasks/script-to-extract-fixed-cves/extract-fixed-cves.sh` | Add report generation function |

#### Details

After processing all entries, generate a markdown report at `tasks/script-to-extract-fixed-cves/report.md` with the following structure:

```markdown
# Fixed CVEs Report

Generated: 2026-06-29

## monitoring-plugin (release-coo-ocp-4.15)

### NPM Vulnerabilities Fixed

| Advisory                                | CVE           | Package | Severity | Title               |
| --------------------------------------- | ------------- | ------- | -------- | ------------------- |
| https://github.com/advisories/GHSA-xxxx | CVE-2024-xxxx | lodash  | high     | Prototype Pollution |

### Go Vulnerabilities Fixed

| ID           | CVE           | Module           | Summary            |
| ------------ | ------------- | ---------------- | ------------------ |
| GO-2024-3321 | CVE-2024-xxxx | golang.org/x/net | HTTP/2 rapid reset |

## monitoring-plugin (release-coo-ocp-4.19)

...

## Summary

| Project           | Branch               | NPM Fixed | Go Fixed | Total |
| ----------------- | -------------------- | --------- | -------- | ----- |
| monitoring-plugin | release-coo-ocp-4.15 | 3         | 1        | 4     |
| ...               | ...                  | ...       | ...      | ...   |
```

The script writes the report incrementally: each project/branch appends its section, then the summary table is generated at the end from the collected
counts.

#### Phase 4 Verification

- Run the full script and verify the generated `report.md` is valid markdown
- Check that every project/branch entry from the spec appears in the report (even if zero CVEs fixed)

---

### Phase 5: Error Handling and Cleanup

**Dependency:** Phase 4 **Parallel with:** None

#### Files Modified

| File                                                       | Change                              |
| ---------------------------------------------------------- | ----------------------------------- |
| `tasks/script-to-extract-fixed-cves/extract-fixed-cves.sh` | Add cleanup trap and error handling |

#### Details

1. **Trap for cleanup**: Ensure the working directory of each submodule is restored even if the script fails mid-execution.

```bash
cleanup() {
  for project_dir in "${TOUCHED_DIRS[@]}"; do
    git -C "$project_dir" checkout - 2>/dev/null || true
    git -C "$project_dir" clean -fd 2>/dev/null || true
  done
  rm -rf "$TMPDIR"
}
trap cleanup EXIT
```

2. **Temp directory for intermediate files**: Use `mktemp -d` for all audit JSON files.

3. **Skip missing tools gracefully**: If `govulncheck` is not installed, log a warning and skip Go analysis (npm is assumed to be present). If `jq` is
   missing, fail with a clear error since it's required.

4. **Skip projects without Go/frontend**: If a project has no `go.mod`, skip govulncheck. If no `package.json` is found, skip npm audit.

5. **Progress logging**: Print progress to stderr so stdout remains clean for piping.

#### Phase 5 Verification

- Verify cleanup works by killing the script mid-run and checking that submodule working trees are restored
- Verify the script handles missing `govulncheck` gracefully

## PR Strategy

| PR | Repository | Branch | Description                                                                 | Dependencies |
| -- | ---------- | ------ | --------------------------------------------------------------------------- | ------------ |
| 1  | ai-sdlc    | main   | Add `extract-fixed-cves.sh` script in `tasks/script-to-extract-fixed-cves/` | None         |

Single PR since the script is self-contained in the task directory and doesn't modify any project submodule code.

## Verification

- **Rollback + npm audit**: Run the script against `monitoring-plugin:release-coo-ocp-4.22` and verify it detects CVEs that were fixed in the recent
  `update-vulnerable-dependencies-26-06-2026` commits.
- **Rollback + govulncheck**: Run the script against `perses-operator:release-coo-1.5` and verify Go vulnerability diff works.
- **Multiple projects in a loop**: Run the full script with all 13 entries and verify the report contains a section for each.
- **Frontend in different directory**: Verify the script correctly finds `web/package.json` for the UI plugins and skips frontend analysis for
  perses-operator.
- **Markdown report**: Verify the generated report is valid markdown with proper tables and summary.

## Risks

| Risk                                                                               | Impact                                          | Mitigation                                                                                 |
| ---------------------------------------------------------------------------------- | ----------------------------------------------- | ------------------------------------------------------------------------------------------ |
| `npm install` at HEAD~3 may fail if Node version is incompatible with old lockfile | Script aborts for that entry, no CVE data       | Use `--ignore-scripts` and `--legacy-peer-deps` flags; catch errors and report as warnings |
| `govulncheck` requires Go to be installed and modules to build                     | Missing Go vulns for some entries               | Check for `govulncheck` availability before running; skip with warning if missing          |
| Checking out branches modifies submodule working trees                             | Uncommitted changes in submodules could be lost | Save and restore original ref; refuse to run if there are uncommitted changes              |
| npm audit JSON format varies between npm versions                                  | Parsing breaks                                  | Use `npm audit --json` which has been stable since npm 7; document minimum npm version     |
| Large number of entries (13) makes full run slow due to `npm install` at each      | Script takes 10+ minutes                        | Add `--project` flag to run for a single project; show progress per entry                  |
