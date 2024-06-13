package getset 

import (
	"strings"
)

// SET key value [NX | XX] [GET] [EX seconds | PX milliseconds | EXAT unix-time-seconds | PXAT unix-time-milliseconds | KEEPTTL]
type paramList struct {
	get            bool
	setInstruction setType
	hasExpiry      bool
	expiry         int64
  keepttl        bool
}

// SET key value [NX | XX] [GET] [EX seconds | PX milliseconds | EXAT unix-time-seconds | PXAT unix-time-milliseconds | KEEPTTL]
func parseParams(params []string) (paramList, error) {
	pList := paramList{}
	for index := 0; index < len(params); index++ {
    command := strings.ToUpper(params[index])
		switch command {
		case "GET":
			pList.get = true
		case "NX":
			pList.setInstruction = NXSetIfKeyNotExist
		case "XX":
			pList.setInstruction = XXOnlySetIfExists
    case "KEEPTTL":
      pList.keepttl = true
		default:
			expiryParseResult, err := tryParseExpiry(command, index, params)
			if err != nil {
				return pList, err
			}

			if expiryParseResult.hasResult {
				pList.hasExpiry = true
				pList.expiry = expiryParseResult.value
				index += 1
			}
		}
	}

	return pList, nil
}
