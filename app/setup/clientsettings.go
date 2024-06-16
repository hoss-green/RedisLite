package setup

import (
	"net"
	"redislite/app/data"
)

type Client struct {
	RedisCommand  data.RedisCommand
	ClientConnection *net.Conn
}
