package cluster

import "obs/internal/component"

var SeedUsers = &component.Component{
	Name: "seed-users",
}

var SeedUsersPermissions = &component.Component{
	Name:      "seed-users-permissions",
	DependsOn: []string{"seed-users"},
}

func init() {
	component.Register(SeedUsers)
	component.Register(SeedUsersPermissions)
}
