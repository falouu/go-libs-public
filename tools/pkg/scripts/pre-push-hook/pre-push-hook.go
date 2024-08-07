package prepushhook

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/falouu/go-libs-public/script"
	"github.com/falouu/go-libs-public/tools/pkg/check"
	log "github.com/sirupsen/logrus"
)

var excludedPatterns = []regexp.Regexp{
	*regexp.MustCompile(`^playground/`),
}

var zeroCommit = regexp.MustCompile(`^0+$`)

func Run() {

	if _, err := script.Init(); err != nil {
		log.Panic(err)
	}

	input := bufio.NewScanner(os.Stdin)
	input.Split(bufio.ScanWords)
	args := []string{}
	for input.Scan() {
		args = append(args, input.Text())
	}

	if len(args) != 4 {
		log.Panicf("Expected 4 argument, received %v: %v", len(args), strings.Join(args, ", "))
	}

	// localRef := args[0]
	localCommit := args[1]
	// remoteRef := args[2]
	remoteCommit := args[3]
	if zeroCommit.MatchString(remoteCommit) {
		remoteCommit = "origin/main"
	}

	log.Debugf(
		"remote commit: %v\n"+
			"local commit: %v\n",
		remoteCommit, localCommit,
	)

	if !shouldRunGoTest(remoteCommit, localCommit) {
		log.Debugf("No go changes - skipping testing")
		os.Exit(0)
	}

	os.Exit(check.Check())
}

func shouldRunGoTest(remoteCommit string, localCommit string) bool {

	if zeroCommit.MatchString(localCommit) {
		log.Debug("local commit is zero, running checks because better be safe")
		return true
	}

	output, err := exec.Command("git", "diff", "--name-only", remoteCommit, localCommit).Output()
	if err != nil {
		if exerr, ok := err.(*exec.ExitError); ok {
			log.Panicf("%v\nStderr:\n%v\n", err, string(exerr.Stderr))
		}
		log.Panic(err)
	}

	filesScanner := bufio.NewScanner(bytes.NewReader(output))
	for filesScanner.Scan() {
		file := filesScanner.Text()
		for _, excludedPattern := range excludedPatterns {
			if !excludedPattern.MatchString(file) {
				return true
			}
		}
	}
	return false
}
