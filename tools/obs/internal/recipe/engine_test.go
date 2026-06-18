package recipe_test

import (
	"fmt"
	"testing"

	"github.com/spf13/pflag"
	"obs/internal/recipe"
)

type prepareStubRecipe struct {
	name  string
	reqs  []recipe.Requirement
	steps []*recipe.Step
}

func (r *prepareStubRecipe) Name() string                              { return r.name }
func (r *prepareStubRecipe) Aliases() []string                         { return nil }
func (r *prepareStubRecipe) Description() string                       { return "" }
func (r *prepareStubRecipe) Flags() *pflag.FlagSet                     { return pflag.NewFlagSet(r.name, pflag.ContinueOnError) }
func (r *prepareStubRecipe) Requirements() []recipe.Requirement        { return r.reqs }
func (r *prepareStubRecipe) Steps(_ *recipe.Config) ([]*recipe.Step, error) { return r.steps, nil }

func TestEngine_Prepare(t *testing.T) {
	eng := recipe.NewEngine()

	seg := recipe.RecipeSegment{
		Recipe: &prepareStubRecipe{
			name: "test",
			steps: []*recipe.Step{
				{Name: "step1"},
				{Name: "step2", DependsOn: []string{"step1"}},
			},
		},
		Flags: pflag.NewFlagSet("test", pflag.ContinueOnError),
	}

	ordered, err := eng.Prepare([]recipe.RecipeSegment{seg}, false, nil)
	if err != nil {
		t.Fatalf("Prepare failed: %v", err)
	}
	if len(ordered) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(ordered))
	}
	if ordered[0].Name != "step1" {
		t.Fatalf("expected step1 first, got %s", ordered[0].Name)
	}
}

func TestEngine_Prepare_RequirementsFail(t *testing.T) {
	eng := recipe.NewEngine()

	seg := recipe.RecipeSegment{
		Recipe: &prepareStubRecipe{
			name: "test",
			reqs: []recipe.Requirement{{Name: "missing", Check: func() error { return fmt.Errorf("not found") }}},
			steps: []*recipe.Step{{Name: "s1"}},
		},
		Flags: pflag.NewFlagSet("test", pflag.ContinueOnError),
	}

	_, err := eng.Prepare([]recipe.RecipeSegment{seg}, false, nil)
	if err == nil {
		t.Fatal("should fail on requirement check")
	}
	if _, ok := err.(*recipe.RequirementsError); !ok {
		t.Fatalf("expected RequirementsError, got %T", err)
	}
}

func TestEngine_Prepare_CircularDeps(t *testing.T) {
	eng := recipe.NewEngine()

	seg := recipe.RecipeSegment{
		Recipe: &prepareStubRecipe{
			name: "test",
			steps: []*recipe.Step{
				{Name: "a", DependsOn: []string{"b"}},
				{Name: "b", DependsOn: []string{"a"}},
			},
		},
		Flags: pflag.NewFlagSet("test", pflag.ContinueOnError),
	}

	_, err := eng.Prepare([]recipe.RecipeSegment{seg}, false, nil)
	if err == nil {
		t.Fatal("circular dependency should return error")
	}
}
