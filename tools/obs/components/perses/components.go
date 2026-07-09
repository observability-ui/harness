package perses

import "obs/internal/component"

var BuildAPI = &component.Component{
	Name: "perses-build",
	Dir:  "projects/perses",
	Config: map[string]string{
		"make-target": "build-api",
	},
}

var Server = &component.Component{
	Name:      "perses-api",
	DependsOn: []string{"perses-build"},
	Dir:       "projects/perses",
	Ports:     []int{8080},
	Config: map[string]string{
		"console-proxy-path": "/api/proxy/plugin/monitoring-console-plugin/perses/",
		"console-proxy-port": "8080",
	},
}

func init() {
	component.Register(BuildAPI)
	component.Register(Server)
}
