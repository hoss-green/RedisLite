package stringparser

import (
	"net"

	"redislite/app/commands/params"
	"redislite/app/data"
	"redislite/app/data/datatypes/kvstring"
	"redislite/app/prototools/protomessages"
	"redislite/app/setup"
)

func set(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
  dataItems := redisCommand.Params
	dataObject := kvstring.KvString{}
	dataObject.Value = dataItems[1]
	err := paramDecoder(dataItems, &dataObject)
	if err != nil {
		return err
	}
	server.DataStore.SetKvString(dataItems[0], dataObject)
  server.RecievedCounter.AddBytes(redisCommand.MessageBytes)
	if server.Settings.Master {
		return protomessages.QuickSendSimpleString(conn, "OK")
	}
	return nil 
}

func paramDecoder(dataItems []string, dataObject *kvstring.KvString) error {
	paramOptions, err := params.ParseParams(2, dataItems)
	if err != nil {
		return err
	}
	if paramOptions.Expiry != 0 {
		dataObject.ExpiryTimeNano = paramOptions.Expiry
	}

	return nil
}
