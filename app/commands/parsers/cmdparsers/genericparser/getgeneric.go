package genericparser

import (
	"net"
	"redislite/app/data"
	"redislite/app/prototools/protomessages"
	"redislite/app/setup"
)

func doDel(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
  result := server.DataStore.DeleteKeys(redisCommand.Params)
	return protomessages.QuickSendInt(conn, result)
}

func doCopy(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
  key := redisCommand.Params[0]
  targetKey := redisCommand.Params[1]
  result := server.DataStore.CopyKey(key, targetKey)
	return protomessages.QuickSendInt(conn, result)
}

func doExists(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
  result := server.DataStore.ExistKeys(redisCommand.Params)
	return protomessages.QuickSendInt(conn, result)
}
