package main

import (
	"log"
	"net"
	"os"

	"redislite/app/setup"
)


func main() {
	server := setup.CreateServer()
  go ProcessEventLoop(server)
	log.SetFlags(0)
	if !server.Settings.Master {
		log.SetPrefix("REPLICA: ")
		initReplica(server)
	} else {
		log.SetPrefix("MASTER: ")
	}
	log.Println("launching server")
	initMaster(server)
}

func initMaster(server *setup.Server) {
  serverSettings := server.Settings
	log.Printf("Initialising Server on: %s\r\n", serverSettings.Host)
	listener, err := net.Listen("tcp", serverSettings.Host)
	if err != nil {
		log.Printf("Failed to bind to port: %d", serverSettings.Port)
		os.Exit(1)
	}

	createMaster(listener, server)
}
