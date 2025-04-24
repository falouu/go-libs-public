package terminalutils

import (
	"errors"
	"fmt"
	"strings"
)

func Confirm(errMsg string) error {
	fmt.Print("Confirm [y/N]: ")
	var answer string
	if _, err := fmt.Scanln(&answer); err != nil {
		return err
	}
	if strings.ToLower(answer) != "y" {
		return errors.New(errMsg)
	}
	return nil
}
