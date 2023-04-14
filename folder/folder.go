package folder

import (
	"fmt"

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
