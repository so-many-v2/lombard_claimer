package filereader

import (
	"bufio"
	"log"
	"os"
)

type FileReader struct{}

func NewFileReader() *FileReader {
	return &FileReader{}
}

func (fr *FileReader) ScanFile(path string) *bufio.Scanner {

	file, err := os.Open(path)

	if err != nil {
		log.Fatalf("error getting wallets from %s : %s", path, err.Error())
	}

	scanner := bufio.NewScanner(file)

	return scanner
}
