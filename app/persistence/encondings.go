package persistence

import (
	"errors"
	"log"
)

type KvEncoding byte
type ByteEncoding byte
type StringLengthEncoding int

type RdbFilePart struct {
	Type     ByteEncoding
	RawBytes []byte
}

type AuxPart struct {
	RawBytes []byte
	Key      string
	Value    string
}

const (
	//None, used to help parsing
	NONE ByteEncoding = 0x00
)


const (
	// 0xFF EOF End of the RDB file
	EOF ByteEncoding = 0xFF
	// 0xFE SELECTDB Database Selector
	SELECTDB ByteEncoding = 0xFE
	// 0xFD EXPIRETIME Expire time in seconds. 4 bytes following are the expiry
	EXPIRETIME ByteEncoding = 0xFD
	// 0xFC EXPIRETIMEMS Expire time in milliseconds. 8 bytes following are the expiry 
	EXPIRETIMEMS ByteEncoding = 0xFC
	// 0xFB RESIZEDB Hash table sizes for the main keyspace and expires, see Resizedb information
	RESIZEDB ByteEncoding = 0xFB
	// 0xFA AUX Auxiliary fields. Arbitrary key-value settings, see Auxiliary fields
	AUX ByteEncoding = 0xFA
)

const (
	//the next six bits represent the length
	SHORTLEN StringLengthEncoding = 0
	//Read one additional byte. The combined 14 bits represent the length
	EXTRABYTE StringLengthEncoding = 1
	//Discard the remaining 6 bits. The next 4 bytes from the stream represent the StringLengthEncoding
	FOURBYTE StringLengthEncoding = 2
	//The next object is encoded in a special format.
	SPEICALFORMAT StringLengthEncoding = 3
)

const (
	// String Encoding
	ENCSTRING KvEncoding = 0x00
	// List Encoding
	//ENCLIST KvEncoding = 1
	//ENCSET KvEncoding = 2 = Set Encoding
	//ENCSORTSET KvEncoding = 3 = Sorted Set Encoding
	//ENCHASH KvEncoding = 4 = Hash Encoding
	//ENCZIP KvEncoding = 9 = Zipmap Encoding
	//ENCZIPLIST KvEncoding = 10 = Ziplist Encoding
	//ENCINSET KvEncoding = 11 = Intset Encoding
	//ENCSORTSETZIP KvEncoding = 12 = Sorted Set in Ziplist Encoding
	//ENCHASHZIP KvEncoding = 13 = Hashmap in Ziplist Encoding (Introduced in RDB version 4)
	//ENCLISTQUICK KvEncoding = 14 = List in Quicklist encoding
)

func CheckIndicatorByte(b byte) (bool, ByteEncoding) {
	switch b {
	case 0xFF:
		log.Println("End of file")
		return true, EOF
	case 0xFE:
		log.Println("Select Db")
		return true, SELECTDB
	case 0xFD:
		return true, EXPIRETIME
	case 0xFC:
		return true, EXPIRETIMEMS
	case 0xFB:
		log.Println("Resize Db")
		return true, RESIZEDB
	case 0xFA:
		// log.Println("Aux")
		return true, AUX
	default:
		return false, NONE
	}
}

func ReadKvPairAndSkip(totalBytes []byte) (string, string, int, error) {
	kvType := totalBytes[0]
	if kvType == byte(ENCSTRING) {
		log.Println("String")
		// kLength, kSkipLength := GetStringLen(totalBytes[1:])
		//  return string(totalBytes[kSkipLength : kSkipLength+kLength]), nil
		k, v, tot := ParseKVString(totalBytes[1:])
		return k, v, tot, nil
	}

	return "", "", 0, errors.New("Unsupported Encoding")

}
func ParseKVString(totalBytes []byte) (string, string, int) {
	key := ""
	value := ""
	kLength, kSkipLength := GetStringLen(totalBytes[:])
	key = string(totalBytes[kSkipLength : kSkipLength+kLength])
	totalKSkip := kSkipLength + kLength
	vLength, vSkipLength := GetStringLen(totalBytes[totalKSkip:])
	value = string(totalBytes[totalKSkip+vSkipLength : totalKSkip+vSkipLength+vLength])

	totalBytesConsumed := totalKSkip + vSkipLength + vLength
	return key, value, totalBytesConsumed
}

func GetStringLen(value []byte) (int, int) {
	b := value[0]
	firstBit := (b>>7)&1 == 1
	secondBit := (b>>6)&1 == 1
	log.Println(firstBit, secondBit)

	sixbits := (b << 2) >> 2
	log.Println("StrLen:", int(sixbits))

	return int(sixbits), 1
}
