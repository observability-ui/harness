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

## Adding New Components

1. Create a directory under `components/` (e.g., `components/my-project/`).
2. Define components as `component.Component` structs in `component.go`.
3. Implement strategy files (e.g., `local.go`) or use built-in strategies via `Config` keys.
4. Register components via `init()` with `component.Register()`.
5. Add a recipe entry in `components/register.go` with `mixer.RegisterRecipe()`.
6. Build: `make obs`

### Built-in Strategy Config Keys

| Config key | Strategy | What it does |
|-----------|----------|-------------|
| `make-target` | LocalMakeRun | Runs `make <target>` in component Dir |
| `npm-cmd` | LocalNPMInstall | Runs `npm <cmd>` (default: install) |
| `compose-file` | PodmanCompose | Runs `podman compose -f <file> up` |
| `oc-args` | OCCommand | Runs `oc <args>` |
| `dockerfile` | ContainerBuild | Builds and pushes container image |
| `console-plugin` | (console strategy) | Registers component as a console plugin |

See `components/monitoring-plugin/component.go` for a complete example.
