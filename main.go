package main

import (
	"os"
	"rcon/cmd"
)

// TODO: remove all panic calls later

func main() {
	yamlDir := os.Getenv("HOME") + "/.config/.mcrcon"

	cmd.InitializeServerCmd(yamlDir)
	cmd.InitializeRootCmd(yamlDir)

	cmd.Execute()
}
