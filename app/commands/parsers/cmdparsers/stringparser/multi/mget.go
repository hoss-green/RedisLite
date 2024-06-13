package multi 

import (
	"net"
	"redislite/app/commands/parsers/utils"
	"redislite/app/data"
	"redislite/app/prototools/protomessages"
	"redislite/app/setup"
)

func MGet(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
	responses := []string{}

	for _, key := range redisCommand.Params {
		dataObject, err := server.DataStore.GetKvString(key)
		if err != nil || utils.Expired(dataObject.ExpiryTimeNano) {
			responses = append(responses, protomessages.BuildNilMsg())
		} else {
			responses = append(responses, protomessages.BuildBulkStringMsg(dataObject.Value))
		}
	}

	return protomessages.SendRespMultiArray(conn, responses)
}
