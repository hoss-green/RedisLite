package persistence 

import (
	"encoding/hex"
	"fmt"
)


func Build() []byte {

	final := buildheader()
	// version := "redis-ver.7.2.0"
	final = appendAux(final, "redis-ver", stringValEncode("7.2.0"))
	final = appendAux(final, "redis-bits", numberValEncode(64))
	final = appendTime(final, "ctime", 64)
	// final = appendAux(final, "ctime", timeValEncode(40))
	// final = appendAux(final, "used-mem", stringValEncode("7.2.0"))

	fmt.Printf("Target: %s\r\n", "524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2")
	fmt.Printf("Actual: %s\r\n", hex.EncodeToString(final))

	fin, err := hex.DecodeString("524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2")
	if err != nil {
		panic("malformed")
	}

	return fin

}

func buildheader() []byte {

	string := "REDIS0011"

	return []byte(fmt.Sprintf("%s", string))
}

func appendAux(rdb []byte, header string, value []byte) []byte {
	headerlen := byte(len(header))
	var aux []byte
	aux = append(aux, byte(AUX))

	auxhead := []byte(fmt.Sprintf("%s", header))
	aux = append(aux, headerlen)
	aux = append(aux, auxhead...)
	aux = append(aux, value...)
	aux = append(rdb, aux...)
	// aux := "redis-bits"

	// return []byte
	return aux

}

func appendTime(rdb []byte, header string, value int) []byte {
	headerlen := byte(len(header))
	var aux []byte
	aux = append(aux, []byte("\r\n")...)
	aux = append(aux, byte(AUX))
	auxhead := []byte(fmt.Sprintf("%s", header))
	aux = append(aux, headerlen)
	aux = append(aux, auxhead...)
  aux = append(aux, byte(value))
	aux = append(rdb, aux...)
  return aux 
}
func stringValEncode(value string) []byte {
	headerval := byte(len(value))
	var aux []byte
	auxval := []byte(fmt.Sprintf("%s", value))
	aux = append(aux, headerval)
	aux = append(aux, auxval...)

	return aux
}

func numberValEncode(value int) []byte {
	var aux []byte

	if value < 64 {
		return append(aux, byte(64))
	}

	return aux
}

func timeValEncode(value int) []byte {
	var aux []byte
	aux = append(aux, 0xFC)
	return append(aux, byte(value))

}
func numberlenencoding() byte {

	return 0x00
}
