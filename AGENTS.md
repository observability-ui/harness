# AI Agent Guide

This is a monorepo harness for the Observability UI team. See [ARCHITECTURE.md](ARCHITECTURE.md) for the project catalog, dependency graph, and
feature delivery stages.

## Setup

Run `make setup` after cloning to install tools and initialize submodules.

## Tools

- **obs** (`./bin/obs`) — CLI for running development and deployment recipes. Build with `make tools`. Run `./bin/obs list` to see available recipes.
  Always use `--non-interactive` or `--output-json` when invoking from an agent. See [tools/obs/AGENTS.md](tools/obs/AGENTS.md) for the full command
  reference.

## Task workflow

Tasks live under `tasks/<name>/` with a structured pipeline:

1. `spec.md` — problem statement, acceptance criteria, related projects
2. `plan.md` — phased implementation plan with file tables and verification
3. `execution.md` — checklist tracking progress through the plan

Use the obsui plugin skills to drive this workflow:

- `/obsui:planner` — create a plan from a spec
- `/obsui:executor` — execute a plan
- `/obsui:bug-diagnostic` — diagnose a bug from a spec
- `/obsui:dev-env` — manage project dev environments
- `/obsui:code-reviewer` — review a PR

## Projects

Projects are git submodules under `projects/`. Each has its own `CLAUDE.md` or `AGENTS.md` with project-specific guidance. Always use relative paths
from the repo root when referencing project files (`projects/<project>/path/to/file`).

Git commands in submodules: `git -C ./projects/<project> <command>`
