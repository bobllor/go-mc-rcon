package cmd

import (
	"errors"
	"fmt"
	"os"
	"rcon/rcon"

	"github.com/spf13/cobra"
)

type serverData struct {
	serverInfo ServerMeta
	command    string
	yamlDir    string
}

var serverFlags = &serverData{}

var serverCmd = &cobra.Command{
	Use:   "server [--name string|--host string --password|-p [string]] [-c string]",
	Short: "Execute a command on a server",
	Long: `Select a server with a tag or by passing in the host address and password 
to execute a command on the server`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := rcon.NewConfig(serverFlags.yamlDir)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if serverFlags.serverInfo.ServerTag != "" {
			err = serverFlags.runServerNameCommand(config)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		// password check is not needed since it is required if host is used.
		if serverFlags.serverInfo.ServerAuth.Host != "" {
			if serverFlags.serverInfo.ServerAuth.Password == "-" {
				password, err := readPassword()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				serverFlags.serverInfo.ServerAuth.Password = password
			}

			err = serverFlags.runServerHostCommand()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
}

// InitializeServer initializes the serverData struct with the arguments given.
func InitializeServerCmd(yamlDirectory string) {
	serverFlags.yamlDir = yamlDirectory

	serverCmd.Flags().StringVarP(
		&serverFlags.serverInfo.ServerTag, "tag", "t",
		"", "The tag of a server, requires an existing YAML entry")
	serverCmd.Flags().StringVar(
		&serverFlags.serverInfo.ServerAuth.Host, "host", "", "Host address of the server")
	serverCmd.Flags().StringVarP(
		&serverFlags.serverInfo.ServerAuth.Password, "password", "p",
		"", `The password for RCON access, pass "-" to prompt for secure input`)
	serverCmd.Flags().StringVarP(
		&serverFlags.command, "command", "c", "", "The command to be executed on the server")

	serverCmd.MarkFlagsOneRequired("host", "tag")
	serverCmd.MarkFlagsRequiredTogether("host", "password")
	serverCmd.MarkFlagsMutuallyExclusive("host", "tag")

	rootCmd.AddCommand(serverCmd)
}

// runServerNameCommand runs the command based off of the server name.
// This requires the given server name to exist in the YAML entry.
func (s *serverData) runServerNameCommand(yamlConfig *rcon.Config) error {
	if _, ok := yamlConfig.Servers[s.serverInfo.ServerTag]; !ok {
		errMsg := fmt.Sprintf("no entries found for %s", s.serverInfo.ServerTag)
		return errors.New(errMsg)
	}

	serverInfo := yamlConfig.Servers[s.serverInfo.ServerTag]

	serverHost := serverInfo.Host
	serverPassword := serverInfo.Password

	rcon := rcon.NewRCON()

	conn, err := rcon.Connect(serverHost)
	if err != nil {
		return err
	}

	err = conn.Authenticate(serverPassword)
	if err != nil {
		return err
	}

	cmdOutput, err := conn.Command(s.command)
	if err != nil {
		return err
	}

	fmt.Println(cmdOutput)
	return nil
}

// runServerNameCommand runs the command with the given server host address and password.
func (s *serverData) runServerHostCommand() error {
	rcon := rcon.NewRCON()

	conn, err := rcon.Connect(s.serverInfo.ServerAuth.Host)
	if err != nil {
		return err
	}

	err = conn.Authenticate(s.serverInfo.ServerAuth.Password)
	if err != nil {
		return err
	}

	cmdOutput, err := conn.Command(s.command)
	if err != nil {
		return err
	}

	fmt.Println(cmdOutput)
	return nil
}
