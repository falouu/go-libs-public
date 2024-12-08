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
	// Starts and wait for completion
	Run() RunError
	// Starts but don't wait for completion. Useful when changing command output to something other than stdout,
	//   using Customize() on builder.
	// `exec.Cmd` is returned to call Wait() on it - but there is no need to call it, if you read the output to
	//   the EOF (in most cases EOF means the process ended)
	Start() (*exec.Cmd, RunError)
}

type RunError interface {
	error
	Command() Command
}

type command struct {
	commandBuilder
}

func (c *command) Run() RunError {
	return c.handleRunError(c.run())
}

func (c *command) Start() (*exec.Cmd, RunError) {
	cmd, err := c.start()
	return cmd, c.handleRunError(err)
}

type commandBuilder struct {
	cmd               string
	args              []string
	preRun            func(c *command) error
	simulate          bool
	confirmationLevel int
	customize         func(*exec.Cmd) error
}

func (b *commandBuilder) Cmd(cmd string, args ...string) Command {
	c := command{*b}
	c.cmd = cmd
	c.args = args
	return &c
}

// the higher the confirmation level, the less likely user will be asked to confirm
func (b *commandBuilder) ConfirmationLevel(level int) CommandBuilder {
	b.confirmationLevel = level
	return b
}

// low level customization. Currently used to stream and process command output in real time
func (b *commandBuilder) Customize(fun func(*exec.Cmd) error) CommandBuilder {
	b.customize = fun
	return b
}

type runError struct {
	error
	command Command
}

func (e *runError) Command() Command {
	return e.command
}

func (c *command) prepareCmd() (*exec.Cmd, error) {
	if err := c.runPreRun(); err != nil {
		return nil, fmt.Errorf("pre run failed: %w", err)
	}
	if c.simulate {
		return nil, nil
	}

	cc := exec.Command(c.cmd, c.args...)
	if c.customize != nil {
		if err := c.customize(cc); err != nil {
			return nil, err
		}
	}
	if cc.Stdout == nil {
		cc.Stdout = os.Stdout
	}
	if cc.Stderr == nil {
		cc.Stderr = os.Stderr
	}
	if cc.Stdin == nil {
		cc.Stdin = os.Stdin
	}
	return cc, nil
}

func (c *command) start() (*exec.Cmd, error) {
	cc, err := c.prepareCmd()
	if err != nil {
		return nil, err
	}
	return cc, cc.Start()
}

func (c *command) run() error {
	cc, err := c.prepareCmd()
	if err != nil {
		return err
	}
	if err := cc.Run(); err != nil {
		return fmt.Errorf("run failed: %w", err)
	}
	return nil
}

func (c *command) handleRunError(err error) RunError {
	if err != nil {
		return &runError{error: err, command: c}
	}
	return nil
}

func (c *command) runPreRun() error {
	if c.preRun != nil {
		return c.preRun(c)
	}
	return nil
}
