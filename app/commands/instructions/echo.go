package instructions

import (
	"redislite/app/prototools/protomessages"
	"net"
)

func Echo(conn net.Conn, text string) error {
	return protomessages.QuickSendSimpleString(conn, text)
}
