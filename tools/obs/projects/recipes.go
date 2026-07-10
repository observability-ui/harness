package projects

import (
	"obs/internal/mixer"

	_ "obs/projects/cluster"
	_ "obs/projects/console"
	_ "obs/projects/logging-plugin"
	_ "obs/projects/monitoring-plugin"
	_ "obs/projects/perses"
)

func init() {
	mixer.RegisterRecipe("start", "monitoring-plugin", []string{"mp"}, []string{
		"mp-install-deps", "mp-frontend", "mp-backend", "console",
	})

	mixer.RegisterRecipe("start", "logging-view-plugin", []string{"lp"}, []string{
		"lp-install-deps", "lp-frontend", "lp-backend", "console", "lp-local-loki",
	})

	mixer.RegisterRecipe("start", "perses", []string{"perses"}, []string{
		"perses-build", "perses-api",
	})

	mixer.RegisterRecipe("deploy", "monitoring-plugin", []string{"mp"}, []string{
		"mp-build-push", "set-mco-unmanaged", "scale-down-mp", "patch-cmo", "scale-up-mp",
	})

	mixer.RegisterRecipe("deploy", "seed-users", []string{"users"}, []string{
		"seed-users",
	})

	mixer.RegisterRecipe("deploy", "seed-users-permissions", []string{"users-perms"}, []string{
		"seed-users", "seed-users-permissions",
	})
}
