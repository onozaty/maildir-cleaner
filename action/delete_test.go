package action

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/onozaty/maildir-cleaner/collector"
	"github.com/onozaty/maildir-cleaner/folder"
	"github.com/onozaty/maildir-cleaner/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDelete(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	allMails := setupMails(t, temp)

	// 奇数を削除対象
	targetMails := []collector.Mail{}
	nonTargetMails := []collector.Mail{}

	for i, mail := range allMails {
		if i%2 == 0 {
			targetMails = append(targetMails, mail)
		} else {
			nonTargetMails = append(nonTargetMails, mail)
		}
	}

	// ACT
	err := Delete(temp, &targetMails)

	// ASSERT
	require.NoError(t, err)

	// 対象のメールが削除されていること
	for _, mail := range targetMails {
		assert.NoFileExists(t, mail.FullPath)
	}

	// 対象外のメールが削除されていないこと
	for _, mail := range nonTargetMails {
		assert.FileExists(t, mail.FullPath)
	}
}

func TestDelete_NotFound(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	// 存在しないメール
	mailPath := filepath.Join(temp, "cur", "a")
	targetMails := []collector.Mail{
		{
			FullPath:   mailPath,
			FolderName: "",
			SubDirName: "cur",
			FileName:   "a",
			Size:       1,
			Time:       time.Now(),
		},
	}

	// ACT
	err := Delete(temp, &targetMails)

	// ASSERT
	// OSによってエラーメッセージが異なるのでファイル名部分だけチェック
	expect := "remove " + mailPath
	assert.Contains(t, err.Error(), expect)
}

func setupMails(t *testing.T, root string) []collector.Mail {

	subscriptionsPath := filepath.Join(root, "subscriptions")
	test.CreateFile(t, subscriptionsPath, "A\nA.B\n&MMYwuTDI-1\n")

	return []collector.Mail{
		// INBOX
		createMailByName(t, root, "", "new", "a"),
		createMailByName(t, root, "", "new", "b"),
		createMailByName(t, root, "", "new", "c"),
		createMailByName(t, root, "", "cur", "d"),
		createMailByName(t, root, "", "cur", "e"),
		createMailByName(t, root, "", "cur", "f"),
		// A
		createMailByName(t, root, "A", "new", "g"),
		createMailByName(t, root, "A", "new", "h"),
		createMailByName(t, root, "A", "new", "i"),
		// A.B
		createMailByName(t, root, "A.B", "cur", "j"),
		createMailByName(t, root, "A.B", "cur", "k"),
		createMailByName(t, root, "A.B", "cur", "l"),
		// テスト1
		createMailByName(t, root, "テスト1", "new", "m"),
		createMailByName(t, root, "テスト1", "cur", "n"),
		createMailByName(t, root, "テスト1", "new", "o"),
	}
}

func createMailByName(t *testing.T, rootDir string, folderName string, sub string, name string) collector.Mail {

	encodedFolderName, _ := folder.EncodeMailFolderName(folderName)
	var physicalFolderName string
	if encodedFolderName == "" {
		physicalFolderName = encodedFolderName
	} else {
		physicalFolderName = "." + encodedFolderName
	}
	folderDir := test.CreateMailFolder(t, rootDir, physicalFolderName)

	size := 1 // サイズは固定

	mailPath, fileName := test.CreateMailByName(t, folderDir, sub, name, size)

	return collector.Mail{
		FullPath:   mailPath,
		FolderName: folderName,
		SubDirName: sub,
		FileName:   fileName,
		Size:       int64(size),
		Time:       time.Now(),
	}
}

func createMailByYearMonth(t *testing.T, rootDir string, folderName string, sub string, year int, month time.Month) collector.Mail {

	encodedFolderName, _ := folder.EncodeMailFolderName(folderName)
	var physicalFolderName string
	if encodedFolderName == "" {
		physicalFolderName = encodedFolderName
	} else {
		physicalFolderName = "." + encodedFolderName
	}
	folderDir := test.CreateMailFolder(t, rootDir, physicalFolderName)

	time := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	size := 1 // サイズは固定

	mailPath, fileName := test.CreateMailByTime(t, folderDir, sub, time, size)

	return collector.Mail{
		FullPath:   mailPath,
		FolderName: folderName,
		SubDirName: sub,
		FileName:   fileName,
		Size:       int64(size),
		Time:       time,
	}
}
