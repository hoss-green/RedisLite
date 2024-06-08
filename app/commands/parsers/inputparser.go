package parsers

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"redislite/app/data"
)

type CommandStringsContainer struct {
	CommandStrings     []string
	CommandStringBytes int64
}

func SeparateByLineBreakAdv(message string) []string {
	items := strings.Split(message, "\r\n")
	itemlen := len(items)
	items = items[0:(itemlen - 1)]
	return items
}

func BreakIntoChunks(items []string) ([]CommandStringsContainer, error) {
	commandstringcontainer := []CommandStringsContainer{}
	arraylen := 0
	bytecount := 0
	breakbytes := len([]byte("/r/n"))
	var commandstringitems []string
	// iterate through all the commands
	for index := 0; index < len(items); index++ {

		var err error
		item := items[index]
		firstchar := ([]rune(items[index]))[0]

		if arraylen == 0 {
      bytecount = 0
			commandstringitems = []string{}
		}

		//if we're not in an array and the first char is a star, then it's an arraay.
		if firstchar == '*' && arraylen == 0 {
			bytecount += len([]byte(item)) + 2 
			arraylen, err = getArrayLen(item)
			if err != nil {
				return []CommandStringsContainer{}, errors.New("Could not parse array length")
			}
			continue
		}

		switch firstchar {
		//REDIS Bulk String
		case '$':
			if index+1 > len(items) {
				return []CommandStringsContainer{}, errors.New("Malformed bulk string, have length, missing string")
			}
			bulkstring, remainderitem, err := handleBulkString(item, items[index+1])
			if err != nil {
				return []CommandStringsContainer{}, errors.New("Malformed bulk string")
			}

			bytecount += len([]byte(items[index])) + breakbytes + len([]byte(items[index+1]))
			commandstringitems = append(commandstringitems, bulkstring)
			if remainderitem == "" {
				index += 1
			} else {
				items[index] = remainderitem
			}
			//REDIS Simple String
		case '+':
			commandstringitems = append(commandstringitems, handleSimpleString(item))
			bytecount += len([]byte(item))
		}

		if arraylen > 0 {
			arraylen -= 1
		}
		if arraylen == 0 {
			commandstringcontainer = append(commandstringcontainer,
				CommandStringsContainer{
					CommandStrings:     commandstringitems,
					CommandStringBytes: int64(bytecount)})
		}
	}

	return commandstringcontainer, nil
}

func getArrayLen(item string) (int, error) {
	arrlen, err := strconv.Atoi(item[1:])
	if err != nil {
		return 0, err
	}

	return arrlen, nil
}

func handleSimpleString(word string) string {
	return string([]rune(word)[1:])
}

func handleBulkString(length string, word string) (string, string, error) {
	strlen, err := strconv.Atoi(length[1:])
	if err != nil {
		return "", "", err
	}

	runeword := []rune(word)
	if strlen < len(runeword) {
		return string(runeword[:strlen]), string(runeword[strlen:]), nil
	}

	//can do check here if we must
	return word, "", nil
}

func ReadString(inputBuffer []byte, inputlen int) string {
	message := fmt.Sprintf(string(inputBuffer[:inputlen]))
	return message
}

func BreakIntoSubString(message string) ([]string, error) {
	if message[0] != '*' {
		return []string{}, errors.New("Unsupported Message")
	}
	substrings := strings.Split(message, "*")
	if len(substrings) == 0 {
		return []string{message}, nil
	}
	if len(substrings) == 2 {
		substring := fmt.Sprintf("*%s", substrings[1])
		return []string{substring}, nil
	}
	for index := 1; index < len(substrings); index++ {
		substrings[index] = fmt.Sprintf("*%s", substrings[index])
		log.Println(substrings[index])
	}

	return substrings, nil
}

func BreakIntoArray(message string) []string {

	if message[0] != '*' {
		return make([]string, 0)
	}
	return strings.Split(message, "\r\n")
}

func ParseIntoRedisCommand(rawcommand string, messageparts []string) (data.RedisCommand, bool) {
	var newCommand data.RedisCommand
	newCommand.RawCommand = rawcommand
	log.Println(fmt.Sprintf("%#v", messageparts))
	newCommand.Command = messageparts[2]
	dataCount := (len(messageparts) - 4) / 2
	dataContainer := make([]string, 0)
	if dataCount > 0 {
		for i, data := range messageparts {
			if i > 3 && i%2 == 0 {
				dataContainer = append(dataContainer, data)
			}
		}
	}

	newCommand.ParamLength = len(dataContainer)
	newCommand.Params = dataContainer
	return newCommand, true
}

