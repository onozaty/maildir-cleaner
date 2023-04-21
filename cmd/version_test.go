package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersionCmd(t *testing.T) {

	// ARRANGE
	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"version",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	require.NoError(t, err)

	result := buf.String()
	assert.Contains(t, result, "Version: dev\nRevision: dev\nOS: ")
}
