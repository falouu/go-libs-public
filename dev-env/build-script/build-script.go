package buildscript

import (
	"github.com/falouu/go-libs-public/dev-env/build-script/api"
	console_shell "github.com/falouu/go-libs-public/dev-env/build-script/console-shell"
	"path/filepath"

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
	return console_shell.Run(spec, rootDir)
}
