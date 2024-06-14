package getset

import (
	"errors"
	// "log"
	"net"
	"redislite/app/commands/parsers/utils"
	"redislite/app/data"
	"redislite/app/data/storage"
	"redislite/app/data/storage/datatyperrors"
	"redislite/app/prototools/protomessages"
	"redislite/app/setup"
)

func doSet(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand, setParams paramList) error {
	dataItems := redisCommand.Params
	dataObject := storage.DataItem{}
	key := dataItems[0]
	dataObject.Value = []byte( dataItems[1])
	keyExists := false
	oldObject := storage.DataItem{}
	var err error
	oldObject, err = server.DataStore.GetKvString(key)
	var tiErr *datatyperrors.WrongtypeError
	if err != nil || utils.Expired(dataObject.ExpiryTimeNano) {
		typeError := errors.As(err, &tiErr)
		keyExists = !typeError
	} else {
		keyExists = true
	}

	if keyExists && setParams.setInstruction == NXSetIfKeyNotExist {
		return protomessages.QuickSendNil(conn)
	}
	if !keyExists && setParams.setInstruction == XXOnlySetIfExists {
		return protomessages.QuickSendNil(conn)
	}

	if setParams.hasExpiry {
		dataObject.ExpiryTimeNano = setParams.expiry
	} else if setParams.keepttl {
    dataObject.ExpiryTimeNano = oldObject.ExpiryTimeNano
  }
	//set all the shit here



	server.DataStore.SetKvString(dataItems[0], dataObject)
	server.RecievedCounter.AddBytes(redisCommand.MessageBytes)
	if setParams.get {
		return protomessages.QuickSendBulkString(conn, string(oldObject.Value))
	}
	return protomessages.QuickSendSimpleString(conn, "OK")
	// if server.Settings.Master {
	// 	return protomessages.QuickSendSimpleString(conn, "OK")
	// }
	// return nil
}
