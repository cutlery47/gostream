package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
)

func RemoveSuffix(str string, sep string) string {
	res := ""

	slice := strings.Split(str, sep)
	// removing the suffix
	slice = (slice)[:len(slice)-1]
	for _, el := range slice {
		res += el
	}

	return res
}

func BufferFile(file *os.File) (*bytes.Buffer, error) {
	if file == nil {
		return nil, errors.New("BufferFile: file is nil")
	}

	meta, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("BufferFile: %v", err)
	}

	buffer := bytes.NewBuffer(make([]byte, meta.Size()))
	file.Read(buffer.Bytes())

	return buffer, nil
}
