package persistence

import (
	"redislite/app/data/datatypes/kvstring"
)

type RdbFile struct {
	RawContents []byte
	KVPairs     map[string]kvstring.KvString
}

func ReadRdbFromFile(filename string, filedir string) RdbFile {
	filebytes, err := readfile(filename, filedir)

	if err != nil {
    return RdbFile {
      KVPairs: map[string]kvstring.KvString{},
    }

	}
	return decodeRdb(filebytes)
}
