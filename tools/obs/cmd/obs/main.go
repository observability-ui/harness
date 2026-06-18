package main

import (
	"obs/internal/cli"
	_ "obs/recipes" // registers all built-in recipes
)

var version = "dev"

func main() {
	cli.Execute(version)
}
