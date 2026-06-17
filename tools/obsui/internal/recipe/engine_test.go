package recipe_test

import (
	"fmt"
	"testing"

	"obsui/internal/recipe"
)

func TestEngine_RequirementsFail(t *testing.T) {
	eng := recipe.NewEngine(nil)

	steps := []*recipe.Step{
		{Name: "step1", Processes: nil},
	}
	reqs := []recipe.Requirement{
		{Name: "always-fail", Check: func() error { return fmt.Errorf("missing tool") }},
	}

	err := eng.CheckRequirements(reqs)
	if err == nil {
		t.Fatal("should fail when requirement check fails")
	}
	_ = steps
}

func TestEngine_ResolveDependencies(t *testing.T) {
	steps := []*recipe.Step{
		{Name: "backend", DependsOn: nil},
		{Name: "frontend", DependsOn: []string{"backend"}},
		{Name: "proxy", DependsOn: []string{"frontend", "backend"}},
	}

	ordered, err := recipe.ResolveDependencies(steps)
	if err != nil {
		t.Fatalf("resolve failed: %v", err)
	}

	indexOf := func(name string) int {
		for i, s := range ordered {
			if s.Name == name {
				return i
			}
		}
		return -1
	}

	if indexOf("backend") > indexOf("frontend") {
		t.Fatal("backend must come before frontend")
	}
	if indexOf("frontend") > indexOf("proxy") {
		t.Fatal("frontend must come before proxy")
	}
}

func TestEngine_CircularDependency(t *testing.T) {
	steps := []*recipe.Step{
		{Name: "a", DependsOn: []string{"b"}},
		{Name: "b", DependsOn: []string{"a"}},
	}

	_, err := recipe.ResolveDependencies(steps)
	if err == nil {
		t.Fatal("circular dependency should return error")
	}
}
