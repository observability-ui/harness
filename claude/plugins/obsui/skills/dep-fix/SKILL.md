---
name: dep-fix
description: Audit and fix vulnerable dependencies (npm + Go) across projects. Use when reviewing CVEs, running security audits, or patching dependency vulnerabilities before a release.
argument-hint: "<project-name> [CVE-ID|package-name]"
allowed-tools: Read, Write, Edit, Bash, Agent, AskUserQuestion
---

# /dep-fix

## Mission

Audit dependencies for known vulnerabilities and apply fixes with verified safety. Supports npm and Go ecosystems across all harness projects.

## Input

$ARGUMENTS must start with a project name (e.g., `monitoring-plugin`). Optionally followed by a CVE ID (e.g., `CVE-2024-1234`) or package name to
target a specific vulnerability. If no project name is provided, ask the user for one before proceeding.

## Error Handling

If any setup command fails (`git pull`, `npm install`, `go mod tidy`), stop and report the error to the user before proceeding. Do not attempt to
fix infrastructure issues — they require user intervention.

## Directory Navigation

The shell CWD persists between Bash tool calls. Preferred pattern: **cd into the target directory once** as a standalone call, then run subsequent
commands without cd. If the agent combines cd with a read-only command like `npm audit`, that is acceptable.

```bash
cd ./projects/<project>/web    # npm — directory containing package.json
cd ./projects/<project>        # Go — directory containing go.mod
```

Then run commands directly:

```bash
npm audit --json
govulncheck ./...
```

**git commands** — use `git -C`, no `cd` needed:

```bash
git -C ./projects/<project> branch --show-current
```

**File reads** — use relative paths from repo root, no `cd`, no leading `./`:

```
projects/<project>/web/package.json
```

## Allowed Commands

Subagents should use these commands. Commands not on this list may trigger permission prompts.

**npm:** `npm audit`, `npm audit fix`, `npm audit --json`, `npm ls`, `npm info`, `npm view`, `npm install`, `npm run`
**audit script:** `npm audit --json > .audit.json 2>/dev/null; echo done` then `node <skill-path>/scripts/audit-summary.js .audit.json`
**Go:** `go get`, `go mod tidy`, `go build`, `go test`, `go list`, `govulncheck`
**git:** `git -C <path> <command>` (always with `-C`)
**shell:** `grep`, `head`, `cat`, `find`

**BANNED — dangerous or unnecessary:**
- `npm audit fix --force` — may downgrade packages to ancient versions
- `--legacy-peer-deps` — silently ignores peer dependency conflicts
- `npm --prefix` — use cd instead
- `node -e`, `node -p` — use `npm ls` to check versions
- `python`, `source`, `nvm` — subagents do not manage runtimes

## Process

### Phase 1: Discovery

1. **Detect ecosystem and locate files**:

```bash
find ./projects/<project> -maxdepth 2 \( -name go.mod -o -name package.json -o -name .nvmrc -o -name .node-version -o -name "Dockerfile*" \) 2>/dev/null
```

   From the output, identify the directory containing `package.json` (npm dir) and `go.mod` (Go dir). These may be in the project root OR a
   subdirectory (e.g., `web/`, `frontend/`, `ui/`). If multiple `package.json` files are found (e.g., `web/` and `tests/`), focus on the main
   application package. Report additional locations to the user.

   Locate all Dockerfiles. Categorize:
   - **Modifiable**: any Dockerfile variant EXCEPT `.art`
   - **NEVER modify**: `Dockerfile.art` — managed by another team

2. **Check prerequisites** — verify only the tools needed for the detected ecosystems:
   - If npm detected: `node --version && npm --version`
   - If Go detected: `govulncheck -version`
   - If a required tool is missing, stop and tell the user what to install.

3. **Switch runtime versions** — if `.nvmrc` or `.node-version` was found, tell the user to run `nvm use` in their terminal. If NOT found, use
   whatever node/npm is in PATH — report the version. For Go, check `.go-version` or the `go` directive in `go.mod`.

4. **Detect branch** — report to user before proceeding:

```bash
git -C ./projects/<project> branch --show-current
```

5. **Update branch and install dependencies**:

```bash
git -C ./projects/<project> pull --ff-only
```

   If pull fails (no tracking info), skip and warn the user.

   - **npm**: cd to the directory containing `package.json` (discovered in step 1 — e.g., `web/`, `frontend/`), then `npm install`
   - **Go**: cd to the directory containing `go.mod`, then `go mod tidy`

6. **Discover verification commands** — use the Read tool to read `projects/<project>/Makefile`, `package.json` scripts, `CLAUDE.md`, and `AGENTS.md`.
   Identify build, lint, and test commands for both ecosystems. Record for Phase 3.

