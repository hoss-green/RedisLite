package instructions

import (
	"log"
	"net"

	"redislite/app/data"
	"redislite/app/prototools/protomessages"
	"redislite/app/setup"
)

func XRange(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
	//                 streamkey starttime    endtime
	//redis-cli XRANGE some_key 1526985054069 1526985054079
	dataItems := redisCommand.Params
	streamName := dataItems[0]
	start := dataItems[1]
	end := dataItems[2]

	streamrange, err := server.DataStore.GetStreamRange(streamName, start, end)
	if err != nil {
		return err
	}

	arrays := &[]string{}
	for _, streamItem := range streamrange {
		title := protomessages.BuildBulkStringMsg(streamItem.Id.ToString())
		kvparray := []string{}
		for _, kvp := range streamItem.Items {
			kvparray = append(kvparray, kvp.Key)
			kvparray = append(kvparray, kvp.Value)
		}

		kvpresp := protomessages.BuildRespArrayMsg(kvparray)
		kvcontainer := protomessages.BuildMultiRespArrayMsg([]string{title, kvpresp})
		log.Printf("ID: %#v", title)

		*arrays = append(*arrays, kvcontainer)

	}
	log.Printf("KVP: %#v", *arrays)
	// log.Println(streamrange)
	response := protomessages.BuildMultiRespArrayMsg(*arrays) + "\r\n"
	log.Printf("Final: %#v", response)
	return protomessages.QuickSendMessage(conn, response)
}
