package b

import (
	"errors"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapErrors(t *testing.T) {
	// given
	cause := errors.New("This is the cause")

	// when
	err := Wrap(cause, "Outer description")

	// then
	assert.Equal(t, "Outer description. Caused by: This is the cause", err.Error())
}

func TestWrapErrorsWithFormat(t *testing.T) {
	// given
	cause := errors.New("This is the cause")

	// when
	err := Wrap(cause, "Outer description with int %d and string %s", 12, "lorem")

	// then
	assert.Equal(t, "Outer description with int 12 and string lorem. Caused by: This is the cause", err.Error())
}

func TestAddStderrToExitError(t *testing.T) {
	// given
	GlobalSettings.AddStderrToExitErrors = true

	cmd := helperCommand(t, "fail_with_stderr")
	_, cause := cmd.Output()
	assert.IsType(t, &exec.ExitError{}, cause)

	// when
	err := Wrap(cause, "Outer description")

	// then
	expectedMessage := strings.Join([]string{
		"Outer description. Caused by: exit status 2",
		"Stderr:",
		"this is stderr of a failed process",
	}, "\n")

	assert.Equal(t, expectedMessage, err.Error())
}

func TestNotAddStderrToExitError(t *testing.T) {
	// given
	GlobalSettings.AddStderrToExitErrors = false

	cmd := helperCommand(t, "fail_with_stderr")
	_, cause := cmd.Output()
	assert.IsType(t, &exec.ExitError{}, cause)

	// when
	err := Wrap(cause, "Outer description")

	// then
	expectedMessage := "Outer description. Caused by: exit status 2"

	assert.Equal(t, expectedMessage, err.Error())
}

func helperCommand(t *testing.T, commandId string) *exec.Cmd {
	buildCmd := exec.Command("go", "build", "helper_process.go")
	buildCmd.Dir = "testdata"
	if err := buildCmd.Run(); err != nil {
		t.Fatal("Failed to build helper command. Caused by: " + err.Error())
	}

	cmd := exec.Command("testdata/helper_process", commandId)
	return cmd
}
