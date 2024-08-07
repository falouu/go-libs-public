/// 2>/dev/null; export APP_CWD="$(pwd)" && cd "$(dirname "$0")" && exec /usr/bin/env go run "$(basename "$0")" "$@"

//go:build main

package main

import (
	"github.com/falouu/go-libs-public/script"
	"github.com/falouu/go-libs-public/tools/pkg/check"
)

// run it only from root dir (scripts)
func main() {
	script.JustRun(func() int { return check.Check() })
}
