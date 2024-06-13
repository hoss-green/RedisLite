package timeparser

import (
	"log"
	"time"
)

const NANOSECOND int = 1000000

func ParseExpiryTime(timecode string, duration int64) (int64, bool) {
	switch timecode {
	case "EX":
		// expiryTimeUnix, err =
    return parseExpSecIntoUnix(duration), true
	case "PX":
		// expiryTimeUnix, err = 
    return parseExpMsIntoUnix(duration), true
	}

	return 0, false 
}

// EX
func parseExpSecIntoUnix(seconds int64) int64 {
	expiryTime := time.Now().UTC().Add(time.Duration(int(seconds) * int(time.Second)))
  log.Println("EX")
	return parseIntoUnixTime(expiryTime)
}

// PX
func parseExpMsIntoUnix(milliseconds int64) int64 {
	expiryTime := time.Now().UTC().Add(time.Duration(int(milliseconds) * int(time.Millisecond)))
  log.Println("PX")
	return parseIntoUnixTime(expiryTime)
}

func parseIntoUnixTime(expiryTime time.Time) int64 {
	return expiryTime.UTC().UnixNano()
}

// // EX
// func parseExpirySec(redisCommand data.ParamOptions, seconds string) (data.ParamOptions, error) {
// 	sec, err := strconv.Atoi(seconds)
// 	if err != nil {
// 		return data.ParamOptions{}, errors.New("found EX but follow up param could not be parsed into seconds")
// 	}
//
// 	expiryTime := time.Now().UTC().Add(time.Duration(sec * int(time.Second)))
// 	return parseExpiry(redisCommand, expiryTime), nil
// }
//
// // PX
// func parseExpiryMs(redisCommand data.ParamOptions, milliseconds string) (data.ParamOptions, error) {
// 	ms, err := strconv.Atoi(milliseconds)
// 	if err != nil {
// 		return data.ParamOptions{}, errors.New("found PX but follow up param could not be parsed into milliseconds")
// 	}
//
// 	expiryTime := time.Now().UTC().Add(time.Duration(ms * int(time.Millisecond)))
// 	return parseExpiry(redisCommand, expiryTime), nil
// }
//
// func parseExpiry(redisCommand data.ParamOptions, expiryTime time.Time) data.ParamOptions {
// 	log.Println("Set expiry time successfully")
// 	redisCommand.Expiry = expiryTime.UTC().UnixNano()
// 	return redisCommand
// }
