# obs

CLI tool for the Observability UI team to run development and deployment recipes.

## Build

```bash
make obs
```

## Usage

```bash
# List available recipes
./bin/obs list

# Start development servers
./bin/obs start monitoring-plugin
./bin/obs start mp con    # multiple recipes

# Deploy to cluster
./bin/obs deploy coo --mode=bundle

# See all options
./bin/obs --help
```

## Modes

- **Interactive** (default): Bubbletea TUI with tabs per process, spinners, keyboard navigation.
- **Non-interactive** (`--non-interactive` or non-TTY): Docker Compose-style prefixed output.
- **JSON** (`--output-json`): Machine-readable JSON lines for CI/agents.
- **Detach** (`--detach`): Start processes in background, exit immediately.

## Architecture

- `cmd/obs/` — Entry point
- `internal/cli/` — Cobra commands
- `internal/recipe/` — Recipe interface, registry, engine
- `internal/process/` — Process lifecycle, ring buffer, port checks
- `internal/ui/` — Bubbletea TUI components
- `internal/runner/` — Interactive and non-interactive runners
- `internal/state/` — .obs/ state directory management
- `internal/output/` — JSON output emitter
- `recipes/` — Recipe implementations
