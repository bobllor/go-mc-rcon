package cmd

import (
	"fmt"
	"os"
	"rcon/rcon"

	"github.com/spf13/cobra"
)

var addFlags = &addData{}

type addData struct {
	serverInfo     ServerMeta
	forceOverwrite bool
}

var addCmd = &cobra.Command{
	Use:   "add [--tag|-t string] [--host string] [--password|-p string]",
	Short: "Add a server entry to the config",
	Run: func(cmd *cobra.Command, args []string) {
		if _, ok := addFlags.serverInfo.Config.Servers[addFlags.serverInfo.ServerTag]; ok && !addFlags.forceOverwrite {
			fmt.Printf(
				`Existing entry found for %s%s`,
				addFlags.serverInfo.ServerTag,
				"\n")
			fmt.Printf(`To force overwrite include the flags -f or --force%s`, "\n")
			fmt.Printf(
				`To remove the existing entry run gorcon remove -t %s%s`,
				addFlags.serverInfo.ServerTag, "\n")
			os.Exit(1)
		}

		if addFlags.serverInfo.ServerAuth.Password == "-" {
			password, err := readPassword()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			addFlags.serverInfo.ServerAuth.Password = password
		}

		err := addFlags.serverInfo.Config.AddServerEntry(
			addFlags.serverInfo.ServerTag, addFlags.serverInfo.ServerAuth.Host,
			addFlags.serverInfo.ServerAuth.Password)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Added entry %s\n", addFlags.serverInfo.ServerTag)
	},
}

func InitializeAddCmd(yamlConfig *rcon.Config) {
	addFlags.serverInfo.Config = yamlConfig

	addCmd.Flags().StringVarP(
		&addFlags.serverInfo.ServerTag, "tag", "t", "", "The identification of the server")
	addCmd.Flags().StringVar(&addFlags.serverInfo.ServerAuth.Host, "host", "", "The host address of the server")
	addCmd.Flags().StringVarP(
		&addFlags.serverInfo.ServerAuth.Password, "password", "p",
		"", `The password for RCON access, pass "-" to type a secure input`)
	addCmd.Flags().BoolVarP(&addFlags.forceOverwrite, "force", "f",
		false, "Overwrites an existing entry without checks")

	addCmd.MarkFlagsOneRequired("tag", "host", "password")
	addCmd.MarkFlagsRequiredTogether("tag", "host", "password")
	rootCmd.AddCommand(addCmd)
}
