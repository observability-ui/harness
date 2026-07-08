package main

import (
	_ "obs/components" // registers components and mixer recipes
	"obs/internal/cli"
)

var version = "dev"

func main() {
	cli.Execute(version)
}
