package main

import (
	"fmt"
	"os"
)

func main() {
	//os.Exit(3)
	//fmt.Fprint(os.Stderr, "this is stderr of a failed process")

	command := os.Args[1]

	switch command {
	case "fail_with_stderr":
		fmt.Fprint(os.Stderr, "this is stderr of a failed process")
		os.Exit(2)
	}
}
