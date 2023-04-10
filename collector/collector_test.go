package collector

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollector(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	expected := []Mail{}

	// INBOX
	{
		mailFolder := createMailFolder(t, temp, "")
		createMailByTime(t, mailFolder, "new", agoDays(0), 1)
		createMailByTime(t, mailFolder, "new", agoDays(1), 1)
		createMailByTime(t, mailFolder, "cur", agoDays(2), 1)
		{
			// 収集対象
			time := agoDays(3)
			mailPath := createMailByTime(t, mailFolder, "cur", time, 1)
			expected = append(expected, Mail{Path: mailPath, FolderName: "", Size: 1, Time: time})
		}
		createMailByTime(t, mailFolder, "tmp", agoDays(4), 1)
		createMailByTime(t, mailFolder, "tmp", agoDays(5), 1)
	}

	// その他フォルダ
	{
		mailFolder := createMailFolder(t, temp, ".A")
		{
			// 収集対象
			time := agoDays(5)
			mailPath := createMailByTime(t, mailFolder, "new", time, 2)
			expected = append(expected, Mail{Path: mailPath, FolderName: "A", Size: 2, Time: time})
		}
		{
			// 収集対象
			time := agoDays(4)
			mailPath := createMailByTime(t, mailFolder, "new", time, 2)
			expected = append(expected, Mail{Path: mailPath, FolderName: "A", Size: 2, Time: time})
		}
		{
			// 収集対象
			time := agoDays(3)
			mailPath := createMailByTime(t, mailFolder, "cur", time, 2)
			expected = append(expected, Mail{Path: mailPath, FolderName: "A", Size: 2, Time: time})
		}
		createMailByTime(t, mailFolder, "cur", agoDays(2), 2)
		createMailByTime(t, mailFolder, "tmp", agoDays(1), 2)
		createMailByTime(t, mailFolder, "tmp", agoDays(0), 2)
	}
	{
		mailFolder := createMailFolder(t, temp, ".B")
		{
			// 収集対象
			time := agoDays(5)
			mailPath := createMailByTime(t, mailFolder, "new", time, 3)
			expected = append(expected, Mail{Path: mailPath, FolderName: "B", Size: 3, Time: time})
		}
	}
	{
		mailFolder := createMailFolder(t, temp, ".C")
		{
			// 収集対象
			time := agoDays(5)
			mailPath := createMailByTime(t, mailFolder, "cur", time, 4)
			expected = append(expected, Mail{Path: mailPath, FolderName: "C", Size: 4, Time: time})
		}
	}
	{
		mailFolder := createMailFolder(t, temp, ".D")
		createMailByTime(t, mailFolder, "tmp", agoDays(5), 5)
	}
	{
		// メールフォルダ以外のフォルダ(先頭に"."無し)
		mailFolder := createMailFolder(t, temp, "a")
		createMailByTime(t, mailFolder, "new", agoDays(5), 6)
		createMailByTime(t, mailFolder, "cur", agoDays(5), 6)
		createMailByTime(t, mailFolder, "tmp", agoDays(5), 6)
	}

	// ACT
	collector := NewCollector(3, "")
	mails, err := collector.Collect(temp)

	// ASSERT
	require.NoError(t, err)
	assert.Equal(t, &expected, mails)
}

func TestCollector_FolderName(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	expected := []Mail{}

	// INBOX
	{
		mailFolder := createMailFolder(t, temp, "")
		createMailByTime(t, mailFolder, "new", agoDays(1), 1)
		{
			// 収集対象
			time := agoDays(10)
			mailPath := createMailByTime(t, mailFolder, "new", time, 1)
			expected = append(expected, Mail{Path: mailPath, FolderName: "", Size: 1, Time: time})
		}
	}

	// a (対象外フォルダ：除外対象のフォルダ)
	{
		mailFolder := createMailFolder(t, temp, ".a")
		createMailByTime(t, mailFolder, "new", agoDays(10), 1)
		createMailByTime(t, mailFolder, "cur", agoDays(10), 1)
	}
	// aa (対象フォルダ：対象外フォルダと前方一致)
	{
		mailFolder := createMailFolder(t, temp, ".aa")
		{
			// 収集対象
			time := agoDays(10)
			mailPath := createMailByTime(t, mailFolder, "new", time, 1)
			expected = append(expected, Mail{Path: mailPath, FolderName: "aa", Size: 1, Time: time})
		}
	}
	// ab (対象外フォルダ：対象外フォルダのサブフォルダ)
	{
		mailFolder := createMailFolder(t, temp, ".a.b")
		createMailByTime(t, mailFolder, "new", agoDays(10), 1)
		createMailByTime(t, mailFolder, "cur", agoDays(10), 1)
	}
	// b (対象フォルダ)
	{
		mailFolder := createMailFolder(t, temp, ".b")
		{
			// 収集対象
			time := agoDays(11)
			mailPath := createMailByTime(t, mailFolder, "cur", time, 1)
			expected = append(expected, Mail{Path: mailPath, FolderName: "b", Size: 1, Time: time})
		}
	}

	// ACT
	collector := NewCollector(10, "a")
	mails, err := collector.Collect(temp)

	// ASSERT
	require.NoError(t, err)
	assert.Equal(t, &expected, mails)
}

