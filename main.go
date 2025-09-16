package main

import (
	"fmt"
	"os"
	"rcon/cmd"
	"rcon/rcon"
)

// TODO: support for windows

func main() {
	yamlDir := os.Getenv("HOME") + "/.config/.mcrcon"

	config, err := rcon.NewConfig(yamlDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd.InitializeServerCmd(config)
	cmd.InitializeRootCmd(config)
	cmd.InitializeAddCmd(config)

	cmd.Execute()
}
