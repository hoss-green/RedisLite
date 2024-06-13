package stringparser

import (
	"fmt"
	"net"

	"redislite/app/commands/parsers/utils"
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
	if err == nil && !utils.Expired(dataObject.ExpiryTimeNano) {
		currenttext = dataObject.Value
		currentexpiry = dataObject.ExpiryTimeNano
	} else {

	}

	dataObject.Value = fmt.Sprintf("%s%s", currenttext, appendvalue)
	dataObject.ExpiryTimeNano = currentexpiry

	server.DataStore.SetKvString(key, dataObject)
	server.RecievedCounter.AddBytes(redisCommand.MessageBytes)
	return protomessages.QuickSendInt(conn, int64(len(dataObject.Value)))
}
