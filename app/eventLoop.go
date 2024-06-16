package main

import (
	"redislite/app/commands/parsers"
	"redislite/app/commands/parsers/utils"
	"redislite/app/setup"
	"time"
)

const LOOP_DELAY_MS = 100

func ProcessEventLoop(server *setup.Server) {
	looptime := getNewLooptime()
	for {
		select {
		case clientChannel := <-server.ClientChannel:
			parsers.ParseCommand(clientChannel.ClientConnection, clientChannel.RedisCommand, server)
		default:
			goto processloop
		}

	processloop:
    if utils.Expired(looptime) {
		server.DataStore.HandleExpired()
      looptime = getNewLooptime()
    }
	}
}

func getNewLooptime() int64 {
	return time.Now().UTC().Add(time.Duration(LOOP_DELAY_MS * time.Millisecond)).UnixNano()
}
