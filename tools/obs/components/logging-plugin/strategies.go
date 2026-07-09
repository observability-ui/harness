package lp

import "obs/internal/strategy"

func init() {
	strategy.Register(InstallDeps.Name, &strategy.LocalNPM{Cmd: "install", Args: []string{"--no-save"}})
}
