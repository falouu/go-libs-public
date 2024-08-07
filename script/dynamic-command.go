package script

import (
	"github.com/alecthomas/kong"
	"github.com/posener/complete"
)

type DynamicCommand struct {
	Cmd        any
	Tags       []string
	Help       string
	Predictors map[string]complete.Predictor
}

type DynamicCommands map[string]DynamicCommand

func (cmds DynamicCommands) toKong() []kong.Option {
	out := []kong.Option{}
	for name, cmd := range cmds {
		out = append(out, kong.DynamicCommand(name, cmd.Help, "", cmd.Cmd, cmd.Tags...))
	}
	return out
}

func (cmds DynamicCommands) Predictors() map[string]complete.Predictor {
	out := map[string]complete.Predictor{}
	for _, cmd := range cmds {
		for key, value := range cmd.Predictors {
			out[key] = value
		}
	}
	return out
}
