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
./bin/obs start mp lp    # shared tasks (console) auto-deduplicated

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

obs uses a **Task/Strategy/Mixer** architecture:

- **Tasks** define WHAT to run and HOW — each task embeds its own `Strategy` along with metadata (name, dependencies, ports, labels)
- **Strategies** implement a unified interface with `Requires()` and `Execute()` methods. Reusable strategies (`MakeRun`, `NPM`, `PodmanCompose`, `ContainerBuild`) live in `internal/strategy/`; project-specific strategies live alongside their tasks
- **Mixer** resolves tasks into an executable step graph (DAG), checks tool requirements, validates required flags, and propagates flag values through `RunContext`
- **Runners** execute the step graph (interactive TUI, non-interactive, or detach)

```
projects/                      — Task definitions per project
  monitoring-plugin/           — MP tasks (frontend, backend, build-push, deploy steps)
  logging-plugin/              — LP tasks (frontend, backend, local loki)
  console/                     — Shared console task + container strategy
  cluster/                     — Cluster operations (seed-users, permissions)
  perses/                      — Perses dashboard tasks (build, API server)
  recipes.go                   — Declarative recipe definitions (task lists)

internal/
  task/                        — Task, Strategy interface, Step, ProcessSpec types
  strategy/                    — Reusable strategy implementations
  mixer/                       — DAG resolution, recipe registry
  runcontext/                  — Centralized state for cross-task communication
  runner/                      — Step execution (interactive, non-interactive, detach)
  process/                     — OS process lifecycle, port management, ring buffer
  ui/                          — Bubbletea TUI (tabs, spinners, main tab tree view)
  cli/                         — Cobra command structure + flag parsing
  state/                       — .obs/ state directory management (detach/attach)
  output/                      — JSON output emitter
```

### Key concepts

**Tasks** combine metadata and strategy in a single struct:
```go
var Frontend = &task.Task{
    Name:      "mp-frontend",
    DependsOn: []string{"mp-install-deps"},
    Dir:       "projects/monitoring-plugin",
    Ports:     []int{9001},
    Labels:    map[string]string{"console-plugin": "monitoring-plugin"},
    Strategy:  strategy.MakeTarget("start-frontend"),
}
```

Tasks can declare **RequiredFlags** for values that must be provided at runtime (e.g., `--image` for deploy). The TUI prompts interactively for missing flags.

**Strategies** implement a unified `Strategy` interface with `Requires()` and `Execute()` methods. Each strategy produces a `Step` with an explicit `Lifecycle` — either `LifecycleOneShot` (runs to completion) or `LifecycleLongRunning` (stays running, readiness via port probing). Every task has an explicit strategy — no magic resolution.

**Recipes** are named task lists:
```go
mixer.RegisterRecipe("start", "monitoring-plugin", []string{"mp"}, []string{
    "mp-install-deps", "mp-frontend", "mp-backend", "console",
})
```

Running `obs start mp lp` expands both recipes. Shared tasks (like `console`) appear once — the Mixer deduplicates them. The console's strategy discovers all active plugins via RunContext and configures itself accordingly.

### Flags

Boolean flags: `--dry-run`, `--non-interactive`, `--detach`, `--output-json`, `--force`. These accept `--flag` or `--flag=true`/`--flag=false` forms.

Custom key-value flags use `--key=value` format (e.g., `--image=quay.io/user/plugin:tag`).

### Signal handling

- First Ctrl+C: graceful shutdown (sends SIGINT to process groups, waits for exit)
- Second Ctrl+C: force shutdown (kills all processes, exits immediately)
