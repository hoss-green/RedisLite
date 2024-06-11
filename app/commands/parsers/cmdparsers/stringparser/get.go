package stringparser

import (
	"net"
	"time"

	"redislite/app/prototools/protomessages"
	"redislite/app/setup"
)

func get(conn net.Conn, server *setup.Server, key string) error {
	dataObject, ok := server.DataStore.GetKvString(key)

	if !ok || expired(dataObject.ExpiryTimeNano) {
		return protomessages.QuickSendNil(conn)
	}

	return protomessages.QuickSendBulkString(conn, dataObject.Value)
}

func expired(expirytime int64) bool {
	return expirytime != 0 && time.Now().UTC().UnixNano() > expirytime
}
