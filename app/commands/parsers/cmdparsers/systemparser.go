package cmdparsers

import (
	"errors"
	"net"
	"strings"

	"redislite/app/commands/errormessages"
	"redislite/app/commands/instructions"
	"redislite/app/data"
	"redislite/app/prototools/protomessages"
	"redislite/app/replication"
	"redislite/app/setup"
)

func ParseSystemCommand(connpointer *net.Conn, redisCommand data.RedisCommand, server *setup.Server) ParserInfo {
	conn := *connpointer
	var err error = nil
	instruction := strings.ToUpper(redisCommand.Command)
	switch instruction {
	case "PING":
		err = instructions.Ping(conn, server, redisCommand)
	case "INFO":
		if redisCommand.ParamLength > 1 {
			err = errors.New(errormessages.IncorrectArgumentsError)
			break
		}
		if redisCommand.ParamLength == 1 {
			err = instructions.Info(conn, server.Settings, redisCommand.Params[0])
		} else {
			err = instructions.Info(conn, server.Settings, "server")
		}
	case "ECHO":
		if redisCommand.ParamLength != 1 {
			err = errors.New(errormessages.IncorrectArgumentsError)
			break
		}
		err = instructions.Echo(conn, redisCommand.Params[0])
	case "CONFIG":
		err = instructions.Config(conn, server, redisCommand)
	case "KEYS":
		err = instructions.Keys(conn, server, redisCommand)
	case "TYPE":
		if redisCommand.ParamLength != 1 {
			err = errors.New(errormessages.IncorrectArgumentsError)
			break
		}
		err = instructions.Type(conn, server, redisCommand.Params[0])
	case "WAIT":
		err = instructions.Wait(conn, server, redisCommand)
	case "COMMAND":
		instructions.Command(conn)
	case "TEST":
		if server.Settings.Master {
			replication.SendToReplicas(server.Replicas, protomessages.BuildRespArrayMsg([]string{"REPLCONF", "GETACK", "*"}))
			return ParserInfo{Executed: true, Err: nil}
		}
		// instructions.Test(conn)
	case "REPLCONF":
		if redisCommand.ParamLength < 1 {
			err = errors.New(errormessages.IncorrectArgumentsError)
			break
		}
		err = instructions.ReplConf(conn, server, redisCommand)
	case "PSYNC":
		replicaConn, err := instructions.Psync(conn, server.Settings)
		if err == nil {
			server.Replicas = append(server.Replicas, replicaConn)
		}
	default:
		return ParserInfo{Executed: false, Err: nil}
	}

	if err != nil {
		return ParserInfo{Executed: false, Err: err}
	}

	return ParserInfo{Executed: true, Err: nil}
}
