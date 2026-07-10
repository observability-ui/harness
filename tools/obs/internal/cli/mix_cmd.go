package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"obs/internal/task"
	"obs/internal/mixer"
	"obs/internal/output"
	"obs/internal/process"
	"obs/internal/runner"
	"obs/internal/state"

	"github.com/spf13/cobra"
)

type runOptions struct {
	tasks          []string
	flagValues     map[string]string
	dryRun         bool
	nonInteractive bool
	detach         bool
	outputJSON     bool
	force          bool
}

func newMixerCmd(command, shortDesc string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                command + " [task...] [flags]",
		Short:              shortDesc,
		DisableFlagParsing: true,
		RunE: func(_ *cobra.Command, args []string) error {
			opts, err := parseArgs(command, args)
			if err != nil {
				return err
			}
			return runMixer(command, opts)
		},
	}
	return cmd
}

func parseArgs(command string, args []string) (runOptions, error) {
	var opts runOptions

	boolFlags := map[string]*bool{
		"dry-run":         &opts.dryRun,
		"non-interactive": &opts.nonInteractive,
		"detach":          &opts.detach,
		"output-json":     &opts.outputJSON,
		"force":           &opts.force,
	}

	var filteredArgs []string
	var setFlags []string
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			return opts, fmt.Errorf("usage: obs %s <recipe|task> [flags]\n\nFlags:\n  --dry-run          Show steps without executing\n  --non-interactive  Disable TUI mode\n  --detach           Run in background\n  --output-json      Output JSON events\n  --force            Free occupied ports\n  --<key>=<value>    Pass custom flag to strategies", command)
		}
		if !strings.HasPrefix(arg, "--") {
			filteredArgs = append(filteredArgs, arg)
			continue
		}
		name := strings.TrimPrefix(arg, "--")
		val := ""
		if eqIdx := strings.Index(name, "="); eqIdx >= 0 {
			val = name[eqIdx+1:]
			name = name[:eqIdx]
		}
		if ptr, ok := boolFlags[name]; ok {
			*ptr = val != "false"
			continue
		}
		if strings.Contains(arg, "=") && name != "" {
			setFlags = append(setFlags, arg)
		} else {
			return opts, fmt.Errorf("unknown flag %q — use --flag=value for custom flags", arg)
		}
	}

	for _, name := range filteredArgs {
		entry, ok := mixer.GetRecipe(command, name)
		if ok {
			opts.tasks = append(opts.tasks, entry.Tasks...)
			continue
		}
		opts.tasks = append(opts.tasks, name)
	}

	if len(opts.tasks) == 0 {
		return opts, fmt.Errorf("no recipe specified — run 'obs list' to see available recipes")
	}

	seen := make(map[string]bool)
	var deduped []string
	for _, n := range opts.tasks {
		if !seen[n] {
			seen[n] = true
			deduped = append(deduped, n)
		}
	}
	opts.tasks = deduped

	opts.flagValues = make(map[string]string)
	for _, flag := range setFlags {
		parts := strings.SplitN(strings.TrimPrefix(flag, "--"), "=", 2)
		if len(parts) == 2 {
			opts.flagValues[parts[0]] = parts[1]
		}
	}

	return opts, nil
}

func printDryRun(steps []*task.Step) {
	fmt.Println("Dry run — steps that would execute:")
	for i, step := range steps {
		fmt.Printf("  %d. %s [%s]\n", i+1, step.Name, step.Lifecycle)
		for _, spec := range step.Processes {
			fmt.Printf("     $ %s %s\n", spec.Command, strings.Join(spec.Args, " "))
			if len(spec.Ports) > 0 {
				fmt.Printf("     ports: %v\n", spec.Ports)
			}
		}
	}
}

func runMixer(command string, opts runOptions) error {
	ctx, cancel := context.WithCancel(context.Background())

	prepare := func() ([]*task.Step, []task.ProjectInfo, error) {
		steps, _, projects, err := mixer.Mix(ctx, opts.tasks, opts.flagValues)
		if err != nil {
			return nil, nil, err
		}
		if opts.force {
			if err := checkAndFreePorts(steps); err != nil {
				return nil, nil, err
			}
		}
		return steps, projects, nil
	}

	if opts.dryRun {
		defer cancel()
		steps, _, err := prepare()
		if err != nil {
			return err
		}
		printDryRun(steps)
		return nil
	}

	mgr := process.NewManager()

	sigCh := make(chan os.Signal, 2)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	if opts.detach {
		go func() {
			<-sigCh
			cancel()
		}()
	} else {
		go func() {
			<-sigCh
			cancel()
			<-sigCh
			mgr.StopAll()
			os.Exit(1)
		}()
		defer cancel()
		defer mgr.StopAll()
	}

	if !opts.nonInteractive && !opts.outputJSON && !opts.detach && runner.IsTerminal() {
		return runner.RunInteractive(ctx, mgr, prepare, opts.flagValues)
	}

	updates := make(chan task.StepUpdate, 100)

	steps, _, err := prepare()
	if err != nil {
		var reqErr *mixer.RequirementsError
		if errors.As(err, &reqErr) {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		return err
	}

	if opts.outputJSON {
		emitter := output.NewJSONEmitter(os.Stdout)
		go func() {
			for u := range updates {
				ev := output.Event{Type: "step_status", Step: u.StepName, Status: u.Status.String()}
				if u.Err != nil {
					ev.Error = u.Err.Error()
				}
				if err := emitter.Emit(ev); err != nil {
					return
				}
			}
		}()
	}

	var r runner.Runner
	if opts.detach {
		r = runner.NewDetach(state.DefaultStateDir())
		if !opts.outputJSON {
			go func() {
				for range updates {
				}
			}()
		}
	} else {
		r = runner.NewNonInteractive(os.Stdout)
		if !opts.outputJSON {
			go func() {
				for u := range updates {
					if u.Err != nil {
						fmt.Fprintf(os.Stderr, "[%s] %s: %v\n", u.Status, u.StepName, u.Err)
					} else {
						fmt.Fprintf(os.Stderr, "[%s] %s\n", u.Status, u.StepName)
					}
				}
			}()
		}
	}

	err = r.Run(ctx, mgr, steps, updates)
	close(updates)
	return err
}

func checkAndFreePorts(steps []*task.Step) error {
	var ports []int
	seen := make(map[int]bool)
	for _, step := range steps {
		for _, spec := range step.Processes {
			for _, p := range spec.Ports {
				if !seen[p] {
					seen[p] = true
					ports = append(ports, p)
				}
			}
		}
	}
	if len(ports) > 0 {
		return process.FreePorts(ports)
	}
	return nil
}
