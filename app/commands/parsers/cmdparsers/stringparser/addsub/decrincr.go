package addsub

import (
	"errors"
	"fmt"
	"net"
	"redislite/app/commands/parsers/utils"
	"redislite/app/data"
	"redislite/app/data/datatypes/kvstring"
	"redislite/app/data/storage/datatyperrors"
	"redislite/app/prototools/protomessages"
	"redislite/app/setup"
	"strconv"
	"strings"
)

func Decr(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
	return addsubtract(conn, server, redisCommand, true, "1")
}

func Decrby(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
	return addsubtract(conn, server, redisCommand, true, redisCommand.Params[1])
}

func Incr(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
	return addsubtract(conn, server, redisCommand, false, "1")
}
func Incrby(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
	return addsubtract(conn, server, redisCommand, false, redisCommand.Params[1])
}

func addsubtract(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand, decr bool, value string) error {
	key := redisCommand.Params[0]
	newvalue, err := strconv.ParseInt(value, 10, 64)
	if strings.HasPrefix(value, "0") || err != nil {
		return protomessages.QuickSendError(conn, "value is not an integer or out of range")
	}

	if decr {
		newvalue *= -1
	}

	dataObject, err := server.DataStore.GetKvString(key)
	var oldvalue int64 = 0
	if err != nil || utils.Expired(dataObject.ExpiryTimeNano) {
    var tiErr *datatyperrors.WrongtypeError
		if errors.As(err, &tiErr) {
      return protomessages.QuickSendError(conn, err.Error())
		}
		dataObject = kvstring.KvString{
			Value: value,
		}
	} else {
		oldvalue, err = strconv.ParseInt(dataObject.Value, 10, 64)
		if err != nil || (len(dataObject.Value) > 1 && strings.HasPrefix(dataObject.Value, "0")) {
			return protomessages.QuickSendError(conn, "value is not an integer or out of range")
		}

		dataObject.Value = fmt.Sprintf("%d", oldvalue+newvalue)
	}

	server.DataStore.SetKvString(key, dataObject)
	server.RecievedCounter.AddBytes(redisCommand.MessageBytes)

	return protomessages.QuickSendInt(conn, oldvalue+newvalue)
}
