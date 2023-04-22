package cmd

import (
	"io"
	"sort"

	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"github.com/onozaty/maildir-cleaner/collector"
)

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
