package cmd

import (
	"fmt"
	"io"

	"github.com/onozaty/maildir-cleaner/action"
	"github.com/onozaty/maildir-cleaner/collector"
	"github.com/spf13/cobra"
)

func newArchiveCmd() *cobra.Command {

	subCmd := &cobra.Command{
		Use:   "archive",
		Short: "Archive old mails",
		RunE: func(cmd *cobra.Command, args []string) error {

			maildirPath, _ := cmd.Flags().GetString("dir")
			age, _ := cmd.Flags().GetInt64("age")
			archiveFolderName, _ := cmd.Flags().GetString("archive-folder")

			// 引数の解析に成功した時点で、エラーが起きてもUsageは表示しない
			cmd.SilenceUsage = true

			return runArchive(
				maildirPath,
				age,
				archiveFolderName,
				cmd.OutOrStdout())
		},
	}

	subCmd.Flags().StringP("dir", "d", "", "User maildir path.")
	subCmd.MarkFlagRequired("dir")
	subCmd.Flags().Int64P("age", "a", 0, "The number of age days to be archived.\nIf you specify 10, mail that has been in the mailbox for more than 10 days since its arrival will be archived.")
	subCmd.MarkFlagRequired("age")

	subCmd.Flags().StringP("archive-folder", "", "Archived", "Archive folder name.")

	return subCmd
}

func runArchive(maildirPath string, age int64, archiveFolderName string, writer io.Writer) error {

	// 対象のメールを収集
	fmt.Fprintf(writer, "Starts searching for the target mails. maildir: %s age: %d\n", maildirPath, age)
	collector := collector.NewCollector(age, "")
	mails, err := collector.Collect(maildirPath)

	if err != nil {
		return err
	}

	fmt.Fprintf(writer, "Completed search. The target mails are listed below.\n")
	renderTargetMails(writer, mails)

	// アーカイブ実施
	fmt.Fprintf(writer, "Starts archiving mails.\n")
	archivedMails, err := action.Archive(maildirPath, mails, archiveFolderName)
	if err != nil {
		return err
	}

	fmt.Fprintf(writer, "Completed archive. The archived mails are listed below.\n")
	renderTargetMails(writer, archivedMails)

	return nil
}
