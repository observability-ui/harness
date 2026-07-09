package cluster

import (
	"context"

	"obs/internal/component"
	"obs/internal/runcontext"
	"obs/internal/strategy"
)

type SeedUsersStrategy struct{}

func (s *SeedUsersStrategy) Name() string        { return "seed-users" }
func (s *SeedUsersStrategy) Requires() []string { return []string{"oc"} }

func (s *SeedUsersStrategy) Run(_ context.Context, comp *component.Component, _ *runcontext.RunContext) (*component.Step, error) {
	return &component.Step{
		Name:      comp.Name,
		DependsOn: comp.DependsOn,
		Processes: []component.ProcessSpec{{
			Name:    "seed-users",
			Command: "bash",
			Args:    []string{"{{path:seed-script}}", "{{path:htpasswd}}", "{{path:oauth}}"},
			Files: map[string]component.FileRef{
				"seed-script": {FS: filesFS, Path: "files/seed-users.sh"},
				"htpasswd":    {FS: filesFS, Path: "files/users.htpasswd"},
				"oauth":       {FS: filesFS, Path: "files/oauth.yaml"},
			},
		}},
	}, nil
}

type PermissionsStrategy struct{}

func (s *PermissionsStrategy) Name() string        { return "seed-users-permissions" }
func (s *PermissionsStrategy) Requires() []string { return []string{"oc"} }

func (s *PermissionsStrategy) Run(_ context.Context, comp *component.Component, _ *runcontext.RunContext) (*component.Step, error) {
	return &component.Step{
		Name:      comp.Name,
		DependsOn: comp.DependsOn,
		Processes: []component.ProcessSpec{{
			Name:    "seed-users-permissions",
			Command: "bash",
			Args:    []string{"{{path:perms-script}}"},
			Files: map[string]component.FileRef{
				"perms-script": {FS: filesFS, Path: "files/permissions.sh"},
			},
		}},
	}, nil
}

func init() {
	strategy.RegisterSelector(func(comp *component.Component, mode string) (strategy.BuildStrategy, strategy.RunStrategy) {
		switch comp.Name {
		case SeedUsers.Name:
			return nil, &SeedUsersStrategy{}
		case SeedUsersPermissions.Name:
			return nil, &PermissionsStrategy{}
		}
		return nil, nil
	})
}
