package cmdparsers

import (
	"errors"
	"net"
	"strings"

	"redislite/app/commands/errormessages"
	"redislite/app/commands/instructions"
	"redislite/app/data"
	"redislite/app/setup"
)

// STREAM COMMAND PARSER
func ParseStreamCommand(connpointer *net.Conn, redisCommand data.RedisCommand, server *setup.Server) ParserInfo {
	conn := *connpointer
	var err error = nil
	instruction := strings.ToUpper(redisCommand.Command)
	switch instruction {
	case "XADD":
		if redisCommand.ParamLength < 2 || redisCommand.ParamLength % 2 != 0{
			err = errors.New(errormessages.IncorrectArgumentsError)
			break
		}
    err = instructions.XAdd(conn, server, redisCommand)
  case "XRANGE":
    if redisCommand.ParamLength != 3 {
			err = errors.New(errormessages.IncorrectArgumentsError)
			break
    }
    err = instructions.XRange(conn, server, redisCommand)
  case "XREAD":
    if redisCommand.ParamLength < 3 {
			err = errors.New(errormessages.IncorrectArgumentsError)
			break
    }
    err = instructions.XRead(conn, server, redisCommand)
	default:
		return ParserInfo{Executed: false, Err: nil}
	}

	if err != nil {
		return ParserInfo{Executed: false, Err: err}
	}

	return ParserInfo{Executed: true, Err: nil}
}
