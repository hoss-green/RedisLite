package protomessages 

import (
	"net"
)

func QuickSendMessage(conn net.Conn, totalMessage string) error {
	_, error := conn.Write([]byte(totalMessage))
	return error
}

func QuickSendSimpleString(conn net.Conn, msg string) error {
	_, error := conn.Write([]byte(BuildSimpleStringMsg(msg)))
	return error
}
func QuickSendInt(conn net.Conn, number int) error {
	_, error := conn.Write([]byte(BuildIntMsg(number)))
	return error
}

func QuickSendNil(conn net.Conn) error {
	_, error := conn.Write([]byte(BuildNilMsg()))
	return error
}

func QuickSendBulkString(conn net.Conn, msg string) error {
	if msg == "" {
		return QuickSendNil(conn)
	}

	_, error := conn.Write([]byte(BuildBulkStringMsg(msg)))
	return error
}

func SendError(conn net.Conn, errmsg string) error {
	_, error := conn.Write([]byte(BuildErrorMsg(errmsg)))
	return error
}

func SendRespArray(conn net.Conn, items []string) error {
	_, error := conn.Write([]byte(BuildRespArrayMsg(items)))
	return error
}

func SendRespMultiArray(conn net.Conn, items []string) error {
	_, error := conn.Write([]byte(BuildMultiRespArrayMsg(items)))
	return error
}
