package action

import (
	"os"
	"path/filepath"

	"github.com/onozaty/maildir-cleaner/collector"
	"github.com/onozaty/maildir-cleaner/folder"
)

type ArchiveFolderNameGenerator interface {
	Generate(mail collector.Mail) string
	BaseName() string
}

type KeepArchiveFolderNameGenerator struct {
	ArchiveFolderBaseName string
}

func (g *KeepArchiveFolderNameGenerator) Generate(mail collector.Mail) string {
	if mail.FolderName == "" {
		// INBOXの場合は、指定したアーカイブフォルダに
		return g.ArchiveFolderBaseName
	}

	return g.ArchiveFolderBaseName + "." + mail.FolderName
}

func (g *KeepArchiveFolderNameGenerator) BaseName() string {
	return g.ArchiveFolderBaseName
}

type YearArchiveFolderNameGenerator struct {
	ArchiveFolderBaseName string
}

func (g *YearArchiveFolderNameGenerator) Generate(mail collector.Mail) string {
	year := mail.Time.UTC().Format("2006")
	return g.ArchiveFolderBaseName + "." + year
}

func (g *YearArchiveFolderNameGenerator) BaseName() string {
	return g.ArchiveFolderBaseName
}

type MonthArchiveFolderNameGenerator struct {
	ArchiveFolderBaseName string
}

func (g *MonthArchiveFolderNameGenerator) Generate(mail collector.Mail) string {
	year := mail.Time.UTC().Format("2006")
	month := mail.Time.UTC().Format("01")
	return g.ArchiveFolderBaseName + "." + year + "." + month
}

func (g *MonthArchiveFolderNameGenerator) BaseName() string {
	return g.ArchiveFolderBaseName
}

func Archive(rootMailFolderPath string, mails *[]collector.Mail, archiveFolderNameGenerator ArchiveFolderNameGenerator) (*[]collector.Mail, error) {
	archivedMails := []collector.Mail{}

	for _, mail := range *mails {
		archiveFolderName := archiveFolderNameGenerator.Generate(mail)
		archivedMail, err := archiveMail(rootMailFolderPath, mail, archiveFolderName)
		if err != nil {
			return nil, err
		}
		archivedMails = append(archivedMails, *archivedMail)
	}

	return &archivedMails, nil
}

func archiveMail(rootMailFolderPath string, mail collector.Mail, archiveFolderName string) (*collector.Mail, error) {

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
