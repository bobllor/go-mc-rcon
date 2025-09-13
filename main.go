package main

import (
	"fmt"
	"rcon/cmd"
	"rcon/rcon"
)

func main() {
	cmd.Execute()
	fmt.Println(cmd.Host)

	fmt.Println(rcon.YamlConfig)
	/*
		network := "tcp"
		address := ""
		rcon := rcon.NewRCON()

		conn, err := rcon.Connect(network, address, time.Second*5)
		if err != nil {
			panic(err)
		}

		password := ""

		err = conn.Authenticate(password)
		if err != nil {
			panic(err)
		}

		cmdOutput, err := conn.Command("")
		if err != nil {
			panic(err)
		}

		fmt.Println(cmdOutput)*/
}
