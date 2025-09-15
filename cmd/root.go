package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type rootData struct {
	yamlDir string
}

var rootFlags = &rootData{}

var rootCmd = &cobra.Command{
	Use: "command",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		fmt.Println("hi")
		os.Exit(1)
	}
}

// InitializeServer initializes the rootData struct with the arguments given.
func InitializeRootCmd(yamlDirectory string) {
	rootFlags.yamlDir = yamlDirectory
}
