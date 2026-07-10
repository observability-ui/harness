package main

import (
	_ "obs/projects" // registers tasks and mixer recipes
	"obs/internal/cli"
)

var version = "dev"

func main() {
	cli.Execute(version)
}
