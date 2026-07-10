package mixer

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"

	"obs/internal/runcontext"
	"obs/internal/task"
)

type RecipeEntry struct {
	Command    string
	Name       string
	Aliases    []string
	Tasks []string
}

var recipes []RecipeEntry

func RegisterRecipe(command, name string, aliases []string, tasks []string) {
	recipes = append(recipes, RecipeEntry{
		Command:    command,
		Name:       name,
		Aliases:    aliases,
		Tasks: tasks,
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

func Mix(ctx context.Context, taskNames []string, flagValues map[string]string) ([]*task.Step, *runcontext.RunContext, []task.ProjectInfo, error) {
	rc := runcontext.New()

	tasks, err := resolveTasks(taskNames)
	if err != nil {
		return nil, nil, nil, err
	}

	names := make([]string, len(tasks))
	for i, t := range tasks {
		names[i] = t.Name
	}
	rc.SetTasks(names)

	projects := task.DetectProjects(tasks)
	for _, t := range tasks {
		if t.Dir == "" && t.Labels["image"] != "" {
			image := os.Getenv(strings.ToUpper(strings.ReplaceAll(t.Name, "-", "_")) + "_IMAGE")
			if image == "" {
				image = t.Labels["image"]
			}
			tag := image
			if idx := strings.LastIndex(image, ":"); idx >= 0 {
				tag = image[idx+1:]
			}
			projects = append(projects, task.ProjectInfo{
				Name: t.Name, Branch: tag, IsImage: true,
			})
		}
	}

	if err := checkRequirements(tasks); err != nil {
		return nil, nil, nil, &RequirementsError{Err: err}
	}

	if err := checkRequiredFlags(tasks, flagValues); err != nil {
		return nil, nil, nil, err
	}

	for k, v := range flagValues {
		rc.Set("_flags", k, v)
	}

	for _, t := range tasks {
		for _, rf := range t.RequiredFlags {
			if v := flagValues[rf.Name]; v != "" {
				rc.Set(t.Name, rf.Name, v)
			}
		}
	}

	var allSteps []*task.Step

	for _, t := range tasks {
		step, err := t.Strategy.Execute(ctx, t, rc)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("execute %q: %w", t.Name, err)
		}
		if step != nil {
			allSteps = append(allSteps, step)
		}

		for _, p := range t.Ports {
			rc.Set(t.Name, "port", strconv.Itoa(p))
		}
	}

	ordered, err := resolveDependencies(allSteps)
	if err != nil {
		return nil, nil, nil, err
	}

	return ordered, rc, projects, nil
}

func checkRequirements(tasks []*task.Task) error {
	seen := make(map[string]bool)
	var missing []string

	for _, t := range tasks {
		for _, tool := range t.Strategy.Requires() {
			if seen[tool] {
				continue
			}
			seen[tool] = true
			if _, err := exec.LookPath(tool); err != nil {
				missing = append(missing, tool)
			}
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("required tools not found:\n  - %s", strings.Join(missing, "\n  - "))
	}
	return nil
}

func checkRequiredFlags(tasks []*task.Task, flagValues map[string]string) error {
	for _, t := range tasks {
		for _, rf := range t.RequiredFlags {
			if flagValues[rf.Name] == "" {
				return &MissingFlagError{Flag: rf.Name, Usage: rf.Usage}
			}
		}
	}
	return nil
}

func resolveTasks(names []string) ([]*task.Task, error) {
	resolved := make(map[string]bool)
	visiting := make(map[string]bool)
	var result []*task.Task

	var visit func(name string) error
	visit = func(name string) error {
		if resolved[name] {
			return nil
		}
		if visiting[name] {
			return fmt.Errorf("circular dependency detected: %s", name)
		}
		visiting[name] = true

		t, ok := task.Get(name)
		if !ok {
			return fmt.Errorf("unknown task %q", name)
		}

		for _, dep := range t.DependsOn {
			if err := visit(dep); err != nil {
				return err
			}
		}

		visiting[name] = false
		resolved[name] = true
		result = append(result, t)
		return nil
	}

	for _, name := range names {
		if err := visit(name); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func resolveDependencies(steps []*task.Step) ([]*task.Step, error) {
	byName := make(map[string]*task.Step, len(steps))
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

	var ordered []*task.Step
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
