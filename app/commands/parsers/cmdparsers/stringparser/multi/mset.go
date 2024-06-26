package multi

import (
	"net"

	"redislite/app/data"
	"redislite/app/data/storage"
	"redislite/app/prototools/protomessages"
	"redislite/app/setup"
)

func MSet(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
	dataItems := redisCommand.Params
	if len(dataItems)%2 != 0 {
		return protomessages.QuickSendError(conn, "ERR wrong number of arguments for command")
	}

	//this is horrible, must fix
	keys := []string{}
	dataObjects := []storage.DataItem{}// []kvstring.KvString{}
	for index := 0; index < len(dataItems)-1; index++ {
		dataObject := storage.DataItem{}//kvstring.KvString{}
		key := dataItems[index]
		dataObject.Value = []byte(dataItems[index+1])
		keys = append(keys, key)
		dataObjects = append(dataObjects, dataObject)
	}
	server.DataStore.SetKvStrings(keys, dataObjects)
	server.RecievedCounter.AddBytes(redisCommand.MessageBytes)

	return protomessages.QuickSendSimpleString(conn, "OK")
}
