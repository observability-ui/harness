package mp

import (
	"context"
	"fmt"

	"obs/internal/component"
	"obs/internal/runcontext"
	"obs/internal/strategy"
)

var BuildPush = &component.Component{
	Name:         "mp-build-push",
	Description:  "Build and push monitoring plugin image",
	Capabilities: []string{"container"},
	Dir:          dir,
	Config: map[string]string{
		"dockerfile": "Dockerfile.dev",
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
	component.Register(BuildPush)
	component.Register(SetMCOUnmanaged)
	component.Register(ScaleDownMP)
	component.Register(PatchCMO)
	component.Register(ScaleUpMP)

	strategy.RegisterSelector(deploySelector)
}

func deploySelector(comp *component.Component, mode string) (strategy.BuildStrategy, strategy.RunStrategy) {
	switch comp.Name {
	case "set-mco-unmanaged":
		return nil, &setMCOStrategy{}
	case "scale-down-mp":
		return nil, &scaleDownMPStrategy{}
	case "patch-cmo":
		return nil, &patchCMOStrategy{}
	case "scale-up-mp":
		return nil, &scaleUpMPStrategy{}
	}
	return nil, nil
}

type setMCOStrategy struct{}

func (s *setMCOStrategy) Name() string        { return "set-mco-unmanaged" }
func (s *setMCOStrategy) Requires() []string { return []string{"oc"} }
func (s *setMCOStrategy) Run(_ context.Context, comp *component.Component, _ *runcontext.RunContext) (*component.Step, error) {
	return &component.Step{
		Name:      comp.Name,
		DependsOn: comp.DependsOn,
		Processes: []component.ProcessSpec{{
			Name:    comp.Name,
			Command: "oc",
			Args:    []string{"patch", "clusterversion", "version", "--type", "json", "-p", "{{content:mco-patch}}"},
			Files: map[string]component.FileRef{
				"mco-patch": {FS: filesFS, Path: "files/set-mco-to-unmanaged.yaml"},
			},
		}},
	}, nil
}

type scaleDownMPStrategy struct{}

func (s *scaleDownMPStrategy) Name() string        { return "scale-down-mp" }
func (s *scaleDownMPStrategy) Requires() []string { return []string{"oc"} }
func (s *scaleDownMPStrategy) Run(_ context.Context, comp *component.Component, _ *runcontext.RunContext) (*component.Step, error) {
	return &component.Step{
		Name:      comp.Name,
		DependsOn: comp.DependsOn,
		Processes: []component.ProcessSpec{
			{
				Name:    "scale down cmo",
				Command: "oc",
				Args:    []string{"scale", "--replicas=0", "-n", "openshift-monitoring", "deployment/cluster-monitoring-operator"},
			},
			{
				Name:    "scale down monitoring plugin",
				Command: "oc",
				Args:    []string{"scale", "--replicas=0", "-n", "openshift-monitoring", "deployment/monitoring-plugin"},
			},
		},
	}, nil
}

type patchCMOStrategy struct{}

func (s *patchCMOStrategy) Name() string        { return "patch-cmo" }
func (s *patchCMOStrategy) Requires() []string { return []string{"oc", "jq", "bash"} }
func (s *patchCMOStrategy) Run(_ context.Context, comp *component.Component, rc *runcontext.RunContext) (*component.Step, error) {
	image := rc.Get("mp-build-push", "image")
	if image == "" {
		return nil, fmt.Errorf("mp-build-push image not found in RunContext")
	}

	return &component.Step{
		Name:      comp.Name,
		DependsOn: comp.DependsOn,
		Processes: []component.ProcessSpec{{
			Name:    comp.Name,
			Command: "bash",
			Args:    []string{"-c", "{{content:patch-cmo-script}}"},
			Env:     map[string]string{"MP_IMAGE": image},
			Files: map[string]component.FileRef{
				"patch-cmo-script": {FS: filesFS, Path: "files/patch-cmo.sh"},
			},
		}},
	}, nil
}

type scaleUpMPStrategy struct{}

func (s *scaleUpMPStrategy) Name() string        { return "scale-up-mp" }
func (s *scaleUpMPStrategy) Requires() []string { return []string{"oc"} }
func (s *scaleUpMPStrategy) Run(_ context.Context, comp *component.Component, _ *runcontext.RunContext) (*component.Step, error) {
	return &component.Step{
		Name:      comp.Name,
		DependsOn: comp.DependsOn,
		Processes: []component.ProcessSpec{
			{
				Name:    "scale up cmo",
				Command: "oc",
				Args:    []string{"scale", "--replicas=1", "-n", "openshift-monitoring", "deployment/cluster-monitoring-operator"},
			},
			{
				Name:    "scale up monitoring plugin",
				Command: "oc",
				Args:    []string{"scale", "--replicas=1", "-n", "openshift-monitoring", "deployment/monitoring-plugin"},
			},
		},
	}, nil
}
