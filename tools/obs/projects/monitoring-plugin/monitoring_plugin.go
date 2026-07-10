package mp

import (
	"context"
	"fmt"

	"obs/internal/runcontext"
	"obs/internal/strategy"
	"obs/internal/task"
)

var dir = "projects/monitoring-plugin"

var InstallDeps = &task.Task{
	Name:     "mp-install-deps",
	Dir:      dir + "/web",
	Strategy: strategy.NPMRun("install", "--no-save"),
}

var Frontend = &task.Task{
	Name:      "mp-frontend",
	DependsOn: []string{"mp-install-deps"},
	Dir:       dir,
	Ports:     []int{9001},
	Strategy:  strategy.MakeTarget("start-frontend"),
}

var Backend = &task.Task{
	Name:      "mp-backend",
	DependsOn: []string{"mp-install-deps"},
	Dir:       dir,
	Ports:     []int{9443},
	Labels:    map[string]string{"console-plugin": "monitoring-plugin"},
	Strategy:  strategy.MakeTarget("start-coo-backend,start-feature-backend"),
}

var BuildPush = &task.Task{
	Name:     "mp-build-push",
	Dir:      dir,
	Strategy: strategy.DockerBuild("Dockerfile.dev"),
	RequiredFlags: []task.RequiredFlag{
		{Name: "image", Usage: "container image to build and push (e.g. quay.io/user/monitoring-plugin:tag)"},
	},
}

var SetMCOUnmanaged = &task.Task{
	Name:      "set-mco-unmanaged",
	DependsOn: []string{"mp-build-push"},
	Strategy:  &setMCOStrategy{},
}

var ScaleDownMP = &task.Task{
	Name:      "scale-down-mp",
	DependsOn: []string{"set-mco-unmanaged"},
	Strategy:  &scaleDownMPStrategy{},
}

var PatchCMO = &task.Task{
	Name:      "patch-cmo",
	DependsOn: []string{"scale-down-mp"},
	Strategy:  &patchCMOStrategy{},
}

var ScaleUpMP = &task.Task{
	Name:      "scale-up-mp",
	DependsOn: []string{"patch-cmo"},
	Strategy:  &scaleUpMPStrategy{},
}

func init() {
	task.Register(InstallDeps)
	task.Register(Frontend)
	task.Register(Backend)
	task.Register(BuildPush)
	task.Register(SetMCOUnmanaged)
	task.Register(ScaleDownMP)
	task.Register(PatchCMO)
	task.Register(ScaleUpMP)
}

type setMCOStrategy struct{}

func (s *setMCOStrategy) Requires() []string { return []string{"oc"} }
func (s *setMCOStrategy) Execute(_ context.Context, t *task.Task, _ *runcontext.RunContext) (*task.Step, error) {
	return &task.Step{
		Name:      t.Name,
		Lifecycle: task.LifecycleOneShot,
		DependsOn: t.DependsOn,
		Processes: []task.ProcessSpec{{
			Name:    t.Name,
			Command: "oc",
			Args:    []string{"patch", "clusterversion", "version", "--type", "json", "-p", "{{content:mco-patch}}"},
			Files: map[string]task.FileRef{
				"mco-patch": {FS: filesFS, Path: "files/set-mco-to-unmanaged.yaml"},
			},
		}},
	}, nil
}

type scaleDownMPStrategy struct{}

func (s *scaleDownMPStrategy) Requires() []string { return []string{"oc"} }
func (s *scaleDownMPStrategy) Execute(_ context.Context, t *task.Task, _ *runcontext.RunContext) (*task.Step, error) {
	return &task.Step{
		Name:      t.Name,
		Lifecycle: task.LifecycleOneShot,
		DependsOn: t.DependsOn,
		Processes: []task.ProcessSpec{
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
func (s *patchCMOStrategy) Execute(_ context.Context, t *task.Task, rc *runcontext.RunContext) (*task.Step, error) {
	image := rc.Get("mp-build-push", "image")
	if image == "" {
		return nil, fmt.Errorf("mp-build-push image not found in RunContext")
	}

	return &task.Step{
		Name:      t.Name,
		Lifecycle: task.LifecycleOneShot,
		DependsOn: t.DependsOn,
		Processes: []task.ProcessSpec{{
			Name:    t.Name,
			Command: "bash",
			Args:    []string{"-c", "{{content:patch-cmo-script}}"},
			Env:     map[string]string{"MP_IMAGE": image},
			Files: map[string]task.FileRef{
				"patch-cmo-script": {FS: filesFS, Path: "files/patch-cmo.sh"},
			},
		}},
	}, nil
}

type scaleUpMPStrategy struct{}

func (s *scaleUpMPStrategy) Requires() []string { return []string{"oc"} }
func (s *scaleUpMPStrategy) Execute(_ context.Context, t *task.Task, _ *runcontext.RunContext) (*task.Step, error) {
	return &task.Step{
		Name:      t.Name,
		Lifecycle: task.LifecycleOneShot,
		DependsOn: t.DependsOn,
		Processes: []task.ProcessSpec{
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
