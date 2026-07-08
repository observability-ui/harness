package mp

import "obs/internal/component"

var dir = "projects/monitoring-plugin"

var InstallDeps = &component.Component{
	Name:        "mp-install-deps",
	Description: "Install monitoring plugin frontend dependencies",
	Dir:         dir + "/web",
	Config: map[string]string{
		"npm-cmd": "install",
	},
}

var Frontend = &component.Component{
	Name:        "mp-frontend",
	Description: "Monitoring plugin frontend dev server",
	DependsOn:   []string{"mp-install-deps"},
	Dir:         dir,
	Outputs:     []component.Output{{Name: "port", Value: "9001"}},
	Config: map[string]string{
		"make-target":    "start-frontend",
		"console-plugin": "monitoring-plugin",
	},
}

var Backend = &component.Component{
	Name:        "mp-backend",
	Description: "Monitoring plugin backend dev server",
	Dir:         dir,
	Outputs:     []component.Output{{Name: "port", Value: "9443"}},
	Config: map[string]string{
		"make-target": "start-feature-backend",
	},
}

func init() {
	component.Register(InstallDeps)
	component.Register(Frontend)
	component.Register(Backend)
}
