package mail

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeMailFolderName(t *testing.T) {

	// ARRANGE
	encodedName := "&MEIwRDBG-"

	// ACT
	name, err := DecodeMailFolderName(encodedName)

	// ASSERT
	require.NoError(t, err)
	assert.Equal(t, "あいう", name)
}

func TestDecodeMailFolderName_Alphabet(t *testing.T) {

	// ARRANGE
	encodedName := "abc"

	// ACT
	name, err := DecodeMailFolderName(encodedName)

	// ASSERT
	require.NoError(t, err)
	assert.Equal(t, "abc", name)
}

func TestDecodeMailFolderName_Invalid(t *testing.T) {

	// ARRANGE
	encodedName := "&A"

	// ACT
	_, err := DecodeMailFolderName(encodedName)

	// ASSERT
	assert.EqualError(t, err, "&A is invalid folder name: utf7: invalid UTF-7")
}

func TestEncodeMailFolderName(t *testing.T) {

	// ARRANGE
	encodedName := "あいう"

	// ACT
	name, err := EncodeMailFolderName(encodedName)

	// ASSERT
	require.NoError(t, err)
	assert.Equal(t, "&MEIwRDBG-", name)
}

func TestEncodeMailFolderName_Alphabet(t *testing.T) {

	// ARRANGE
	encodedName := "abc"

	// ACT
	name, err := EncodeMailFolderName(encodedName)

	// ASSERT
	require.NoError(t, err)
	assert.Equal(t, "abc", name)
}
