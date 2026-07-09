package mixer

import (
	"context"
	"fmt"
	"os/exec"
	"slices"
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

func ListRecipes(command string) []RecipeEntry {
	var result []RecipeEntry
	for _, r := range recipes {
		if r.Command == command {
			result = append(result, r)
		}
	}
	return result
}

func ListAllRecipes() []RecipeEntry {
	return recipes
}

type Mixer struct{}

func New() *Mixer {
	return &Mixer{}
}

type RequirementsError struct {
	Err error
}

func (e *RequirementsError) Error() string { return e.Err.Error() }

type MissingFlagError struct {
	Flag  string
	Usage string
}

func (e *MissingFlagError) Error() string {
	return fmt.Sprintf("--%s is required — %s", e.Flag, e.Usage)
}

func (m *Mixer) Mix(ctx context.Context, componentNames []string, mode string, flagValues map[string]string) ([]*component.Step, *runcontext.RunContext, error) {
	rc := runcontext.New(mode)

	components, err := resolveComponents(componentNames)
	if err != nil {
		return nil, nil, err
	}

	if err := m.checkRequirements(components, mode); err != nil {
		return nil, nil, &RequirementsError{Err: err}
	}

	if err := m.checkRequiredFlags(components, flagValues); err != nil {
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
		build, run := strategy.Select(comp, mode)

		if build != nil {
			step, err := build.Build(ctx, comp, rc)
			if err != nil {
				return nil, nil, fmt.Errorf("build %q: %w", comp.Name, err)
			}
			if step != nil {
				allSteps = append(allSteps, step)
			}
		}

		if run != nil {
			step, err := run.Run(ctx, comp, rc)
			if err != nil {
				return nil, nil, fmt.Errorf("run %q: %w", comp.Name, err)
			}
			if step != nil {
				allSteps = append(allSteps, step)
			}
		}

		for _, out := range comp.Outputs {
			if out.Value != "" {
				rc.Set(comp.Name, out.Name, out.Value)
			}
		}
	}

	ordered, err := resolveDependencies(allSteps)
	if err != nil {
		return nil, nil, err
	}

	return ordered, rc, nil
}

func (m *Mixer) checkRequirements(components []*component.Component, mode string) error {
	seen := make(map[string]bool)
	var missing []string

	for _, comp := range components {
		build, run := strategy.Select(comp, mode)
		if build != nil {
			for _, tool := range build.Requires() {
				if seen[tool] {
					continue
				}
				seen[tool] = true
				if _, err := exec.LookPath(tool); err != nil {
					missing = append(missing, tool)
				}
			}
		}
		if run != nil {
			for _, tool := range run.Requires() {
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

func (m *Mixer) checkRequiredFlags(components []*component.Component, flagValues map[string]string) error {
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
