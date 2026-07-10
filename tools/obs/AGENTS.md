# obs — Agent Guide

obs is a CLI tool that orchestrates development and deployment workflows for Observability UI projects using a task/strategy architecture.

## Quick Reference

```bash
# List available recipes
./bin/obs list

# Start a recipe (interactive TUI)
./bin/obs start mp

# Start multiple recipes (shared tasks are deduplicated)
./bin/obs start mp lp

# Start in non-interactive mode (for agents/CI)
./bin/obs start mp --non-interactive

# Deploy recipes
./bin/obs deploy mp --image=quay.io/user/monitoring-plugin:tag

# Dry run (show what would happen without executing)
./bin/obs start mp --dry-run

# JSON output (for machine parsing)
./bin/obs start mp --non-interactive --output-json

# Detach (background processes)
./bin/obs start mp --detach

# Force (kill processes on busy ports)
./bin/obs start mp --force

# Override console image via environment variable
CONSOLE_IMAGE=quay.io/my/console:dev ./bin/obs start mp

# Check status of running processes
./bin/obs status

# Attach to running processes
./bin/obs attach

# Clean up
./bin/obs cleanup [--force]
```

## For Agents

- Always use `--non-interactive` or `--output-json` when invoking from an agent.
- Use `--dry-run` to preview what a recipe will do before running it.
- Use `--force` to kill processes on busy ports before starting.
- Exit codes: 0 = success, 1 = failure, 2 = requirements not met.
- JSON output emits one JSON object per line with fields: `type`, `step`, `status`, `error`.
- Pass `--key=value` flags for recipe-specific configuration (e.g., `--image=...` for deploy).
- Boolean flags accept `--flag`, `--flag=true`, or `--flag=false` forms.
- Unknown `--flag` args without `=` are rejected with a helpful error.
- First Ctrl+C triggers graceful shutdown; second Ctrl+C force-kills all processes.

## Adding a New Project

1. Create a directory under `projects/` (e.g., `projects/my-project/`).
2. Define tasks as `task.Task` structs with an explicit `Strategy` in a single Go file.
3. Register tasks via `init()` with `task.Register()`.
4. Add a recipe entry in `projects/recipes.go` with `mixer.RegisterRecipe()`.
5. Add a blank import in `projects/recipes.go`: `_ "obs/projects/my-project"`.
6. Build: `make obs`

### Example project file

```go
package myproject

import (
    "obs/internal/strategy"
    "obs/internal/task"
)

var Build = &task.Task{
    Name:     "my-build",
    Dir:      "projects/my-project",
    Strategy: strategy.MakeTarget("build"),
}

var Server = &task.Task{
    Name:      "my-server",
    DependsOn: []string{"my-build"},
    Dir:       "projects/my-project",
    Ports:     []int{8080},
    Strategy:  strategy.MakeTarget("run"),
}

func init() {
    task.Register(Build)
    task.Register(Server)
}
```

### Task struct fields

| Field           | Purpose                                                                              |
| --------------- | ------------------------------------------------------------------------------------ |
| `Name`          | Unique identifier used in recipes, DependsOn, and RunContext                         |
| `DependsOn`     | List of task names that must complete before this one starts                          |
| `Dir`           | Working directory for process execution                                              |
| `Ports`         | TCP ports the task listens on (used for readiness probing and plugin discovery)       |
| `Labels`        | Metadata for cross-task discovery (e.g., `console-plugin`, `console-proxy-path`)     |
| `RequiredFlags` | Flags the user must provide (prompted interactively if missing)                      |
| `Strategy`      | How to execute this task — a built-in or custom `Strategy` implementation            |

### Built-in strategies

| Constructor                        | Strategy         | What it does                                    |
| ---------------------------------- | ---------------- | ----------------------------------------------- |
| `strategy.MakeTarget("target")`    | `MakeRun`        | Runs `make <target>` in task Dir                |
| `strategy.NPMRun("cmd", args...)`  | `NPM`            | Runs `npm <cmd> <args>` in task Dir             |
| `strategy.Compose("file")`         | `PodmanCompose`  | Runs `podman compose -f <file> up`              |
| `strategy.DockerBuild("file")`     | `ContainerBuild` | Builds and pushes container image via `podman`  |

### Labels for console integration

| Label key            | Purpose                                             |
| -------------------- | --------------------------------------------------- |
| `console-plugin`     | Registers task as a console dynamic plugin          |
| `console-proxy-path` | Registers a proxy endpoint on the console           |
| `console-proxy-port` | Port for the proxy endpoint (defaults to `8080`)    |

### Strategy interface

All strategies implement a unified `Strategy` interface:

```go
type Strategy interface {
    Requires() []string
    Execute(ctx context.Context, t *task.Task, rc *runcontext.RunContext) (*task.Step, error)
}
```

Each strategy's `Execute()` returns a `*task.Step` with an explicit `Lifecycle` field:

- `LifecycleOneShot` — process runs to completion (builds, deploys, setup scripts)
- `LifecycleLongRunning` — process stays running; readiness determined by port probing (dev servers, compose stacks, containers)

### Available recipes

**Start (development):**

- `monitoring-plugin` (alias: `mp`) — Frontend + backend dev servers + console
- `logging-view-plugin` (alias: `lp`) — Frontend + backend dev servers + console + local Loki
- `perses` (alias: `perses`) — Build Perses API + run server

**Deploy (cluster):**

- `monitoring-plugin` (alias: `mp`) — Build/push image + patch cluster (requires `--image=...`)
- `seed-users` (alias: `users`) — Create HTPasswd users on cluster
- `seed-users-permissions` (alias: `users-perms`) — Create users + set RBAC permissions

See `projects/recipes.go` and individual `projects/*/` files for full definitions.
