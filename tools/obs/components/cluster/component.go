package cluster

import "obs/internal/component"

var SeedUsers = &component.Component{
	Name:         "seed-users",
	Description:  "Create test users with standard RBAC roles",
	Capabilities: []string{"cluster"},
	Config: map[string]string{
		"oc-args": "apply -f -",
	},
}

var ScaleDownCMO = &component.Component{
	Name:         "scale-down-cmo",
	Description:  "Scale down cluster monitoring operator",
	Capabilities: []string{"cluster"},
	Config: map[string]string{
		"oc-args": "scale --replicas=0 -n openshift-monitoring deployment/cluster-monitoring-operator",
	},
}

var ScaleUpCMO = &component.Component{
	Name:         "scale-up-cmo",
	Description:  "Scale up cluster monitoring operator",
	Capabilities: []string{"cluster"},
	Config: map[string]string{
		"oc-args": "scale --replicas=1 -n openshift-monitoring deployment/cluster-monitoring-operator",
	},
}

func init() {
	component.Register(SeedUsers)
	component.Register(ScaleDownCMO)
	component.Register(ScaleUpCMO)
}
