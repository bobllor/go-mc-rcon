package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"rcon/rcon"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

type ServerMeta struct {
	ServerTag  string
	ServerAuth rcon.Server
	Config     *rcon.Config
}

type rootData struct {
	ServerInfo ServerMeta
}

var rootFlags = &rootData{}

var rootCmd = &cobra.Command{
	Use: "command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("test")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
	}
}

// InitializeServer initializes the rootData struct with the arguments given.
func InitializeRootCmd(yamlConfig *rcon.Config) {
	rootFlags.ServerInfo.Config = yamlConfig
}

func readPassword() (string, error) {
	stdin := int(syscall.Stdin)

	oldState, err := term.GetState(stdin)
	if err != nil {
		return "", err
	}
	defer term.Restore(stdin, oldState)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		for _ = range ch {
			term.Restore(stdin, oldState)
			os.Exit(1)
		}
	}()

	fmt.Print("Password: ")
	pwBytes, err := term.ReadPassword(stdin)
	if err != nil {
		return "", err
	}

	// formatting, yes...
	fmt.Println()

	return string(pwBytes), nil
}
