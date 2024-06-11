package systemparser 

import (
	"errors"
	"net"
	"strings"

	"redislite/app/commands/errormessages"
	"redislite/app/commands/parsers/parserentities"
	"redislite/app/data"
	"redislite/app/setup"
)

func ParseSystemCommand(connpointer *net.Conn, redisCommand data.RedisCommand, server *setup.Server) parserentities.ParserInfo {
	conn := *connpointer
	var err error = nil
	instruction := strings.ToUpper(redisCommand.Command)
	switch instruction {
	case "PING":
		err = ping(conn, server, redisCommand)
	case "INFO":
		if redisCommand.ParamLength > 1 {
			err = errors.New(errormessages.IncorrectArgumentsError)
			break
		}
		if redisCommand.ParamLength == 1 {
			err = info(conn, server.Settings, redisCommand.Params[0])
		} else {
			err = info(conn, server.Settings, "server")
		}
	case "ECHO":
		if redisCommand.ParamLength != 1 {
			err = errors.New(errormessages.IncorrectArgumentsError)
			break
		}
		err = echo(conn, redisCommand.Params[0])
	case "CONFIG":
		err = config(conn, server, redisCommand)
	case "KEYS":
		err = keys(conn, server, redisCommand)
	case "TYPE":
		if redisCommand.ParamLength != 1 {
			err = errors.New(errormessages.IncorrectArgumentsError)
			break
		}
		err = Type(conn, server, redisCommand.Params[0])
	case "WAIT":
		err = wait(conn, server, redisCommand)
	case "COMMAND":
		command(conn)
	// case "TEST":
	// 	if server.Settings.Master {
	// 		replication.SendToReplicas(server.Replicas, protomessages.BuildRespArrayMsg([]string{"REPLCONF", "GETACK", "*"}))
	// 		return parserentities.ParserInfo{Executed: true, Err: nil}
	// 	}
	case "REPLCONF":
		if redisCommand.ParamLength < 1 {
			err = errors.New(errormessages.IncorrectArgumentsError)
			break
		}
		err = replconf(conn, server, redisCommand)
	case "PSYNC":
		replicaConn, err := psync(conn, server.Settings)
		if err == nil {
			server.Replicas = append(server.Replicas, replicaConn)
		}
	default:
		return parserentities.ParserInfo{Executed: false, Err: nil}
	}

	if err != nil {
		return parserentities.ParserInfo{Executed: false, Err: err}
	}

	return parserentities.ParserInfo{Executed: true, Err: nil}
}
