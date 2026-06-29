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

// --- Provider tests ---

type stubProvider struct {
	name  string
	steps []*recipe.Step
}

func (p *stubProvider) Name() string { return p.name }
func (p *stubProvider) Provide(needs []recipe.StepNeed, cfg *recipe.Config) ([]*recipe.Step, error) {
	return p.steps, nil
}

type needfulStubRecipe struct {
	prepareStubRecipe
	needs []recipe.StepNeed
}

func (r *needfulStubRecipe) Needs() []recipe.StepNeed { return r.needs }

func TestEngine_Prepare_ProviderMergesNeeds(t *testing.T) {
	eng := recipe.NewEngine()

	var receivedNeeds []recipe.StepNeed
	recipe.RegisterProvider(&stubProvider{
		name:  "console",
		steps: []*recipe.Step{{Name: "start-console"}},
	})
	// Override with a provider that captures needs
	recipe.RegisterProvider(&captureProvider{
		name:  "console",
		needs: &receivedNeeds,
		steps: []*recipe.Step{{Name: "start-console"}},
	})

	seg1 := recipe.RecipeSegment{
		Recipe: &needfulStubRecipe{
			prepareStubRecipe: prepareStubRecipe{name: "mp", steps: []*recipe.Step{{Name: "mp-frontend"}}},
			needs:             []recipe.StepNeed{{Provider: "console", Config: map[string]string{"plugin": "monitoring-plugin"}}},
		},
		Flags: pflag.NewFlagSet("mp", pflag.ContinueOnError),
	}
	seg2 := recipe.RecipeSegment{
		Recipe: &needfulStubRecipe{
			prepareStubRecipe: prepareStubRecipe{name: "lp", steps: []*recipe.Step{{Name: "lp-frontend"}}},
			needs:             []recipe.StepNeed{{Provider: "console", Config: map[string]string{"plugin": "logging-view-plugin"}}},
		},
		Flags: pflag.NewFlagSet("lp", pflag.ContinueOnError),
	}

	ordered, err := eng.Prepare([]recipe.RecipeSegment{seg1, seg2}, true, nil)
	if err != nil {
		t.Fatalf("Prepare failed: %v", err)
	}

	if len(receivedNeeds) != 2 {
		t.Fatalf("expected provider called with 2 needs, got %d", len(receivedNeeds))
	}

	hasConsole := false
	for _, s := range ordered {
		if s.Name == "start-console" {
			hasConsole = true
		}
	}
	if !hasConsole {
		t.Fatal("expected provider-generated start-console step in output")
	}
}

type captureProvider struct {
	name  string
	needs *[]recipe.StepNeed
	steps []*recipe.Step
}

func (p *captureProvider) Name() string { return p.name }
func (p *captureProvider) Provide(needs []recipe.StepNeed, cfg *recipe.Config) ([]*recipe.Step, error) {
	*p.needs = needs
	return p.steps, nil
}

func TestEngine_Prepare_NoNeedsRecipeUnchanged(t *testing.T) {
	eng := recipe.NewEngine()

	seg := recipe.RecipeSegment{
		Recipe: &prepareStubRecipe{
			name:  "test",
			steps: []*recipe.Step{{Name: "s1"}, {Name: "s2", DependsOn: []string{"s1"}}},
		},
		Flags: pflag.NewFlagSet("test", pflag.ContinueOnError),
	}

	ordered, err := eng.Prepare([]recipe.RecipeSegment{seg}, true, nil)
	if err != nil {
		t.Fatalf("Prepare failed: %v", err)
	}
	if len(ordered) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(ordered))
	}
}

func TestEngine_Prepare_ProviderStepsInDependencyResolution(t *testing.T) {
	eng := recipe.NewEngine()

	recipe.RegisterProvider(&stubProvider{
		name:  "infra",
		steps: []*recipe.Step{{Name: "setup-infra"}},
	})

	seg := recipe.RecipeSegment{
		Recipe: &needfulStubRecipe{
			prepareStubRecipe: prepareStubRecipe{
				name: "app",
				steps: []*recipe.Step{{Name: "start-app", DependsOn: []string{"setup-infra"}}},
			},
			needs: []recipe.StepNeed{{Provider: "infra"}},
		},
		Flags: pflag.NewFlagSet("app", pflag.ContinueOnError),
	}

	ordered, err := eng.Prepare([]recipe.RecipeSegment{seg}, true, nil)
	if err != nil {
		t.Fatalf("Prepare failed: %v", err)
	}
	if len(ordered) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(ordered))
	}
	if ordered[0].Name != "setup-infra" {
		t.Fatalf("expected setup-infra first, got %s", ordered[0].Name)
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
