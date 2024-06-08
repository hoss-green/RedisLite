package instructions

import (
	"net"
	"redislite/app/prototools/protomessages"
	"redislite/app/setup"
)

func Type(conn net.Conn, server *setup.Server, key string) error {
	typename := server.DataStore.KvType(key)

	return protomessages.QuickSendSimpleString(conn, typename)
}

