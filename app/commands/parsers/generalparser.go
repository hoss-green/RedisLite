package parsers

import (
	"fmt"
	"log"
	"net"
	"redislite/app/commands/parsers/cmdparsers/genericparser"
	"redislite/app/commands/parsers/cmdparsers/stringparser"
	"redislite/app/commands/parsers/cmdparsers/systemparser"
	"redislite/app/commands/parsers/parserentities"
	"redislite/app/data"
	"redislite/app/prototools/protomessages"
	"redislite/app/setup"
	"strings"
)

func ParseCommand(connpointer *net.Conn, redisCommand data.RedisCommand, server *setup.Server) bool {
	parserInfo := parserentities.ParserInfo{}
	instruction := strings.ToUpper(redisCommand.Command)
	if instruction == "EXIT" || instruction == "QUIT" {
		return false
	}

  parserInfo = systemparser.ParseSystemCommand(connpointer, redisCommand, server)
	if parserInfo.Executed {
		return true
	}
	if parserInfo.Err != nil {
		goto parser
	}
	parserInfo = genericparser.ParseGenericCommand(connpointer, redisCommand, server)
	if parserInfo.Executed {
		return true
	}
	if parserInfo.Err != nil {
		goto parser
	}
	parserInfo = stringparser.ParseStringCommand(connpointer, redisCommand, server)
	if parserInfo.Executed {
		return true
	}
	// if parserInfo.Err != nil {
	// 	goto parser
	// }
	// parserInfo = streamparser.ParseStreamCommand(connpointer, redisCommand, server)
	// if parserInfo.Executed {
	// 	return true
	// }
parser:

	if parserInfo.Err != nil {
		log.Printf("Command Error: %s.\n\r", instruction)
		err := parserInfo.Err
		protomessages.QuickSendError(*connpointer, fmt.Sprintf("%s", err))
		return true //keep looping anyway
	}

	unrec := fmt.Sprintf("Command not recognised: %s.", instruction)
	log.Println(unrec)
	protomessages.QuickSendError(*connpointer, unrec)
	return true
}
