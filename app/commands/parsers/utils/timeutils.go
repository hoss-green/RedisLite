package utils 

import "time"

func Expired(expirytime int64) bool {
	return expirytime != 0 && time.Now().UTC().UnixNano() > expirytime
}
