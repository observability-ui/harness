package recipe

import (
	"fmt"

	"github.com/spf13/pflag"
)

type RecipeSegment struct {
	Recipe Recipe
	Flags  *pflag.FlagSet
}

func ParseRecipeArgs(reg *Registry, command string, args []string) ([]RecipeSegment, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("no recipe specified — run 'obs list' to see available recipes")
	}

	var segments []RecipeSegment
	var currentRecipe Recipe
	var currentFlags []string

	flush := func() error {
		if currentRecipe == nil {
			return nil
		}
		fs := copyFlagSet(currentRecipe.Flags())
		if err := fs.Parse(currentFlags); err != nil {
			return fmt.Errorf("invalid flags for %q: %w", currentRecipe.Name(), err)
		}
		segments = append(segments, RecipeSegment{Recipe: currentRecipe, Flags: fs})
		currentFlags = nil
		return nil
	}

	for _, arg := range args {
		if arg == "" {
			continue
		}
		// Flags start with -
		if len(arg) > 0 && arg[0] == '-' {
			if currentRecipe == nil {
				return nil, fmt.Errorf("flag %q before any recipe name", arg)
			}
			currentFlags = append(currentFlags, arg)
			continue
		}
		// Try to look up as a recipe
		rec, ok := reg.Lookup(command, arg)
		if !ok {
			return nil, fmt.Errorf("unknown recipe %q for command %q — run 'obs list' to see available recipes", arg, command)
		}
		if err := flush(); err != nil {
			return nil, err
		}
		currentRecipe = rec
	}

	if currentRecipe == nil {
		return nil, fmt.Errorf("no valid recipe found in arguments")
	}
	if err := flush(); err != nil {
		return nil, err
	}

	return segments, nil
}

func copyFlagSet(src *pflag.FlagSet) *pflag.FlagSet {
	dst := pflag.NewFlagSet(src.Name(), pflag.ContinueOnError)
	src.VisitAll(func(f *pflag.Flag) {
		dst.AddFlag(f)
	})
	return dst
}
