package genericparser

import (
	"errors"
	"net"
	"strings"

	"redislite/app/commands/errormessages"
	"redislite/app/commands/parsers/parserentities"
	"redislite/app/data"
	"redislite/app/setup"
)

// STREAM COMMAND PARSER
func ParseGenericCommand(connpointer *net.Conn, redisCommand data.RedisCommand, server *setup.Server) parserentities.ParserInfo {
	conn := *connpointer
	var err error = nil
	instruction := strings.ToUpper(redisCommand.Command)
	switch instruction {
	case "DEL":
		if redisCommand.ParamLength < 1 {
			err = errors.New(errormessages.IncorrectArgumentsError)
			break
		}
		err = doDel(conn, server, redisCommand)
	case "COPY":
		if redisCommand.ParamLength != 2 {
			err = errors.New(errormessages.IncorrectArgumentsError)
			break
		}
		err = doCopy(conn, server, redisCommand)
	case "EXISTS":
		if redisCommand.ParamLength < 1 {
			err = errors.New(errormessages.IncorrectArgumentsError)
			break
		}
		err = doExists(conn, server, redisCommand)
	default:
		return parserentities.ParserInfo{Executed: false, Err: nil}
	}

	if err != nil {
		return parserentities.ParserInfo{Executed: false, Err: err}
	}

	return parserentities.ParserInfo{Executed: true, Err: nil}
}
