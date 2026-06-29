# Execution: Script to Extract Fixed CVEs

> Results are annotated inline: `-- **value**` for discovered values, `-- **passes/FAILED**` for verification.

## Phase 1: Core Script — Git Operations and Project Loop
Depends on: nothing | Parallel with: none | Type: implementation | Projects: ai-sdlc (tasks/)

- [x] Create `tasks/script-to-extract-fixed-cves/extract-fixed-cves.sh` with project/branch config array
- [x] Add repo-name → local-dir mapping function (`get_local_dir`)
- [x] Add remote resolution function (`resolve_remote`)
- [x] Add git checkout/restore logic with original ref tracking -- **uses temp file instead of associative array for bash 3.x compat**
- [x] Fix `.git` check to use `-e` instead of `-d` for submodule gitlinks
- [x] Add `git checkout -- .` before branch switches to reset npm install side effects

### Phase 1 Verification
- [x] `bash -n tasks/script-to-extract-fixed-cves/extract-fixed-cves.sh` — **passes**

## Phase 2: NPM Audit Diff Logic
Depends on: Phase 1 | Parallel with: none (same file as Phase 3) | Type: implementation | Projects: ai-sdlc (tasks/)

- [x] Add frontend directory detection function (`detect_frontend_dir`)
- [x] Add npm audit at HEAD and HEAD~N
- [x] Add `extract_npm_urls` function using jq to extract advisory URLs
- [x] Add `extract_npm_detail` function for report enrichment
- [x] Add diff logic: extract URLs from HEAD and HEAD~N, use `comm -23` to find fixed

### Phase 2 Verification
- [x] jq expression correctly parses npm audit JSON format -- **passes, tested against distributed-tracing and troubleshooting-panel**

## Phase 3: govulncheck Diff Logic
Depends on: Phase 1 | Parallel with: none (same file as Phase 2) | Type: implementation | Projects: ai-sdlc (tasks/)

- [x] Add govulncheck run at HEAD and HEAD~N
- [x] Add `extract_go_ids` function using jq to extract OSV IDs from finding entries
- [x] Add `extract_go_detail` function for report enrichment (aliases, module, summary)
- [x] Add diff logic mirroring NPM approach

### Phase 3 Verification
- [x] jq expression correctly parses govulncheck JSON format -- **passes, found 20-32 Go vulns fixed per branch**

## Phase 4: Markdown Report Generation
Depends on: Phase 2, Phase 3 | Parallel with: none | Type: implementation | Projects: ai-sdlc (tasks/)

- [x] Add report header generation with date
- [x] Add per-project/branch section generation (NPM + Go tables)
- [x] Add summary table generation at end with status column

### Phase 4 Verification
- [x] Generated report.md is valid markdown with proper tables -- **passes**

## Phase 5: Error Handling and Cleanup
Depends on: Phase 4 | Parallel with: none | Type: implementation | Projects: ai-sdlc (tasks/)

- [x] Add cleanup trap (restore git state with `checkout -- .` + ref restore, remove temp files)
- [x] Add temp directory creation (`mktemp -d`)
- [x] Add tool availability checks (jq required, govulncheck optional, npm required)
- [x] Add dirty-worktree guard (refuse to run with uncommitted changes)
- [x] Add progress logging to stderr
- [x] Add skip for too-few-commits, missing remote, fetch failure

### Phase 5 Verification
- [x] Script handles missing `govulncheck` gracefully -- **confirmed: HAS_GOVULNCHECK flag skips with warning**
- [x] Cleanup restores working trees correctly after normal completion -- **verified: projects back on main, clean**

## End-to-End Verification
- [x] Run against `troubleshooting-panel:release-coo-ocp-4.19` — **20 Go vulns fixed detected**
- [x] Run against `distributed-tracing:release-coo-ocp-4.{12,15,19,22}` — **32, 28, 28, 0 Go vulns fixed**
- [x] Report contains proper markdown tables with CVE IDs, modules, summaries
- [x] Frontend detected in `web/` for UI plugins
- [x] Projects restored to main branch with clean working tree after script completes
- [x] Full parallel run of all 13 entries — **5 min, 362% CPU, all entries ok**

---

## Phase 6: Parallel Execution by Project
> Added post-initial implementation to reduce wall-clock time from ~25 min to ~5 min.

- [x] Extract entry processing into `process_entry` function with index-based output files
- [x] Add `process_project` function to run entries for the same project sequentially
- [x] Replace main for-loop with parallel dispatch: one background job per unique project
- [x] Write report sections to `$WORK_DIR/sections/NNNN.md` and summaries to `$WORK_DIR/summaries/NNNN.txt` for ordered assembly
- [x] Move refs tracking to per-project files (`$WORK_DIR/refs/${local_dir}.txt`) to avoid concurrent writes
- [x] Update cleanup to iterate over `$WORK_DIR/refs/*.txt`
- [x] Prefix stderr with `[$local_dir]` for readable interleaved output
- [x] Fix `pids[@]` unbound variable when filter matches no projects

