package mp

import "obs/internal/component"

var dir = "projects/monitoring-plugin"

var InstallDeps = &component.Component{
	Name: "mp-install-deps",
	Dir:  dir + "/web",
}

var Frontend = &component.Component{
	Name:      "mp-frontend",
	DependsOn: []string{"mp-install-deps"},
	Dir:       dir,
	Ports:     []int{9001},
	Config: map[string]string{
		"make-target": "start-frontend",
	},
}

var Backend = &component.Component{
	Name:      "mp-backend",
	DependsOn: []string{"mp-install-deps"},
	Dir:       dir,
	Ports: []int{9443},
	Config: map[string]string{
		"make-target":    "start-feature-backend",
		"console-plugin": "monitoring-plugin",
	},
}

var BuildPush = &component.Component{
	Name: "mp-build-push",
	Dir:  dir,
	Config: map[string]string{
		"dockerfile": "Dockerfile.dev",
	},
	RequiredFlags: []component.RequiredFlag{
		{Name: "image", Usage: "container image to build and push (e.g. quay.io/user/monitoring-plugin:tag)"},
	},
}

var SetMCOUnmanaged = &component.Component{
	Name:      "set-mco-unmanaged",
	DependsOn: []string{"mp-build-push"},
}

var ScaleDownMP = &component.Component{
	Name:      "scale-down-mp",
	DependsOn: []string{"set-mco-unmanaged"},
}

var PatchCMO = &component.Component{
	Name:      "patch-cmo",
	DependsOn: []string{"scale-down-mp"},
}

var ScaleUpMP = &component.Component{
	Name:      "scale-up-mp",
	DependsOn: []string{"patch-cmo"},
}

func init() {
	component.Register(InstallDeps)
	component.Register(Frontend)
	component.Register(Backend)
	component.Register(BuildPush)
	component.Register(SetMCOUnmanaged)
	component.Register(ScaleDownMP)
	component.Register(PatchCMO)
	component.Register(ScaleUpMP)
}
