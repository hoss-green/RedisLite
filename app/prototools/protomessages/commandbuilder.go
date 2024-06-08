package protomessages 

import (
	"fmt"
)

func BuildSimpleStringMsg(msg string) string {
	return "+" + msg + "\r\n"
}

func BuildIntMsg(number int) string {
	return fmt.Sprintf(":%d\r\n", number)
}

func BuildNilMsg() string {
	return "$-1\r\n"
}

func BuildBulkStringMsg(msg string) string {
	return "$" + fmt.Sprint(len(msg)) + "\r\n" + msg + "\r\n"
}

func BuildErrorMsg(errmsg string) string {
	return fmt.Sprintf("-%s\r\n", errmsg)
}

func BuildRespArrayMsg(items []string) string {
	arraystring := fmt.Sprintf("*%d\r\n", len(items))
	for _, item := range items {
		arraystring = fmt.Sprintf("%s%s", arraystring, BuildBulkStringMsg(item))
	}

	return arraystring
}

func BuildMultiRespArrayMsg(items []string) string {
	arraystring := fmt.Sprintf("*%d\r\n", len(items))
	for _, item := range items {
		arraystring = fmt.Sprintf("%s%s", arraystring, item)
	}
	return arraystring
}
