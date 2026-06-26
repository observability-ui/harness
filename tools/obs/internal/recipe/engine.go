package recipe

import (
	"fmt"
	"strings"
)

type Engine struct{}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Prepare(segments []RecipeSegment, dryRun bool, checkPorts func([]int) error) ([]*Step, error) {
	var allReqs []Requirement
	var allSteps []*Step

	for _, seg := range segments {
		allReqs = append(allReqs, seg.Recipe.Requirements(seg.Flags)...)

		cfg := &Config{Flags: seg.Flags, DryRun: dryRun}
		steps, err := seg.Recipe.Steps(cfg)
		if err != nil {
			return nil, fmt.Errorf("recipe %q: %w", seg.Recipe.Name(), err)
		}
		allSteps = append(allSteps, steps...)
	}

	providerSteps, err := e.resolveProviders(segments)
	if err != nil {
		return nil, err
	}
	allSteps = append(allSteps, providerSteps...)

	if !dryRun {
		if err := e.checkRequirements(allReqs); err != nil {
			return nil, &RequirementsError{Err: err}
		}
		if checkPorts != nil {
			if err := checkPortAvailability(allSteps, checkPorts); err != nil {
				return nil, &RequirementsError{Err: err}
			}
		}
	}

	return resolveDependencies(allSteps)
}

func (e *Engine) resolveProviders(segments []RecipeSegment) ([]*Step, error) {
	grouped := make(map[string][]StepNeed)
	for _, seg := range segments {
		nr, ok := seg.Recipe.(NeedfulRecipe)
		if !ok {
			continue
		}
		for _, need := range nr.Needs() {
			grouped[need.Provider] = append(grouped[need.Provider], need)
		}
	}

	var steps []*Step
	for providerName, needs := range grouped {
		provider, ok := GetProvider(providerName)
		if !ok {
			return nil, fmt.Errorf("unknown step provider %q", providerName)
		}
		provided, err := provider.Provide(needs, &Config{})
		if err != nil {
			return nil, fmt.Errorf("provider %q: %w", providerName, err)
		}
		steps = append(steps, provided...)
	}
	return steps, nil
}

func checkPortAvailability(steps []*Step, checkPorts func([]int) error) error {
	seen := make(map[int]bool)
	var ports []int
	for _, step := range steps {
		for _, spec := range step.Processes {
			for _, p := range spec.Ports {
				if !seen[p] {
					seen[p] = true
					ports = append(ports, p)
				}
			}
		}
	}
	if len(ports) == 0 {
		return nil
	}
	if err := checkPorts(ports); err != nil {
		return fmt.Errorf("port availability check failed:\n  - %v", err)
	}
	return nil
}

type RequirementsError struct {
	Err error
}

func (e *RequirementsError) Error() string { return e.Err.Error() }

func (e *Engine) checkRequirements(reqs []Requirement) error {
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

func resolveDependencies(steps []*Step) ([]*Step, error) {
	byName := make(map[string]*Step, len(steps))
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
