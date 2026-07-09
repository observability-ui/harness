package lp

import "obs/internal/component"

var dir = "projects/logging-view-plugin"

var InstallDeps = &component.Component{
	Name:        "lp-install-deps",
	Description: "Install logging plugin frontend dependencies",
	Dir:         dir + "/web",
	Config: map[string]string{
		"npm-cmd": "install",
	},
}

var Frontend = &component.Component{
	Name:        "lp-frontend",
	Description: "Logging plugin frontend dev server",
	DependsOn:   []string{"lp-install-deps"},
	Dir:         dir,
	Outputs:     []component.Output{{Name: "port", Value: "9001"}},
	Config: map[string]string{
		"make-target": "start-frontend",
	},
}

var Backend = &component.Component{
	Name:        "lp-backend",
	Description: "Logging plugin backend dev server",
	DependsOn:   []string{"lp-frontend"},
	Dir:         dir,
	Outputs:     []component.Output{{Name: "port", Value: "9002"}},
	Config: map[string]string{
		"make-target":    "start-backend",
		"console-plugin": "logging-view-plugin",
	},
}

var LocalLoki = &component.Component{
	Name:        "lp-local-loki",
	Description: "Local Loki instance via podman compose",
	Dir:         dir,
	Outputs:     []component.Output{{Name: "port", Value: "3100"}},
	Config: map[string]string{
		"compose-file": "hack/docker-compose/docker-compose.test.yml",
	},
}

func init() {
	component.Register(InstallDeps)
	component.Register(Frontend)
	component.Register(Backend)
	component.Register(LocalLoki)
}
