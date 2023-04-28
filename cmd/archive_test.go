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

	// アーカイブ対象(10日以上経過)
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

	// アーカイブ対象外(10日未満 or tmp)
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

	// subscriptionsに登録されていること
	assert.Equal(t, "A\nA.B\n&MMYwuTDI-1\nArchived\nArchived.A\nArchived.A.B\nArchived.&MMYwuTDI-1\n", test.ReadFile(t, subscriptionsPath))

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

func TestArchiveCmd_ArchivePatternYear(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	// アーカイブ対象(10日以上経過)
	targetMails := []collector.Mail{
		// INBOX
		createMailByYearMonth(t, temp, "", "new", 2020, 1),
		createMailByYearMonth(t, temp, "", "cur", 2020, 12),
		// A
		createMailByYearMonth(t, temp, "A", "new", 2021, 1),
		createMailByYearMonth(t, temp, "A", "cur", 2021, 1),
		// A.B
		createMailByYearMonth(t, temp, "A.B", "cur", 2020, 1),
		createMailByYearMonth(t, temp, "A.B", "cur", 2020, 2),
		// テスト1
		createMailByYearMonth(t, temp, "テスト1", "cur", 2020, 11),
	}
	archivedFolderNames := []string{
		"Archived.2020",
		"Archived.2020",
		"Archived.2021",
		"Archived.2021",
		"Archived.2020",
		"Archived.2020",
		"Archived.2020",
	}

	// アーカイブ対象外
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

	subscriptionsPath := filepath.Join(temp, "subscriptions")
	test.CreateFile(t, subscriptionsPath, "A\nA.B\n&MMYwuTDI-1\n")

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"archive",
		"-d", temp,
		"-a", "10",
		"--archive-pattern", "year",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	require.NoError(t, err)

	// 対象のメールが元のフォルダに無い＆移動先にあること
	for i, mail := range targetMails {
		assert.NoFileExists(t, mail.FullPath)

		archivedFolderName := archivedFolderNames[i]

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

	// subscriptionsに登録されていること
	assert.Equal(t, "A\nA.B\n&MMYwuTDI-1\nArchived.2020\nArchived.2021\n", test.ReadFile(t, subscriptionsPath))

	// 標準出力の内容確認
	result := buf.String()
	expected := fmt.Sprintf(`Starts searching for the target mails. maildir: %s age: %d
Completed search. The target mails are listed below.
+---------+-----------------+------------------+
| Name    | Number of mails | Total size(byte) |
+---------+-----------------+------------------+
|         |               2 |            4,053 |
| A       |               2 |            4,044 |
| A.B     |               2 |            4,043 |
| テスト1 |               1 |            2,031 |
+---------+-----------------+------------------+
|   Total |               7 |           14,171 |
+---------+-----------------+------------------+
Starts archiving mails.
Completed archive. The archived mails are listed below.
+---------------+-----------------+------------------+
| Name          | Number of mails | Total size(byte) |
+---------------+-----------------+------------------+
| Archived.2020 |               5 |           10,127 |
| Archived.2021 |               2 |            4,044 |
+---------------+-----------------+------------------+
|         Total |               7 |           14,171 |
+---------------+-----------------+------------------+
`, temp, 10)
	assert.Equal(t, expected, result)
}

func TestArchiveCmd_ArchivePatternMonth(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	// アーカイブ対象(10日以上経過)
	targetMails := []collector.Mail{
		// INBOX
		createMailByYearMonth(t, temp, "", "new", 2020, 1),
		createMailByYearMonth(t, temp, "", "cur", 2020, 12),
		// A
		createMailByYearMonth(t, temp, "A", "new", 2021, 1),
		createMailByYearMonth(t, temp, "A", "cur", 2021, 1),
		// A.B
		createMailByYearMonth(t, temp, "A.B", "cur", 2020, 1),
		createMailByYearMonth(t, temp, "A.B", "cur", 2020, 2),
		// テスト1
		createMailByYearMonth(t, temp, "テスト1", "cur", 2020, 11),
	}
	archivedFolderNames := []string{
		"Archived.2020.01",
		"Archived.2020.12",
		"Archived.2021.01",
		"Archived.2021.01",
		"Archived.2020.01",
		"Archived.2020.02",
		"Archived.2020.11",
	}

	// アーカイブ対象外
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

	subscriptionsPath := filepath.Join(temp, "subscriptions")
	test.CreateFile(t, subscriptionsPath, "A\nA.B\n&MMYwuTDI-1\nArchived.2020.01\n")

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"archive",
		"-d", temp,
		"-a", "10",
		"--archive-pattern", "month",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	require.NoError(t, err)

	// 対象のメールが元のフォルダに無い＆移動先にあること
	for i, mail := range targetMails {
		assert.NoFileExists(t, mail.FullPath)

		archivedFolderName := archivedFolderNames[i]

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

	// subscriptionsに登録されていること
	assert.Equal(t, "A\nA.B\n&MMYwuTDI-1\nArchived.2020.01\nArchived.2020.12\nArchived.2021.01\nArchived.2020.02\nArchived.2020.11\n", test.ReadFile(t, subscriptionsPath))

	// 標準出力の内容確認
	result := buf.String()
	expected := fmt.Sprintf(`Starts searching for the target mails. maildir: %s age: %d
Completed search. The target mails are listed below.
+---------+-----------------+------------------+
| Name    | Number of mails | Total size(byte) |
+---------+-----------------+------------------+
|         |               2 |            4,053 |
| A       |               2 |            4,044 |
| A.B     |               2 |            4,043 |
| テスト1 |               1 |            2,031 |
+---------+-----------------+------------------+
|   Total |               7 |           14,171 |
+---------+-----------------+------------------+
Starts archiving mails.
Completed archive. The archived mails are listed below.
+------------------+-----------------+------------------+
| Name             | Number of mails | Total size(byte) |
+------------------+-----------------+------------------+
| Archived.2020.01 |               2 |            4,042 |
| Archived.2020.02 |               1 |            2,022 |
| Archived.2020.11 |               1 |            2,031 |
| Archived.2020.12 |               1 |            2,032 |
| Archived.2021.01 |               2 |            4,044 |
+------------------+-----------------+------------------+
|            Total |               7 |           14,171 |
+------------------+-----------------+------------------+
`, temp, 10)
	assert.Equal(t, expected, result)
}

func TestArchiveCmd_IgnoreArchiveFolder(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	// アーカイブ対象(1日以上経過)
	targetMails := []collector.Mail{
		// INBOX
		createMailByDays(t, temp, "", "new", 1),
		createMailByDays(t, temp, "", "cur", 2),
		// A
		createMailByDays(t, temp, "A", "new", 3),
		// A.B
		createMailByDays(t, temp, "A.B", "cur", 4),
		createMailByDays(t, temp, "A.B", "new", 5),
	}

	// アーカイブ対象外(経過しているが、既にアーカイブフォルダ配下にあるもの)
	nonTargetMails := []collector.Mail{
		createMailByDays(t, temp, "Archived", "new", 11),
		createMailByDays(t, temp, "Archived.A", "cur", 12),
		createMailByDays(t, temp, "Archived.A.B", "new", 13),
	}

	subscriptionsPath := filepath.Join(temp, "subscriptions")
	test.CreateFile(t, subscriptionsPath, "A\nA.B\nArchived\nArchived.A\nArchived.A.B\n")

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"archive",
		"-d", temp,
		"-a", "1",
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

	// subscriptionsに登録されていること
	assert.Equal(t, "A\nA.B\nArchived\nArchived.A\nArchived.A.B\n", test.ReadFile(t, subscriptionsPath))

	// 標準出力の内容確認
	result := buf.String()
	expected := fmt.Sprintf(`Starts searching for the target mails. maildir: %s age: %d
Completed search. The target mails are listed below.
+-------+-----------------+------------------+
| Name  | Number of mails | Total size(byte) |
+-------+-----------------+------------------+
|       |               2 |                3 |
| A     |               1 |                3 |
| A.B   |               2 |                9 |
+-------+-----------------+------------------+
| Total |               5 |               15 |
+-------+-----------------+------------------+
Starts archiving mails.
Completed archive. The archived mails are listed below.
+--------------+-----------------+------------------+
| Name         | Number of mails | Total size(byte) |
+--------------+-----------------+------------------+
| Archived     |               2 |                3 |
| Archived.A   |               1 |                3 |
| Archived.A.B |               2 |                9 |
+--------------+-----------------+------------------+
|        Total |               5 |               15 |
+--------------+-----------------+------------------+
`, temp, 1)
	assert.Equal(t, expected, result)
}

func TestArchiveCmd_ArchiveFileNameMultibyte(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	// アーカイブ対象(5日以上)
	targetMails := []collector.Mail{
		// INBOX
		createMailByDays(t, temp, "", "new", 100),
		// A
		createMailByDays(t, temp, "A", "new", 7),
		// A.B
		createMailByDays(t, temp, "A.B", "cur", 6),
		// テスト1
		createMailByDays(t, temp, "テスト1", "cur", 5),
	}

	// アーカイブ対象外(5日未満 or アーカイブフォルダ配下)
	nonTargetMails := []collector.Mail{
		// INBOX
		createMailByDays(t, temp, "", "new", 4),
		// テスト1
		createMailByDays(t, temp, "テスト1", "cur", 3),
		// アーカイブ
		createMailByDays(t, temp, "アーカイブ", "cur", 100),
		// アーカイブ.A
		createMailByDays(t, temp, "アーカイブ.A", "cur", 1000),
	}

	subscriptionsPath := filepath.Join(temp, "subscriptions")
	test.CreateFile(t, subscriptionsPath, "A\nA.B\n&MMYwuTDI-1\n&MKIw,DCrMKQw1g-\n&MKIw,DCrMKQw1g-.A\n")

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"archive",
		"-d", temp,
		"--archive-folder", "アーカイブ",
		"-a", "5",
		"--archive-pattern", "keep",
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
			archivedFolderName = "アーカイブ"
		} else {
			archivedFolderName = "アーカイブ" + "." + mail.FolderName
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

	// subscriptionsに登録されていること
	assert.Equal(t, "A\nA.B\n&MMYwuTDI-1\n&MKIw,DCrMKQw1g-\n&MKIw,DCrMKQw1g-.A\n&MKIw,DCrMKQw1g-.A.B\n&MKIw,DCrMKQw1g-.&MMYwuTDI-1\n", test.ReadFile(t, subscriptionsPath))

	// 標準出力の内容確認
	result := buf.String()
	expected := fmt.Sprintf(`Starts searching for the target mails. maildir: %s age: %d
Completed search. The target mails are listed below.
+---------+-----------------+------------------+
| Name    | Number of mails | Total size(byte) |
+---------+-----------------+------------------+
|         |               1 |              100 |
| A       |               1 |                7 |
| A.B     |               1 |                6 |
| テスト1 |               1 |                5 |
+---------+-----------------+------------------+
|   Total |               4 |              118 |
+---------+-----------------+------------------+
Starts archiving mails.
Completed archive. The archived mails are listed below.
+--------------------+-----------------+------------------+
| Name               | Number of mails | Total size(byte) |
+--------------------+-----------------+------------------+
| アーカイブ         |               1 |              100 |
| アーカイブ.A       |               1 |                7 |
| アーカイブ.A.B     |               1 |                6 |
| アーカイブ.テスト1 |               1 |                5 |
+--------------------+-----------------+------------------+
|              Total |               4 |              118 |
+--------------------+-----------------+------------------+
`, temp, 5)
	assert.Equal(t, expected, result)
}

func TestArchiveCmd_Empty(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	// アーカイブ対象外のみ
	nonTargetMails := []collector.Mail{
		// INBOX
		createMailByDays(t, temp, "", "new", 10),
		createMailByDays(t, temp, "", "cur", 9),
		createMailByDays(t, temp, "", "tmp", 11),
		// A
		createMailByDays(t, temp, "A", "new", 10),
		// A.B
		createMailByDays(t, temp, "A.B", "cur", 1),
		// テスト1
		createMailByDays(t, temp, "テスト1", "cur", 9),
	}

	subscriptionsPath := filepath.Join(temp, "subscriptions")
	test.CreateFile(t, subscriptionsPath, "A\nA.B\n&MMYwuTDI-1\n")

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"archive",
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

	// subscriptionsが変わっていないこと
	assert.Equal(t, "A\nA.B\n&MMYwuTDI-1\n", test.ReadFile(t, subscriptionsPath))

	// 標準出力の内容確認
	result := buf.String()
	expected := fmt.Sprintf(`Starts searching for the target mails. maildir: %s age: %d
Completed search. There were no target mails.
`, temp, 11)
	assert.Equal(t, expected, result)
}

func TestArchiveCmd_SubscriptionsNotFound(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	// アーカイブ対象
	createMailByDays(t, temp, "", "new", 100)

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
	require.EqualError(t, err, "subscriptions file not found: currently only dovecot is supported")
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

func TestArchiveCmd_InvalidArchivePattern(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"archive",
		"-d", temp,
		"-a", "10",
		"--archive-pattern", "xxx",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	// ACT
	err := rootCmd.Execute()

	// ASSERT
	require.EqualError(t, err, "invalid archive-pattern 'xxx'")
}
