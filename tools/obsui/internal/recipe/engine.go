package recipe

import (
	"fmt"
	"strings"
)

type Engine struct {
	Manager any // process.Manager, stored as any to avoid import cycle
}

func NewEngine(mgr any) *Engine {
	return &Engine{Manager: mgr}
}

func (e *Engine) CheckRequirements(reqs []Requirement) error {
	var failures []string
	for _, r := range reqs {
		if err := r.Check(); err != nil {
			failures = append(failures, fmt.Sprintf("  - %s: %v", r.Name, err))
		}
	}
	if len(failures) > 0 {
		return fmt.Errorf("requirements not met:\n%s", strings.Join(failures, "\n"))
	}
	return nil
}

func ResolveDependencies(steps []*Step) ([]*Step, error) {
	byName := make(map[string]*Step, len(steps))
	for _, s := range steps {
		byName[s.Name] = s
	}

	// Kahn's algorithm
	inDegree := make(map[string]int, len(steps))
	dependents := make(map[string][]string) // name -> steps that depend on it
	for _, s := range steps {
		if _, exists := inDegree[s.Name]; !exists {
			inDegree[s.Name] = 0
		}
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

	var ordered []*Step
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
