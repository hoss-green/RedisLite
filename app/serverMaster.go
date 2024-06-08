package main

import (
	"log"
	"net"
	"os"

	"redislite/app/setup"
)

func createMaster(listener net.Listener, server *setup.Server) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		} else {
			log.Println("Connection received from: ", conn.RemoteAddr().String())
		}

		go handleInboundConnection(&conn, server)
	}
}

