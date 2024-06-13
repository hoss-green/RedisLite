package stringparser

import (
	"net"
	"redislite/app/commands/parsers/utils"
	"redislite/app/data"
	"redislite/app/data/datatypes/kvstring"
	"redislite/app/prototools/protomessages"
	"redislite/app/setup"
)

func getset(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
	key := redisCommand.Params[0]
	dataObject, ok := server.DataStore.GetKvString(key)
	response := ""
	if !ok || utils.Expired(dataObject.ExpiryTimeNano) {
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