func TestCollector_TimeNotIncluded(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	expected := []Mail{}

	// INBOX
	{
		mailFolder := createMailFolder(t, temp, "")
		createMailByTime(t, mailFolder, "new", agoDays(1).Add(time.Second*2), 1)
		createMailByName(t, mailFolder, "new", "abc", 1) // 日付の情報含まない
		{
			// 収集対象
			time := agoDays(1).Add(time.Second * (-2))
			mailPath := createMailByTime(t, mailFolder, "cur", time, 1)
			expected = append(expected, Mail{Path: mailPath, FolderName: "", Size: 1, Time: time})
		}
	}

	// ACT
	collector := NewCollector(1, "")
	mails, err := collector.Collect(temp)

	// ASSERT
	require.NoError(t, err)
	assert.Equal(t, &expected, mails)
}

func TestCollector_SkipSubFolder(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	expected := []Mail{}

	// INBOX
	{
		mailFolder := createMailFolder(t, temp, "")
		{
			// 収集対象
			time := agoDays(2)
			mailPath := createMailByTime(t, mailFolder, "cur", time, 1)
			expected = append(expected, Mail{Path: mailPath, FolderName: "", Size: 1, Time: time})
		}
	}

	// その他フォルダ
	{
		// Maildirとしてあるべきフォルダ無し
		// -> エラーとならずにスキップされること
		createDir(t, temp, ".a")
	}

	// ACT
	collector := NewCollector(2, "")
	mails, err := collector.Collect(temp)

	// ASSERT
	require.NoError(t, err)
	assert.Equal(t, &expected, mails)
}

func TestCollector_InvalidFolderName(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	// INBOX
	{
		mailFolder := createMailFolder(t, temp, "")
		createMailByTime(t, mailFolder, "cur", agoDays(10), 1)
	}

	// その他フォルダ
	{
		// フォルダ名としておかしなもの(修正UTF-7としてデコードできないもの)
		mailFolder := createMailFolder(t, temp, ".&A")
		createMailByTime(t, mailFolder, "cur", agoDays(10), 1)
	}

	// ACT
	collector := NewCollector(2, "")
	_, err := collector.Collect(temp)

	// ASSERT
	assert.EqualError(t, err, "&A is invalid folder name: utf7: invalid UTF-7")
}

func TestCollector_RootFolderNotFound(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	rootMailFolderPath := filepath.Join(temp, "xx") // 存在しないフォルダ

	// ACT
	collector := NewCollector(2, "")
	_, err := collector.Collect(rootMailFolderPath)

	// ASSERT
	require.Error(t, err)
	// OSによってエラーメッセージが異なるのでファイル名部分だけチェック
	expect := "open " + rootMailFolderPath
	assert.Contains(t, err.Error(), expect)
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

func createMailFolder(t *testing.T, parentDir string, folderName string) string {

	folderDir := parentDir

	if folderName != "" {
		// INBOX以外
		folderDir = createDir(t, parentDir, folderName)
	}

	createDir(t, folderDir, "tmp")
	createDir(t, folderDir, "new")
	createDir(t, folderDir, "cur")

	return folderDir
}

func createMailByTime(t *testing.T, folderDir string, sub string, time time.Time, size int) string {

	return createMailByName(t, folderDir, sub, fmt.Sprintf("%d", time.Unix()), size)
}

func createMailByName(t *testing.T, folderDir string, sub string, name string, size int) string {

	mailPath := filepath.Join(folderDir, sub, name)
	err := os.MkdirAll(filepath.Dir(mailPath), 0777)
	require.NoError(t, err)

	// サイズだけあっていればよいので中身は適当に
	createFile(t, mailPath, strings.Repeat("x", size))

	return mailPath
}

func agoDays(days int64) time.Time {
	return time.Unix(time.Now().Unix()-(days*24*60*60), 0)
}
