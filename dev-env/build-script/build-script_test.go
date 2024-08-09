package buildscript

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateBashInitFile(t *testing.T) {

	// err := startShell([]string{})

	resultFile, err := createBashInitFile([]string{})
	defer os.Remove(resultFile)

	require.NoError(t, err)
	assert.NotEmpty(t, resultFile)

	content, err := os.ReadFile(resultFile)
	require.NoError(t, err)
	assert.NotEmpty(t, content)
}
