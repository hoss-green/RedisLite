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
)

func IncrbyFloat(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
	return addsubtractfloat(conn, server, redisCommand, false, redisCommand.Params[1])
}

func addsubtractfloat(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand, decr bool, value string) error {
	key := redisCommand.Params[0]
	var newvalue float64
	newvalue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return protomessages.QuickSendError(conn, "value is not a valid float")
	}

	if decr {
		newvalue *= -1
	}

	dataObject, err := server.DataStore.GetKvString(key)
	var oldvalue float64 = 0
	if err != nil || utils.Expired(dataObject.ExpiryTimeNano) {
		var tiErr *datatyperrors.WrongtypeError
		if errors.As(err, &tiErr) {
			return protomessages.QuickSendError(conn, err.Error())
		}
		dataObject = kvstring.KvString{
			Value: value,
		}
	} else {
		oldvalue, err = strconv.ParseFloat(dataObject.Value, 64)
		if err != nil {
			return protomessages.QuickSendError(conn, "value is not a valid float")
		}

		dataObject.Value = fmt.Sprintf("%f", oldvalue+newvalue)
	}

	server.DataStore.SetKvString(key, dataObject)
	server.RecievedCounter.AddBytes(redisCommand.MessageBytes)

	return protomessages.QuickSendBulkString(conn, fmt.Sprintf("%f", oldvalue+newvalue))
	// return protomessages.QuickSendInt(conn, oldvalue+newvalue)
}
