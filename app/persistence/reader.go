package persistence

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func readfile(filename string, filedir string) ([]byte, error) {
	file := filename
	if filedir != "" {
		file = filepath.Join(filedir, filename)
	}
	log.Println(fmt.Sprintf("Reading rdb file from %s", file))
	content, err := os.ReadFile(file)

	if err != nil {
		log.Println("Could not read rdb file")
    return []byte{}, errors.New("File not found error")
	}

	fmt.Println(string(content))

	return content, nil
}
