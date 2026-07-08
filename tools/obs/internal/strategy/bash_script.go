package strategy

import (
	"context"
	"strings"

	"obs/internal/component"
	"obs/internal/runcontext"
)

type BashScript struct{}

func (s *BashScript) Name() string        { return "bash-script" }
func (s *BashScript) Requires() []string { return []string{"bash"} }

func (s *BashScript) Run(_ context.Context, comp *component.Component, _ *runcontext.RunContext) (*component.Step, error) {
	env := make(map[string]string)
	for k, v := range comp.Config {
		if after, ok := strings.CutPrefix(k, "env."); ok {
			env[after] = v
		}
	}

	spec := component.ProcessSpec{
		Name:    comp.Name,
		Command: "bash",
		Args:    []string{"-c", comp.Config["script-content"]},
		Env:     env,
	}

	if comp.Config["script-file"] != "" && comp.Config["files-fs"] != "" {
		spec.Args = []string{"-c", "{{content:" + comp.Name + "-script}}"}
		spec.Files = map[string]component.FileRef{}
	}

	return &component.Step{
		Name:      comp.Name,
		DependsOn: comp.DependsOn,
		Processes: []component.ProcessSpec{spec},
	}, nil
}
