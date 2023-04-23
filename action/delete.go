package action

import (
	"os"

	"github.com/onozaty/maildir-cleaner/collector"
)

func Delete(rootMailFolderPath string, mails *[]collector.Mail) error {
	for _, mail := range *mails {
		if err := os.Remove(mail.FullPath); err != nil {
			return err
		}
	}

	return nil
}
