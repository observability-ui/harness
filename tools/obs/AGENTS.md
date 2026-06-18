# obs — Agent Guide

obs is a CLI tool for running development and deployment recipes for Observability UI projects.

## Quick Reference

```bash
# List available recipes
./bin/obs list

# Start a recipe (interactive TUI)
./bin/obs start <recipe> [flags]

# Start in non-interactive mode (for agents/CI)
./bin/obs start <recipe> --non-interactive

# Deploy recipes
./bin/obs deploy <recipe> [flags]

# Multiple recipes at once
./bin/obs start mp con

# Per-recipe flags
./bin/obs start mp --version=4.18 con --version=4.18

# Dry run (show what would happen without executing)
./bin/obs start mp --dry-run

# JSON output (for machine parsing)
./bin/obs start mp --non-interactive --output-json

# Detach (background processes)
./bin/obs start mp --detach

# Force (kill processes on busy ports)
./bin/obs start mp --force

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
- Exit codes: 0 = success, 1 = recipe failure, 2 = requirements not met.
- JSON output emits one JSON object per line with fields: `type`, `step`, `status`, `error`.

## Adding New Recipes

1. Create a new directory under `tools/obs/recipes/` (e.g., `recipes/coo/`).
2. Implement the `recipe.Recipe` interface in a `start.go` and/or `deploy.go` file.
3. Register it in `recipes/register.go`.
4. Build: `make obs`

See `recipes/mp/start.go` for a complete example.
