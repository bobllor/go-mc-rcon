package main

import (
	"fmt"
	"os"
	"rcon/cmd"
	"rcon/rcon"
)

// TODO: remove all panic calls later

func main() {
	cmd.Execute()
	yamlDir := os.Getenv("HOME") + "/.config/.mcrcon"

	config, err := rcon.NewConfig(yamlDir)
	if err != nil {
		panic(err)
	}

	// TODO: fix the name, it will use the name given by the flag instead (--name or -n).
	serverInfo, err := config.GetServer("default")
	if err != nil {
		panic(err)
	}

	rcon := rcon.NewRCON()

	conn, err := rcon.Connect(serverInfo.Host)
	if err != nil {
		panic(err)
	}

	err = conn.Authenticate(serverInfo.RCON_Password)
	if err != nil {
		panic(err)
	}

	// TODO: fix commmand, will use the command given by the flag (-c).
	cmdOutput, err := conn.Command("deop Notch")
	if err != nil {
		panic(err)
	}

	fmt.Println(cmdOutput)
}
