package mp

import (
	"obs/internal/recipe"

	"github.com/spf13/pflag"
)

type DeployMonitoringPlugin struct{}

func (r *DeployMonitoringPlugin) Name() string      { return "monitoring-plugin" }
func (r *DeployMonitoringPlugin) Aliases() []string { return []string{"mp"} }
func (r *DeployMonitoringPlugin) Description() string {
	return "Deploy monitoring plugin image in a running OpenShift cluster"
}

func (r *DeployMonitoringPlugin) Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("monitoring-plugin", pflag.ContinueOnError)
	fs.String("image", "", "container image to build and push (e.g. quay.io/user/monitoring-plugin:tag)")
	return fs
}

func (r *DeployMonitoringPlugin) Requirements(flags *pflag.FlagSet) []recipe.Requirement {
	return []recipe.Requirement{
		recipe.RequireFlag(flags, "image", "container image to build and push"),
		recipe.RequireOCLogin(),
		recipe.RequirePodman(),
		recipe.RequireJQ(),
	}
}

func (r *DeployMonitoringPlugin) Steps(cfg *recipe.Config) ([]*recipe.Step, error) {
	dir := "projects/monitoring-plugin"

	image, _ := cfg.Flags.GetString("image")

	return []*recipe.Step{
		{
			Name:      "build-push-mp-image",
			DependsOn: []string{},
			Processes: []recipe.ProcessSpec{
				{
					Name:    "build mp image",
					Command: "podman",
					Args:    []string{"build", "-f", "Dockerfile.dev", "--platform=linux/amd64", "-t", image},
					Dir:     dir,
				},
				{
					Name:    "push mp image",
					Command: "podman",
					Args:    []string{"push", image},
					Dir:     dir,
				},
			},
		},
		{
			Name:      "set-mco-to-unmanaged",
			DependsOn: []string{"build-push-mp-image"},
			Processes: []recipe.ProcessSpec{
				{
					Name:    "set mco to unmanaged",
					Command: "oc",
					Args:    []string{"patch", "clusterversion", "version", "--type", "json", "-p", "{{content:mco-patch}}"},
					Files: map[string]recipe.FileRef{
						"mco-patch": {FS: filesFS, Path: "files/set-mco-to-unmanaged.yaml"},
					},
				},
			},
		},
		{
			Name:      "scale-down-cmo",
			DependsOn: []string{"set-mco-to-unmanaged"},
			Processes: []recipe.ProcessSpec{
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
		},
		{
			Name:      "patch-cmo",
			DependsOn: []string{"scale-down-cmo"},
			Processes: []recipe.ProcessSpec{
				{
					Name:    "patch cmo",
					Command: "bash",
					Args:    []string{"-c", "{{content:patch-cmo-script}}"},
					Env:     map[string]string{"MP_IMAGE": image},
					Files: map[string]recipe.FileRef{
						"patch-cmo-script": {FS: filesFS, Path: "files/patch-cmo.sh"},
					},
				},
			},
		},
		{
			Name:      "scale-up-cmo",
			DependsOn: []string{"patch-cmo"},
			Processes: []recipe.ProcessSpec{
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
		},
	}, nil
}
