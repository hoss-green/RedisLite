package main

import (
	"encoding/hex"
	"errors"
	"io"
	"log"
	"net"
	"os"

	"redislite/app/commands/parsers"
	"redislite/app/data"
	"redislite/app/setup"
)

func handleInboundConnection(conn *net.Conn, server *setup.Server) {
	defer (*conn).Close()
	loop := true
	firstrun := !server.Settings.Master
	for {
		inputlen, inputBuffer := readcomm(conn)
		if inputBuffer == nil {
			break
		}

		if firstrun {
			log.Println("First Run")
			inputlen, inputBuffer = readRdbFile(inputlen, *inputBuffer)
			firstrun = false
			if inputlen == 0 {
				log.Println("oops")
				continue
			}
		}

		message := parsers.ReadString(*inputBuffer, inputlen)
		log.Printf("Mesrec: %#v\r\n", message)
		items := parsers.SeparateByLineBreakAdv(message)
		log.Printf("MesrecItems: %#v\r\n", items)

		cscontainers, err := parsers.BreakIntoChunks(items)
		log.Printf("Containers in message: %d", len(cscontainers))
		if err != nil {
			log.Println(err)
		}
		for _, cscontainer := range cscontainers {
			log.Printf("CSC: %#v\r\n", cscontainer)
			items := cscontainer.CommandStrings
			params := []string{}
			paramlen := 0
			if len(items) > 1 {
				params = items[1:]
				paramlen = len(params)
			}
			redisCommand := data.RedisCommand{
				Command:      items[0],
				Params:       params,
				ParamLength:  paramlen,
				MessageBytes: cscontainer.CommandStringBytes,
			}

			log.Printf("CSC: %#v\r\n", redisCommand)
			loop = parsers.ParseCommand(conn, redisCommand, server)
			if !loop {
				break

			}
		}

		if !loop {
			break
		}
	}

	log.Println("Connection Closed")
}

func readRdbFile(inputlen int, inputBuffer []byte) (int, *[]byte) {
	bytearray := inputBuffer[:inputlen]
	rdb := hex.EncodeToString(bytearray)
	log.Println(rdb)
	for index, databyte := range bytearray {
		if databyte == 0xFF {
			log.Println("Termination byte set, success")
			cutbuffer := bytearray[index+1 : inputlen]
			log.Printf("size:%d\r\n", len(cutbuffer))
			if len(cutbuffer) == 8 {
				log.Println("only bytes remaining were the checksum")
				return 0, &[]byte{}
			}
			cutbuffer = cutbuffer[8:]
			return len(cutbuffer), &cutbuffer
		}
	}

	return inputlen, &inputBuffer
	// if bytearray[inputlen - 1] != 0xFF {
	// }

}

func readcomm(conn *net.Conn) (int, *[]byte) {
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
