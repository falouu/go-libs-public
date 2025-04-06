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

	script, err := script.Init()
	if err != nil {
		return err
	}

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

type FuncRequirement func(env RequirementEnvironment) (*RequirementResult, error)

func (f FuncRequirement) Info() *RequirementInfo {
	return nil
}
func (f FuncRequirement) Ensure(env RequirementEnvironment) (*RequirementResult, error) {
	return f(env)
}
