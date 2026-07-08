package console

import "obs/internal/component"

var Console = &component.Component{
	Name:         "console",
	Description:  "OpenShift console serving registered plugins",
	Capabilities: []string{"local", "container"},
	Outputs:      []component.Output{{Name: "port", Value: "9000"}},
	Config: map[string]string{
		"image":    "quay.io/openshift/origin-console:latest",
		"platform": "linux/amd64",
	},
}

func init() {
	component.Register(Console)
}
