package script

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/falouu/go-libs-public/b"
	log "github.com/sirupsen/logrus"
	"github.com/willabides/kongplete"
)

type Script interface {
	// directory where go script is located
	Dir() string
	// directory from where go script was called
	Cwd() string
	// run command as defined by passed Kong CLI
	RunCommand(context any) error
}

type Options struct {
	CLI             interface{}
	Description     string
	DynamicCommands DynamicCommands
}

func InitCustom(modifyOptions func(options *Options)) (Script, error) {
	options := defaultOptions
	modifyOptions(&options)
	return initScript(&options)
}

func Init() (Script, error) {
	return initScript(&defaultOptions)
}

func JustRun(fun runFunc) {
	if _, err := Init(); err != nil {
		log.Panic(err)
	}
	os.Exit(fun())
}

type runFunc func() (exitCode int)

var defaultOptions = Options{
	CLI:             &struct{}{},
	DynamicCommands: DynamicCommands{},
}

var shebang = `///bin/true; export APP_CWD="$(pwd)" && cd "$(dirname "$0")" && exec /usr/bin/env go run "$0" "$@"`

func initScript(options *Options) (Script, error) {

	file, err := getScriptFilepath()
	if err != nil {
		return nil, b.Wrap(err, "Cannot determine script filepath")
	}

	cwd, err := fixCwd()
	if err != nil {
		return nil, err
	}

	var cli struct {
		kong.Plugins
		InstallCompletions InstallCompletions `short:"i" help:"install shell completions"`
		LogLevel           string             `help:"Set log level"`
	}

	script := script{
		filepath: file,
		cwd:      cwd,
	}

	cli.Plugins = kong.Plugins{options.CLI}
	kongOptions := []kong.Option{
		kong.Description(options.Description),
		kong.Bind(&script),
		kong.UsageOnError(),
	}
	kongOptions = append(kongOptions, options.DynamicCommands.toKong()...)
	kongParser := kong.Must(&cli, kongOptions...)

	// if err := kongCtx.Validate(); err != nil {
	// 	return nil, err
	// }

	kongplete.Complete(kongParser, kongplete.WithPredictors(options.DynamicCommands.Predictors()))

	script.kongCtx, err = kongParser.Parse(os.Args[1:])
	kongParser.FatalIfErrorf(err)

	if cli.LogLevel != "" {
		level, err := log.ParseLevel(cli.LogLevel)
		if err != nil {
			return nil, err
		}
		log.SetLevel(level)
	}

	log.Debug("script path: " + file)
	log.Debug("cwd: ", cwd)
	log.Debug("log level: " + log.GetLevel().String())

	return &script, nil
}

type script struct {
	filepath string
	cwd      string
	kongCtx  *kong.Context
}

func (s *script) Dir() string {
	return filepath.Dir(s.filepath)
}

func (s *script) Cwd() string {
	return s.cwd
}

func (s *script) RunCommand(context any) error {
	return s.kongCtx.Run(context)
}

type InstallCompletions bool

func (f InstallCompletions) BeforeReset(script *script) error {
	name := filepath.Base(os.Args[0])
	path := resolveHome(script.filepath)
	fmt.Println("# To get command and params completions, execute:\n" +
		"alias " + name + "='" + path + "'\n" +
		"complete -C \"" + path + "\" " + name + "\n" +
		"# Add it to .bashrc for permanent effect")
	os.Exit(0)
	return nil
}

func getScriptFilepath() (string, error) {
	rpc := make([]uintptr, 100)
	n := runtime.Callers(1, rpc)
	if n < 1 {
		return "", errors.New("no frames returned")
	}

	frames := runtime.CallersFrames(rpc)

	//more := true
	//f, more := frames.Next()

	more := true
	for more {
		var f runtime.Frame
		f, more = frames.Next()
		if strings.HasPrefix(f.Function, "main.") {
			return f.File, nil
		}
	}

	return "", errors.New("cannot find frame with 'main' package")
}

var goRunBuildDirDetectPattern = regexp.MustCompile("go-build[0-9]")

func fixCwd() (string, error) {
	cwd, cwdFound := os.LookupEnv("APP_CWD")
	if !cwdFound {
		origCwd, err := os.Getwd()
		if err != nil {
			return "", b.Wrap(err, "Cannot obtain cwd and APP_CWD env var is missing")
		}
		if goRunBuildDirDetectPattern.MatchString(os.Args[0]) {
			return "", fmt.Errorf("It looks like you run *.go script directly without building a binary. "+
				"You can do that, but you have to set APP_CWD env var, which is not found. "+
				"To fix that, add proper first line in your script:\n\n%v\n\nset executable flag:\n\n"+
				"chmod +x <your script>.go\n\n"+
				"and run it directly:\n\n./<your script>.go", shebang)
		}
		cwd = origCwd
	}
	if err := os.Chdir(cwd); err != nil {
		return "", b.Wrap(err, "Cannot set working directory to directory provided by APP_CWD=%v", cwd)
	}
	return cwd, nil
}

func resolveHome(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	rel, err := filepath.Rel(home, path)
	if err != nil {
		return path
	}
	if strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return path
	}
	return "${HOME}/" + rel
}
