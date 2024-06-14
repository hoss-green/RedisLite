package streamparser 

// import (
// 	"net"
// 	"strings"
// 	"time"
//
// 	"redislite/app/data"
// 	"redislite/app/data/datatypes/kvstring"
// 	"redislite/app/prototools/protomessages"
// 	"redislite/app/setup"
// )
//
// func xadd(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
// 	dataItems := redisCommand.Params
// 	streamName := dataItems[0]
// 	key := dataItems[1]
//
// 	kvstrings := []kvstring.KvString{}
// 	for index := 2; index < len(dataItems); index += 2 {
// 		value := dataItems[index+1]
// 		kvstrings = append(kvstrings, kvstring.KvString{Key: dataItems[index], Value: value})
// 	}
//
// 	kvsd, err := server.DataStore.AppendKvStream(streamName, key, kvstrings)
// 	if err != nil {
// 		return err
// 	}
//
// 	channelsend := func() {
// 		server.StreamChannel <- setup.KvStreamChannel{
// 			StreamName: strings.ToUpper(streamName), Data: kvsd, AddedTime: time.Now().UTC().UnixNano()}
// 	}
// 	go channelsend()
//
// 	// log.Println("Ok: ", key)
// 	return protomessages.QuickSendSimpleString(conn, kvsd.Id.ToString())
// }
