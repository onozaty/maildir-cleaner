package folder

import (
	"path/filepath"
	"testing"

	"github.com/onozaty/maildir-cleaner/test"
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
	test.CreateFile(t, subscriptionsPath, "X\n")

	// ACT
	folderPath, err := Setup(temp, "AAA")

	// ASSERT
	require.NoError(t, err)

	expectedFolderPath := filepath.Join(temp, ".AAA")
	assert.Equal(t, expectedFolderPath, folderPath)
	assert.Equal(t, "X\nAAA\n", test.ReadFile(t, subscriptionsPath))

	assert.DirExists(t, expectedFolderPath)
	assert.DirExists(t, filepath.Join(expectedFolderPath, "cur"))
	assert.DirExists(t, filepath.Join(expectedFolderPath, "new"))
	assert.DirExists(t, filepath.Join(expectedFolderPath, "tmp"))
}

func TestSetup_AlreadyExists(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	subscriptionsPath := filepath.Join(temp, "subscriptions")
	test.CreateFile(t, subscriptionsPath, "&MEIwRDBG-\nAAA")

	expectedFolderPath := test.CreateDir(t, temp, ".&MEIwRDBG-")
	test.CreateDir(t, expectedFolderPath, "cur")
	test.CreateDir(t, expectedFolderPath, "new")
	test.CreateDir(t, expectedFolderPath, "tmp")

	// ACT
	folderPath, err := Setup(temp, "あいう")

	// ASSERT
	require.NoError(t, err)

	assert.Equal(t, expectedFolderPath, folderPath)
	assert.Equal(t, "&MEIwRDBG-\nAAA", test.ReadFile(t, subscriptionsPath))

	assert.DirExists(t, expectedFolderPath)
	assert.DirExists(t, filepath.Join(expectedFolderPath, "cur"))
	assert.DirExists(t, filepath.Join(expectedFolderPath, "new"))
	assert.DirExists(t, filepath.Join(expectedFolderPath, "tmp"))
}

func TestSetup_ParentFolders(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	subscriptionsPath := filepath.Join(temp, "subscriptions")
	test.CreateFile(t, subscriptionsPath, "X\n")

	// ACT
	folderPath, err := Setup(temp, "X.Y.Z.テスト")

	// ASSERT
	require.NoError(t, err)

	expectedFolderPath := filepath.Join(temp, ".X.Y.Z.&MMYwuTDI-")
	assert.Equal(t, expectedFolderPath, folderPath)
	assert.Equal(t, "X\nX.Y\nX.Y.Z\nX.Y.Z.&MMYwuTDI-\n", test.ReadFile(t, subscriptionsPath))

	assert.DirExists(t, expectedFolderPath)
	assert.DirExists(t, filepath.Join(expectedFolderPath, "cur"))
	assert.DirExists(t, filepath.Join(expectedFolderPath, "new"))
	assert.DirExists(t, filepath.Join(expectedFolderPath, "tmp"))

	// 親フォルダも生成されていることをチェック
	{
		expectedParentFolderPath := filepath.Join(temp, ".X")
		assert.DirExists(t, expectedParentFolderPath)
		assert.DirExists(t, filepath.Join(expectedParentFolderPath, "cur"))
		assert.DirExists(t, filepath.Join(expectedParentFolderPath, "new"))
		assert.DirExists(t, filepath.Join(expectedParentFolderPath, "tmp"))
	}
	{
		expectedParentFolderPath := filepath.Join(temp, ".X.Y")
		assert.DirExists(t, expectedParentFolderPath)
		assert.DirExists(t, filepath.Join(expectedParentFolderPath, "cur"))
		assert.DirExists(t, filepath.Join(expectedParentFolderPath, "new"))
		assert.DirExists(t, filepath.Join(expectedParentFolderPath, "tmp"))
	}
	{
		expectedParentFolderPath := filepath.Join(temp, ".X.Y.Z")
		assert.DirExists(t, expectedParentFolderPath)
		assert.DirExists(t, filepath.Join(expectedParentFolderPath, "cur"))
		assert.DirExists(t, filepath.Join(expectedParentFolderPath, "new"))
		assert.DirExists(t, filepath.Join(expectedParentFolderPath, "tmp"))
	}
}

func TestSetup_SubscriptionsEmpty(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	subscriptionsPath := filepath.Join(temp, "subscriptions")
	test.CreateFile(t, subscriptionsPath, "")

	// ACT
	folderPath, err := Setup(temp, "A.B")

	// ASSERT
	require.NoError(t, err)

	expectedFolderPath := filepath.Join(temp, ".A.B")
	assert.Equal(t, expectedFolderPath, folderPath)
	assert.Equal(t, "A\nA.B\n", test.ReadFile(t, subscriptionsPath))

	assert.DirExists(t, expectedFolderPath)
	assert.DirExists(t, filepath.Join(expectedFolderPath, "cur"))
	assert.DirExists(t, filepath.Join(expectedFolderPath, "new"))
	assert.DirExists(t, filepath.Join(expectedFolderPath, "tmp"))

	// 親フォルダも作成されていること
	expectedParentFolderPath := filepath.Join(temp, ".A")
	assert.DirExists(t, expectedParentFolderPath)
	assert.DirExists(t, filepath.Join(expectedParentFolderPath, "cur"))
	assert.DirExists(t, filepath.Join(expectedParentFolderPath, "new"))
	assert.DirExists(t, filepath.Join(expectedParentFolderPath, "tmp"))
}

func TestSetup_SubscriptionsNoLineBreak(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	subscriptionsPath := filepath.Join(temp, "subscriptions")
	test.CreateFile(t, subscriptionsPath, "AAA\nBBB")

	// ACT
	folderPath, err := Setup(temp, "AA")

	// ASSERT
	require.NoError(t, err)

	expectedFolderPath := filepath.Join(temp, ".AA")
	assert.Equal(t, expectedFolderPath, folderPath)
	assert.Equal(t, "AAA\nBBB\nAA\n", test.ReadFile(t, subscriptionsPath))

	assert.DirExists(t, expectedFolderPath)
	assert.DirExists(t, filepath.Join(expectedFolderPath, "cur"))
	assert.DirExists(t, filepath.Join(expectedFolderPath, "new"))
	assert.DirExists(t, filepath.Join(expectedFolderPath, "tmp"))
}

func TestSetup_SubscriptionsNotFound(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	// subscriptions無し

	// ACT
	_, err := Setup(temp, "AA")

	// ASSERT
	require.EqualError(t, err, "subscriptions file not found: currently only dovecot is supported")
}

func TestSetup_RootDirNotFound(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	rootMailFolderPath := filepath.Join(temp, "xxxx") // 存在しないフォルダ

	// ACT
	_, err := Setup(rootMailFolderPath, "AA")

	// ASSERT
	require.Error(t, err)
	// OSによってエラーメッセージが異なるのでファイル名部分だけチェック
	expect := "mkdir " + filepath.Join(rootMailFolderPath, ".AA")
	assert.Contains(t, err.Error(), expect)
}
