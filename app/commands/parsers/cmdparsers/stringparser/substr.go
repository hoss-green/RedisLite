package stringparser

import (
	"errors"
	"net"
	"redislite/app/commands/parsers/utils"
	"redislite/app/data"
	"redislite/app/data/storage/datatyperrors"
	"redislite/app/prototools/protomessages"
	"redislite/app/setup"
	"strconv"
)

func substr(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
	params := redisCommand.Params
	key := params[0]

	from, fromerr := strconv.ParseInt(params[1], 10, 64)
	to, toerr := strconv.ParseInt(params[2], 10, 64)
	if fromerr != nil || toerr != nil {
		return errors.New("value is not an integer or out of range")
	}
	//includes 0 0, no equality check is needed
	if from < 0 && to > 0 {
		return protomessages.QuickSendEmptyString(conn)
	}

	dataObject, err := server.DataStore.GetKvString(key)
	if err != nil || utils.Expired(dataObject.ExpiryTimeNano) {
    var tiErr *datatyperrors.WrongtypeError
		if errors.As(err, &tiErr) {
      return protomessages.QuickSendError(conn, tiErr.Error())
		}

		return protomessages.QuickSendEmptyString(conn)
	}

	// absto := abs(to)
	// absfrom := abs(from)

	maxlen := int64(len(dataObject.Value))

	if from <= maxlen && from <= to && to > 0 {
		if from < 0 {
			from = 0
		}
		return protomessages.QuickSendBulkString(conn, dataObject.Value[from:limitmax(maxlen, to + 1)])
	}
	if from < 0 && to < 0 && from <= to {
		if from < -maxlen {
			from = -maxlen
		}
		return protomessages.QuickSendBulkString(conn, dataObject.Value[maxlen+from:maxlen+to+1])
	}

	return protomessages.QuickSendEmptyString(conn)
}
func limitmax(limit int64, num int64) int64 {
	if num < limit {
		return num
	}

	return limit
}

func abs(num int64) int64 {
	if num < 0 {
		return -num
	}
	return num
}
