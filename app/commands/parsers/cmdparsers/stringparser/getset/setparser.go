package getset 

import (
	"errors"
	"log"
	"redislite/app/commands/parsers/cmdparsers/timeparser"
	"strconv"
	"strings"
)

// SET key value [NX | XX] [GET] [EX seconds | PX milliseconds | EXAT unix-time-seconds | PXAT unix-time-milliseconds | KEEPTTL]
type paramList struct {
	get            bool
	setInstruction setType
	hasExpiry      bool
	expiry         int64
}

// SET key value [NX | XX] [GET] [EX seconds | PX milliseconds | EXAT unix-time-seconds | PXAT unix-time-milliseconds | KEEPTTL]
func parseParams(params []string) (paramList, error) {
  log.Println("Params: ", params)
  log.Println("ParamsLen: ", len(params))
	pList := paramList{}
	for index := 0; index < len(params); index++ {
    command := strings.ToUpper(params[index])
		switch command {
		case "GET":
			pList.get = true
		case "NX":
			pList.setInstruction = NXSetIfKeyNotExist
		case "XX":
			pList.setInstruction = XXOnlySetIfExists
		default:
			expiryParseResult, err := tryParseExpiry(command, index, params)
			if err != nil {
				return pList, err
			}

			if expiryParseResult.hasResult {
				pList.hasExpiry = true
				pList.expiry = expiryParseResult.value
				index += 1
			}
		}
	}

	return pList, nil
}

type expiryParseResult struct {
	value     int64
	hasResult bool
}

func tryParseExpiry(command string, index int, params []string) (expiryParseResult, error) {
	if index > len(params) {
		return expiryParseResult{value: 0, hasResult: false}, errors.New("syntax error")
	}

	duration, err := strconv.ParseInt(params[index+1], 10, 64)
	log.Println("DURATION", duration)
	if err != nil {
		return expiryParseResult{value: 0, hasResult: false}, errors.New("value is not an integer or out of range")
	}

	if duration < 0 {
		return expiryParseResult{value: 0, hasResult: false}, errors.New("invalid expire time in 'set' command")
	}

	timeValue, timeHasValue := timeparser.ParseExpiryTime(command, duration)
	return expiryParseResult{value: timeValue, hasResult: timeHasValue}, nil
}
