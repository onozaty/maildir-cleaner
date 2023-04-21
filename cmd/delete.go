package cmd

import (
	"fmt"
	"io"
	"sort"

	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"github.com/onozaty/maildir-cleaner/action"
	"github.com/onozaty/maildir-cleaner/collector"
	"github.com/spf13/cobra"
)

func newDeleteCmd() *cobra.Command {

	subCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete old mails",
		RunE: func(cmd *cobra.Command, args []string) error {

			maildirPath, _ := cmd.Flags().GetString("dir")
			age, _ := cmd.Flags().GetInt64("age")

			// 引数の解析に成功した時点で、エラーが起きてもUsageは表示しない
			cmd.SilenceUsage = true

			return runDelete(
				maildirPath,
				age,
				cmd.OutOrStdout())
		},
	}

	subCmd.Flags().StringP("dir", "d", "", "User maildir path.")
	subCmd.MarkFlagRequired("dir")
	subCmd.Flags().Int64P("age", "a", 0, "The number of age days to be deleted.\nIf you specify 10, mail that has been in the mailbox for more than 10 days since its arrival will be deleted.")
	subCmd.MarkFlagRequired("age")

	return subCmd
}

func runDelete(maildirPath string, age int64, writer io.Writer) error {

	// 対象のメールを収集
	fmt.Fprintf(writer, "Starts searching for the target mails. maildir: %s age: %d\n", maildirPath, age)
	collector := collector.NewCollector(age, "")
	mails, err := collector.Collect(maildirPath)

	if err != nil {
		return err
	}

	fmt.Fprintf(writer, "Completed search. The target mails are listed below.\n")
	renderTargetMails(writer, mails)

	// 削除実施
	fmt.Fprintf(writer, "Starts deleting mails.\n")
	if err := action.Delete(maildirPath, mails); err != nil {
		return err
	}
	fmt.Fprintf(writer, "Completed deletion.\n")

	return nil
}

func renderTargetMails(writer io.Writer, mails *[]collector.Mail) {

	aggregateResults := aggregateMails(mails)
	allMailCount := int64(0)
	allMailSize := int64(0)

	table := tablewriter.NewWriter(writer)
	table.SetAutoFormatHeaders(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT})
	table.SetHeader([]string{"Name", "Number of mails", "Total size(byte)"})

	for _, result := range aggregateResults {
		table.Append(
			[]string{result.FolderName, humanize.Comma(result.Count), humanize.Comma(result.TotalSize)})

		allMailCount += result.Count
		allMailSize += result.TotalSize
	}

	table.SetFooterAlignment(tablewriter.ALIGN_RIGHT)
	table.SetFooter([]string{"Total", humanize.Comma(allMailCount), humanize.Comma(allMailSize)})

	table.Render()
}

func aggregateMails(mails *[]collector.Mail) []aggregateResult {

	// フォルダ名毎に集計
	aggregateResultsMap := map[string]*aggregateResult{}

	for _, mail := range *mails {
		result := aggregateResultsMap[mail.FolderName]
		if result == nil {
			result = &aggregateResult{
				FolderName: mail.FolderName,
				Count:      0,
				TotalSize:  0,
			}
			aggregateResultsMap[mail.FolderName] = result
		}

		result.Count++
		result.TotalSize += mail.Size
	}

	aggregateResults := []aggregateResult{}
	for _, result := range aggregateResultsMap {
		aggregateResults = append(aggregateResults, *result)
	}

	// フォルダ名でソートして返す
	sort.Slice(aggregateResults, func(i, j int) bool {
		return aggregateResults[i].FolderName < aggregateResults[j].FolderName
	})

	return aggregateResults
}

type aggregateResult struct {
	FolderName string
	Count      int64
	TotalSize  int64
}
