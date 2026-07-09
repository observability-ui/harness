package console

import "obs/internal/component"

var Console = &component.Component{
	Name:    "console",
	Ports: []int{9000},
	Config: map[string]string{
		"image":    "quay.io/openshift/origin-console:latest",
		"platform": "linux/amd64",
	},
}

func init() {
	component.Register(Console)
}
