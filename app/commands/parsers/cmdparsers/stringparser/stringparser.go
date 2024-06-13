package stringparser

import (
	"errors"
	"net"
	"strings"

	"redislite/app/commands/errormessages"
	"redislite/app/commands/parsers/cmdparsers/stringparser/addsub"
	"redislite/app/commands/parsers/cmdparsers/stringparser/doget"
	"redislite/app/commands/parsers/cmdparsers/stringparser/doset"
	"redislite/app/commands/parsers/cmdparsers/stringparser/multi"
	"redislite/app/commands/parsers/parserentities"
	"redislite/app/data"
	"redislite/app/setup"
)

// STRING COMMAND PARSER
func ParseStringCommand(connpointer *net.Conn, redisCommand data.RedisCommand, server *setup.Server) parserentities.ParserInfo {
	conn := *connpointer
	var err error = nil
	instruction := strings.ToUpper(redisCommand.Command)
	switch instruction {
	case "APPEND":
		if redisCommand.ParamLength != 2 {
			err = errors.New(errormessages.IncorrectArgumentsError)
			break
		}
		err = sAppend(conn, server, redisCommand)
	case "DECR":
		if redisCommand.ParamLength != 1 {
			err = errors.New(errormessages.IncorrectArgumentsError)
			break
		}
		err = addsub.Decr(conn, server, redisCommand)
	case "DECRBY":
		if redisCommand.ParamLength != 2 {
			err = errors.New(errormessages.IncorrectArgumentsError)
			break
		}
		err = addsub.Decrby(conn, server, redisCommand)
	case "GET":
		if redisCommand.ParamLength != 1 {
			err = errors.New(errormessages.IncorrectArgumentsError)
			break
		}
		err = doget.Get(conn, server, redisCommand)
	case "GETRANGE":
		if redisCommand.ParamLength != 3 {
			err = errors.New(errormessages.IncorrectArgumentsError)
			break
		}
		err = substr(conn, server, redisCommand)
	case "GETDEL":
		if redisCommand.ParamLength != 1 {
			err = errors.New(errormessages.IncorrectArgumentsError)
			break
		}
		err = doget.GetDel(conn, server, redisCommand)
	case "GETSET":
		if redisCommand.ParamLength != 2 {
			err = errors.New(errormessages.IncorrectArgumentsError)
			break
		}
		err = doget.GetSet(conn, server, redisCommand)
	case "INCR":
		if redisCommand.ParamLength != 1 {
			err = errors.New(errormessages.IncorrectArgumentsError)
			break
		}
		err = addsub.Incr(conn, server, redisCommand)
	case "INCRBY":
		if redisCommand.ParamLength != 2 {
			err = errors.New(errormessages.IncorrectArgumentsError)
			break
		}
		err = addsub.Incrby(conn, server, redisCommand)
	case "MGET":
		if redisCommand.ParamLength < 1 {
			err = errors.New(errormessages.IncorrectArgumentsError)
			break
		}
		err = multi.MGet(conn, server, redisCommand)
	case "MSET":
		if redisCommand.ParamLength < 2 {
			err = errors.New(errormessages.IncorrectArgumentsError)
			break
		}
		err = multi.MSet(conn, server, redisCommand)
	case "SET":
		if redisCommand.ParamLength < 2 {
			err = errors.New(errormessages.IncorrectArgumentsError)
			break
		}
		err = doset.Set(conn, server, redisCommand)
	case "STRLEN":
		if redisCommand.ParamLength != 1 {
			err = errors.New(errormessages.IncorrectArgumentsError)
			break
		}
		err = strlen(conn, server, redisCommand)
	case "SUBSTR":
		if redisCommand.ParamLength != 3 {
			err = errors.New(errormessages.IncorrectArgumentsError)
			break
		}
		err = substr(conn, server, redisCommand)
	default:
		return parserentities.ParserInfo{Executed: false, Err: nil}
	}

	if err != nil {
		return parserentities.ParserInfo{Executed: false, Err: err}
	}

	return parserentities.ParserInfo{Executed: true, Err: nil}
}
