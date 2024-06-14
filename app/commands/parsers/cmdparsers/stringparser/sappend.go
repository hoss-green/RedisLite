package stringparser

import (
	"fmt"
	"net"
	"redislite/app/data"
	"redislite/app/prototools/protomessages"
	"redislite/app/setup"
)

func sAppend(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
	dataItems := redisCommand.Params
	key := dataItems[0]
	appendvalue := dataItems[1]
	dataObject, err := server.DataStore.GetKvString(key)

	currenttext := ""
	var currentexpiry int64 = 0
	if err != nil {
	} else {
		currenttext = string(dataObject.Value)
		currentexpiry = dataObject.ExpiryTimeNano

	}

	dataObject.Value = []byte(fmt.Sprintf("%s%s", currenttext, appendvalue))
	dataObject.ExpiryTimeNano = currentexpiry

	server.DataStore.SetKvString(key, dataObject)
	server.RecievedCounter.AddBytes(redisCommand.MessageBytes)
	return protomessages.QuickSendInt(conn, int64(len(dataObject.Value)))
}
