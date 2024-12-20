package utils

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

func RemoveSuffix(str string, sep string) string {
	res := ""

	slice := strings.Split(str, sep)
	// string doesn't have sep
	if len(slice) == 1 {
		return slice[0]
	}

	// removing the suffix
	slice = slice[:len(slice)-1]

	// assembling the leftovers
	for _, el := range slice {
		res += el
	}

	return res
}

func BufferReader(reader io.Reader) (*bytes.Buffer, error) {
	if reader == nil {
		return nil, errors.New("BufferReader: reader is nil")
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	return buf, nil
}
