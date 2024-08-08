package b

import (
	"github.com/falouu/go-libs-public/b/internal/fs"
)

// Read file contents.
// If file exists and there is no error:
//
//	content = <file content>
//	isExists = true
//	err = nil
//
// If file doesn't exist:
//
//	isExist = false
//	err = nil
//
// If there was an error:
//
//	err != nil
func ReadFileIfExists(path string) (content []byte, isExist bool, err error) {
	return fs.ReadFileIfExists(path)
}

func IsFile(path string) bool {
	return fs.IsFile(path)
}
