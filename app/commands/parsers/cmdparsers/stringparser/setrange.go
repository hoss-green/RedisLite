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

func setrange(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
	params := redisCommand.Params
	key := params[0]
	replacementString := []rune(params[2])
	replacementStringLen := len(params[2])

	from, fromerr := strconv.ParseInt(params[1], 10, 64)
	if fromerr != nil {
		return errors.New("value is not an integer or out of range")
	}

	if from < 0 {
		return protomessages.QuickSendError(conn, "offset is out of range")
	}
	if from+int64(replacementStringLen) > 536870911 {
		return protomessages.QuickSendError(conn, "offset too big")
	}

	dataObject, err := server.DataStore.GetKvString(key)
	if err != nil || utils.Expired(dataObject.ExpiryTimeNano) {
		var tiErr *datatyperrors.WrongtypeError
		if errors.As(err, &tiErr) {
			return protomessages.QuickSendError(conn, tiErr.Error())
		}
	}

	currentVal := []rune(dataObject.Value)
	currentValLen := int64(len(currentVal))

	newArrayLength := from + int64(replacementStringLen)
	//grow array
	if newArrayLength >= currentValLen {
		//create padding object first
		// paddedstring := []rune("")
		if from > currentValLen {
			// padding := from - currentValLen
			// for range padding {
			// 	paddedstring = append(paddedstring, rune(0x00))
			// }
			//
			for index := currentValLen; index < newArrayLength; index++ {
				currentVal = append(currentVal, rune(0x00))
				currentValLen = int64(len(currentVal))
			}
		} else {
			emptyarray := make([]rune, from+int64(replacementStringLen)-currentValLen)
			currentVal = append(currentVal, emptyarray...)
			currentValLen = int64(len(currentVal))

		}
	}

	count := 0
	for index := from; index < from+int64(replacementStringLen); index++ {
		currentVal[index] = replacementString[count]
		count += 1
	}

	dataObject.Value = string(currentVal)
	server.DataStore.SetKvString(dataObject.Key, dataObject)

	return protomessages.QuickSendInt(conn, currentValLen)
}
