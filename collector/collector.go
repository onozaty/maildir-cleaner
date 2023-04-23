package collector

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/onozaty/maildir-cleaner/folder"
)

type Mail struct {
	FullPath   string
	FolderName string // エンコード前のメールフォルダ名(先頭の"."は付かない)
	SubDirName string // curなど
	FileName   string
	Size       int64
	Time       time.Time
}

type Collector struct {
	target func(Mail) bool
}

func NewCollector(ageOfDays int64, ignoreFolderName string) *Collector {

	// 現在日時 - 経過日
	targetMaxTime := time.Now().AddDate(0, 0, -int(ageOfDays))

	return &Collector{
		target: func(mail Mail) bool {

			if ignoreFolderName != "" &&
				(mail.FolderName == ignoreFolderName || strings.HasPrefix(mail.FolderName, ignoreFolderName+".")) {
				// 対象外のフォルダ名と一致(サブフォルダも考慮)
				return false
			}

			// 日時が取れなかった場合(=0)は対象外
			return mail.Time.Unix() != 0 && mail.Time.Before(targetMaxTime)
		},
	}
}

func (c *Collector) Collect(rootMailFolderPath string) (*[]Mail, error) {

	collectedMails := []Mail{}

	// ルート(INBOX)
	mails, err := c.collectMailFolder("", rootMailFolderPath, false)
	if err != nil {
		return nil, err
	}
	collectedMails = append(collectedMails, *mails...)

	// その他メールフォルダ
	entries, err := os.ReadDir(rootMailFolderPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		// ディレクトリの先頭が"."になっているものがメールフォルダ
		if entry.IsDir() && strings.HasPrefix(entry.Name(), ".") {
			mailFolderName, err := folder.DecodeMailFolderName(entry.Name()[1:]) // 先頭の"."は除く
			if err != nil {
				return nil, err
			}

			// その他メールフォルダは作成直後にcurフォルダなどが無いことがあるので無かったらスキップするように設定
			mails, err := c.collectMailFolder(mailFolderName, filepath.Join(rootMailFolderPath, entry.Name()), true)
			if err != nil {
				return nil, err
			}
			collectedMails = append(collectedMails, *mails...)
		}
	}

	// フォルダ名+ファイル名でソート
	// (順番が必ず同じになるように)
	sort.Slice(collectedMails, func(i, j int) bool {
		if collectedMails[i].FolderName == collectedMails[j].FolderName {
			return collectedMails[i].FileName < collectedMails[j].FileName
		}
		return collectedMails[i].FolderName < collectedMails[j].FolderName
	})

	return &collectedMails, nil
}

func (c *Collector) collectMailFolder(mailFolderName string, mailFolderPath string, skipSubdirMissing bool) (*[]Mail, error) {

	collectedMails := []Mail{}

	// tmpにあるのは配送中のものなので対象から除いておく
	for _, subName := range []string{"new", "cur"} {
		subDir := filepath.Join(mailFolderPath, subName)
		if _, err := os.Stat(subDir); os.IsNotExist(err) && skipSubdirMissing {
			// サブディレクトリが無いことを無視する場合はスキップ
			continue
		}

		mails, err := c.collectMails(mailFolderName, filepath.Join(mailFolderPath, subName))
		if err != nil {
			return nil, err
		}

		collectedMails = append(collectedMails, *mails...)
	}

	return &collectedMails, nil
}

func (c *Collector) collectMails(mailFolderName string, dirPath string) (*[]Mail, error) {

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	collectedMails := []Mail{}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			return nil, err
		}

		mail := Mail{
			FullPath:   filepath.Join(dirPath, info.Name()),
			FolderName: mailFolderName,
			SubDirName: filepath.Base(dirPath),
			FileName:   info.Name(),
			Size:       info.Size(),
			Time:       MailTime(info.Name()),
		}

		if c.target(mail) {
			collectedMails = append(collectedMails, mail)
		}
	}

	return &collectedMails, nil
}
