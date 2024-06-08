package instructions

import (
	"net"

	"redislite/app/data"
	"redislite/app/prototools/protomessages"
	"redislite/app/setup"
)

func Keys(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
	keys := []string{}
	for k := range server.Rdb.KVPairs {
		keys = append(keys, k)
	}
	return protomessages.SendRespArray(conn, keys)
}
