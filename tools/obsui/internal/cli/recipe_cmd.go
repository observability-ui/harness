package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"obsui/internal/output"
	"obsui/internal/process"
	"obsui/internal/recipe"
	"obsui/internal/runner"
	"obsui/internal/state"
)

func newRecipeCmd(command, shortDesc string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                command + " [recipe...] [flags]",
		Short:              shortDesc,
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRecipes(cmd, command, args)
		},
	}
	return cmd
}

func runRecipes(cmd *cobra.Command, command string, args []string) error {
	// Parse global flags from args before filtering
	var filteredArgs []string
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
		case "--help", "-h":
			// Skip help flag - it shouldn't reach here but just in case
		default:
			filteredArgs = append(filteredArgs, arg)
		}
	}

	segments, err := recipe.ParseRecipeArgs(recipe.DefaultRegistry, command, filteredArgs)
	if err != nil {
		return err
	}

	// Collect all requirements and steps
	var allReqs []recipe.Requirement
	var allSteps []*recipe.Step

	eng := recipe.NewEngine(nil)

	for _, seg := range segments {
		allReqs = append(allReqs, seg.Recipe.Requirements()...)

		cfg := &recipe.Config{Flags: seg.Flags, DryRun: dryRun}
		steps, err := seg.Recipe.Steps(cfg)
		if err != nil {
			return fmt.Errorf("recipe %q: %w", seg.Recipe.Name(), err)
		}
		allSteps = append(allSteps, steps...)
	}

	// Check requirements
	if err := eng.CheckRequirements(allReqs); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	// Check ports
	for _, step := range allSteps {
		for _, spec := range step.Processes {
			if err := process.CheckPorts(spec.Ports); err != nil {
				return err
			}
		}
	}

	// Resolve dependencies
	ordered, err := recipe.ResolveDependencies(allSteps)
	if err != nil {
		return err
	}

	if dryRun {
		fmt.Println("Dry run — steps that would execute:")
		for i, step := range ordered {
			fmt.Printf("  %d. %s\n", i+1, step.Name)
			for _, spec := range step.Processes {
				fmt.Printf("     $ %s %s\n", spec.Command, joinArgs(spec.Args))
				if len(spec.Ports) > 0 {
					fmt.Printf("     ports: %v\n", spec.Ports)
				}
			}
		}
		return nil
	}

	// Set up signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	mgr := process.NewManager()
	defer mgr.StopAll()

	updates := make(chan runner.StepUpdate, 100)

	// JSON output mode
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

	// Detach mode
	if detach {
		store := state.NewStore(state.DefaultStateDir())
		lock := state.NewLock(state.DefaultStateDir())
		if err := lock.Acquire(); err != nil {
			return err
		}
		// Don't release lock — background processes need it

		r := runner.NewNonInteractive(io.Discard)
		go r.Run(ctx, mgr, ordered, updates)

		// Wait briefly for processes to start, then save state and exit
		time.Sleep(500 * time.Millisecond)

		var procs []state.ProcessState
		for _, p := range mgr.All() {
			procs = append(procs, state.ProcessState{
				Name:   p.Spec.Name,
				PID:    p.PID(),
				Status: "running",
			})
			store.WritePID(p.Spec.Name, p.PID())
		}
		var recipeNames []string
		for _, seg := range segments {
			recipeNames = append(recipeNames, seg.Recipe.Name())
		}
		store.Save(&state.RunState{Recipes: recipeNames, Processes: procs})

		fmt.Println("Processes started in background:")
		state.PrintStatus(&state.RunState{Processes: procs})
		return nil
	}

	// Select runner
	var r runner.Runner
	useInteractive := !nonInteractive && !outputJSON && runner.IsTerminal()
	if useInteractive {
		r = runner.NewInteractive()
	} else {
		r = runner.NewNonInteractive(os.Stdout)
		// Only print updates in non-interactive mode
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

	return r.Run(ctx, mgr, ordered, updates)
}

func joinArgs(args []string) string {
	result := ""
	for i, a := range args {
		if i > 0 {
			result += " "
		}
		result += a
	}
	return result
}
