package systemparser 

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"redislite/app/data"
	"redislite/app/prototools/protomessages"
	"redislite/app/setup"
)

func replconf(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
	command := strings.ToUpper(redisCommand.Params[0])
	switch command {
	case "LISTENING-PORT":
		if len(redisCommand.Params) != 2 {
			return errors.New("listening-port requires a port number")
		}

		port, err := strconv.Atoi(redisCommand.Params[1])
		if err != nil {
			return errors.New("listening-port should be a number")
		}

		server.Settings.MasterPort = port
		return protomessages.QuickSendSimpleString(conn, "OK")
	case "CAPA":
		return protomessages.QuickSendSimpleString(conn, "OK")
	case "ACK":
		if server.Settings.Master {
			log.Println("ACK RECIEVED")
			server.AckChannel <- true
		}
		return nil
	case "GETACK":
		log.Println("GETACK RECIEVED")
		recievedBytes := fmt.Sprintf("%d", server.RecievedCounter.Total())
		log.Println("recievedBytes: ", server.RecievedCounter.Total())
		response := protomessages.BuildRespArrayMsg([]string{"REPLCONF", "ACK", recievedBytes})
		server.RecievedCounter.AddBytes(redisCommand.MessageBytes)
		return protomessages.QuickSendMessage(conn, response)
	default:
		return errors.New(fmt.Sprintf("Command not recognised or malformed: %s", redisCommand.Params[0]))
	}

}
