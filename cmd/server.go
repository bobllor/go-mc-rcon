package cmd

import "rcon/rcon"

var Host string

// NOTE: maybe having a "server" subcommand?

func init() {
	serversMap := rcon.YamlConfig.Servers
	network := serversMap[rcon.YamlConfig.Default_Server].Host

	/*
		TODO:
		 	1. password flag, required if choosing --host (does not work alone)
		 	2. name flag, if used then it will go through the yaml config instead
			3. list server entries
			4. add server entry
		 	5. remove server entry
		 	6.
	*/

	rootCmd.PersistentFlags().StringVar(&Host, "host", network, "Set the server's address")
}
