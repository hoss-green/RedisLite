package instructions

import (
	"net"

	"redislite/app/data"
	"redislite/app/prototools/protomessages"
	"redislite/app/setup"
)

func Ping(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
	server.RecievedCounter.AddBytes(redisCommand.MessageBytes)
	if !server.Settings.Master {
		return nil
	}
	return protomessages.QuickSendSimpleString(conn, "PONG")
}
