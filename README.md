# Observability UI AI SDLC Harness

AI-driven software development lifecycle harness for the Observability UI team. This repo provides structured context, task tracking, and automation
for using AI coding agents across the team's project portfolio.

## How it works

Each task follows a three-document workflow:

1. **`spec.md`** — Problem statement, related projects/branches, and acceptance criteria.
2. **`plan.md`** — Step-by-step breakdown an AI agent can execute against.
3. **`execution.md`** — Progress tracking with checkboxes and notes captured during execution.

Tasks live in `tasks/`. The `projects/` directory contains git submodules for every repo in scope, giving agents direct access to source code.

## Projects

See [ARCHITECTURE.md](ARCHITECTURE.md) for the full project catalog and system architecture.

## Repository layout

```
tasks/                  # Active tasks (spec + plan + execution)
completed/              # Archived completed tasks
projects/               # Git submodules for all in-scope repos
tools/obs/              # obs CLI source (Go)
bin/                    # Built tools (obs + dprint) — gitignored
claude/plugins/obsui/   # Claude Code plugin (skills for planning, execution, debugging, dev environments, code review)
```

## Setup

```sh
git clone --recurse-submodules https://github.com/observability-ui/harness/
make setup    # install tools, build obs CLI, and reset submodules to their configured branches
```

## Tools

### obs CLI

The `obs` CLI runs development and deployment recipes for projects in the harness. Built with `make tools`, the binary lands in `bin/obs`.

```sh
obs list                              # list available recipes
obs start mp                          # start monitoring plugin (frontend + backend + console)
obs start mp --force                  # kill processes on busy ports, then start
obs --dry-run start mp                # show what would run without executing
obs deploy coo                        # deploy cluster observability operator
obs status                            # show running processes
obs cleanup                           # stop all processes
```

Runs in interactive mode (TUI with tabs per process) by default, falls back to non-interactive (docker-compose style) in CI or with
`--non-interactive`. See [tools/obs/README.md](tools/obs/README.md) for the full reference.

### AI agent skills

The [obsui plugin](claude/plugins/obsui/) provides skills for AI-assisted development:

| Skill                   | Purpose                                           |
| ----------------------- | ------------------------------------------------- |
| `/obsui:planner`        | Create an implementation plan from a spec         |
| `/obsui:executor`       | Execute a plan with parallel agents               |
| `/obsui:bug-diagnostic` | Diagnose a bug from a spec                        |
| `/obsui:dev-env`        | Manage project dev environments via the obsui CLI |
| `/obsui:code-reviewer`  | Multi-angle PR review                             |

## Resetting projects

After working on tasks, submodules may have checked-out branches or uncommitted changes. Run `make reset-projects` to reset all submodules back to the
branches defined in `.gitmodules` at the latest remote HEAD.

## Markdown formatting

All markdown is formatted with [dprint](https://dprint.dev/) (150 char line width). Run `make lint` to format and `make check` to validate.
