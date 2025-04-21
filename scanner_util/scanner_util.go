package scanner_util

import (
	"bufio"
	"encoding/binary"
	"github.com/pkg/errors"
	"io"
)

const maxBufferSize = 100 << 20 // 100MB

func splitLengthPrefixU32(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, bufio.ErrFinalToken
	}

	if len(data) < 4 {
		return 0, nil, nil
	}

	length := int(binary.LittleEndian.Uint32(data[:4]))

	if length > maxBufferSize {
		return 0, nil, errors.New("too large, max 100MB")
	}

	if len(data) < length+4 {
		return 0, nil, nil
	}

	return length + 4, data[4 : length+4], nil
}

func NewLengthPrefixedReader(r io.Reader, maxTokenSize int) *bufio.Scanner {
	scanner := bufio.NewScanner(r)
	scanner.Split(splitLengthPrefixU32)
	scanner.Buffer(make([]byte, 0, 102400), maxBufferSize)
	return scanner
}
