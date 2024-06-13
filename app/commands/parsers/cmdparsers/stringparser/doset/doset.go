package doset

import (
	"errors"
	// "log"
	"net"
	"redislite/app/commands/parsers/utils"
	"redislite/app/data"
	"redislite/app/data/datatypes/kvstring"
	"redislite/app/data/storage/datatyperrors"
	"redislite/app/prototools/protomessages"
	"redislite/app/setup"
)

func doSet(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand, setParams paramList) error {
	dataItems := redisCommand.Params
	dataObject := kvstring.KvString{}
	key := dataItems[0]
	dataObject.Value = dataItems[1]
	keyExists := false
	oldObject := kvstring.KvString{}
	if setParams.get || setParams.setInstruction != NormalSetType {
		var err error
		oldObject, err = server.DataStore.GetKvString(key)
		var tiErr *datatyperrors.WrongtypeError
		if err != nil || utils.Expired(dataObject.ExpiryTimeNano) {
			typeError := errors.As(err, &tiErr)
			keyExists = !typeError
		} else {
			keyExists = true
		}

	}

	if keyExists && setParams.setInstruction == NXSetIfKeyNotExist {
		return protomessages.QuickSendNil(conn)
	}
	if !keyExists && setParams.setInstruction == XXOnlySetIfExists {
		return protomessages.QuickSendNil(conn)
	}

	if setParams.hasExpiry {
		dataObject.ExpiryTimeNano = setParams.expiry
	}
	//set all the shit here

	server.DataStore.SetKvString(dataItems[0], dataObject)
	server.RecievedCounter.AddBytes(redisCommand.MessageBytes)
	if setParams.get {
		return protomessages.QuickSendBulkString(conn, oldObject.Value)
	}
	return protomessages.QuickSendSimpleString(conn, "OK")
	// if server.Settings.Master {
	// 	return protomessages.QuickSendSimpleString(conn, "OK")
	// }
	// return nil
}
