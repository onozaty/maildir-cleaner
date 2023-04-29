package collector

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/onozaty/maildir-cleaner/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollector(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	expected := []Mail{}

	// INBOX
	{
		mailFolder := test.CreateMailFolder(t, temp, "")
		test.CreateMailByTime(t, mailFolder, "new", test.AgoDays(t, 0), 1)
		test.CreateMailByTime(t, mailFolder, "new", test.AgoDays(t, 1), 1)
		test.CreateMailByTime(t, mailFolder, "cur", test.AgoDays(t, 2), 1)
		{
			// 収集対象
			time := test.AgoDays(t, 3)
			mailPath, fileName := test.CreateMailByTime(t, mailFolder, "cur", time, 1)
			expected = append(expected, Mail{
				FullPath:   mailPath,
				FolderName: "",
				SubDirName: "cur",
				FileName:   fileName,
				Size:       1,
				Time:       time,
			})
		}
		test.CreateMailByTime(t, mailFolder, "tmp", test.AgoDays(t, 4), 1)
		test.CreateMailByTime(t, mailFolder, "tmp", test.AgoDays(t, 5), 1)
	}

	// その他フォルダ
	{
		mailFolder := test.CreateMailFolder(t, temp, ".A")
		{
			// 収集対象
			time := test.AgoDays(t, 5)
			mailPath, fileName := test.CreateMailByTime(t, mailFolder, "new", time, 2)
			expected = append(expected, Mail{
				FullPath:   mailPath,
				FolderName: "A",
				SubDirName: "new",
				FileName:   fileName,
				Size:       2,
				Time:       time,
			})
		}
		{
			// 収集対象
			time := test.AgoDays(t, 4)
			mailPath, fileName := test.CreateMailByTime(t, mailFolder, "new", time, 2)
			expected = append(expected, Mail{
				FullPath:   mailPath,
				FolderName: "A",
				SubDirName: "new",
				FileName:   fileName,
				Size:       2,
				Time:       time,
			})
		}
		{
			// 収集対象
			time := test.AgoDays(t, 3)
			mailPath, fileName := test.CreateMailByTime(t, mailFolder, "cur", time, 2)
			expected = append(expected, Mail{
				FullPath:   mailPath,
				FolderName: "A",
				SubDirName: "cur",
				FileName:   fileName,
				Size:       2,
				Time:       time,
			})
		}
		test.CreateMailByTime(t, mailFolder, "cur", test.AgoDays(t, 2), 2)
		test.CreateMailByTime(t, mailFolder, "tmp", test.AgoDays(t, 1), 2)
		test.CreateMailByTime(t, mailFolder, "tmp", test.AgoDays(t, 0), 2)
	}
	{
		mailFolder := test.CreateMailFolder(t, temp, ".B")
		{
			// 収集対象
			time := test.AgoDays(t, 5)
			mailPath, fileName := test.CreateMailByTime(t, mailFolder, "new", time, 3)
			expected = append(expected, Mail{
				FullPath:   mailPath,
				FolderName: "B",
				SubDirName: "new",
				FileName:   fileName,
				Size:       3,
				Time:       time,
			})
		}
	}
	{
		mailFolder := test.CreateMailFolder(t, temp, ".C")
		{
			// 収集対象
			time := test.AgoDays(t, 5)
			mailPath, fileName := test.CreateMailByTime(t, mailFolder, "cur", time, 4)
			expected = append(expected, Mail{
				FullPath:   mailPath,
				FolderName: "C",
				SubDirName: "cur",
				FileName:   fileName,
				Size:       4,
				Time:       time,
			})
		}
	}
	{
		mailFolder := test.CreateMailFolder(t, temp, ".D")
		test.CreateMailByTime(t, mailFolder, "tmp", test.AgoDays(t, 5), 5)
	}
	{
		// メールフォルダ以外のフォルダ(先頭に"."無し)
		mailFolder := test.CreateMailFolder(t, temp, "a")
		test.CreateMailByTime(t, mailFolder, "new", test.AgoDays(t, 5), 6)
		test.CreateMailByTime(t, mailFolder, "cur", test.AgoDays(t, 5), 6)
		test.CreateMailByTime(t, mailFolder, "tmp", test.AgoDays(t, 5), 6)
	}

	// ACT
	collector := NewCollector(3)
	mails, err := collector.Collect(temp)

	// ASSERT
	require.NoError(t, err)
	assert.Equal(t, &expected, mails)
}

