package cmd

import (
	"bytes"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/onozaty/maildir-cleaner/collector"
	"github.com/onozaty/maildir-cleaner/folder"
	"github.com/onozaty/maildir-cleaner/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArchiveCmd(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	// 削除対象(10日以上経過)
	targetMails := []collector.Mail{
		// INBOX
		createMail(t, temp, "", "new", 10),
		createMail(t, temp, "", "new", 11),
		createMail(t, temp, "", "cur", 12),
		createMail(t, temp, "", "cur", 13),
		// A
		createMail(t, temp, "A", "new", 1000),
		// A.B
		createMail(t, temp, "A.B", "cur", 10),
		createMail(t, temp, "A.B", "new", 10),
		// テスト1
		createMail(t, temp, "テスト1", "cur", 11),
		createMail(t, temp, "テスト1", "cur", 12),
	}

	// 削除対象外(10日未満 or tmp)
	nonTargetMails := []collector.Mail{
		// INBOX
		createMail(t, temp, "", "new", 1),
		createMail(t, temp, "", "new", 9),
		createMail(t, temp, "", "cur", 9),
		createMail(t, temp, "", "tmp", 10),
		// A
		createMail(t, temp, "A", "new", 9),
		// A.B
		createMail(t, temp, "A.B", "cur", 1),
		// テスト1
		createMail(t, temp, "テスト1", "cur", 9),
	}

	subscriptionsPath := filepath.Join(temp, "subscriptions")
	test.CreateFile(t, subscriptionsPath, "A\nA.B\n&MMYwuTDI-1\n")

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"archive",
		"-d", temp,
		"-a", "10",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	require.NoError(t, err)

	// 対象のメールが元のフォルダに無い＆移動先にあること
	for _, mail := range targetMails {
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
	}

	// 対象外のメールが削除されていないこと
	for _, mail := range nonTargetMails {
		assert.FileExists(t, mail.FullPath)
	}

	// 標準出力の内容確認
	result := buf.String()
	expected := fmt.Sprintf(`Starts searching for the target mails. maildir: %s age: %d
Completed search. The target mails are listed below.
+---------+-----------------+------------------+
| Name    | Number of mails | Total size(byte) |
+---------+-----------------+------------------+
|         |               4 |               46 |
| A       |               1 |            1,000 |
| A.B     |               2 |               20 |
| テスト1 |               2 |               23 |
+---------+-----------------+------------------+
|   Total |               9 |            1,089 |
+---------+-----------------+------------------+
Starts archiving mails.
Completed archive. The archived mails are listed below.
+------------------+-----------------+------------------+
| Name             | Number of mails | Total size(byte) |
+------------------+-----------------+------------------+
| Archived         |               4 |               46 |
| Archived.A       |               1 |            1,000 |
| Archived.A.B     |               2 |               20 |
| Archived.テスト1 |               2 |               23 |
+------------------+-----------------+------------------+
|            Total |               9 |            1,089 |
+------------------+-----------------+------------------+
`, temp, 10)
	assert.Equal(t, expected, result)
}

func TestArchiveCmd_MaildirNotFound(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	rootMailFolderPath := filepath.Join(temp, "xx") // 存在しないフォルダ

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"archive",
		"-d", rootMailFolderPath,
		"-a", "10",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	require.Error(t, err)
	// OSによってエラーメッセージが異なるのでファイル名部分だけチェック
	expect := "open " + rootMailFolderPath
	assert.Contains(t, err.Error(), expect)
}
