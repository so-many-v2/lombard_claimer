package filereader

import (
	"bufio"
	"fmt"
	"os"
)

type FileReader struct{}

func NewFileReader() *FileReader {
	return &FileReader{}
}

func (fr *FileReader) ScanFile(path string) ([]string, error) {

	file, err := os.Open(path)

	if err != nil {
		return nil, fmt.Errorf("error getting wallets from %s : %s", path, err.Error())
	}

	scanner := bufio.NewScanner(file)

	output := []string{}

	for scanner.Scan() {
		output = append(output, scanner.Text())
	}

	return output, nil
}
