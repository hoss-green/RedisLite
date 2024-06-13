package doget

import (
	"errors"
	"net"
	"redislite/app/commands/parsers/utils"
	"redislite/app/data"
	"redislite/app/data/storage/datatyperrors"
	"redislite/app/prototools/protomessages"
	"redislite/app/setup"
)

func doGet(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand, doDelete bool) error {
	key := redisCommand.Params[0]
	dataObject, err := server.DataStore.GetKvString(key)

	if err != nil || utils.Expired(dataObject.ExpiryTimeNano) {
		var tiErr *datatyperrors.WrongtypeError
		if errors.As(err, &tiErr) {
			return protomessages.QuickSendError(conn, tiErr.Error())
		}
		return protomessages.QuickSendNil(conn)
	}
	value := dataObject.Value
	if doDelete {
		server.DataStore.DelKvString(key)
	}
	return protomessages.QuickSendBulkString(conn, value)
}
