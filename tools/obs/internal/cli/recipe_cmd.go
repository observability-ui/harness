package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"obs/internal/output"
	"obs/internal/process"
	"obs/internal/recipe"
	"obs/internal/runner"
	"obs/internal/state"
)

func newRecipeCmd(command, shortDesc string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                command + " [recipe...] [flags]",
		Short:              shortDesc,
		DisableFlagParsing: true,
		RunE: func(_ *cobra.Command, args []string) error {
			return runRecipes(command, args)
		},
	}
	return cmd
}

func runRecipes(command string, args []string) error {
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
		case "--force":
			force = true
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

	portCheck := process.CheckPorts
	if force {
		portCheck = process.FreePorts
	}

	eng := recipe.NewEngine()
	prepare := func() ([]*recipe.Step, error) {
		return eng.Prepare(segments, dryRun, portCheck)
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

	updates := make(chan recipe.StepUpdate, 100)

	// Interactive mode handles requirements inside the TUI
	if !nonInteractive && !outputJSON && !detach && !dryRun && runner.IsTerminal() {
		return runner.RunInteractive(ctx, mgr, prepare, updates)
	}

	// Non-interactive: prepare upfront, exit on failure
	ordered, err := prepare()
	if err != nil {
		if reqErr, ok := err.(*recipe.RequirementsError); ok {
			fmt.Fprintln(os.Stderr, reqErr)
			os.Exit(2)
		}
		return err
	}

	if dryRun {
		fmt.Println("Dry run — steps that would execute:")
		for i, step := range ordered {
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

	// Select runner
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

	return r.Run(ctx, mgr, ordered, updates)
}

