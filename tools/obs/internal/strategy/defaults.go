package strategy

import "obs/internal/component"

func resolveByConfig(comp *component.Component) []Strategy {
	if comp.Config["dockerfile"] != "" {
		return []Strategy{&ContainerBuild{}}
	}

	if comp.Config["make-target"] != "" {
		return []Strategy{&LocalMakeRun{}}
	}

	if comp.Config["compose-file"] != "" {
		return []Strategy{&PodmanCompose{}}
	}

	return nil
}
