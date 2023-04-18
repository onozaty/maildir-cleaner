package folder

import (
	"io/fs"
	"os"
	"path/filepath"
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

func TestSetup(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	subscriptionsPath := filepath.Join(temp, "subscriptions")
	createFile(t, subscriptionsPath, "X\n")

	// ACT
	folderPath, err := Setup(temp, "AAA")

	// ASSERT
	require.NoError(t, err)

	expectedFolderPath := filepath.Join(temp, ".AAA")
	assert.Equal(t, expectedFolderPath, folderPath)
	assert.Equal(t, "X\nAAA\n", readFile(t, subscriptionsPath))

	assert.DirExists(t, expectedFolderPath)
	assert.DirExists(t, filepath.Join(expectedFolderPath, "cur"))
	assert.DirExists(t, filepath.Join(expectedFolderPath, "new"))
	assert.DirExists(t, filepath.Join(expectedFolderPath, "tmp"))
}

func TestSetup_AlreadyExists(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	subscriptionsPath := filepath.Join(temp, "subscriptions")
	createFile(t, subscriptionsPath, "&MEIwRDBG-\nAAA")

	expectedFolderPath := createDir(t, temp, ".&MEIwRDBG-")
	createDir(t, expectedFolderPath, "cur")
	createDir(t, expectedFolderPath, "new")
	createDir(t, expectedFolderPath, "tmp")

	// ACT
	folderPath, err := Setup(temp, "あいう")

	// ASSERT
	require.NoError(t, err)

	assert.Equal(t, expectedFolderPath, folderPath)
	assert.Equal(t, "&MEIwRDBG-\nAAA", readFile(t, subscriptionsPath))

	assert.DirExists(t, expectedFolderPath)
	assert.DirExists(t, filepath.Join(expectedFolderPath, "cur"))
	assert.DirExists(t, filepath.Join(expectedFolderPath, "new"))
	assert.DirExists(t, filepath.Join(expectedFolderPath, "tmp"))
}

func createDir(t *testing.T, parent string, name string) string {

	dir := filepath.Join(parent, name)
	err := os.Mkdir(dir, 0777)
	require.NoError(t, err)

	return dir
}

func createFile(t *testing.T, path string, content string) fs.FileInfo {

	file, err := os.Create(path)
	require.NoError(t, err)

	_, err = file.Write([]byte(content))
	require.NoError(t, err)

	info, err := file.Stat()
	require.NoError(t, err)

	err = file.Close()
	require.NoError(t, err)

	return info
}

func readFile(t *testing.T, path string) string {

	bo, err := os.ReadFile(path)
	require.NoError(t, err)

	return string(bo)
}
