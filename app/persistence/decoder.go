package persistence

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"log"
	"redislite/app/data/storage"
	"redislite/app/data/storage/datatypes"
	"time"
)

func decodeRdb(rdbFileBytes []byte) RdbFile {
	rdbFile := RdbFile{
		DataItems: make(map[string]storage.DataItem),
	}
	log.Println(hex.EncodeToString(rdbFileBytes))
	fileLen := len(rdbFileBytes)
	log.Println("Bytes in file:", fileLen)
	//first five bytes should be "REDIS"
	redisheader := string(rdbFileBytes[:5])
	log.Println("Redis Header (REDIS): ", redisheader)
	//next four bytes should be a code
	rdbversion := string(rdbFileBytes[5:9])
	log.Println("RDB Version: ", rdbversion)

	checksum := rdbFileBytes[len(rdbFileBytes)-8:]
	log.Println("Checksum: ", checksum)

rdbloop:
	for index := 9; index < fileLen; index++ {
		currentbyte := rdbFileBytes[index]
		_, indicator := CheckIndicatorByte(currentbyte)
		section := []byte{}
		switch indicator {
		case EOF:
			break rdbloop
		case AUX:
			index, section = readBytes(index, rdbFileBytes)
			decodeAux(section)
			continue
		case SELECTDB:
			log.Println("DBSElECTOR")
			index += 1
			dbselector := rdbFileBytes[index]
			log.Println("Db Selector:", dbselector)
			continue
		case RESIZEDB:
			log.Println("RESIZEDB")
			// index, _ = readBytes(index, rdbFileBytes)
			//cross this bridge later
			index += 2
			continue
		default:
			log.Println("Reading Key")
			var expiry int64 = 0
			switch indicator {
			case EXPIRETIME:
				secBytes := rdbFileBytes[index+1 : index+5]
				expiryBuffer := bytes.NewReader(secBytes)
				var nowVar int32
				err := binary.Read(expiryBuffer, binary.LittleEndian, &nowVar)
				if err != nil {
					log.Fatal("FD Time set incorrectly")
				}
				expiry = int64(nowVar)
				index += 5
			case EXPIRETIMEMS:
				msBytes := rdbFileBytes[index+1 : index+9]
				expiryBuffer := bytes.NewReader(msBytes)
				err := binary.Read(expiryBuffer, binary.LittleEndian, &expiry)
				if err != nil {
					log.Fatal("FE Time set incorrectly")
				}
				index += 9
			}
			key, val, totLen, err := ReadKvPairAndSkip(rdbFileBytes[index:])
			if err != nil {
				log.Println("Could not parse: ", err.Error())
			} else {
				exptime := time.Unix(0, int64(expiry)*int64(time.Millisecond)).UTC()
				t := time.Now().UTC()
				log.Printf("CUR:  %d | %s \r\n", t.UnixNano(), t.String())
				log.Printf("EXP:  %d | %s \r\n", exptime.UnixNano(), exptime.String())
				log.Println("EXP: ", exptime)
				log.Println("KEY: ", key)
				log.Println("VAL: ", val)
				rdbFile.DataItems[key] = storage.CreateDataItem(key, datatypes.DATA_TYPE_STRING, []byte(val), exptime.UnixNano())
				// storage.DataItem{
				// 	Key:   key,
				// 	Value: val,
				// 	// ExpiryTimeNano: exptime.UnixNano(),
				// }
				index += totLen
			}
		}

		// log.Printf("sec: %d", len(section))
	}

	return rdbFile
}

func decodeAux(section []byte) {
	k, v, _, _ := ReadKvPairAndSkip(section)

	log.Printf("Key %s, Value %s", k, v)
}

func decodeString(chunk []byte) {
	// stringEnconding :=

}

func readBytes(startIndex int, filebytes []byte) (int, []byte) {
	fileLen := len(filebytes)
	sectionBytes := []byte{}
	counter := 0
	for newIndex := startIndex + 1; newIndex < fileLen; newIndex++ {
		currentbyte := filebytes[newIndex]
		isIndicator, _ := CheckIndicatorByte(currentbyte)
		if isIndicator {
			counter = newIndex - 1
			break
		}
		sectionBytes = append(sectionBytes, currentbyte)
	}
	log.Printf("Read from %d to %d", startIndex, counter)
	return counter, sectionBytes
}
