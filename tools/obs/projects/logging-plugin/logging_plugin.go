package lp

import (
	"obs/internal/strategy"
	"obs/internal/task"
)

var dir = "projects/logging-view-plugin"

var InstallDeps = &task.Task{
	Name:     "lp-install-deps",
	Dir:      dir + "/web",
	Strategy: strategy.NPMRun("install", "--no-save"),
}

var Frontend = &task.Task{
	Name:      "lp-frontend",
	DependsOn: []string{"lp-install-deps"},
	Dir:       dir,
	Ports:     []int{9001},
	Strategy:  strategy.MakeTarget("start-frontend"),
}

var Backend = &task.Task{
	Name:      "lp-backend",
	DependsOn: []string{"lp-frontend"},
	Dir:       dir,
	Ports:     []int{9002},
	Labels:    map[string]string{"console-plugin": "logging-view-plugin"},
	Strategy:  strategy.MakeTarget("start-backend"),
}

var LocalLoki = &task.Task{
	Name:     "lp-local-loki",
	Dir:      dir,
	Ports:    []int{3100},
	Strategy: strategy.Compose("hack/docker-compose/docker-compose.test.yml"),
}

func init() {
	task.Register(InstallDeps)
	task.Register(Frontend)
	task.Register(Backend)
	task.Register(LocalLoki)
}
