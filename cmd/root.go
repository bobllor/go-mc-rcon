package cmd

import (
	"rcon/rcon"

	"github.com/spf13/cobra"
)

type ServerMeta struct {
	ServerTag  string
	ServerAuth rcon.Server
}

type rootData struct {
	yamlDir string
}

var rootFlags = &rootData{}

var rootCmd = &cobra.Command{
	Use: "command",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
	}
}

// InitializeServer initializes the rootData struct with the arguments given.
func InitializeRootCmd(yamlDirectory string) {
	rootFlags.yamlDir = yamlDirectory
}
