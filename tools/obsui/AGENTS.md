# obsui — Agent Guide

obsui is a CLI tool for running development and deployment recipes for Observability UI projects.

## Quick Reference

```bash
# List available recipes
./bin/obsui list

# Start a recipe (interactive TUI)
./bin/obsui start <recipe> [flags]

# Start in non-interactive mode (for agents/CI)
./bin/obsui start <recipe> --non-interactive

# Deploy recipes
./bin/obsui deploy <recipe> [flags]

# Multiple recipes at once
./bin/obsui start mp con

# Per-recipe flags
./bin/obsui start mp --version=4.18 con --version=4.18

# Dry run (show what would happen without executing)
./bin/obsui start mp --dry-run

# JSON output (for machine parsing)
./bin/obsui start mp --non-interactive --output-json

# Detach (background processes)
./bin/obsui start mp --detach

# Check status of running processes
./bin/obsui status

# Attach to running processes
./bin/obsui attach mp

# Clean up
./bin/obsui cleanup [--force]
```

## For Agents

- Always use `--non-interactive` or `--output-json` when invoking from an agent.
- Use `--dry-run` to preview what a recipe will do before running it.
- Exit codes: 0 = success, 1 = recipe failure, 2 = requirements not met.
- JSON output emits one JSON object per line with fields: `type`, `step`, `status`, `error`.

## Adding New Recipes

1. Create a new directory under `tools/obsui/recipes/` (e.g., `recipes/deploycoo/`).
2. Implement the `recipe.Recipe` interface in a `recipe.go` file.
3. Register it in `recipes/register.go`.
4. Build: `make obsui`

See `recipes/startmp/recipe.go` for a complete example.
