package stringparser

import (
	"net"
	"redislite/app/commands/parsers/utils"
	"redislite/app/data"
	"redislite/app/prototools/protomessages"
	"redislite/app/setup"
)

func getdel(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
	key := redisCommand.Params[0]
	dataObject, ok := server.DataStore.GetKvString(key)
	if !ok || utils.Expired(dataObject.ExpiryTimeNano) {
		return protomessages.QuickSendNil(conn)
	}

	value := dataObject.Value
	server.DataStore.DelKvString(key)
	return protomessages.QuickSendBulkString(conn, value)
}
