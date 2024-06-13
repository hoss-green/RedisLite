package stringparser

import (
	"errors"
	"net"
	"redislite/app/commands/parsers/utils"
	"redislite/app/data"
	"redislite/app/data/datatypes/kvstring"
	"redislite/app/data/storage/datatyperrors"
	"redislite/app/prototools/protomessages"
	"redislite/app/setup"
)

func getset(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
	key := redisCommand.Params[0]
	dataObject, err := server.DataStore.GetKvString(key)
	response := ""
	if err != nil || utils.Expired(dataObject.ExpiryTimeNano) {
		var tiErr *datatyperrors.WrongtypeError
		if errors.As(err, &tiErr) {
			return protomessages.QuickSendError(conn, tiErr.Error())
		}
		response = protomessages.BuildNilMsg()
		dataObject = kvstring.KvString{
			Value: redisCommand.Params[1],
		}
	} else {
		response = protomessages.BuildBulkStringMsg(dataObject.Value)
		dataObject.Value = redisCommand.Params[1]
	}

	server.DataStore.SetKvString(key, dataObject)
	server.RecievedCounter.AddBytes(redisCommand.MessageBytes)

	return protomessages.QuickSendMessage(conn, response)
	// if server.Settings.Master { return protomessages.QuickSendMessage(conn, response)
	// }
	// return nil

}
