package main

import (
	"fmt"
	rcon "rcon/core"
	"time"
)

func main() {
	network := ""
	address := ""

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

	fmt.Println(cmdOutput)
}
