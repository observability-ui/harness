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
func (r *prepareStubRecipe) Requirements(_ *pflag.FlagSet) []recipe.Requirement { return r.reqs }
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

func TestEngine_Prepare_AutoPortRequirements(t *testing.T) {
	eng := recipe.NewEngine()

	checkedPorts := make(map[int]bool)
	portChecker := func(ports []int) error {
		for _, p := range ports {
			checkedPorts[p] = true
		}
		return nil
	}

	seg := recipe.RecipeSegment{
		Recipe: &prepareStubRecipe{
			name: "test",
			steps: []*recipe.Step{
				{
					Name: "s1",
					Processes: []recipe.ProcessSpec{
						{Name: "frontend", Ports: []int{9001}},
					},
				},
				{
					Name: "s2",
					Processes: []recipe.ProcessSpec{
						{Name: "backend", Ports: []int{9443}},
						{Name: "console", Ports: []int{9000}},
					},
				},
			},
		},
		Flags: pflag.NewFlagSet("test", pflag.ContinueOnError),
	}

	_, err := eng.Prepare([]recipe.RecipeSegment{seg}, false, portChecker)
	if err != nil {
		t.Fatalf("Prepare failed: %v", err)
	}

	for _, port := range []int{9001, 9443, 9000} {
		if !checkedPorts[port] {
			t.Errorf("port %d was not checked", port)
		}
	}
}

func TestEngine_Prepare_AutoPortRequirements_Dedup(t *testing.T) {
	eng := recipe.NewEngine()

	callCount := 0
	portChecker := func(ports []int) error {
		callCount++
		return nil
	}

	seg := recipe.RecipeSegment{
		Recipe: &prepareStubRecipe{
			name: "test",
			steps: []*recipe.Step{
				{Name: "s1", Processes: []recipe.ProcessSpec{{Name: "a", Ports: []int{9001}}}},
				{Name: "s2", Processes: []recipe.ProcessSpec{{Name: "b", Ports: []int{9001}}}},
			},
		},
		Flags: pflag.NewFlagSet("test", pflag.ContinueOnError),
	}

	_, err := eng.Prepare([]recipe.RecipeSegment{seg}, false, portChecker)
	if err != nil {
		t.Fatalf("Prepare failed: %v", err)
	}
	if callCount != 1 {
		t.Errorf("expected port 9001 checked once, got %d calls", callCount)
	}
}

func TestEngine_Prepare_AutoPortRequirements_Fail(t *testing.T) {
	eng := recipe.NewEngine()

	portChecker := func(ports []int) error {
		return fmt.Errorf("port %d is already in use", ports[0])
	}

	seg := recipe.RecipeSegment{
		Recipe: &prepareStubRecipe{
			name: "test",
			steps: []*recipe.Step{
				{Name: "s1", Processes: []recipe.ProcessSpec{{Name: "a", Ports: []int{9001}}}},
			},
		},
		Flags: pflag.NewFlagSet("test", pflag.ContinueOnError),
	}

	_, err := eng.Prepare([]recipe.RecipeSegment{seg}, false, portChecker)
	if err == nil {
		t.Fatal("expected error for busy port")
	}
	if _, ok := err.(*recipe.RequirementsError); !ok {
		t.Fatalf("expected RequirementsError, got %T: %v", err, err)
	}
}

func TestEngine_Prepare_RequireFlag_Missing(t *testing.T) {
	eng := recipe.NewEngine()

	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	fs.String("image", "", "container image")

	seg := recipe.RecipeSegment{
		Recipe: &prepareStubRecipe{
			name:  "test",
			reqs:  []recipe.Requirement{recipe.RequireFlag(fs, "image", "container image to deploy")},
			steps: []*recipe.Step{{Name: "s1"}},
		},
		Flags: fs,
	}

	_, err := eng.Prepare([]recipe.RecipeSegment{seg}, false, nil)
	if err == nil {
		t.Fatal("should fail when required flag is not set")
	}
	if _, ok := err.(*recipe.RequirementsError); !ok {
		t.Fatalf("expected RequirementsError, got %T", err)
	}
}

func TestEngine_Prepare_RequireFlag_Set(t *testing.T) {
	eng := recipe.NewEngine()

	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	fs.String("image", "", "container image")
	fs.Parse([]string{"--image=quay.io/my/image:latest"})

	seg := recipe.RecipeSegment{
		Recipe: &prepareStubRecipe{
			name:  "test",
			reqs:  []recipe.Requirement{recipe.RequireFlag(fs, "image", "container image to deploy")},
			steps: []*recipe.Step{{Name: "s1"}},
		},
		Flags: fs,
	}

	_, err := eng.Prepare([]recipe.RecipeSegment{seg}, false, nil)
	if err != nil {
		t.Fatalf("should pass when required flag is set: %v", err)
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
