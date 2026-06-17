package main

import "obsui/internal/cli"

var version = "dev"

func main() {
	cli.Execute(version)
}
