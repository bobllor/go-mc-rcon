package cmd

import (
	"errors"
	"fmt"
	"rcon/rcon"

	"github.com/spf13/cobra"
)

// TODO: remove all panics

type serverData struct {
	serverInfo ServerMeta
	command    string
	yamlDir    string
}

var serverFlags = &serverData{}

var serverCmd = &cobra.Command{
	Use:   "server [--name string|--host string --password|-p [string]] [-c string]",
	Short: "connect to a server and run a command",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := rcon.NewConfig(serverFlags.yamlDir)
		if err != nil {
			panic(err)
		}

		if serverFlags.serverInfo.ServerTag != "" {
			err = serverFlags.runServerNameCommand(config)
			if err != nil {
				panic(err)
			}
		}

		// password check is not needed since it is required if host is used.
		if serverFlags.serverInfo.ServerAuth.Host != "" {
			err = serverFlags.runServerHostCommand()
			if err != nil {
				panic(err)
			}
		}
	},
}

// InitializeServer initializes the serverData struct with the arguments given.
func InitializeServerCmd(yamlDirectory string) {
	serverFlags.yamlDir = yamlDirectory

	serverCmd.Flags().StringVarP(
		&serverFlags.serverInfo.ServerTag, "tag", "t",
		"", "the tag of a server, requires a YAML entry")
	serverCmd.Flags().StringVar(
		&serverFlags.serverInfo.ServerAuth.Host, "host", "", "host address of the server, requires password")
	serverCmd.Flags().StringVarP(
		&serverFlags.serverInfo.ServerAuth.Password, "password", "p",
		"", "password of the server, pass - to prompt for secure input")
	serverCmd.Flags().StringVarP(
		&serverFlags.command, "command", "c", "", "the command sent to the server through RCON")

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
