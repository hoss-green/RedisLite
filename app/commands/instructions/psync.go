package instructions

import (
	"errors"
	"fmt"
	"net"

	"redislite/app/persistence"
	"redislite/app/setup"
)

func Psync(conn net.Conn, serverSettings setup.ServerSettings) (net.Conn, error) {
	println("Configuring Psync")
	response := fmt.Sprintf("+FULLRESYNC %s %d\r\n", serverSettings.MasterReplId, serverSettings.MasterReplIdOffset)
	conn.Write([]byte(response))
	//decode db
	SendDbToReplica(conn)
	return conn, nil
}

func SendDbToReplica(conn net.Conn) error {
	rdb := persistence.Build()
	_, err := conn.Write(append([]byte(fmt.Sprintf("$%d\r\n", len(rdb))), rdb...))
	if err != nil {
		fmt.Println("Failed to Send Db to Replica")
		return errors.New("Failed to Send Db to Replica")
	}

	return nil
}
