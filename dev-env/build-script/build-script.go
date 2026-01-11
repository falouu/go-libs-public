package buildscript

import (
	"path/filepath"

	"github.com/falouu/go-libs-public/dev-env/build-script/api"
	console_shell "github.com/falouu/go-libs-public/dev-env/build-script/console-shell"

	"github.com/falouu/go-libs-public/script"
)

type Specification = api.Specification
type Requirement = api.Requirement
type RequirementEnvironment = api.RequirementEnvironment
type RequirementResult = api.RequirementResult
type RequirementInfo = api.RequirementInfo

func Run(spec *Specification) error {
	return RunWithContext(func(ctx Context) (*Specification, error) {
		return spec, nil
	})
}

type SpecSupplier func(context Context) (*Specification, error)

// Prefer over Run() if the context is needed to initialize Specification (e.g. access to script instance)
func RunWithContext(spec SpecSupplier) error {
	script, err := script.Init()
	if err != nil {
		return err
	}
	ctx := Context{Script: script}
	s, err := spec(ctx)
	if err != nil {
		return err
	}
	return run(s, script)
}

func run(spec *Specification, script script.Script) error {
	rootDir := script.Dir()
	if spec.RootDir != "" {
		if filepath.IsAbs(spec.RootDir) {
			rootDir = spec.RootDir
		} else {
			rootDir = filepath.Join(rootDir, spec.RootDir)
		}
	}
	spec.RootDir = rootDir
	return console_shell.Run(spec)
}

type Context struct {
	Script script.Script
}

type FuncRequirement func(env RequirementEnvironment) (*RequirementResult, error)

func (f FuncRequirement) Info() *RequirementInfo {
	return nil
}
func (f FuncRequirement) Ensure(env RequirementEnvironment) (*RequirementResult, error) {
	return f(env)
}
