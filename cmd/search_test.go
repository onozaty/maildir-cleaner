package cmd

import (
	"bytes"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/onozaty/maildir-cleaner/collector"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchCmd(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	// 対象(10日以上経過)
	targetMails := []collector.Mail{
		// INBOX
		createMailByDays(t, temp, "", "new", 10),
		createMailByDays(t, temp, "", "new", 11),
		createMailByDays(t, temp, "", "cur", 12),
		createMailByDays(t, temp, "", "cur", 13),
		// A
		createMailByDays(t, temp, "A", "new", 1000),
		// A.B
		createMailByDays(t, temp, "A.B", "cur", 10),
		createMailByDays(t, temp, "A.B", "new", 10),
		// テスト1
		createMailByDays(t, temp, "テスト1", "cur", 11),
		createMailByDays(t, temp, "テスト1", "cur", 12),
	}

	// 対象外(10日未満 or tmp)
	nonTargetMails := []collector.Mail{
		// INBOX
		createMailByDays(t, temp, "", "new", 1),
		createMailByDays(t, temp, "", "new", 9),
		createMailByDays(t, temp, "", "cur", 9),
		createMailByDays(t, temp, "", "tmp", 10),
		// A
		createMailByDays(t, temp, "A", "new", 9),
		// A.B
		createMailByDays(t, temp, "A.B", "cur", 1),
		// テスト1
		createMailByDays(t, temp, "テスト1", "cur", 9),
	}

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"search",
		"-d", temp,
		"-a", "10",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	require.NoError(t, err)

	// 対象/対象外のメールが削除されていないこと
	for _, mail := range targetMails {
		assert.FileExists(t, mail.FullPath)
	}
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
`, temp, 10)
	assert.Equal(t, expected, result)
}

func TestSearchCmd_ExcludeFolder(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	// 対象
	targetMails := []collector.Mail{
		// INBOX
		createMailByDays(t, temp, "", "new", 1),
		createMailByDays(t, temp, "", "new", 2),
		createMailByDays(t, temp, "", "cur", 3),
		createMailByDays(t, temp, "", "cur", 4),
		// A
		createMailByDays(t, temp, "A", "new", 5),
		// A.B
		createMailByDays(t, temp, "A.B", "cur", 6),
		createMailByDays(t, temp, "A.B", "new", 7),
		// テスト1
		createMailByDays(t, temp, "テスト1", "cur", 8),
		createMailByDays(t, temp, "テスト1", "cur", 9),
	}

	// 対象外(除外フォルダとして指定)
	nonTargetMails := []collector.Mail{
		// B
		createMailByDays(t, temp, "B", "new", 11),
		// B.A
		createMailByDays(t, temp, "B.A", "cur", 12),
		// テスト2
		createMailByDays(t, temp, "テスト2", "cur", 13),
	}

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"search",
		"-d", temp,
		"-a", "1",
		"--exclude-folder", "B",
		"--exclude-folder", "テスト2",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	require.NoError(t, err)

	// 対象/対象外のメールが削除されていないこと
	for _, mail := range targetMails {
		assert.FileExists(t, mail.FullPath)
	}
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
|         |               4 |               10 |
| A       |               1 |                5 |
| A.B     |               2 |               13 |
| テスト1 |               2 |               17 |
+---------+-----------------+------------------+
|   Total |               9 |               45 |
+---------+-----------------+------------------+
`, temp, 1)
	assert.Equal(t, expected, result)
}

func TestSearchCmd_Empty(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	// 対象外のみ
	nonTargetMails := []collector.Mail{
		// INBOX
		createMailByDays(t, temp, "", "new", 1),
		createMailByDays(t, temp, "", "new", 9),
		createMailByDays(t, temp, "", "cur", 10),
		createMailByDays(t, temp, "", "tmp", 11),
		// A
		createMailByDays(t, temp, "A", "new", 9),
		// A.B
		createMailByDays(t, temp, "A.B", "cur", 1),
		// テスト1
		createMailByDays(t, temp, "テスト1", "cur", 9),
	}

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"search",
		"-d", temp,
		"-a", "11",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	require.NoError(t, err)

	// 対象外のメールが削除されていないこと
	for _, mail := range nonTargetMails {
		assert.FileExists(t, mail.FullPath)
	}

	// 標準出力の内容確認
	result := buf.String()
	expected := fmt.Sprintf(`Starts searching for the target mails. maildir: %s age: %d
Completed search. There were no target mails.
`, temp, 11)
	assert.Equal(t, expected, result)
}

func TestSearch_MaildirNotFound(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	rootMailFolderPath := filepath.Join(temp, "xx") // 存在しないフォルダ

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"search",
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
