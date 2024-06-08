package parsers

import (
	"fmt"
	"log"
	"net"
	"strings"

	"redislite/app/commands/parsers/cmdparsers"
	"redislite/app/data"
	"redislite/app/prototools/protomessages"
	"redislite/app/setup"
)

func ParseCommand(connpointer *net.Conn, redisCommand data.RedisCommand, server *setup.Server) bool {
	parserInfo := cmdparsers.ParserInfo{}
	instruction := strings.ToUpper(redisCommand.Command)
	// log.Println("Recieved: ", instruction)
	if instruction == "EXIT" || instruction == "QUIT" {
		return false
	}
	parserInfo = cmdparsers.ParseSystemCommand(connpointer, redisCommand, server)
	if parserInfo.Executed {
		return true
	}
	parserInfo = cmdparsers.ParseStringCommand(connpointer, redisCommand, server)
	if parserInfo.Executed {
		return true
	}
	parserInfo = cmdparsers.ParseStreamCommand(connpointer, redisCommand, server)
	if parserInfo.Executed {
		return true
	}

	if parserInfo.Err != nil {
    log.Printf("Command Error: %s.\n\r", instruction)
    err := parserInfo.Err
		protomessages.SendError(*connpointer, fmt.Sprintf("%s", err))
		return true //keep looping anyway
	}

	unrec := fmt.Sprintf("Command not recognised: %s.\n", instruction)
	log.Println(unrec)
	protomessages.SendError(*connpointer, unrec)
	return true
}
