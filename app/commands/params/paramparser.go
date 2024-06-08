package params

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"redislite/app/data"
)

const NANOSECOND int = 1000000

func ParseParams(ignorecount int, params []string) (data.ParamOptions, error) {
	paramOptions := data.ParamOptions{}
	slice := params[ignorecount:]
  if len(slice) == 0 {
    return paramOptions, nil
  }
	log.Printf("Params: %v\r\n", slice)
	for i := 0; i < len(slice); i++ {
		log.Printf("Param: %v\r\n", slice[i])
		var err error = nil
		hasnext := i+1 < len(slice)
		//params that require a follow up value
		if hasnext {
			i += 1
			switch param := strings.ToUpper(slice[i-1]); param {
			case "EX":
				paramOptions, err = parseExpirySec(paramOptions, slice[i])
			case "PX":
				paramOptions, err = parseExpiryMs(paramOptions, slice[i])
			}
			if err != nil {
				return data.ParamOptions{}, err
			}
			continue
		}

		switch param := strings.ToUpper(slice[i]); param {
		case "NX":
		case "XX":
		default:
			err = errors.New("Paramater String Malformed")
		}

		if err != nil {
			return data.ParamOptions{}, err
		}
	}

	return paramOptions, nil
}

// EX
func parseExpirySec(redisCommand data.ParamOptions, seconds string) (data.ParamOptions, error) {
	sec, err := strconv.Atoi(seconds)
	if err != nil {
		return data.ParamOptions{}, errors.New("found EX but follow up param could not be parsed into seconds")
	}

	expiryTime := time.Now().UTC().Add(time.Duration(sec * int(time.Second)))
	return parseExpiry(redisCommand, expiryTime), nil
}

// PX
func parseExpiryMs(redisCommand data.ParamOptions, milliseconds string) (data.ParamOptions, error) {
	ms, err := strconv.Atoi(milliseconds)
	if err != nil {
		return data.ParamOptions{}, errors.New("found PX but follow up param could not be parsed into milliseconds")
	}

	expiryTime := time.Now().UTC().Add(time.Duration(ms * int(time.Millisecond)))
	return parseExpiry(redisCommand, expiryTime), nil
}

func parseExpiry(redisCommand data.ParamOptions, expiryTime time.Time) data.ParamOptions {
	log.Println("Set expiry time successfully")
	redisCommand.Expiry = expiryTime.UTC().UnixNano()
	return redisCommand
}
