# obsui

CLI tool for the Observability UI team to run development and deployment recipes.

## Build

```bash
make obsui
```

## Usage

```bash
# List available recipes
./bin/obsui list

# Start development servers
./bin/obsui start monitoring-plugin
./bin/obsui start mp con    # multiple recipes

# Deploy to cluster
./bin/obsui deploy coo --mode=bundle

# See all options
./bin/obsui --help
```

## Modes

- **Interactive** (default): Bubbletea TUI with tabs per process, spinners, keyboard navigation.
- **Non-interactive** (`--non-interactive` or non-TTY): Docker Compose-style prefixed output.
- **JSON** (`--output-json`): Machine-readable JSON lines for CI/agents.
- **Detach** (`--detach`): Start processes in background, exit immediately.

## Architecture

- `cmd/obsui/` — Entry point
- `internal/cli/` — Cobra commands
- `internal/recipe/` — Recipe interface, registry, engine
- `internal/process/` — Process lifecycle, ring buffer, port checks
- `internal/ui/` — Bubbletea TUI components
- `internal/runner/` — Interactive and non-interactive runners
- `internal/state/` — .obsui/ state directory management
- `internal/output/` — JSON output emitter
- `recipes/` — Recipe implementations
