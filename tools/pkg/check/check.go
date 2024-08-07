package check

import "github.com/falouu/go-libs-public/shell"

func Check() (exitCode int) {
	return shell.Run("go build ./... && go test -v ./...")
}
