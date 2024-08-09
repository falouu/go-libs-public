package console_shell

import (
	_ "embed"
	"github.com/falouu/go-libs-public/dev-env/build-script/api"
	"github.com/sirupsen/logrus"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

type Specification = api.Specification

func Run(spec *Specification) error {
	rootDir := spec.RootDir
	buildDir := filepath.Join(rootDir, ".build")
	if err := os.MkdirAll(buildDir, 0700); err != nil {
		return err
	}
	log.Debug("buildDir: ", buildDir)

	requirementEnvironment := requirementEnvironment{buildDir: buildDir}
	addToPath := []string{}
	appendToBashRc := []*api.Template{}

	for i := range spec.Requirements {
		req := spec.Requirements[i]

		result, err := req.Ensure(&requirementEnvironment)
		if err != nil {
			return err
		}

		addToPath = append(addToPath, result.AddToPath...)
		appendToBashRc = append(appendToBashRc, result.BashRcAppend)
	}
	return startShell(addToPath, appendToBashRc, spec.PromptText)
}

func startShell(addToPath []string, appendToBashRc []*api.Template, promptText string) error {

	bashPath, err := exec.LookPath("bash")
	if err != nil {
		return err
	}

	bashInitFile, err := createBashInitFile(addToPath, appendToBashRc, promptText)
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

func createBashInitFile(addToPath []string, append []*api.Template, promptText string) (path string, err error) {
	bashInitFile, err := os.CreateTemp("", "go-build-script-bash-init-")
	if err != nil {
		return "", err
	}

	templateText := bashInitFileTemplate
	extra := map[string]any{}
	for _, appendTemplate := range append {
		templateText += "\n" + appendTemplate.Template
		maps.Copy(extra, appendTemplate.Args)
	}

	initFileTmpl, err := template.New("").Parse(templateText)
	if err != nil {
		return "", err
	}

	addToPathString := strings.Join(addToPath, ":")

	if promptText == "" {
		promptText = "buildscript"
	}
	err = initFileTmpl.Execute(bashInitFile, struct {
		AddToPath  string
		PromptText string
		Extra      map[string]any
	}{
		AddToPath:  addToPathString,
		Extra:      extra,
		PromptText: promptText,
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
