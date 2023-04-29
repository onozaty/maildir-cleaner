package cmd

import (
	"fmt"
	"io"

	"github.com/onozaty/maildir-cleaner/action"
	"github.com/onozaty/maildir-cleaner/collector"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func newArchiveCmd() *cobra.Command {

	subCmd := &cobra.Command{
		Use:   "archive",
		Short: "Archive old mails",
		RunE: func(cmd *cobra.Command, args []string) error {

			maildirPath, _ := cmd.Flags().GetString("dir")
			age, _ := cmd.Flags().GetInt64("age")

			archiveFolderNameGenerator, err := newArchiveFolderNameGenerator(cmd.Flags())
			if err != nil { // 許可されていなパラメータの可能性あり
				return err
			}

			excludeFolderNames, _ := cmd.Flags().GetStringArray("exclude-folder")

			// 引数の解析に成功した時点で、エラーが起きてもUsageは表示しない
			cmd.SilenceUsage = true

			return runArchive(
				maildirPath,
				age,
				archiveFolderNameGenerator,
				excludeFolderNames,
				cmd.OutOrStdout())
		},
	}

	subCmd.Flags().StringP("dir", "d", "", "User maildir path.")
	subCmd.MarkFlagRequired("dir")
	subCmd.Flags().Int64P("age", "a", 0, "The number of age days to be archived.\nIf you specify 10, mail that has been in the mailbox for more than 10 days since its arrival will be archived.")
	subCmd.MarkFlagRequired("age")

	subCmd.Flags().StringP("archive-folder", "", "Archived", "Archive folder name.")
	subCmd.Flags().StringP("archive-pattern", "", "keep", "Archive pattern. can be specified: keep, year, month")
	subCmd.Flags().StringArrayP("exclude-folder", "", []string{}, "The name of the folder to exclude.")

	return subCmd
}

func runArchive(maildirPath string, age int64, archiveFolderNameGenerator action.ArchiveFolderNameGenerator, excludeFolderNames []string, writer io.Writer) error {

	// 対象のメールを収集
	fmt.Fprintf(writer, "Starts searching for the target mails. maildir: %s age: %d\n", maildirPath, age)
	collector := collector.NewCollector(
		age,
		// アーカイブフォルダも対象外に
		append(excludeFolderNames, archiveFolderNameGenerator.BaseName())...)
	mails, err := collector.Collect(maildirPath)

	if err != nil {
		return err
	}

	if len(*mails) == 0 {
		// アーカイブ対象無し
		fmt.Fprintf(writer, "Completed search. There were no target mails.\n")
		return nil
	}

	fmt.Fprintf(writer, "Completed search. The target mails are listed below.\n")
	renderTargetMails(writer, mails)

	// アーカイブ実施
	fmt.Fprintf(writer, "Starts archiving mails.\n")
	archivedMails, err := action.Archive(maildirPath, mails, archiveFolderNameGenerator)
	if err != nil {
		return err
	}

	fmt.Fprintf(writer, "Completed archive. The archived mails are listed below.\n")
	renderTargetMails(writer, archivedMails)

	return nil
}

func newArchiveFolderNameGenerator(f *pflag.FlagSet) (action.ArchiveFolderNameGenerator, error) {

	archiveFolderName, _ := f.GetString("archive-folder")
	archivePattern, _ := f.GetString("archive-pattern")

	switch archivePattern {
	case "keep":
		return &action.KeepArchiveFolderNameGenerator{
			ArchiveFolderBaseName: archiveFolderName,
		}, nil
	case "year":
		return &action.YearArchiveFolderNameGenerator{
			ArchiveFolderBaseName: archiveFolderName,
		}, nil
	case "month":
		return &action.MonthArchiveFolderNameGenerator{
			ArchiveFolderBaseName: archiveFolderName,
		}, nil
	default:
		return nil, fmt.Errorf("invalid archive-pattern '%s'", archivePattern)
	}
}
