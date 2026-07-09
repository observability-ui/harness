package mixer

import (
	"context"
	"fmt"
	"os/exec"
	"slices"
	"strconv"
	"strings"

	"obs/internal/component"
	"obs/internal/runcontext"
	"obs/internal/strategy"
)

type RecipeEntry struct {
	Command    string
	Name       string
	Aliases    []string
	Components []string
}

var recipes []RecipeEntry

func RegisterRecipe(command, name string, aliases []string, components []string) {
	recipes = append(recipes, RecipeEntry{
		Command:    command,
		Name:       name,
		Aliases:    aliases,
		Components: components,
	})
}

func GetRecipe(command, nameOrAlias string) (*RecipeEntry, bool) {
	for i := range recipes {
		r := &recipes[i]
		if r.Command != command {
			continue
		}
		if r.Name == nameOrAlias {
			return r, true
		}
		if slices.Contains(r.Aliases, nameOrAlias) {
			return r, true
		}
	}
	return nil, false
}

func ListAllRecipes() []RecipeEntry {
	cp := make([]RecipeEntry, len(recipes))
	copy(cp, recipes)
	return cp
}

func ResetRecipes() {
	recipes = nil
}

type RequirementsError struct {
	Err error
}

func (e *RequirementsError) Error() string  { return e.Err.Error() }
func (e *RequirementsError) Unwrap() error { return e.Err }

type MissingFlagError struct {
	Flag  string
	Usage string
}

func (e *MissingFlagError) Error() string {
	return fmt.Sprintf("--%s is required — %s", e.Flag, e.Usage)
}

func Mix(ctx context.Context, componentNames []string, flagValues map[string]string) ([]*component.Step, *runcontext.RunContext, error) {
	rc := runcontext.New()

	components, err := resolveComponents(componentNames)
	if err != nil {
		return nil, nil, err
	}

	names := make([]string, len(components))
	for i, c := range components {
		names[i] = c.Name
	}
	rc.SetComponents(names)

	if err := checkRequirements(components); err != nil {
		return nil, nil, &RequirementsError{Err: err}
	}

	if err := checkRequiredFlags(components, flagValues); err != nil {
		return nil, nil, err
	}

	for k, v := range flagValues {
		rc.Set("_flags", k, v)
	}

	for _, comp := range components {
		for _, rf := range comp.RequiredFlags {
			if v := flagValues[rf.Name]; v != "" {
				rc.Set(comp.Name, rf.Name, v)
			}
		}
	}

	var allSteps []*component.Step

	for _, comp := range components {
		strategies := strategy.Resolve(comp)
		for _, s := range strategies {
			step, err := s.Execute(ctx, comp, rc)
			if err != nil {
				return nil, nil, fmt.Errorf("execute %q: %w", comp.Name, err)
			}
			if step != nil {
				allSteps = append(allSteps, step)
			}
		}

		for _, p := range comp.Ports {
			rc.Set(comp.Name, "port", strconv.Itoa(p))
		}
	}

	ordered, err := resolveDependencies(allSteps)
	if err != nil {
		return nil, nil, err
	}

	return ordered, rc, nil
}

func checkRequirements(components []*component.Component) error {
	seen := make(map[string]bool)
	var missing []string

	for _, comp := range components {
		strategies := strategy.Resolve(comp)
		for _, s := range strategies {
			for _, tool := range s.Requires() {
				if seen[tool] {
					continue
				}
				seen[tool] = true
				if _, err := exec.LookPath(tool); err != nil {
					missing = append(missing, tool)
				}
			}
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("required tools not found:\n  - %s", strings.Join(missing, "\n  - "))
	}
	return nil
}

func checkRequiredFlags(components []*component.Component, flagValues map[string]string) error {
	for _, comp := range components {
		for _, rf := range comp.RequiredFlags {
			if flagValues[rf.Name] == "" {
				return &MissingFlagError{Flag: rf.Name, Usage: rf.Usage}
			}
		}
	}
	return nil
}

func resolveComponents(names []string) ([]*component.Component, error) {
	resolved := make(map[string]bool)
	visiting := make(map[string]bool)
	var result []*component.Component

	var visit func(name string) error
	visit = func(name string) error {
		if resolved[name] {
			return nil
		}
		if visiting[name] {
			return fmt.Errorf("circular dependency detected: %s", name)
		}
		visiting[name] = true

		comp, ok := component.Get(name)
		if !ok {
			return fmt.Errorf("unknown component %q", name)
		}

		for _, dep := range comp.DependsOn {
			if err := visit(dep); err != nil {
				return err
			}
		}

		visiting[name] = false
		resolved[name] = true
		result = append(result, comp)
		return nil
	}

	for _, name := range names {
		if err := visit(name); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func resolveDependencies(steps []*component.Step) ([]*component.Step, error) {
	byName := make(map[string]*component.Step, len(steps))
	for _, s := range steps {
		if _, exists := byName[s.Name]; exists {
			return nil, fmt.Errorf("duplicate step name %q", s.Name)
		}
		byName[s.Name] = s
	}

	inDegree := make(map[string]int, len(steps))
	dependents := make(map[string][]string)
	for _, s := range steps {
		for _, dep := range s.DependsOn {
			if _, ok := byName[dep]; !ok {
				return nil, fmt.Errorf("step %q depends on unknown step %q", s.Name, dep)
			}
			inDegree[s.Name]++
			dependents[dep] = append(dependents[dep], s.Name)
		}
	}

	var queue []string
	for _, s := range steps {
		if inDegree[s.Name] == 0 {
			queue = append(queue, s.Name)
		}
	}

	var ordered []*component.Step
	for len(queue) > 0 {
		name := queue[0]
		queue = queue[1:]
		ordered = append(ordered, byName[name])
		for _, dep := range dependents[name] {
			inDegree[dep]--
			if inDegree[dep] == 0 {
				queue = append(queue, dep)
			}
		}
	}

	if len(ordered) != len(steps) {
		return nil, fmt.Errorf("circular dependency detected among steps")
	}
	return ordered, nil
}
