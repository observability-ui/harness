---
name: bug-diagnostic
description: Diagnose a bug from a spec, produce a diagnostic.md with root cause analysis and a plan.md for the executor to implement the fix.
allowed-tools: Read, Bash(find:*), Bash(grep:*), Bash(rg:*), Bash(git log:*), Bash(git diff:*), Bash(git show:*), Bash(git branch:*), Bash(git checkout:*), Bash(git tag:*), Bash(git -C:*), Bash(wc:*), Bash(ls:*), Bash(npm test:*), Bash(npm run:*), Bash(make:*), Bash(go test:*), Bash(./bin/obs *), LSP, Agent
---

## Input

$ARGUMENTS is a task folder name. The folder must exist under `tasks/` and contain a `spec.md` file describing the bug.

The spec should include:

- **Description** of the bug (what is broken, when it happens)
- **Reproduction steps** (commands, user actions, or test cases that trigger it)
- **Expected vs. actual behavior**
- **Related projects and branches** (which codebases are affected)
- **Hints** (optional — error messages, stack traces, suspected areas)

## Prerequisites

The projects referenced in the spec live as git submodules under `projects/`. The repository's `.claude/settings.json` already includes
`additionalDirectories` for the project submodules, so file reads and bash commands against those paths will not trigger permission prompts.

**Path rules — ALWAYS use relative paths from the repo root:**

- File reads: `projects/<project>/path/to/file` (relative, no leading `./`)
- Bash find/grep/ls: `./projects/<project>/...`
- Git commands in submodules: `git -C ./projects/<project> <command>` (e.g., `git -C ./projects/perses log --oneline -5`)
- NEVER use absolute paths or `cd /absolute/path && git ...` — these trigger permission prompts for untrusted hooks
- The permission allowlist matches `cd ./projects/* && git *`, so if you must use `cd`, always use the relative form:
  `cd ./projects/<project> && git ...`

**Iron law: NO FIXES WITHOUT ROOT CAUSE INVESTIGATION FIRST.**

Do not propose a fix, write a plan, or modify any code until the root cause is identified and documented in the diagnostic. Resist the urge to "just
try something" — trace the data flow first.

## Steps

### 1. Check out the affected branch

The bug may only exist on a specific branch or version. This step ensures every project is on the correct branch **before** any reproduction or
investigation begins.

**1a. Read the spec and system context**

Read these files in order:

```
tasks/$ARGUMENTS/spec.md
ARCHITECTURE.md
```

For each project listed in the spec's "Related projects and branches" section, use the Read tool with relative paths to read these files if they
exist:

- `projects/<project>/CLAUDE.md`
- `projects/<project>/AGENTS.md`
- `projects/<project>/README.md`

Run all project reads in parallel across projects to minimize round-trips.

**1b. Determine the target branch per project**

For each project in scope, check what branch is currently checked out and what branches/tags are available:

```bash
git -C ./projects/<project> branch --show-current && echo "---" && git -C ./projects/<project> log --oneline -5
git -C ./projects/<project> tag --sort=-creatordate | head -10
```

Now compare against the spec:

- If the spec's "Related projects and branches" section **specifies a branch or tag** for the project → use that.
- If the spec **does not specify a branch** for a project → ask the user before proceeding. Use AskUserQuestion with the current branch, recent tags,
  and `main` as options:

```
Which branch/version of <project> has the bug?
Options:
- <current branch> (currently checked out)
- main
- <latest tag> (latest release)
- [Other — user types a branch or tag]
```

Do NOT assume the current branch is correct. A version-specific bug on a release tag will not reproduce on `main`.

**1c. Switch to the target branch**

For each project, check out the confirmed branch:

```bash
git -C ./projects/<project> checkout <target-branch>
```

Verify the checkout succeeded and record the exact commit:

```bash
git -C ./projects/<project> log --oneline -1
```

**1d. Identify scope**

After reading and checking out branches, identify:

