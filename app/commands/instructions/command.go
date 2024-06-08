package instructions

import (
	"fmt"
	"net"
)

func Command(conn net.Conn) {
	fmt.Printf("Client connected from ip: %s.\n", conn.RemoteAddr().String())
	conn.Write([]byte("+\r\n"))
}
