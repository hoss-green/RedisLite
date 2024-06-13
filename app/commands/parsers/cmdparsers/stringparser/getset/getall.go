package getset

import (
	"errors"
	"net"
	"redislite/app/commands/parsers/utils"
	"redislite/app/data"
	"redislite/app/data/storage/datatyperrors"
	"redislite/app/prototools/protomessages"
	"redislite/app/setup"
)

func doGet(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand, getParamList getParamList) error {
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
	if getParamList.delete {
		server.DataStore.DelKvString(key)
	}

  if getParamList.persist {
    dataObject.ExpiryTimeNano = 0
    server.DataStore.SetKvString(key, dataObject)
  }

  if getParamList.hasExpiry {
    dataObject.ExpiryTimeNano = getParamList.expiry
    server.DataStore.SetKvString(key, dataObject)
  }

	return protomessages.QuickSendBulkString(conn, value)
}
