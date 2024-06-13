package stringparser
//
// import (
// 	"fmt"
// 	"net"
// 	"redislite/app/commands/parsers/utils"
// 	"redislite/app/data"
// 	"redislite/app/data/datatypes/kvstring"
// 	"redislite/app/prototools/protomessages"
// 	"redislite/app/setup"
// 	"strconv"
// 	"strings"
// )
//
// func decrby(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand, decr bool, value string) error {
//   return nil
// }
//
// func incrby(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand, decr bool, value string) error {
//   return nil
// }
//
//
// func addsubtract(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand, decr bool, value string) error {
// 	key := redisCommand.Params[0]
// 	// value := redisCommand.Params[1]
// 	intvalue, err := strconv.ParseInt(value, 10, 64)
// 	if strings.HasPrefix(value, "0") || err != nil {
// 		return protomessages.SendError(conn, "value is not an integer or out of range")
// 	}
//
// 	dataObject, exists := server.DataStore.GetKvString(key)
// 	var oldvalue int64 = 0
// 	if !exists || utils.Expired(dataObject.ExpiryTimeNano) {
// 		dataObject = kvstring.KvString{
// 			Value: redisCommand.Params[1],
// 		}
// 	} else {
// 		oldvalue, err = strconv.ParseInt(dataObject.Value, 10, 64)
// 		if err != nil {
// 			return protomessages.SendError(conn, "value is not an integer or out of range")
// 		}
//     
// 		dataObject.Value = fmt.Sprintf("%d", oldvalue+intvalue)
// 	}
//
// 	server.DataStore.SetKvString(key, dataObject)
//
// 	server.RecievedCounter.AddBytes(redisCommand.MessageBytes)
//
// 	return protomessages.QuickSendInt(conn, oldvalue+intvalue)
// }
//
//
// // func createnumber(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
// // 	dataObject := kvstring.KvString{
// // 		Value: redisCommand.Params[1],
// // 	}
// // 	server.DataStore.SetKvString(redisCommand.Params[0], dataObject)
// // 	server.RecievedCounter.AddBytes(redisCommand.MessageBytes)
// // 	if server.Settings.Master {
// // 		return protomessages.QuickSendSimpleString(conn, "OK")
// // 	}
// // 	return nil
// // }
