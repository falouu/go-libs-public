package console_shell

import (
	_ "embed"
	"github.com/falouu/go-libs-public/dev-env/build-script/api"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

type Specification = api.Specification

func Run(spec *Specification, rootDir string) error {
	buildDir := filepath.Join(rootDir, ".build")
	if err := os.MkdirAll(buildDir, 0700); err != nil {
		return err
	}
	log.Debug("buildDir: ", buildDir)

	requirementEnvironment := requirementEnvironment{buildDir: buildDir}
	addToPath := []string{}

	for i := range spec.Requirements {
		req := spec.Requirements[i]

		result, err := req.Ensure(&requirementEnvironment)
		if err != nil {
			return err
		}

		addToPath = append(addToPath, result.AddToPath...)
	}
	return startShell(addToPath)
}

func startShell(addToPath []string) error {

	bashPath, err := exec.LookPath("bash")
	if err != nil {
		return err
	}

	bashInitFile, err := createBashInitFile(addToPath)
	if err != nil {
		return err
	}
	defer os.Remove(bashInitFile)

	bashArgs := []string{"--init-file", bashInitFile, "-i"}
	log.Debugf("exec: %v %v", bashPath, strings.Join(bashArgs, " "))
	// return syscall.Exec(bashPath, bashArgs, os.Environ())
	cmd := exec.Command(bashPath, bashArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func createBashInitFile(addToPath []string) (path string, err error) {
	bashInitFile, err := os.CreateTemp("", "go-build-script-bash-init-")
	if err != nil {
		return "", err
	}

	initFileTmpl, err := template.New("").Parse(bashInitFileTemplate)
	if err != nil {
		return "", err
	}

	addToPathString := strings.Join(addToPath, ":")

	err = initFileTmpl.Execute(bashInitFile, struct {
		AddToPath string
	}{
		AddToPath: addToPathString,
	})
	if err != nil {
		return "", err
	}
	return bashInitFile.Name(), bashInitFile.Sync()
}

type requirementEnvironment struct {
	buildDir string
}

func (e *requirementEnvironment) BuildDir() string {
	return e.buildDir
}

//go:embed bash-init-file.sh
var bashInitFileTemplate string

var log = logrus.WithField("src", "build-script")
