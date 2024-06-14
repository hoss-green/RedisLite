package persistence

import "redislite/app/data/storage"


type RdbFile struct {
	RawContents []byte
	DataItems     map[string]storage.DataItem
}

func ReadRdbFromFile(filename string, filedir string) RdbFile {
	filebytes, err := readfile(filename, filedir)

	if err != nil {
    return RdbFile {
      DataItems: map[string]storage.DataItem{},
    }

	}
	return decodeRdb(filebytes)
}
