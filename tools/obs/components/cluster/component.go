package cluster

import "obs/internal/component"

var SeedUsers = &component.Component{
	Name:         "seed-users",
	Description:  "Create test users in the cluster (22 users, password = username)",
	Capabilities: []string{"cluster"},
}

var SeedUsersPermissions = &component.Component{
	Name:         "seed-users-permissions",
	Description:  "Assign RBAC roles to test users (monitoring, logging, perses)",
	DependsOn:    []string{"seed-users"},
	Capabilities: []string{"cluster"},
}

func init() {
	component.Register(SeedUsers)
	component.Register(SeedUsersPermissions)
}
