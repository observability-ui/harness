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

## Prerequisites

Verify required tools are installed for the detected ecosystems:

```bash
node --version && npm --version  # required for npm projects
govulncheck -version             # required for Go projects
```

If a required tool is missing, stop and tell the user what to install before proceeding.

## Process

### Phase 1: Discovery

1. **Detect ecosystem and locate files** — find `go.mod`, `package.json`, `.nvmrc`, and `.node-version`. These may be in the project root OR a
   subdirectory (e.g., `web/`, `frontend/`, `ui/`). Record the directory where each file lives — this is the `<npm-dir>` or `<go-dir>` used in all
   subsequent commands.

```bash
find ./projects/<project> -maxdepth 2 -name go.mod -o -name package.json -o -name .nvmrc -o -name .node-version 2>/dev/null
```

2. **Switch runtime versions** — check if the project pins a specific runtime version and switch to it before any other commands. Different versions
   may have different vulnerabilities or available fixes.

   - **Node**: nvm is a shell function, not a binary — it must be sourced before use. Every Bash call that needs node/npm must initialize nvm first.
     Look for `.nvmrc` or `.node-version` in the same directory as `package.json` OR in the project root. If found, `cd` to the directory containing
     the nvm config file and run `nvm use`:

     ```bash
     export NVM_DIR="$HOME/.nvm" && [ -s "$NVM_DIR/nvm.sh" ] && . "$NVM_DIR/nvm.sh" && cd ./projects/<project>/<npm-dir> && nvm use
     ```

     If the version is not installed, install it first:

     ```bash
     export NVM_DIR="$HOME/.nvm" && [ -s "$NVM_DIR/nvm.sh" ] && . "$NVM_DIR/nvm.sh" && cd ./projects/<project>/<npm-dir> && nvm install && nvm use
     ```

     IMPORTANT: because each Bash tool call runs in a fresh shell, you must prepend the nvm initialization to EVERY subsequent command that uses
     node or npm throughout the entire skill execution. Always `cd` to the `<npm-dir>` so nvm picks up the correct `.nvmrc`:

     ```bash
     export NVM_DIR="$HOME/.nvm" && [ -s "$NVM_DIR/nvm.sh" ] && . "$NVM_DIR/nvm.sh" && cd ./projects/<project>/<npm-dir> && nvm use &>/dev/null && <your npm command>
     ```

   - **Go**: look for `.go-version` in the `<go-dir>` or the `go` directive in `go.mod`. If `gvm` is available (`command -v gvm`), run
     `gvm use <version>`. If the version is not installed, tell the user to install it with `gvm install <version>`.

3. **Detect branch** — report to user before proceeding (different branches = different product versions):

```bash
git -C ./projects/<project> branch --show-current
```

4. **Update branch and install dependencies** — pull the latest changes and install dependencies to ensure the analysis reflects the current state:

```bash
git -C ./projects/<project> pull --ff-only
```

- **npm**: `cd ./projects/<project>/<npm-dir> && npm install`
- **Go**: `cd ./projects/<project> && go mod tidy`

5. **Discover verification commands** — read the project's `Makefile`, `package.json` scripts, `CLAUDE.md`, and `AGENTS.md` to identify build, lint,
   and test commands. Record them for Phase 3.

6. **Run audit**:
   - **npm**: `cd ./projects/<project>/<npm-dir> && npm audit` (add `--json` for parseable output)
   - **Go**: `cd ./projects/<project> && govulncheck ./...`
   - If a specific CVE or package was given as argument, filter results to that target

7. **Present findings** as a table: package, severity, CVE, current version, fixed version (if known), direct vs transitive.

If no vulnerabilities found, report clean and stop.

### Phase 2: Fix

Apply fixes using this escalation strategy. Try each level in order — stop at the first that resolves the vulnerability.

#### npm escalation

1. **Direct update**: the vulnerable package is a direct dependency in `package.json` → check `npm info <pkg> versions` for a patched version → update
   `package.json` + `npm install`.

2. **Parent dependency update**: the vulnerability is in a transitive dependency → use `npm ls <vulnerable-pkg>` to trace the dependency chain back to
   the direct (parent) dependency in `package.json`. Check if the parent has a newer version (`npm info <parent-pkg> versions`) that pulls in a fixed
   version of the transitive dep. If so, update the parent in `package.json` + `npm install`. This is preferred over `npm audit fix` because it
   produces a cleaner, more intentional update.

