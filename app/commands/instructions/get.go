package instructions

import (
	"net"

	"redislite/app/prototools/protomessages"
	"redislite/app/setup"
)

func Get(conn net.Conn, server *setup.Server, key string) error {
	dataObject, ok := server.DataStore.GetKvString(key)

	if !ok || expired(dataObject.ExpiryTimeNano) {
		return protomessages.QuickSendNil(conn)
	}

	return protomessages.QuickSendBulkString(conn, dataObject.Value)
}
