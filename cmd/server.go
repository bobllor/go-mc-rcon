package cmd

var Host string

// NOTE: maybe having a "server" subcommand?

func init() {
	/*
		TODO:
		 	1. password flag, required if choosing --host (does not work alone)
		 	2. name flag, if used then it will go through the yaml config instead
			3. list server entries
			4. add server entry
		 	5. remove server entry
		 	6.
	*/

	rootCmd.PersistentFlags().StringVar(&Host, "host", "", "Set the server's address")
}
