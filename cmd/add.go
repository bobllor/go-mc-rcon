package cmd

import (
	"fmt"
	"os"
	"rcon/rcon"

	"github.com/spf13/cobra"
)

var addFlags = &ServerMeta{}

var addCmd = &cobra.Command{
	Use:   "add [--tag|-t string] [--host string] [--password|-p string]",
	Short: "Add a server entry to the config",
	Run: func(cmd *cobra.Command, args []string) {
		if addFlags.ServerAuth.Password == "-" {
			password, err := readPassword()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			addFlags.ServerAuth.Password = password
		}

		err := addFlags.Config.AddServerEntry(
			addFlags.ServerTag, addFlags.ServerAuth.Host, addFlags.ServerAuth.Password)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Added entry %s\n", addFlags.ServerTag)
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
