# obs

CLI tool for the Observability UI team to orchestrate development and deployment workflows.

## Build

```bash
make obs
```

## Usage

```bash
# List available recipes
./bin/obs list

# Start development servers
./bin/obs start mp
./bin/obs start mp lp    # shared components (console) auto-deduplicated

# Deploy to cluster
./bin/obs deploy mp --image=quay.io/user/monitoring-plugin:tag

# Override console image
CONSOLE_IMAGE=quay.io/my/console:dev ./bin/obs start mp

# See all options
./bin/obs --help
```

## Modes

- **Interactive** (default): Bubbletea TUI with tabs per process, spinners, keyboard navigation.
- **Non-interactive** (`--non-interactive` or non-TTY): Prefixed output per process.
- **JSON** (`--output-json`): Machine-readable JSON lines for CI/agents.
- **Detach** (`--detach`): Start processes in background, exit immediately.

## Architecture

obs uses a **Component/Strategy/Mixer** architecture:

- **Components** define WHAT can run (name, dependencies, outputs, config)
- **Strategies** define HOW to run each component (local process, container, helm, etc.)
- **Mixer** resolves components into an executable step graph (DAG)
- **Runners** execute the step graph (interactive TUI or non-interactive)

```
components/                    — Component definitions + strategies (WHAT + HOW)
  monitoring-plugin/           — MP components (frontend, backend) + deploy strategies
  logging-plugin/              — LP components (frontend, backend, local loki)
  console/                     — Shared console component + container strategy
  cluster/                     — Cluster operations (seed-users, scale CMO)
  register.go                  — Declarative recipe definitions (component lists)

internal/
  component/                   — Component and Step types
  strategy/                    — Strategy interfaces + built-in strategies
  mixer/                       — DAG resolution, strategy selection, recipe registry
  runcontext/                  — Centralized state for cross-component communication
  runner/                      — Step execution (interactive, non-interactive, detach)
  process/                     — OS process lifecycle, port management
  ui/                          — Bubbletea TUI (tabs, spinners, main tab tree view)
  cli/                         — Cobra command structure
  state/                       — .obs/ state directory management
  output/                      — JSON output emitter
```

### Key concepts

**Components** are declarative structs — no imperative code:
```go
var Frontend = &component.Component{
    Name:      "mp-frontend",
    DependsOn: []string{"mp-install-deps"},
    Dir:       "projects/monitoring-plugin",
    Outputs:   []component.Output{{Name: "port", Value: "9001"}},
    Config:    map[string]string{"make-target": "start-frontend"},
}
```

**Strategies** are selected by config keys. A component with `make-target` automatically uses the `LocalMakeRun` strategy. Custom strategies can be registered for specific components (e.g., the console's container strategy).

**Recipes** are named component lists:
```go
mixer.RegisterRecipe("start", "monitoring-plugin", []string{"mp"}, []string{
    "mp-install-deps", "mp-frontend", "mp-backend", "console",
})
```

Running `obs start mp lp` expands both recipes. Shared components (like `console`) appear once — the Mixer deduplicates them. The console's strategy discovers all active plugins via RunContext and configures itself accordingly.
