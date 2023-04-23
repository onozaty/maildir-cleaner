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

func TestArchive(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	allMails := setupMails(t, temp)

	// 奇数を移動対象に
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
	resultArchiveMails, err := Archive(temp, &targetMails, "Archived")

	// ASSERT
	require.NoError(t, err)

	// 対象のメールが元のフォルダに無い＆移動先にあること
	for i, mail := range targetMails {
		assert.NoFileExists(t, mail.FullPath)

		var archivedFolderName string
		if mail.FolderName == "" {
			archivedFolderName = "Archived"
		} else {
			archivedFolderName = "Archived" + "." + mail.FolderName
		}

		encodedFolderName, _ := folder.EncodeMailFolderName(archivedFolderName)
		archivedMailPath := filepath.Join(
			temp,
			"."+encodedFolderName,
			mail.SubDirName,
			mail.FileName)

		assert.FileExists(t, archivedMailPath)
		assert.NoFileExists(t, mail.FullPath)

		// 戻りの内容とあっていること
		resultArchiveMail := (*resultArchiveMails)[i]
		assert.Equal(t, archivedMailPath, resultArchiveMail.FullPath)
		assert.Equal(t, archivedFolderName, resultArchiveMail.FolderName)
		assert.Equal(t, mail.SubDirName, resultArchiveMail.SubDirName)
		assert.Equal(t, mail.FileName, resultArchiveMail.FileName)
		assert.Equal(t, mail.Size, resultArchiveMail.Size)
		assert.Equal(t, mail.Time, resultArchiveMail.Time)
	}

	// 対象外のメールが削除されていないこと
	for _, mail := range nonTargetMails {
		assert.FileExists(t, mail.FullPath)
	}
}

func TestArchive_ArchiveFolderBaseNameMultibyte(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	allMails := setupMails(t, temp)

	// 奇数を移動対象に
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
	resultArchiveMails, err := Archive(temp, &targetMails, "第1.第2") // マルチバイト

	// ASSERT
	require.NoError(t, err)

	// 対象のメールが元のフォルダに無い＆移動先にあること
	for i, mail := range targetMails {
		assert.NoFileExists(t, mail.FullPath)

		var archivedFolderName string
		if mail.FolderName == "" {
			archivedFolderName = "第1.第2"
		} else {
			archivedFolderName = "第1.第2" + "." + mail.FolderName
		}

		encodedFolderName, _ := folder.EncodeMailFolderName(archivedFolderName)
		archivedMailPath := filepath.Join(
			temp,
			"."+encodedFolderName,
			mail.SubDirName,
			mail.FileName)

		assert.FileExists(t, archivedMailPath)
		assert.NoFileExists(t, mail.FullPath)

		// 戻りの内容とあっていること
		resultArchiveMail := (*resultArchiveMails)[i]
		assert.Equal(t, archivedMailPath, resultArchiveMail.FullPath)
		assert.Equal(t, archivedFolderName, resultArchiveMail.FolderName)
		assert.Equal(t, mail.SubDirName, resultArchiveMail.SubDirName)
		assert.Equal(t, mail.FileName, resultArchiveMail.FileName)
		assert.Equal(t, mail.Size, resultArchiveMail.Size)
		assert.Equal(t, mail.Time, resultArchiveMail.Time)
	}

	// 対象外のメールが削除されていないこと
	for _, mail := range nonTargetMails {
		assert.FileExists(t, mail.FullPath)
	}
}

func TestArchive_MailNotFound(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	subscriptionsPath := filepath.Join(temp, "subscriptions")
	test.CreateFile(t, subscriptionsPath, "")

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
	_, err := Archive(temp, &targetMails, "Archived")

	// ASSERT
	// OSによってエラーメッセージが異なるのでファイル名部分だけチェック
	expect := "rename " + mailPath
	assert.Contains(t, err.Error(), expect)
}

func TestArchive_RootNotFound(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	rootMailFolderPath := filepath.Join(temp, "xxxx") // 存在しないフォルダ

	mailPath := filepath.Join(rootMailFolderPath, "cur", "a")
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
	_, err := Archive(rootMailFolderPath, &targetMails, "Archived")

	// ASSERT
	// OSによってエラーメッセージが異なるのでファイル名部分だけチェック
	expect := "mkdir " + rootMailFolderPath
	assert.Contains(t, err.Error(), expect)
}
