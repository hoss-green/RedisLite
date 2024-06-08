package main

import (
	// "encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"redislite/app/setup"
)

func initReplica(server *setup.Server) {
	conn := createMasterTCPConn(server.Settings)
	err := handshake(conn, server.Settings)
	if err != nil {
		log.Println("Error Handshaking with Master")
		os.Exit(1)
	}

	log.Println("Connecting to Master")
	go handleInboundConnection(conn, server)
	log.Println("Connected to Master")
}



func createMasterTCPConn(serverSettings setup.ServerSettings) *net.Conn {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverSettings.MasterHost, serverSettings.MasterPort))
	if err != nil {
		panic("Could not connect to port")
	}
	return &conn
}

func handshake(conn *net.Conn, serverSettings setup.ServerSettings) error {
	var err error
	response := ""
	{
		response, err = sendMessage(conn, "*1\r\n$4\r\nPING\r\n")
		if strings.ToUpper(response) != "+PONG" {
			return errors.New("PING FAILED")
		}
		log.Print("Response from ping: ", response)
	}

	{
		response, err = sendMessage(conn, fmt.Sprintf("*3\r\n$8\r\nREPLCONF\r\n$14\r\nlistening-port\r\n$4\r\n%d\r\n", serverSettings.Port))
		if strings.ToUpper(response) != "+OK" {
			return errors.New("REPLCONF failed on settings listening port")
		}
		log.Print("Response from REPLCONF listening port: ", response)
	}

	{
		response, err = sendMessage(conn, fmt.Sprintf("*3\r\n$8\r\nREPLCONF\r\n$4\r\ncapa\r\n$6\r\npsync2\r\n"))
		if strings.ToUpper(response) != "+OK" {
			return errors.New("REPLCONF failed to set config")
		}
		log.Print("Response from REPLCONF config: ", response)
	}

	{
		response, err = sendMessage(conn, fmt.Sprintf("*3\r\n$5\r\nPSYNC\r\n$1\r\n?\r\n$2\r\n-1\r\n"))
		if !strings.HasPrefix(strings.ToUpper(response), "+FULLRESYNC") {
			return errors.New("PSYC failed to init config")
		}
		log.Print("Response from PSYNC config: ", response)
	}

	return err
}


func readrdbcomm(conn *net.Conn) (int, *[]byte) {
	inputBuffer := make([]byte, 1024)
	inputlen, err := (*conn).Read(inputBuffer)
	if inputlen == 0 {
		return 0, nil
	}
	if errors.Is(err, io.EOF) {
		log.Println("Redis client connection closed")
		return 0, nil
	} else if err != nil {
		log.Println("Unrecoverable error: ", err.Error())
		os.Exit(1)
	}

	return inputlen, &inputBuffer
}

func sendMessage(conn *net.Conn, message string) (string, error) {
	_, err := fmt.Fprintf(*conn, message)
	if err != nil {
		panic("Could not send message")
	}
	readbuffer := make([]byte, 1024)
	_, readErr := (*conn).Read(readbuffer)
	if readErr != nil {
		panic("Could not read response message")
	}
	resp := strings.Split(string(readbuffer), "\r\n")
	return resp[0], nil
}

func parseReplicaParams(replicaOf string) (string, int, error) {
	var parts = strings.Split(replicaOf, " ")
	if len(parts) != 2 {
		return "", 0, errors.New("Wrong number of parameters")
	}

	var host string
	hostip := net.ParseIP(parts[0])
	if hostip != nil {
		host = string(hostip.String())
		fmt.Printf("host ip set, using: %s \r\n", host)
	} else {
		host = parts[0]
		fmt.Printf("host is not an ip address, attempting to use hostname: %s \r\n", host)
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, errors.New("Port is not a number")
	}

	return host, port, nil

}
