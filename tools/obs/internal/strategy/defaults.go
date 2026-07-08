package strategy

import "obs/internal/component"

func init() {
	RegisterSelector(defaultSelector)
}

func defaultSelector(comp *component.Component, mode string) (BuildStrategy, RunStrategy) {
	if comp.Config["dockerfile"] != "" {
		return &ContainerBuild{}, nil
	}

	if comp.Config["npm-cmd"] != "" {
		return &LocalNPMInstall{}, nil
	}

	if comp.Config["make-target"] != "" {
		return nil, &LocalMakeRun{}
	}

	if comp.Config["compose-file"] != "" {
		return nil, &PodmanCompose{}
	}

	if comp.Config["oc-args"] != "" {
		return nil, &OCCommand{}
	}

	if comp.Config["script-content"] != "" {
		return nil, &BashScript{}
	}

	return nil, nil
}
