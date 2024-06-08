package replication

import (
	"fmt"
	"log"
	"net"
)

func SendToReplicas(replicas []net.Conn, command string) {
	if replicas == nil || len(replicas) == 0 {
		return
	}
	log.Printf("Command: %s\n\r", command)
	log.Printf("Sending to %d replicas\n\r", len(replicas))
	for index, replica := range replicas {
		log.Printf("Sending to replica: #%d\n\r", index+1)
		info := fmt.Sprintf("%s", command)
		// replica.Write([]byte(command))
		replica.Write([]byte(info))
	}
}
