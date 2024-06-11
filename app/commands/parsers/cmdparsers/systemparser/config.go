package systemparser 

import (
	"net"
	"redislite/app/data"
	"redislite/app/prototools/protomessages"
	"redislite/app/setup"
	"strings"
)

func config(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
	switch strings.ToUpper(redisCommand.Params[0]) {
	case "GET":
		switch strings.ToUpper(redisCommand.Params[1]) {
		case "DIR":
			protomessages.SendRespArray(conn, []string{"dir", server.Settings.Dir})
		case "DBFILENAME":
			protomessages.SendRespArray(conn, []string{"dbfilename", server.Settings.DbFilename})
		}
	}

	return nil
}
