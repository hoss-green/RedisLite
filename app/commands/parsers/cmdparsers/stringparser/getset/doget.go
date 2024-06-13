package getset

import (
	"net"
	"redislite/app/data"
	"redislite/app/setup"
)

func Get(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
	return doGet(conn, server, redisCommand, getParamList{
    delete: false,
  })
}

func GetDel(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
	return doGet(conn, server, redisCommand, getParamList{
    delete: true,
  })
}

func GetSet(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
	return doSet(conn, server, redisCommand, paramList{
		get: true,
	})
}

func GetEx(conn net.Conn, server *setup.Server, redisCommand data.RedisCommand) error {
	getPList := getParamList{}
	if len(redisCommand.Params) == 1 {
		return doGet(conn, server, redisCommand, getParamList{})
	}
	if len(redisCommand.Params) > 1 {
		var err error
		getPList, err = parseGetExParams(redisCommand.Params[1:])
		if err != nil {
			return err
		}
	}

	return doGet(conn, server, redisCommand, getPList)
}
