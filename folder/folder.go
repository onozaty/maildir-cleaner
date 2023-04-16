package folder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/emersion/go-imap/utf7"
)

func DecodeMailFolderName(encodedName string) (string, error) {
	decoder := utf7.Encoding.NewDecoder()
	decodedName, err := decoder.String(encodedName)

	if err != nil {
		return "", fmt.Errorf("%s is invalid folder name: %w", encodedName, err)
	}
	return decodedName, nil
}

func EncodeMailFolderName(decodedName string) (string, error) {
	encoder := utf7.Encoding.NewEncoder()
	encodedName, err := encoder.String(decodedName)

	if err != nil {
		return "", fmt.Errorf("%s is invalid folder name: %w", decodedName, err)
	}
	return encodedName, nil
}

func Setup(rootMailFolderPath string, decodedFolderName string) (string, error) {
	encodedFolderName, err := EncodeMailFolderName(decodedFolderName)
	if err != nil {
		return "", err
	}

	folderPath := filepath.Join(rootMailFolderPath, "."+encodedFolderName)

	_, err = os.Stat(folderPath)
	if os.IsNotExist(err) {
		// メールフォルダに対応する物理ディレクトリの作成
		if err := os.Mkdir(folderPath, 0777); err != nil {
			return "", err
		}

		for _, subName := range []string{"new", "cur", "tmp"} {
			subDir := filepath.Join(folderPath, subName)
			if err := os.Mkdir(subDir, 0777); err != nil {
				return "", err
			}
		}

		// メールフォルダを購読状態に
		subscriptionsPath := filepath.Join(rootMailFolderPath, "subscriptions")
		_, err = os.Stat(subscriptionsPath)
		if os.IsExist(err) {

			bytes, err := os.ReadFile(subscriptionsPath)
			if err != nil {
				return "", err
			}
			contents := string(bytes)

			file, err := os.OpenFile(subscriptionsPath, os.O_APPEND|os.O_WRONLY, 0600)
			if err != nil {
				return "", err
			}
			defer file.Close()

			// 末尾に改行が無かったら付与
			if !strings.HasPrefix(contents, "\n") {
				file.WriteString("\n")
			}

			file.WriteString(encodedFolderName + "\n")
		}

	}

	return folderPath, nil
}
