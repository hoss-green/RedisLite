package getset 

import (
	"net"
	"redislite/app/data"
	"redislite/app/setup"
)


func Set(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
	setParamsList := paramList{}
	if len(redisCommand.Params) > 1 {
		var err error
		setParamsList, err = parseParams(redisCommand.Params[2:])
		if err != nil {
			return err
		}
	}
	return doSet(conn, server, redisCommand, setParamsList)
}

