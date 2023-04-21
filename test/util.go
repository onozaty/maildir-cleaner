package test

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CreateDir(t *testing.T, parent string, name string) string {

	dir := filepath.Join(parent, name)
	err := os.Mkdir(dir, 0777)
	if !os.IsExist(err) {
		require.NoError(t, err)
	}

	return dir
}

func CreateFile(t *testing.T, path string, content string) fs.FileInfo {

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

func ReadFile(t *testing.T, path string) string {

	bo, err := os.ReadFile(path)
	require.NoError(t, err)

	return string(bo)
}

func CreateMailFolder(t *testing.T, parentDir string, physicalFolderName string) string {

	folderDir := parentDir

	if physicalFolderName != "" {
		// INBOX以外
		folderDir = CreateDir(t, parentDir, physicalFolderName)
	}

	CreateDir(t, folderDir, "tmp")
	CreateDir(t, folderDir, "new")
	CreateDir(t, folderDir, "cur")

	return folderDir
}

func CreateMailByTime(t *testing.T, folderDir string, sub string, time time.Time, size int) (string, string) {

	return CreateMailByName(t, folderDir, sub, fmt.Sprintf("%d", time.Unix()), size)
}

func CreateMailByName(t *testing.T, folderDir string, sub string, name string, size int) (string, string) {

	mailPath := filepath.Join(folderDir, sub, name)
	err := os.MkdirAll(filepath.Dir(mailPath), 0777)
	require.NoError(t, err)

	// サイズだけあっていればよいので中身は適当に
	CreateFile(t, mailPath, strings.Repeat("x", size))

	return mailPath, name
}

func AgoDays(t *testing.T, days int) time.Time {
	return time.Unix(time.Now().Unix()-(int64(days)*24*60*60), 0)
}
