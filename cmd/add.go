package cmd

import (
	"github.com/spf13/cobra"
)

var addFlags = &ServerMeta{}

var addCmd = &cobra.Command{
	Use:   "add [--tag|-t string] [--host string] [--password|-p string]",
	Short: "add a server entry to the config",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func InitializeAddCmd() {
	addCmd.Flags().StringVarP(
		&addFlags.ServerTag, "tag", "t", "", "the server's tag for identification")
	addCmd.Flags().StringVar(&addFlags.ServerAuth.Host, "host", "", "the host address of the server")
	addCmd.Flags().StringVarP(
		&addFlags.ServerAuth.Password, "password", "p", "", "the password of the server")

	addCmd.MarkFlagsRequiredTogether("tag", "host", "password")
	rootCmd.AddCommand(addCmd)
}
