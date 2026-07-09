package perses

import "obs/internal/component"

var BuildAPI = &component.Component{
	Name:        "perses-build",
	Description: "Build Perses API server",
	Dir:         "projects/perses",
	Config: map[string]string{
		"make-target": "build-api",
	},
}

var Server = &component.Component{
	Name:        "perses-api",
	Description: "Perses API server",
	DependsOn:   []string{"perses-build"},
	Dir:         "projects/perses",
	Outputs:     []component.Output{{Name: "port", Value: "8080"}},
	Config: map[string]string{
		"console-proxy-path": "/api/proxy/plugin/monitoring-console-plugin/perses/",
		"console-proxy-port": "8080",
	},
}

func init() {
	component.Register(BuildAPI)
	component.Register(Server)
}
