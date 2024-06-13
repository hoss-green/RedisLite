package getset

import (
	"errors"
	"strings"
)

type getParamList struct {
	delete    bool
	hasExpiry bool
	expiry    int64
	persist   bool
}

func parseGetExParams(params []string) (getParamList, error) {
	pList := getParamList{}
	for index := 0; index < len(params); index++ {
		command := strings.ToUpper(params[index])
		switch command {
		case "PERSIST":
			pList.persist = true
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

	if pList.hasExpiry && pList.persist {
		return pList, errors.New("syntax error")
	}

	return pList, nil
}
