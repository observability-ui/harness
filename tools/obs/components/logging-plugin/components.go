package lp

import "obs/internal/component"

var dir = "projects/logging-view-plugin"

var InstallDeps = &component.Component{
	Name: "lp-install-deps",
	Dir:  dir + "/web",
}

var Frontend = &component.Component{
	Name:      "lp-frontend",
	DependsOn: []string{"lp-install-deps"},
	Dir:       dir,
	Ports:     []int{9001},
	Config: map[string]string{
		"make-target": "start-frontend",
	},
}

var Backend = &component.Component{
	Name:      "lp-backend",
	DependsOn: []string{"lp-frontend"},
	Dir:       dir,
	Ports:     []int{9002},
	Config: map[string]string{
		"make-target":    "start-backend",
		"console-plugin": "logging-view-plugin",
	},
}

var LocalLoki = &component.Component{
	Name:    "lp-local-loki",
	Dir:     dir,
	Ports: []int{3100},
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
