package main

import (
	"github.com/cluion/paq/internal/cli"
	_ "github.com/cluion/paq/internal/provider"
)

func main() {
	cli.RegisterProviderCommands()
	cli.Execute()
}
