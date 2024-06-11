package systemparser 

import (
	"redislite/app/prototools/protomessages"
	"net"
)

func echo(conn net.Conn, text string) error {
	return protomessages.QuickSendSimpleString(conn, text)
}
