package action

import (
	"os"
	"path/filepath"

	"github.com/onozaty/maildir-cleaner/collector"
	"github.com/onozaty/maildir-cleaner/folder"
)

func Delete(rootMailFolderPath string, mails *[]collector.Mail) error {
	for _, mail := range *mails {
		if err := os.Remove(mail.FullPath); err != nil {
			return err
		}
	}

	return nil
}

func Archive(rootMailFolderPath string, mails *[]collector.Mail, archiveFolderBaseName string) (*[]collector.Mail, error) {
	archivedMails := []collector.Mail{}

	for _, mail := range *mails {
		archivedMail, err := archiveMail(rootMailFolderPath, mail, archiveFolderBaseName)
		if err != nil {
			return nil, err
		}
		archivedMails = append(archivedMails, *archivedMail)
	}

	return &archivedMails, nil
}

func archiveMail(rootMailFolderPath string, mail collector.Mail, archiveFolderBaseName string) (*collector.Mail, error) {

	archiveFolderName := joinArchiveFolderName(archiveFolderBaseName, mail.FolderName)

	archiveFolderPath, err := folder.Setup(rootMailFolderPath, archiveFolderName)
	if err != nil {
		return nil, err
	}

	archiveMailPath := filepath.Join(archiveFolderPath, mail.SubDirName, mail.FileName)
	if err := os.Rename(mail.FullPath, archiveMailPath); err != nil {
		return nil, err
	}

	return &collector.Mail{
		FullPath:   archiveMailPath,
		FolderName: archiveFolderName,
		SubDirName: mail.SubDirName,
		FileName:   mail.FileName,
		Size:       mail.Size,
		Time:       mail.Time,
	}, nil
}

func joinArchiveFolderName(archiveFolderBaseName string, beforeFolderName string) string {
	if beforeFolderName == "" {
		// INBOXの場合は、指定したアーカイブフォルダに
		return archiveFolderBaseName
	}

	return archiveFolderBaseName + "." + beforeFolderName
}
