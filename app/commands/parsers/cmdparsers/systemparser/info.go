package systemparser 

import (
	"fmt"
	"net"
	"strings"

	"redislite/app/prototools/protomessages"
	"redislite/app/setup"
)

func info(conn net.Conn, serverSettings setup.ServerSettings, section string) error {
	infoParams := make(map[string]string)
	infoParams["role"] = getrole(serverSettings)
	infoParams["port"] = fmt.Sprintf("%d", serverSettings.Port)
	infoParams["master_replid"] = serverSettings.MasterReplId
	infoParams["master_repl_offset"] = fmt.Sprintf("%d", serverSettings.MasterReplIdOffset)
	section = reformatSection(section)
	info := createBulkString(section, infoParams)
	protomessages.QuickSendBulkString(conn, string(info))
	return nil
}

func getrole(serverSettings setup.ServerSettings) string {
	if serverSettings.Master {
		return "master"
	}

	return "slave"
}

func createBulkString(sectionHeader string, stringParams map[string]string) []rune {
	bulkstring := sectionHeader
	for k, v := range stringParams {
		kv := fmt.Sprintf("%s:%s\r\n", k, v)
		bulkstring = bulkstring + kv
	}
	bulkrunes := []rune(bulkstring)
	return bulkrunes
}

func reformatSection(section string) string {
	if len(section) == 0 {
		return section
	}
	runes := []rune(strings.ToLower(section))
	runes[0] = []rune(strings.ToUpper(string(runes[0])))[0]
	return fmt.Sprintf("# %s\r\n", string(runes))
}
