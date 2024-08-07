package b

import (
	"fmt"
	"os/exec"
)

type settings struct {
	AddStderrToExitErrors bool
}

var GlobalSettings = settings{
	AddStderrToExitErrors: false,
}

type ErrorWithCause struct {
	Cause                 error
	addStderrToExitErrors bool
	OuterMsg              string
}

func (e *ErrorWithCause) Error() string {

	msg := fmt.Sprintf("%s. Caused by: %s", e.OuterMsg, e.Cause.Error())
	if e.addStderrToExitErrors {
		if exitErr, ok := e.Cause.(*exec.ExitError); ok {
			msg += "\nStderr:\n" + string(exitErr.Stderr)
		}
	}

	return msg
}

func Wrap(err error, msg string, args ...interface{}) error {
	return &ErrorWithCause{
		Cause:                 err,
		addStderrToExitErrors: GlobalSettings.AddStderrToExitErrors,
		OuterMsg:              fmt.Sprintf(msg, args...),
	}
}
