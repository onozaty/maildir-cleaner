package folder

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

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

func Setup(rootMailFolderPath string, folderName string) (string, error) {
	encodedFolderName, err := EncodeMailFolderName(folderName)
	if err != nil {
		return "", err
	}

	// メールフォルダに対応するディレクトリが無かったら作成
	folderPath := filepath.Join(rootMailFolderPath, "."+encodedFolderName)
	if err := ensureDir(folderPath); err != nil {
		return "", err
	}

	for _, subName := range []string{"new", "cur", "tmp"} {
		subDir := filepath.Join(folderPath, subName)
		if err := ensureDir(subDir); err != nil {
			return "", err
		}
	}

	// メールフォルダを購読状態に
	if err := subscribe(rootMailFolderPath, encodedFolderName); err != nil {
		return "", err
	}

	return folderPath, nil
}

func subscribe(rootMailFolderPath string, encodedFolderName string) error {
	// subscriptions に該当のフォルダを追加
	// TODO: Dovecot 以外への対応
	subscriptionsPath := filepath.Join(rootMailFolderPath, "subscriptions")
	if isNotExist(subscriptionsPath) {
		return fmt.Errorf("subscriptions file not found: currently only dovecot is supported")
	}

	// 購読済みかチェック
	file, err := os.OpenFile(subscriptionsPath, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == encodedFolderName {
			// 既に購読済みなので何もしない
			return nil
		}
	}

	// 末尾が改行でなければ、改行を追加したうえでフォルダを追加
	fileStat, err := file.Stat()
	if err != nil {
		return err
	}

	if fileStat.Size() != 0 {
		b := make([]byte, 1)
		if _, err := file.ReadAt(b, fileStat.Size()-1); err != nil {
			return err
		}
		if b[0] != '\n' {
			file.WriteString("\n")
		}
	}

	file.WriteString(encodedFolderName + "\n")
	return nil
}

func ensureDir(dirPath string) error {
	if isNotExist(dirPath) {
		err := os.Mkdir(dirPath, 0777)
		if err != nil {
			return err
		}
	}

	return nil
}

func isNotExist(dirPath string) bool {
	_, err := os.Stat(dirPath)
	return os.IsNotExist(err)
}
