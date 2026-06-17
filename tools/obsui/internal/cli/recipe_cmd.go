package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"obsui/internal/process"
	"obsui/internal/recipe"
	"obsui/internal/runner"
)

func newRecipeCmd(command, shortDesc string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                command + " [recipe...] [flags]",
		Short:              shortDesc,
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRecipes(command, args)
		},
	}
	return cmd
}

func runRecipes(command string, args []string) error {
	segments, err := recipe.ParseRecipeArgs(recipe.DefaultRegistry, command, args)
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
	r := runner.NewNonInteractive(os.Stdout)

	go func() {
		for u := range updates {
			if u.Err != nil {
				fmt.Fprintf(os.Stderr, "[%s] %s: %v\n", u.Status, u.StepName, u.Err)
			} else {
				fmt.Fprintf(os.Stderr, "[%s] %s\n", u.Status, u.StepName)
			}
		}
	}()

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
