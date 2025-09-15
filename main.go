package main

import (
	"os"
	"rcon/cmd"
)

// TODO: support for windows

func main() {
	yamlDir := os.Getenv("HOME") + "/.config/.mcrcon"

	cmd.InitializeServerCmd(yamlDir)
	cmd.InitializeRootCmd(yamlDir)
	cmd.InitializeAddCmd()

	cmd.Execute()
}
