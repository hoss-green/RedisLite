package doget 

import (
	"net"
	"redislite/app/data"
	"redislite/app/setup"
)


func Get(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error { 
  return doGet(conn, server, redisCommand, false)
}

func GetDel(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error { 
  return doGet(conn, server, redisCommand, false)
}
