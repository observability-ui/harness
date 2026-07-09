# obs — Agent Guide

obs is a CLI tool that orchestrates development and deployment workflows for Observability UI projects using a component/strategy architecture.

## Quick Reference

```bash
# List available recipes
./bin/obs list

# Start a recipe (interactive TUI)
./bin/obs start mp

# Start multiple recipes (shared components are deduplicated)
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

## Adding New Components

1. Create a directory under `components/` (e.g., `components/my-project/`).
2. Define components as `component.Component` structs in `components.go`.
3. Optionally define custom strategies in `strategies.go`, or rely on built-in strategies via `Config` keys.
4. Register components via `init()` with `component.Register()`.
5. Register strategies via `init()` with `strategy.Register(componentName, &MyStrategy{})` if using custom strategies.
6. Add a recipe entry in `components/register.go` with `mixer.RegisterRecipe()`.
7. Build: `make obs`

### Component struct fields

| Field           | Purpose                                                                              |
| --------------- | ------------------------------------------------------------------------------------ |
| `Name`          | Unique identifier used in recipes, DependsOn, and RunContext                         |
| `DependsOn`     | List of component names that must complete before this one starts                    |
| `Dir`           | Working directory for process execution                                              |
| `Ports`         | TCP ports the component listens on (used for readiness probing and plugin discovery) |
| `Config`        | Strategy selection hints and configuration values                                    |
| `RequiredFlags` | Flags the user must provide (prompted interactively if missing)                      |

### Built-in Strategy Config Keys

| Config key           | Strategy           | What it does                                    |
| -------------------- | ------------------ | ----------------------------------------------- |
| `make-target`        | LocalMakeRun       | Runs `make <target>` in component Dir           |
| `compose-file`       | PodmanCompose      | Runs `podman compose -f <file> up`              |
| `dockerfile`         | ContainerBuild     | Builds and pushes container image via `podman`  |
| `console-plugin`     | (console strategy) | Registers component as a console dynamic plugin |
| `console-proxy-path` | (console strategy) | Registers a proxy endpoint on the console       |

`LocalNPM` is registered explicitly per component (not via config keys) with `Cmd` and `Args` fields, e.g., `&strategy.LocalNPM{Cmd: "install", Args: []string{"--no-save"}}`.

### Strategy interface

All strategies implement a unified `Strategy` interface:

```go
type Strategy interface {
    Requires() []string
    Execute(ctx context.Context, comp *component.Component, rc *runcontext.RunContext) (*component.Step, error)
}
```

Each strategy's `Execute()` returns a `*component.Step` with an explicit `Lifecycle` field:

- `LifecycleOneShot` — process runs to completion (builds, deploys, setup scripts)
- `LifecycleLongRunning` — process stays running; readiness determined by port probing (dev servers, compose stacks, containers)

Strategies are bound to components by name via `strategy.Register(componentName, &MyStrategy{})` in each component package's `init()`. Components
without explicit registrations are resolved by config keys (e.g., `make-target` → `LocalMakeRun`). The mixer calls `strategy.Resolve(comp)` to get the
strategies for a component and executes each one.

### Available recipes

**Start (development):**

- `monitoring-plugin` (alias: `mp`) — Frontend + backend dev servers + console
- `logging-view-plugin` (alias: `lp`) — Frontend + backend dev servers + console + local Loki
- `perses` (alias: `perses`) — Build Perses API + run server

**Deploy (cluster):**

- `monitoring-plugin` (alias: `mp`) — Build/push image + patch cluster (requires `--image=...`)
- `seed-users` (alias: `users`) — Create HTPasswd users on cluster
- `seed-users-permissions` (alias: `users-perms`) — Create users + set RBAC permissions

See `components/register.go` and individual `components/*/components.go` files for full definitions.
