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
		"make-target": "start-frontend",
	},
}

var Backend = &component.Component{
	Name:        "mp-backend",
	Description: "Monitoring plugin backend dev server",
	Dir:         dir,
	Outputs:     []component.Output{{Name: "port", Value: "9443"}},
	Config: map[string]string{
		"make-target":    "start-feature-backend",
		"console-plugin": "monitoring-plugin",
	},
}

var BuildPush = &component.Component{
	Name:         "mp-build-push",
	Description:  "Build and push monitoring plugin image",
	Capabilities: []string{"container"},
	Dir:          dir,
	Config: map[string]string{
		"dockerfile": "Dockerfile.dev",
	},
	RequiredFlags: []component.RequiredFlag{
		{Name: "image", Usage: "container image to build and push (e.g. quay.io/user/monitoring-plugin:tag)"},
	},
}

var SetMCOUnmanaged = &component.Component{
	Name:         "set-mco-unmanaged",
	Description:  "Set MCO to unmanaged",
	DependsOn:    []string{"mp-build-push"},
	Capabilities: []string{"cluster"},
}

var ScaleDownMP = &component.Component{
	Name:         "scale-down-mp",
	Description:  "Scale down CMO and monitoring plugin",
	DependsOn:    []string{"set-mco-unmanaged"},
	Capabilities: []string{"cluster"},
}

var PatchCMO = &component.Component{
	Name:         "patch-cmo",
	Description:  "Patch CMO with custom monitoring plugin image",
	DependsOn:    []string{"scale-down-mp"},
	Capabilities: []string{"cluster"},
}

var ScaleUpMP = &component.Component{
	Name:         "scale-up-mp",
	Description:  "Scale up CMO and monitoring plugin",
	DependsOn:    []string{"patch-cmo"},
	Capabilities: []string{"cluster"},
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
