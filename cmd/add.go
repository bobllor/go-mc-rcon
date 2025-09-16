package cmd

import (
	"rcon/rcon"

	"github.com/spf13/cobra"
)

var addFlags = &ServerMeta{}

var addCmd = &cobra.Command{
	Use:   "add [--tag|-t string] [--host string] [--password|-p string]",
	Short: "Add a server entry to the config",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func InitializeAddCmd(yamlConfig *rcon.Config) {
	addFlags.Config = yamlConfig

	addCmd.Flags().StringVarP(
		&addFlags.ServerTag, "tag", "t", "", "The identification of the server")
	addCmd.Flags().StringVar(&addFlags.ServerAuth.Host, "host", "", "The host address of the server")
	addCmd.Flags().StringVarP(
		&addFlags.ServerAuth.Password, "password", "p",
		"", `The password for RCON access, pass "-" to type a secure input`)

	addCmd.MarkFlagsOneRequired("tag", "host", "password")
	addCmd.MarkFlagsRequiredTogether("tag", "host", "password")
	rootCmd.AddCommand(addCmd)
}
