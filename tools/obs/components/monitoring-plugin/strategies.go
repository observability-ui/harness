package mp

import (
	"context"
	"fmt"

	"obs/internal/component"
	"obs/internal/runcontext"
	"obs/internal/strategy"
)

func init() {
	strategy.Register(InstallDeps.Name, &strategy.LocalNPM{Cmd: "install", Args: []string{"--no-save"}})
	strategy.Register(SetMCOUnmanaged.Name, &setMCOStrategy{})
	strategy.Register(ScaleDownMP.Name, &scaleDownMPStrategy{})
	strategy.Register(PatchCMO.Name, &patchCMOStrategy{})
	strategy.Register(ScaleUpMP.Name, &scaleUpMPStrategy{})
}

type setMCOStrategy struct{}

func (s *setMCOStrategy) Requires() []string { return []string{"oc"} }
func (s *setMCOStrategy) Execute(_ context.Context, comp *component.Component, _ *runcontext.RunContext) (*component.Step, error) {
	return &component.Step{
		Name:      comp.Name,
		Lifecycle: component.LifecycleOneShot,
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

func (s *scaleDownMPStrategy) Requires() []string { return []string{"oc"} }
func (s *scaleDownMPStrategy) Execute(_ context.Context, comp *component.Component, _ *runcontext.RunContext) (*component.Step, error) {
	return &component.Step{
		Name:      comp.Name,
		Lifecycle: component.LifecycleOneShot,
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

func (s *patchCMOStrategy) Requires() []string { return []string{"oc", "jq", "bash"} }
func (s *patchCMOStrategy) Execute(_ context.Context, comp *component.Component, rc *runcontext.RunContext) (*component.Step, error) {
	image := rc.Get("mp-build-push", "image")
	if image == "" {
		return nil, fmt.Errorf("mp-build-push image not found in RunContext")
	}

	return &component.Step{
		Name:      comp.Name,
		Lifecycle: component.LifecycleOneShot,
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

func (s *scaleUpMPStrategy) Requires() []string { return []string{"oc"} }
func (s *scaleUpMPStrategy) Execute(_ context.Context, comp *component.Component, _ *runcontext.RunContext) (*component.Step, error) {
	return &component.Step{
		Name:      comp.Name,
		Lifecycle: component.LifecycleOneShot,
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