### Phase 6 Verification
- [x] Single-project test (`troubleshooting-panel`) — **passes, sequential within project**
- [x] Full 13-entry run — **5 min, 362% CPU, all entries ok, correct report order**
- [x] All projects restored to main with clean working trees

## Phase 7: Fix NPM Audit JSON Corruption
> `npm install` stdout (e.g., `up to date, audited 1700 packages`) was written into the same file as `npm audit --json`, corrupting the JSON. jq silently returned nothing, so all npm reports were empty.

- [x] Root cause: `(cd "$dir" && npm install 2>/dev/null && npm audit --json)` piped both commands' stdout into one file
- [x] Fix: redirect npm install stdout to `/dev/null` (`npm install >/dev/null 2>&1`)
- [x] Verified: monitoring-plugin 4.22 now reports 23 npm fixes (was 0)

### Phase 7 Verification
- [x] monitoring-plugin: 20, 20, 23 npm fixes on 4.15, 4.19, 4.22 — **passes**

## Phase 8: Remove `npm install` from HEAD Audit
> `npm audit` reads `package-lock.json` directly — `npm install` is unnecessary and can modify the lockfile, contaminating the HEAD vs HEAD~N comparison.

- [x] Remove `npm install` from the HEAD audit pipeline (keep only `npm audit --json`)
- [x] Remove `--ignore-scripts --legacy-peer-deps` flags (user request)
- [x] Verified: `npm audit --json` produces identical results with and without `node_modules` (11 URLs both ways)

### Phase 8 Verification
- [x] monitoring-plugin results consistent with manual `npm audit` — **passes**

## Phase 9: Map GHSA Advisories to CVE IDs
> NPM audit reports advisory URLs (GHSA-xxxx) but errata/release notes need CVE IDs. GitHub Advisory API returns the mapping via `gh api /advisories/GHSA-xxxx`.

- [x] Add `gh` CLI availability check (`HAS_GH` flag, graceful fallback to "N/A")
- [x] Add `fetch_cve_for_ghsa` function using `gh api` + `jq .cve_id`
- [x] Add CVE column to NPM report table: `| Advisory | CVE | Package | Severity | Title |`
- [x] Rate limit confirmed: 5000 req/hr, worst case ~520 calls — **no concern**

### Phase 9 Verification
- [x] troubleshooting-panel 4.19: 48 npm advisories, all CVE-mapped — **passes**
- [x] Advisories without CVE show "N/A" (e.g., GHSA-442j, GHSA-7rx3) — **passes**

## Phase 10: Add Commit Messages to Report
> Each report section now includes the commits being analyzed (HEAD~N..HEAD), giving context on what changes introduced the fixes.

- [x] Capture `git log --oneline $remote/$branch~$ROLLBACK..$remote/$branch` before checkout
- [x] Add "Commits analyzed" subsection at top of each project/branch report section

### Phase 10 Verification
- [x] perses-operator report shows 5 commit messages — **passes**

## Phase 11: YAML Output with Konflux Component Mapping
> Errata tooling needs a YAML file mapping each fixed CVE to its Konflux component name. Component names come from a hardcoded mapping derived from `component-mapping.md`.

- [x] Add `--yaml` flag parsing (coexists with optional filter argument)
- [x] Add `get_component` function: hardcoded `(project, branch) → component-name` mapping for all 13 entries
- [x] Cache CVE lookups in `$etmp/cve_cache.txt` during report generation to avoid duplicate API calls
- [x] Write per-entry YAML fragments to `$WORK_DIR/yaml/NNNN.yaml`
- [x] Assemble `report.yaml` from fragments in entry order after all jobs complete
- [x] Go CVEs: extract CVE-* aliases from `extract_go_detail` output
- [x] NPM CVEs: read from cache file (populated during report table generation)

### Phase 11 Verification
- [x] `--yaml perses-operator`: 20 Go CVEs mapped to `perses-operator-1-5` — **passes**
- [x] YAML format matches spec: `cves:\n  - key: CVE-xxx\n    component: xxx` — **passes**
- [x] Without `--yaml`: no report.yaml generated — **passes**

---

## Summary

**Status:** Complete (all 11 phases done)

### Script capabilities
- Parallel processing of 5 projects (entries for same project run sequentially)
- Diffs npm audit + govulncheck at HEAD vs HEAD~N to find fixed CVEs
- Markdown report with commit messages, NPM table (with CVE via GitHub API), Go table, and summary
- Optional `--yaml` output mapping CVEs to Konflux component names
- Graceful handling of missing tools, dirty worktrees, missing branches
- Cleanup trap restores all project working trees on exit

### Usage
```bash
./extract-fixed-cves.sh                          # all projects, markdown only
./extract-fixed-cves.sh monitoring-plugin         # filter by project name
./extract-fixed-cves.sh --yaml                    # all projects, markdown + YAML
./extract-fixed-cves.sh --yaml troubleshooting    # filter + YAML
```