func TestCollector_ExcludeFolderName(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	expected := []Mail{}

	// INBOX
	{
		mailFolder := test.CreateMailFolder(t, temp, "")
		test.CreateMailByTime(t, mailFolder, "new", test.AgoDays(t, 1), 1)
		{
			// 収集対象
			time := test.AgoDays(t, 10)
			mailPath, fileName := test.CreateMailByTime(t, mailFolder, "new", time, 1)
			expected = append(expected, Mail{
				FullPath:   mailPath,
				FolderName: "",
				SubDirName: "new",
				FileName:   fileName,
				Size:       1,
				Time:       time,
			})
		}
	}

	// a (対象外フォルダ：除外対象のフォルダ)
	{
		mailFolder := test.CreateMailFolder(t, temp, ".a")
		test.CreateMailByTime(t, mailFolder, "new", test.AgoDays(t, 10), 1)
		test.CreateMailByTime(t, mailFolder, "cur", test.AgoDays(t, 10), 1)
	}
	// aa (対象フォルダ：対象外フォルダと前方一致)
	{
		mailFolder := test.CreateMailFolder(t, temp, ".aa")
		{
			// 収集対象
			time := test.AgoDays(t, 10)
			mailPath, fileName := test.CreateMailByTime(t, mailFolder, "new", time, 1)
			expected = append(expected, Mail{
				FullPath:   mailPath,
				FolderName: "aa",
				SubDirName: "new",
				FileName:   fileName,
				Size:       1,
				Time:       time,
			})
		}
	}
	// ab (対象外フォルダ：対象外フォルダのサブフォルダ)
	{
		mailFolder := test.CreateMailFolder(t, temp, ".a.b")
		test.CreateMailByTime(t, mailFolder, "new", test.AgoDays(t, 10), 1)
		test.CreateMailByTime(t, mailFolder, "cur", test.AgoDays(t, 10), 1)
	}
	// b (対象フォルダ)
	{
		mailFolder := test.CreateMailFolder(t, temp, ".b")
		{
			// 収集対象
			time := test.AgoDays(t, 11)
			mailPath, fileName := test.CreateMailByTime(t, mailFolder, "cur", time, 1)
			expected = append(expected, Mail{
				FullPath:   mailPath,
				FolderName: "b",
				SubDirName: "cur",
				FileName:   fileName,
				Size:       1,
				Time:       time,
			})
		}
	}

	// ACT
	collector := NewCollector(10, "a")
	mails, err := collector.Collect(temp)

	// ASSERT
	require.NoError(t, err)
	assert.Equal(t, &expected, mails)
}

