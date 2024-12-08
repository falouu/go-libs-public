package shell

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"testing"
	"time"

	textio "github.com/falouu/go-libs-public/text/io"
	"github.com/stretchr/testify/require"
)

func TestReadCommandOutputRealtime(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error

	builder := commandBuilder{}

	var stdoutReader *bufio.Reader
	var stdin io.WriteCloser

	builder.Customize(func(c *exec.Cmd) (err error) {
		reader, err := c.StdoutPipe()
		if err != nil {
			return err
		}
		stdoutReader = bufio.NewReader(reader)
		stdin, err = c.StdinPipe()
		return err
	})
	cmd := builder.Cmd("go", "run", "./testdata/helper_process.go")
	_, err = cmd.Start()
	require.NoError(t, err)

	readingDone := make(chan bool)

	// first line
	var line, text string
	go func() {
		defer func() { readingDone <- true }()
		line, text, err = textio.ReadLineUntilText(stdoutReader, []string{"\n"})
	}()

	select {
	case <-ctx.Done():
		require.Fail(t, "timeout")
	case <-readingDone:
	}

	require.NoError(t, err)
	require.Equal(t, "Printing first line immiedatelly. Waiting for confirmation...", line)
	require.Equal(t, "", text)

	// second line
	_, err = stdin.Write([]byte("\n"))
	require.NoError(t, err)

	go func() {
		defer func() { readingDone <- true }()
		line, text, err = textio.ReadLineUntilText(stdoutReader, []string{"..."})
	}()

	select {
	case <-ctx.Done():
		require.Fail(t, "timeout")
	case <-readingDone:
	}

	require.NoError(t, err)
	require.Equal(t, "Printing half of the second line and waiting...", line)
	require.Equal(t, "...", text)

	var textBytes []byte
	go func() {
		defer func() { readingDone <- true }()
		// reading to the end to check if it's good enough method of waiting for process,
		// instead of cmd.Wait()
		textBytes, err = io.ReadAll(stdoutReader)
	}()

	// resuming writing only AFTER io.ReadAll(), to make sure it's the proper way of reading it to
	// the end
	_, err = stdin.Write([]byte("\n"))
	require.NoError(t, err)

	_, err = stdin.Write([]byte("\n"))
	require.NoError(t, err)

	select {
	case <-ctx.Done():
		require.Fail(t, "timeout")
	case <-readingDone:
	}

	require.NoError(t, err)
	require.Equal(t, " then printing the rest.\n", string(textBytes))
	// stdoutReader.Close()
}
