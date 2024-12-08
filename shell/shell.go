package shell

import (
	"errors"
	"fmt"
	"strings"
)

type Shell interface {
	Cmd(cmd string, args ...string) Command
	CmdBuilder() CommandBuilder
	Run(cmd string, args ...string) RunError
	// show confirmation prompt and wait for user input
	Confirm() error
}

type ShellConfig struct {
	PrintCommands     bool
	Simulate          bool
	ConfirmationLevel int
}

func New() Shell {
	return NewCustom(func(config *ShellConfig) {})
}

func NewCustom(modify func(config *ShellConfig)) Shell {
	config := defaultConfig
	modify(&config)

	if config.Simulate {
		config.PrintCommands = true
	}

	return &shell{config: &config}
}

type shell struct {
	config *ShellConfig
}

var defaultConfig = ShellConfig{
	PrintCommands:     false,
	Simulate:          false,
	ConfirmationLevel: 1,
}

func (s *shell) CmdBuilder() CommandBuilder {
	return &commandBuilder{preRun: s.preRun, simulate: s.config.Simulate, confirmationLevel: 1}
}

func (s *shell) Cmd(cmd string, args ...string) Command {
	return s.CmdBuilder().Cmd(cmd, args...)
}

func (s *shell) Run(cmd string, args ...string) RunError {
	return s.Cmd(cmd, args...).Run()
}

func (s *shell) preRun(c *command) error {
	shouldConfirm := c.confirmationLevel <= s.config.ConfirmationLevel
	shouldPrint := shouldConfirm || s.config.PrintCommands || c.simulate
	if shouldPrint {
		prefix := ">"
		if c.simulate {
			prefix = ">[dry-run]"
		}
		fmt.Println(prefix, c.cmd, c.args)
	}
	if shouldConfirm {
		return s.Confirm()
	}
	return nil
}

func (s *shell) Confirm() error {
	fmt.Print("Confirm [y/N]: ")
	var answer string
	fmt.Scanln(&answer)
	if strings.ToLower(answer) != "y" {
		return errors.New("command aborted")
	}
	return nil
}
