package cluster

import (
	"context"

	"obs/internal/runcontext"
	"obs/internal/task"
)

var SeedUsers = &task.Task{
	Name:     "seed-users",
	Strategy: &SeedUsersStrategy{},
}

var SeedUsersPermissions = &task.Task{
	Name:      "seed-users-permissions",
	DependsOn: []string{"seed-users"},
	Strategy:  &PermissionsStrategy{},
}

func init() {
	task.Register(SeedUsers)
	task.Register(SeedUsersPermissions)
}

type SeedUsersStrategy struct{}

func (s *SeedUsersStrategy) Requires() []string { return []string{"oc"} }

func (s *SeedUsersStrategy) Execute(_ context.Context, t *task.Task, _ *runcontext.RunContext) (*task.Step, error) {
	return &task.Step{
		Name:      t.Name,
		Lifecycle: task.LifecycleOneShot,
		DependsOn: t.DependsOn,
		Processes: []task.ProcessSpec{{
			Name:    "seed-users",
			Command: "bash",
			Args:    []string{"{{path:seed-script}}", "{{path:htpasswd}}", "{{path:oauth}}"},
			Files: map[string]task.FileRef{
				"seed-script": {FS: filesFS, Path: "files/seed-users.sh"},
				"htpasswd":    {FS: filesFS, Path: "files/users.htpasswd"},
				"oauth":       {FS: filesFS, Path: "files/oauth.yaml"},
			},
		}},
	}, nil
}

type PermissionsStrategy struct{}

func (s *PermissionsStrategy) Requires() []string { return []string{"oc"} }

func (s *PermissionsStrategy) Execute(_ context.Context, t *task.Task, _ *runcontext.RunContext) (*task.Step, error) {
	return &task.Step{
		Name:      t.Name,
		Lifecycle: task.LifecycleOneShot,
		DependsOn: t.DependsOn,
		Processes: []task.ProcessSpec{{
			Name:    "seed-users-permissions",
			Command: "bash",
			Args:    []string{"{{path:perms-script}}"},
			Files: map[string]task.FileRef{
				"perms-script": {FS: filesFS, Path: "files/permissions.sh"},
			},
		}},
	}, nil
}
