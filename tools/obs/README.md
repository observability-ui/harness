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

- **Components** define WHAT can run (name, dependencies, ports, config)
- **Strategies** define HOW to run each component via a unified `Strategy` interface with `Requires()` and `Execute()` methods
- **Mixer** resolves components into an executable step graph (DAG), checks tool requirements, validates required flags, and propagates flag values through `RunContext`
- **Runners** execute the step graph (interactive TUI, non-interactive, or detach)

```
components/                    — Component definitions + strategies (WHAT + HOW)
  monitoring-plugin/           — MP components (frontend, backend, build-push) + deploy strategies
  logging-plugin/              — LP components (frontend, backend, local loki)
  console/                     — Shared console component + container strategy
  cluster/                     — Cluster operations (seed-users, permissions)
  perses/                      — Perses dashboard components (build, API server)
  register.go                  — Declarative recipe definitions (component lists)

internal/
  component/                   — Component, Step (with Lifecycle), ProcessSpec types + helpers
  strategy/                    — Strategy interface + built-in strategies (defaults.go)
  mixer/                       — DAG resolution, strategy selection, recipe registry
  runcontext/                  — Centralized state for cross-component communication
  runner/                      — Step execution (interactive, non-interactive, detach)
  process/                     — OS process lifecycle, port management, ring buffer
  ui/                          — Bubbletea TUI (tabs, spinners, main tab tree view)
  cli/                         — Cobra command structure + flag parsing
  state/                       — .obs/ state directory management (detach/attach)
  output/                      — JSON output emitter
```

### Key concepts

**Components** are declarative structs — no imperative code:
```go
var Frontend = &component.Component{
    Name:      "mp-frontend",
    DependsOn: []string{"mp-install-deps"},
    Dir:       "projects/monitoring-plugin",
    Ports:     []int{9001},
    Config:    map[string]string{"make-target": "start-frontend"},
}
```

Components can declare **RequiredFlags** for values that must be provided at runtime (e.g., `--image` for deploy). The TUI prompts interactively for missing flags.

**Strategies** implement a unified `Strategy` interface with `Requires()` and `Execute()` methods. Each strategy produces a `Step` with an explicit `Lifecycle` — either `LifecycleOneShot` (runs to completion) or `LifecycleLongRunning` (stays running, readiness via port probing). Strategies are bound to components by name via `strategy.Register()`. Components without explicit registrations are resolved automatically by config keys (e.g., `make-target` → `LocalMakeRun`).

**Recipes** are named component lists:
```go
mixer.RegisterRecipe("start", "monitoring-plugin", []string{"mp"}, []string{
    "mp-install-deps", "mp-frontend", "mp-backend", "console",
})
```

Running `obs start mp lp` expands both recipes. Shared components (like `console`) appear once — the Mixer deduplicates them. The console's strategy discovers all active plugins via RunContext and configures itself accordingly.

### Flags

Boolean flags: `--dry-run`, `--non-interactive`, `--detach`, `--output-json`, `--force`. These accept `--flag` or `--flag=true`/`--flag=false` forms.

Custom key-value flags use `--key=value` format (e.g., `--image=quay.io/user/plugin:tag`).

### Signal handling

- First Ctrl+C: graceful shutdown (sends SIGINT to process groups, waits for exit)
- Second Ctrl+C: force shutdown (kills all processes, exits immediately)
