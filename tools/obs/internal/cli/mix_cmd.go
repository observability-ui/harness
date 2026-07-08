package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"obs/internal/component"
	"obs/internal/mixer"
	"obs/internal/output"
	"obs/internal/process"
	"obs/internal/runner"
	"obs/internal/state"

	"github.com/spf13/cobra"
)

func newMixerCmd(command, shortDesc string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                command + " [component...] [flags]",
		Short:              shortDesc,
		DisableFlagParsing: true,
		RunE: func(_ *cobra.Command, args []string) error {
			return runMixer(command, args)
		},
	}
	return cmd
}

func runMixer(command string, args []string) error {
	var filteredArgs []string
	var setFlags []string
	for _, arg := range args {
		switch arg {
		case "--dry-run":
			dryRun = true
		case "--non-interactive":
			nonInteractive = true
		case "--detach":
			detach = true
		case "--output-json":
			outputJSON = true
		case "--force":
			force = true
		default:
			if strings.HasPrefix(arg, "--") && strings.Contains(arg, "=") {
				setFlags = append(setFlags, arg)
			} else {
				filteredArgs = append(filteredArgs, arg)
			}
		}
	}

	var componentNames []string
	for _, name := range filteredArgs {
		entry, ok := mixer.GetRecipe(command, name)
		if ok {
			componentNames = append(componentNames, entry.Components...)
			continue
		}
		componentNames = append(componentNames, name)
	}

	if len(componentNames) == 0 {
		return fmt.Errorf("no recipe specified — run 'obs list' to see available recipes")
	}

	seen := make(map[string]bool)
	var deduped []string
	for _, n := range componentNames {
		if !seen[n] {
			seen[n] = true
			deduped = append(deduped, n)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	m := mixer.New()

	prepare := func() ([]*component.Step, error) {
		steps, rc, err := m.Mix(ctx, deduped, "local")
		if err != nil {
			return nil, err
		}
		for _, flag := range setFlags {
			parts := strings.SplitN(strings.TrimPrefix(flag, "--"), "=", 2)
			if len(parts) == 2 {
				rc.Set("_flags", parts[0], parts[1])
			}
		}
		if force {
			if err := checkAndFreePorts(steps); err != nil {
				return nil, err
			}
		}
		return steps, nil
	}

	if dryRun {
		steps, err := prepare()
		if err != nil {
			return err
		}
		fmt.Println("Dry run — steps that would execute:")
		for i, step := range steps {
			fmt.Printf("  %d. %s\n", i+1, step.Name)
			for _, spec := range step.Processes {
				fmt.Printf("     $ %s %s\n", spec.Command, strings.Join(spec.Args, " "))
				if len(spec.Ports) > 0 {
					fmt.Printf("     ports: %v\n", spec.Ports)
				}
			}
		}
		return nil
	}

	mgr := process.NewManager()
	defer mgr.StopAll()

	updates := make(chan component.StepUpdate, 100)

	if !nonInteractive && !outputJSON && !detach && runner.IsTerminal() {
		return runner.RunInteractive(ctx, mgr, prepare, updates)
	}

	steps, err := prepare()
	if err != nil {
		if _, ok := err.(*mixer.RequirementsError); ok {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		return err
	}

	if outputJSON {
		emitter := output.NewJSONEmitter(os.Stdout)
		go func() {
			for u := range updates {
				ev := output.Event{Type: "step_status", Step: u.StepName, Status: u.Status.String()}
				if u.Err != nil {
					ev.Error = u.Err.Error()
				}
				emitter.Emit(ev)
			}
		}()
	}

	var r runner.Runner
	if detach {
		r = runner.NewDetach(state.DefaultStateDir())
	} else {
		r = runner.NewNonInteractive(os.Stdout)
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

	return r.Run(ctx, mgr, steps, updates)
}

func checkAndFreePorts(steps []*component.Step) error {
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
