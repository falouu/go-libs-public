/// 2>/dev/null; export APP_CWD="$(pwd)" && cd "$(dirname "$0")" && exec /usr/bin/env go run "$(basename "$0")" "$@"

//go:build main

package main

import prepushhook "github.com/falouu/go-libs-public/tools/pkg/scripts/pre-push-hook"

func main() {
	prepushhook.Run()
}
