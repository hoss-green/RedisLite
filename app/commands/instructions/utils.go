package instructions

import "time"

func expired(expirytime int64) bool {
	return expirytime != 0 && time.Now().UTC().UnixNano() > expirytime
}
