package cmd

import (
	"fmt"
	"io"

	"github.com/onozaty/maildir-cleaner/collector"
	"github.com/spf13/cobra"
)

func newSearchCmd() *cobra.Command {

	subCmd := &cobra.Command{
		Use:   "search",
		Short: "Search old mails",
		RunE: func(cmd *cobra.Command, args []string) error {

			maildirPath, _ := cmd.Flags().GetString("dir")
			age, _ := cmd.Flags().GetInt64("age")
			excludeFolderNames, _ := cmd.Flags().GetStringArray("exclude-folder")

			// 引数の解析に成功した時点で、エラーが起きてもUsageは表示しない
			cmd.SilenceUsage = true

			return runSearch(
				maildirPath,
				age,
				excludeFolderNames,
				cmd.OutOrStdout())
		},
	}

	subCmd.Flags().StringP("dir", "d", "", "User maildir path.")
	subCmd.MarkFlagRequired("dir")
	subCmd.Flags().Int64P("age", "a", 0, "The number of age days to be displayed.\nIf you specify 10, mail that has been in the mailbox for more than 10 days since its arrival will be displayed.")
	subCmd.MarkFlagRequired("age")
	subCmd.Flags().StringArrayP("exclude-folder", "", []string{}, "The name of the folder to exclude.")

	return subCmd
}

func runSearch(maildirPath string, age int64, excludeFolderNames []string, writer io.Writer) error {

	// 対象のメールを収集
	fmt.Fprintf(writer, "Starts searching for the target mails. maildir: %s age: %d\n", maildirPath, age)
	collector := collector.NewCollector(age, excludeFolderNames...)
	mails, err := collector.Collect(maildirPath)

	if err != nil {
		return err
	}

	if len(*mails) == 0 {
		// 対象無し
		fmt.Fprintf(writer, "Completed search. There were no target mails.\n")
		return nil
	}

	fmt.Fprintf(writer, "Completed search. The target mails are listed below.\n")
	renderTargetMails(writer, mails)

	return nil
}
