package instructions

import (
	"errors"
	"redislite/app/data"
	"redislite/app/prototools/protomessages"
	"redislite/app/setup"
	"log"
	"net"
	"strconv"
	"time"
)

func Wait(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
	log.Println("WAIT start")
	replicacount, err := strconv.Atoi(redisCommand.Params[0])
	if err != nil {
		return errors.New("Replica Count incorrect for Wait Command")
	}
	timerdurationMS, err := strconv.Atoi(redisCommand.Params[1])
	if err != nil {
		return errors.New("Timer parameter incorrect for Wait Command")
	}
	log.Printf("Waiting for %d replicas for %d milliseconds", replicacount, timerdurationMS)
	server.RecievedCounter.AddBytes(redisCommand.MessageBytes)
	if len(server.Replicas) == 0 ||
		server.DataStore.CountKvString() == 0 {
		return protomessages.QuickSendInt(conn, len(server.Replicas))
	}

	var acks int
	for index, replica := range server.Replicas {
		checkAck(replica, index)
	}

	timer := time.After(time.Duration(timerdurationMS) * time.Millisecond)
loop:
	for acks < replicacount {
		select {
		case <-server.AckChannel:
			acks += 1
			log.Println("Recieved ack")
		case <-timer:
			log.Println("Timer Expired")
			break loop
		}
	}
	return protomessages.QuickSendInt(conn, acks)
}

func checkAck(replica net.Conn, index int) {
	getack := protomessages.BuildRespArrayMsg([]string{"REPLCONF", "GETACK", "*"})
	log.Printf("Sending check GETACK to replica: #%d\n\r", index+1)
	replica.Write([]byte(getack))
}
