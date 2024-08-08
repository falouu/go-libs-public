package fs

import (
	"errors"
	"os"
)

func ReadFileIfExists(path string) (content []byte, isExist bool, err error) {
	bytes, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return bytes, true, nil
}

func IsFile(path string) bool {
	stat, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	if err != nil {
		panic(err)
	}
	return stat.Mode().IsRegular()
}
