package stringparser

import (
	"errors"
	"net"

	"redislite/app/data"
	"redislite/app/data/storage/datatyperrors"
	"redislite/app/prototools/protomessages"
	"redislite/app/setup"
)

func strlen(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
  key := redisCommand.Params[0]
	dataObject, err := server.DataStore.GetKvString(key)

	if err != nil  {
		var tiErr *datatyperrors.WrongtypeError
		if errors.As(err, &tiErr) {
			return protomessages.QuickSendError(conn, tiErr.Error())
		}
		return protomessages.QuickSendInt(conn, 0)
	}

		return protomessages.QuickSendInt(conn, int64(len(string(dataObject.Value))))
}