3. **npm audit fix**: if no parent update resolves it → run `npm audit fix`. NEVER use `--force` — it can introduce breaking major version changes.
   This patches the lock file for safe transitive fixes.

4. **Override**: last resort — no upgrade path exists. Check whether the vulnerable package is only reachable through `devDependencies` (use
   `npm ls <vulnerable-pkg>` and trace the root ancestor). If it is dev-only, recommend skipping the override — it does not affect production
   payloads and overriding it risks breaking CI, tests, or build tooling. Report it as a known dev-only vulnerability and move on.
   If the vulnerability affects production dependencies, add an `overrides` entry in `package.json` pinning the transitive dep to the fixed version.
   Warn the user this is risky and may cause incompatibilities. Confirm with the user before applying. Always run full verification (Phase 3) after
   applying overrides.

IMPORTANT: the dev-only skip applies ONLY to overrides (step 4). Dev dependencies should still be fixed through steps 1-3 (direct update, parent
update, npm audit fix) when possible — these are safe and keep the dependency tree clean.

#### Go escalation

1. **Module update**: the vulnerability is in a third-party module → check for an updated version with `go list -m -versions <module>` or the
   advisory's fixed version → apply with `go get <module>@<version> && go mod tidy`.

2. **Go version upgrade**: the vulnerability is in the Go standard library → check the advisory for the fixed Go version. Update the `go` directive in
   `go.mod`, NOT the `toolchain` directive — the `go` directive is the minimum version requirement that govulncheck checks against and that gets
   enforced when others build the module. The `toolchain` directive only controls which toolchain binary is downloaded locally.
   Warn the user this may affect CI/CD pipelines, container base images, and other build systems. Confirm with the user before applying.
   Always run full verification (Phase 3).

After each fix, re-run `npm audit` or `govulncheck` to confirm the vulnerability is resolved before proceeding.

### Phase 3: Verify

Run the project's build, lint, and test commands discovered in Phase 1.

For override fixes or major version bumps: ALL THREE (build + lint + test) must pass. If any fail, attempt up to 2 fix cycles. If still failing,
revert the change, report the issue, and suggest the user investigate manually with the versions that were attempted.

### Phase 4: Cross-branch

After fixing on the current branch, ask the user:

1. "Should I check which other release branches are affected by this vulnerability?"
2. "Should I apply this fix to other affected release branches?"

If yes: for each target branch, checkout → re-audit to confirm vulnerability exists → apply same fix strategy → verify → return to original branch
when done.

## Safety Rules

- NEVER run `npm audit fix --force`
- NEVER use `--legacy-peer-deps` — it silently ignores peer dependency conflicts that may cause runtime errors
- When updating versions in `package.json`, always use the caret range `^x.y.z` — never use tilde `~x.y.z` or exact `x.y.z`
- NEVER edit `package-lock.json` or `go.sum` directly — let the package manager regenerate them
- Overrides are LAST RESORT — always try direct update, parent update, and audit fix first
- Always confirm override changes with the user before applying
- Always verify overrides with build + lint + test
- ALWAYS use relative paths from the repo root — never absolute paths or `cd /absolute/path && ...` as these trigger permission prompts
- File reads: `projects/<project>/path/to/file` (relative, no leading `./`)
- Bash find/grep/ls: `./projects/<project>/...`
- Git commands in submodules: `git -C ./projects/<project> <command>` (e.g., `git -C ./projects/monitoring-plugin log --oneline -5`)
- If you must use `cd`, always use the relative form: `cd ./projects/<project> && ...`

## Report

When complete, use this template exactly:

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
|         |     |          |        |                  |               |

Methods: `direct update` | `parent update (<parent-pkg>)` | `npm audit fix` | `override` | `module update` | `go version upgrade`

### Skipped Vulnerabilities

| Package | CVE | Severity | Reason |
|---------|-----|----------|--------|
|         |     |          |        |

Reasons: `dev-only dependency` | `override too risky` | `no fix available`

### Unfixed Vulnerabilities

| Package | CVE | Severity | Blocked By | Recommendation |
|---------|-----|----------|------------|----------------|
|         |     |          |            |                |

### Verification

| Check | Status | Notes |
|-------|--------|-------|
| Build |        |       |
| Lint  |        |       |
| Test  |        |       |

### Cross-branch Status

| Branch | Affected | Action Taken |
|--------|----------|--------------|
|        |          |              |
```

Omit sections that have no entries (e.g., if nothing was skipped, omit "Skipped Vulnerabilities"). Always include Summary and Verification.
