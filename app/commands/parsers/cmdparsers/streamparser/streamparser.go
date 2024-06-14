package streamparser

// import (
// 	"errors"
// 	"net"
// 	"strings"
//
// 	"redislite/app/commands/errormessages"
// 	"redislite/app/commands/parsers/parserentities"
// 	"redislite/app/data"
// 	"redislite/app/setup"
// )
//
// // STREAM COMMAND PARSER
// func ParseStreamCommand(connpointer *net.Conn, redisCommand data.RedisCommand, server *setup.Server) parserentities.ParserInfo { 
// 	conn := *connpointer
// 	var err error = nil
// 	instruction := strings.ToUpper(redisCommand.Command)
// 	switch instruction {
// 	case "XADD":
// 		if redisCommand.ParamLength < 2 || redisCommand.ParamLength % 2 != 0{
// 			err = errors.New(errormessages.IncorrectArgumentsError)
// 			break
// 		}
//     err = xadd(conn, server, redisCommand)
//   case "XRANGE":
//     if redisCommand.ParamLength != 3 {
// 			err = errors.New(errormessages.IncorrectArgumentsError)
// 			break
//     }
//     err = xrange(conn, server, redisCommand)
//   case "XREAD":
//     if redisCommand.ParamLength < 3 {
// 			err = errors.New(errormessages.IncorrectArgumentsError)
// 			break
//     }
//     err = xread(conn, server, redisCommand)
// 	default:
// 		return parserentities.ParserInfo{Executed: false, Err: nil}
// 	}
//
// 	if err != nil {
// 		return parserentities.ParserInfo{Executed: false, Err: err}
// 	}
//
// 	return parserentities.ParserInfo{Executed: true, Err: nil}
// }
