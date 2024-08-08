package fs

import (
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileExists(t *testing.T) {
	// given
	tempFile, err := os.CreateTemp("", "")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	err = os.WriteFile(tempFile.Name(), []byte("file content"), 0600)
	require.NoError(t, err)

	err = tempFile.Close()
	require.NoError(t, err)
	existingFile := tempFile.Name()

	// when
	content, isExist, err := ReadFileIfExists(existingFile)

	// then
	assert.NoError(t, err)
	assert.True(t, isExist)
	assert.Equal(t, "file content", string(content))
}

func TestFileNotExists(t *testing.T) {
	// given
	tempFile, err := os.CreateTemp("", "")
	require.NoError(t, err)

	err = tempFile.Close()
	require.NoError(t, err)

	notExistingFile := tempFile.Name()
	err = os.Remove(notExistingFile)
	require.NoError(t, err)

	// when
	_, isExist, err := ReadFileIfExists(notExistingFile)

	// then
	assert.NoError(t, err)
	assert.False(t, isExist)
}

func TestFileNotPermitted(t *testing.T) {
	// given
	tempFile, err := os.CreateTemp("", "")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	onlyWrite := fs.FileMode(0200)
	err = tempFile.Chmod(onlyWrite)
	require.NoError(t, err)

	err = tempFile.Close()
	require.NoError(t, err)

	notPermittedFile := tempFile.Name()

	// when
	_, _, err = ReadFileIfExists(notPermittedFile)

	// then
	assert.ErrorIs(t, err, os.ErrPermission)
}
