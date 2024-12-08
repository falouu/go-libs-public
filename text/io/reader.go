package io

import (
	"bufio"
	"strings"
)

func ReadLineUntilText(r *bufio.Reader, texts []string) (line string, foundText string, err error) {
	var byt byte
	bytes := []byte{}
	for {
		byt, err = r.ReadByte()
		if err != nil {
			return
		}
		if byt == '\n' {
			if len(bytes) > 0 && bytes[len(bytes)-1] == '\r' {
				bytes = bytes[0 : len(bytes)-1]
			}
			return string(bytes), "", nil
		}
		bytes = append(bytes, byt)
		str := string(bytes)

		for _, text := range texts {
			if strings.HasSuffix(str, text) {
				return str, text, nil
			}
		}
	}
}