func TestCollector_ExcludeFolderNames(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()

	expected := []Mail{}

	// INBOX
	{
		mailFolder := test.CreateMailFolder(t, temp, "")
		test.CreateMailByTime(t, mailFolder, "new", test.AgoDays(t, 1), 1)
		{
			// 収集対象
			time := test.AgoDays(t, 10)
			mailPath, fileName := test.CreateMailByTime(t, mailFolder, "new", time, 1)
			expected = append(expected, Mail{
				FullPath:   mailPath,
				FolderName: "",
				SubDirName: "new",
				FileName:   fileName,
				Size:       1,
				Time:       time,
			})
		}
	}

	// a (対象外フォルダ：除外対象のフォルダ)
	{
		mailFolder := test.CreateMailFolder(t, temp, ".a")
		test.CreateMailByTime(t, mailFolder, "new", test.AgoDays(t, 10), 1)
		test.CreateMailByTime(t, mailFolder, "cur", test.AgoDays(t, 10), 1)
	}
	// aa (対象フォルダ：対象外フォルダと前方一致)
	{
		mailFolder := test.CreateMailFolder(t, temp, ".aa")
		{
			// 収集対象
			time := test.AgoDays(t, 10)
			mailPath, fileName := test.CreateMailByTime(t, mailFolder, "new", time, 1)
			expected = append(expected, Mail{
				FullPath:   mailPath,
				FolderName: "aa",
				SubDirName: "new",
				FileName:   fileName,
				Size:       1,
				Time:       time,
			})
		}
	}
	// ab (対象外フォルダ：対象外フォルダのサブフォルダ)
	{
		mailFolder := test.CreateMailFolder(t, temp, ".a.b")
		test.CreateMailByTime(t, mailFolder, "new", test.AgoDays(t, 10), 1)
		test.CreateMailByTime(t, mailFolder, "cur", test.AgoDays(t, 10), 1)
	}

	// b (対象フォルダ：対象外フォルダの一部)
	{
		mailFolder := test.CreateMailFolder(t, temp, ".b")
		{
			// 収集対象
			time := test.AgoDays(t, 11)
			mailPath, fileName := test.CreateMailByTime(t, mailFolder, "new", time, 1)
			expected = append(expected, Mail{
				FullPath:   mailPath,
				FolderName: "b",
				SubDirName: "new",
				FileName:   fileName,
				Size:       1,
				Time:       time,
			})
		}
	}
	// bb (対象外フォルダ)
	{
		mailFolder := test.CreateMailFolder(t, temp, ".bb")
		test.CreateMailByTime(t, mailFolder, "cur", test.AgoDays(t, 11), 1)
	}

	// ACT
	collector := NewCollector(10, "a", "bb")
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
		mailFolder := test.CreateMailFolder(t, temp, "")
		test.CreateMailByTime(t, mailFolder, "new", test.AgoDays(t, 1).Add(time.Second*2), 1)
		test.CreateMailByName(t, mailFolder, "new", "abc", 1) // 日付の情報含まない
		{
			// 収集対象
			time := test.AgoDays(t, 1).Add(time.Second * (-2))
			mailPath, fileName := test.CreateMailByTime(t, mailFolder, "cur", time, 1)
			expected = append(expected, Mail{
				FullPath:   mailPath,
				FolderName: "",
				SubDirName: "cur",
				FileName:   fileName,
				Size:       1,
				Time:       time,
			})
		}
	}

	// ACT
	collector := NewCollector(1)
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
		mailFolder := test.CreateMailFolder(t, temp, "")
		{
			// 収集対象
			time := test.AgoDays(t, 2)
			mailPath, fileName := test.CreateMailByTime(t, mailFolder, "cur", time, 1)
			expected = append(expected, Mail{
				FullPath:   mailPath,
				FolderName: "",
				SubDirName: "cur",
				FileName:   fileName,
				Size:       1,
				Time:       time,
			})
		}
	}

	// その他フォルダ
	{
		// Maildirとしてあるべきフォルダ無し
		// -> エラーとならずにスキップされること
		test.CreateDir(t, temp, ".a")
	}

	// ACT
	collector := NewCollector(2)
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
		mailFolder := test.CreateMailFolder(t, temp, "")
		test.CreateMailByTime(t, mailFolder, "cur", test.AgoDays(t, 10), 1)
	}

	// その他フォルダ
	{
		// フォルダ名としておかしなもの(修正UTF-7としてデコードできないもの)
		mailFolder := test.CreateMailFolder(t, temp, ".&A")
		test.CreateMailByTime(t, mailFolder, "cur", test.AgoDays(t, 10), 1)
	}

	// ACT
	collector := NewCollector(2)
	_, err := collector.Collect(temp)

	// ASSERT
	assert.EqualError(t, err, "&A is invalid folder name: utf7: invalid UTF-7")
}

func TestCollector_RootFolderNotFound(t *testing.T) {

	// ARRANGE
	temp := t.TempDir()
	rootMailFolderPath := filepath.Join(temp, "xx") // 存在しないフォルダ

	// ACT
	collector := NewCollector(2)
	_, err := collector.Collect(rootMailFolderPath)

	// ASSERT
	require.Error(t, err)
	// OSによってエラーメッセージが異なるのでファイル名部分だけチェック
	expect := "open " + rootMailFolderPath
	assert.Contains(t, err.Error(), expect)
}