- Which repositories are in scope and on which branch/version
- What the bug symptoms are (error messages, incorrect behavior, test failures)
- What the spec's reproduction steps are
- Any hints about suspected root cause or affected areas

**1e. Check dev environment readiness**

Before attempting reproduction, check what recipes and environments are available, and whether processes are already running.

1. Ensure the CLI is built, then discover available recipes and check running state:

```bash
make obs 2>/dev/null
./bin/obs list
./bin/obs status
```

2. Use `--dry-run` to preview what a recipe would do without executing:

```bash
./bin/obs --dry-run start <recipe-alias>
./bin/obs --dry-run deploy <recipe-alias>
```

3. Based on the bug and available recipes, decide what to set up:

   - **Bug requires a running dev server** (UI behavior, frontend rendering, API responses): Start the relevant recipe. Always use
     `--non-interactive` when running from an agent. Use `--force` to kill any processes on busy ports:

     ```bash
     ./bin/obs --non-interactive --force start <recipe-alias>
     ```

   - **Bug requires a deployed component on a cluster** (operator behavior, plugin loading, CRD reconciliation): Deploy the relevant
     recipe:

     ```bash
     ./bin/obs --non-interactive deploy <recipe-alias>
     ```

   - **Already running:** `./bin/obs status` shows active processes — note them and proceed to reproduction.

   - **No matching recipe exists:** Skip this step and proceed to reproduction using the project's own commands (Makefile targets,
     npm scripts, go test).

   - **Bug reproduces via unit tests or code inspection alone:** Skip environment setup entirely.

4. After reproduction and investigation are complete, clean up any environments you started:

```bash
./bin/obs cleanup --force
```

**Principle: suggest, never block.** Many bugs can be reproduced through unit tests or code inspection without a full dev environment.
Do not refuse to proceed if no recipe exists or if the dev environment is not running. The obs tool is a convenience, not a gate.

**Agent-specific notes:**
- Always use `--non-interactive` or `--output-json` — interactive mode requires a TTY.
- Exit codes: 0 = success, 1 = recipe failure, 2 = requirements not met (e.g., `oc` not logged in).
- Use `--dry-run` before running to understand what commands and ports a recipe uses.
- Use `--force` to automatically kill processes on busy ports instead of failing.

### 2. Reproduce the bug

Attempt to reproduce the bug using the steps from the spec. The goal is to see the failure firsthand and capture exact output.

**Run reproduction steps:**

- Execute the commands or test cases from the spec
- Use the project's own test/build commands (from CLAUDE.md, Makefile, package.json):

```bash
# Check available commands
grep -E '^[a-zA-Z_-]+:' ./projects/<project>/Makefile 2>/dev/null | head -20
grep -A 30 '"scripts"' ./projects/<project>/package.json 2>/dev/null
```

- Capture exact error messages, stack traces, and failing test output
- Note the environment: branch, commit, any relevant configuration

**Command rules:** NEVER run `npx` commands directly — they trigger permission prompts. Always use Makefile targets or npm scripts like
`npm run build`, `npm run test`, `npm run lint` or `npm run type-check`. Go commands (`go test`, `go build`) are fine.

**If reproduction fails:**

- Try variations (different inputs, different order of steps)
- Check whether the spec's branch/commit is checked out
- If the bug cannot be reproduced after reasonable effort, ask the user for clarification using AskUserQuestion before proceeding. Do NOT guess or
  skip ahead.

**Document what you observed** — exact commands run, exact output, and how it differs from expected behavior. This goes into the diagnostic.

### 3. Investigate root cause

Follow a structured investigation. Work backward from the symptom to the cause.

**3a. Trace the error to its source**

Start at the error message or failing assertion and trace backward through the code:

```bash
# Find where the error is raised
grep -rn "error message text" projects/<project>/src/ --include="*.ts"
grep -rn "error message text" projects/<project>/ --include="*.go"

# Check recent changes to affected files
git -C ./projects/<project> log --oneline -20 -- path/to/affected/file
git -C ./projects/<project> diff HEAD~5 -- path/to/affected/file
```

