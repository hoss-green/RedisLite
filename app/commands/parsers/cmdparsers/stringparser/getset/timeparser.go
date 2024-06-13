package getset 

import (
	"errors"
	"log"
	"redislite/app/commands/parsers/cmdparsers/timeparser"
	"strconv"
)

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
