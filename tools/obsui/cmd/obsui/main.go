package main

import (
	"obsui/internal/cli"
	_ "obsui/recipes" // registers all built-in recipes
)

var version = "dev"

func main() {
	cli.Execute(version)
}
