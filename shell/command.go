package shell

import (
	"fmt"
	"os"
	"os/exec"
)

type CommandBuilder interface {
	Cmd(cmd string, args ...string) Command
	ConfirmationLevel(level int) CommandBuilder
}

type Command interface {
	Run() RunError
}

type RunError interface {
	error
	Command() Command
}

func (c *command) Run() RunError {
	if err := c.run(); err != nil {
		return &runError{error: err, command: c}
	}
	return nil
}

type command struct {
	commandBuilder
}

type commandBuilder struct {
	cmd               string
	args              []string
	preRun            func(c *command) error
	simulate          bool
	confirmationLevel int
}

func (b *commandBuilder) Cmd(cmd string, args ...string) Command {
	c := command{*b}
	c.cmd = cmd
	c.args = args
	return &c
}

func (b *commandBuilder) ConfirmationLevel(level int) CommandBuilder {
	b.confirmationLevel = level
	return b
}

type runError struct {
	error
	command Command
}

func (e *runError) Command() Command {
	return e.command
}

func (c *command) run() error {
	if err := c.preRun(c); err != nil {
		return fmt.Errorf("pre run failed: %w", err)
	}
	if c.simulate {
		return nil
	}

	cc := exec.Command(c.cmd, c.args...)
	cc.Stdout = os.Stdout
	cc.Stderr = os.Stderr
	cc.Stdin = os.Stdin
	if err := cc.Run(); err != nil {
		return fmt.Errorf("run failed: %w", err)
	}
	return nil
}