Use LSP when available:

- `goToDefinition` — follow function calls from the error site
- `findReferences` — find all callers of a broken function
- `incomingCalls` / `outgoingCalls` — trace the call chain
- `hover` — check type signatures at suspicious points

**3b. Multi-repo bugs**

When the bug spans multiple repositories, launch parallel Explore agents (one per repo) to investigate simultaneously. Each agent should report:

- The relevant code paths in that repo
- Recent changes that could have introduced the bug
- How the repo interacts with other affected repos (API contracts, shared types)

Synthesize their findings to identify where the contract is broken.

**3c. Compare against working state**

- Find similar code that works correctly and compare
- Check the last known good commit: `git -C ./projects/<project> log --oneline --all -- path/to/file`
- Use `git -C ./projects/<project> diff <good-commit> <bad-commit> -- path/to/file` to isolate what changed

**3d. Form and test hypotheses**

1. Form a single hypothesis based on the evidence gathered
2. Test it minimally — read one file, run one command, check one output
3. If confirmed, proceed to Step 4
4. If refuted, record the hypothesis and evidence, then form a new one

**Red flags — restart the investigation if you notice:**

- Proposing a fix before tracing the full data flow
- "Just try changing X and see if it works"
- Third failed hypothesis in a row — step back and question your assumptions
- Each "fix" reveals a new problem in a different place (symptom of wrong root cause)

### 4. Clarify and confirm

Before writing artifacts, present your findings to the user:

1. State the root cause you identified (one sentence)
2. Show the key evidence (file:line references, command output)
3. Describe the fix approach at a high level

Ask targeted questions using AskUserQuestion:

- **Fix scope** — should the fix be minimal (patch the symptom) or structural (address the underlying design issue)?
- **Acceptable trade-offs** — performance vs. correctness, backward compatibility constraints
- **Testing expectations** — unit tests sufficient, or integration/E2E tests needed?

Wait for the user's answers before proceeding to Steps 5 and 6.

### 5. Write diagnostic.md

Save to `tasks/$ARGUMENTS/diagnostic.md`. This document is the evidence record — it must stand on its own without requiring someone to re-run the
investigation.

Use the diagnostic template below.

### 6. Write plan.md

Save to `tasks/$ARGUMENTS/plan.md` using the plan template below. This plan must be in the exact format the executor skill expects — same sections,
same table structures, same phase conventions.

**Plan authoring rules:**

- The Problem section should reference the diagnostic: `See tasks/$ARGUMENTS/diagnostic.md for full root cause analysis.`
- Current State table must include the buggy components with their current (broken) behavior
- Each phase's Details section should include code snippets for non-obvious fixes
- Verification section must include the reproduction case — after the fix, the original bug must not reproduce

**Detail calibration:**

- **Code snippets:** Include for type signature changes, API contract changes, non-obvious logic, tricky merge patterns
- **Line references:** Include when the exact insertion/modification point matters
- **Prose:** Use for straightforward config changes, import updates
- **Files Modified table:** Required for every phase that modifies files

**Parallel execution annotations:**

Each phase must declare its dependency and whether it can run in parallel with other phases. The constraint: only one agent should modify a given file
at a time. Phases touching different repos or non-overlapping files can run in parallel via separate agents.

**Self-review before saving:**

1. **Root cause coverage** — the plan addresses the root cause from the diagnostic, not just the symptom
2. **Regression test** — at least one phase adds or modifies a test that would have caught this bug
3. **Dependency ordering and parallelism** — phases reference correct dependencies, parallel phases don't modify overlapping files
4. **File path accuracy** — every path exists in the codebase or is marked as a new file
5. **Reproduction case in verification** — the spec's reproduction steps appear in the Verification section with "should no longer reproduce" as
   expected outcome

## Diagnostic template