7. **Dispatch ecosystem agents** — if the project has BOTH npm and Go, spawn two subagents in parallel. If only one ecosystem, run Phase 2 directly.

   Each subagent receives:
   - The ecosystem (npm or Go) and its resolved directory path
   - The CVE or package filter (if any)
   - The audit-summary script path (npm only)
   - The allowed commands and banned commands from this skill
   - The verification commands from step 6
   - **Navigation**: "cd to the directory as your first call. Then run all commands without cd."

   Each subagent returns:
   1. Vulnerabilities found (table)
   2. Fixes applied (old → new versions)
   3. Recommendations needing user confirmation (overrides, Go version upgrades)
   4. Skipped vulnerabilities with reason
   5. Remaining vulnerability count

   Subagents must NOT apply overrides (npm) or Go version upgrades — return as recommendations only.
   Go subagents must check `go` directive before and after each `go get` — flag implicit bumps.

### Phase 2: Audit and Fix

#### npm audit and fix

1. **Audit, fix, re-audit** — run as separate Bash calls:

   a. `npm audit --json > .audit.json 2>/dev/null; echo done`
   b. `node <skill-path>/scripts/audit-summary.js .audit.json` — produces a table with: package, declared version, installed version, vuln range,
      severity, prod/dev, parent chain, fix version, advisory IDs. If installed version is outside the vuln range, it's marked as a false positive.
   c. `npm audit fix` (NEVER `--force`)
   d. Re-run steps a-b to see what remains.
   e. When analysis is complete, clean up: `rm -f .audit.json`

2. **For remaining vulnerabilities**, apply fixes in escalation order:

   a. **Direct update**: vulnerable package is a direct dep → `npm info <pkg> versions` to find a patched version → update `package.json` +
      `npm install`. If blocked by peer deps, note the constraining package.

   b. **Parent dependency update**: transitive dep → use `npm ls <pkg> --depth=3` (only if the audit summary lacks detail) to find the parent →
      check if parent has an update that pulls in the fix.

   c. **Override**: last resort. If the vuln is dev-only, recommend skipping — it doesn't affect production. If production, recommend an `overrides`
      entry but DO NOT APPLY without user confirmation.

   The dev-only skip applies ONLY to overrides (step c). Dev deps should still be fixed through steps a-b and npm audit fix.

3. **Align declared versions**: when `package.json` was modified (steps a or b), bump the base version to match the installed version
   (e.g., `^4.17.11` → `^4.17.21`). Always use `^x.y.z`. Skip after npm audit fix (lock file only).

#### Go audit and fix

1. **Audit**: `govulncheck ./...`. Report both go.mod directive and local toolchain version. Categorize:
   - **"Standard library"** → Go version upgrade (step b)
   - **"Module"** → module update (step a)

2. **Fix**:

   a. **Module update**: `go get <module>@<version> && go mod tidy`. After each `go get`, check if the `go` directive changed — if so, flag as
      implicit Go version upgrade needing confirmation.

   b. **Go version upgrade**: find the HIGHEST fix version across all stdlib CVEs. Update `go` directive (NOT `toolchain`). Confirm with user first.

3. **Re-audit**: `govulncheck ./...` to confirm.

### Phase 3: Update Dockerfiles and Verify

**Dockerfiles** — if Go or Node version changed, scan modifiable Dockerfiles (excluding `.art`) for `FROM` lines with the old version. Show the user
what will change before applying. Skip if no version changed.

**Verify** — run ALL verification commands from step 6 (both ecosystems). For overrides or major bumps: build + lint + test must ALL pass. Up to 2
retry cycles. If still failing, revert and report.

## Safety Rules

- Use `^x.y.z` for versions — never `~` or exact
- NEVER edit `package-lock.json` or `go.sum` directly
- Overrides are LAST RESORT — confirm with user before applying
- NEVER modify `Dockerfile.art` files

## Report

When complete, use this template:

```
## Dependency Audit Report

**Project:** <project-name>
**Branch:** <branch-name>
**Date:** <YYYY-MM-DD>
**Ecosystems:** npm | Go | both

### Summary

| Severity | Found | Fixed | Skipped | Unfixed |
|----------|-------|-------|---------|---------|
| Critical |       |       |         |         |
| High     |       |       |         |         |
| Moderate |       |       |         |         |
| Low      |       |       |         |         |

### Fixed Vulnerabilities

| Package | CVE | Severity | Method | Previous Version | Fixed Version |
|---------|-----|----------|--------|------------------|---------------|

Methods: `direct update` | `parent update (<parent-pkg>)` | `npm audit fix` | `override` | `module update` | `go version upgrade`

### Skipped Vulnerabilities

| Package | CVE | Severity | Reason |
|---------|-----|----------|--------|

Reasons: `dev-only dependency` | `override too risky` | `no fix available` | `false positive`

### Unfixed Vulnerabilities

| Package | CVE | Severity | Blocked By | Recommendation |
|---------|-----|----------|------------|----------------|

### Dockerfile Updates

| Dockerfile | Old Base Image | New Base Image |
|------------|---------------|----------------|

### Verification

| Check | Status | Notes |
|-------|--------|-------|
| Build |        |       |
| Lint  |        |       |
| Test  |        |       |
```

Omit empty sections. Always include Summary and Verification.