```
# Diagnostic: [Bug Name]

## Bug Summary

[One-paragraph description: what is broken, what the user experiences, and the severity/impact.]

## Reproduction

| Step | Command / Action | Expected | Actual |
| ---- | ---------------- | -------- | ------ |
| 1    | [command or action] | [what should happen] | [what actually happens] |
| ...  | ...              | ...      | ...    |

### Environment

- **Branch:** [branch name per project]
- **Commit:** [short SHA per project]
- **Relevant config:** [any environment-specific settings]

### Error Output
```

[Exact error message, stack trace, or failing test output — verbatim, not paraphrased]

```
## Investigation

### Hypothesis 1: [short description]

**Evidence:** [what was checked — files read, commands run, LSP queries]
**Result:** Confirmed | Refuted
**Details:** [what was found and why it confirms or rules out this hypothesis]

### Hypothesis N: [short description]

[Repeat for each hypothesis tested. Include refuted hypotheses — they narrow the search space and prevent re-investigation.]

## Root Cause

[Clear, precise explanation of why the bug occurs. Reference specific file:line locations. Explain the mechanism — what triggers it, what state
becomes incorrect, and why the current code produces wrong output.]

### Affected Components

| Component | File / Location | Impact |
| --------- | --------------- | ------ |
| [name]    | `project/path/to/file.ext:line` | [How this component is affected] |
| ...       | ...             | ...    |

### Contributing Factors

[Environment, configuration, timing, or data conditions that contribute to the bug. Omit this section if the bug is purely a code defect with no
environmental factors.]

## Fix Strategy

[High-level approach to fixing the bug. Explain why this approach was chosen over alternatives. Mention any trade-offs (e.g., "minimal patch now,
structural fix in follow-up" vs. "fix the root design issue").]
```

## Plan template

```
# Plan: Fix — [Bug Name]

## Problem

[Why this fix is needed. Reference the diagnostic: "See `tasks/$ARGUMENTS/diagnostic.md` for full root cause analysis."
Summarize the root cause in 2-3 sentences. Link upstream issues if relevant.]

## Current State

| Component | File / Location | Current Behavior |
| --------- | --------------- | ---------------- |
| [name]    | `project/path/to/file.ext:line` | [What it does now — the buggy behavior] |
| ...       | ...             | ...              |

## Changes

### Phase 1: [Name]

**Dependency:** None
**Parallel with:** None | Phase N (when touching different repos/files)

#### Files Modified

| File | Change |
| ---- | ------ |
| `project/path/to/file.ext` | [Brief description of what changes] |
| ...  | ...    |

#### Details

[Detailed description of the changes. Include code snippets for type changes and non-obvious logic. Include line references when the exact point matters.]

##### [Sub-section for complex changes within this phase]

[For phases with multiple independent changes, use sub-sections.]

#### Phase 1 Verification

- [Specific command and expected output]
- [Manual check if automated verification is not possible]

### Phase 2: [Name]

**Dependency:** Phase 1
**Parallel with:** Phase 3 (different repo)

[Same structure as Phase 1]

...

## PR Strategy

| PR | Repository | Branch | Description | Dependencies |
| -- | ---------- | ------ | ----------- | ------------ |
| 1  | [repo]     | [branch] | [what this PR contains] | None |
| 2  | [repo]     | [branch] | [what this PR contains] | PR 1 merged |
| ...| ...        | ...    | ...         | ...          |

[If all changes fit in a single PR, use one row. For multi-repo tasks, list PRs in merge order. Note which can be reviewed in parallel.]

## Verification

[End-to-end verification mapped to the spec's acceptance criteria and the original reproduction case.]

- [Original reproduction case] - should no longer reproduce after fix
- [Acceptance criterion] - [how to verify]
- ...

## Risks

| Risk | Impact | Mitigation |
| ---- | ------ | ---------- |
| [What could go wrong] | [What breaks] | [How to prevent or recover] |
| ...  | ...    | ...        |
```
