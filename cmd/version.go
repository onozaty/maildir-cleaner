package cmd

import (
	"runtime"

	"github.com/spf13/cobra"
)

var (
	Revision = "dev"
	Version  = "dev"
)

func newVersionCmd() *cobra.Command {

	subCmd := &cobra.Command{
		Use:   "version",
		Short: "Show maildir-cleaner command version",
		RunE: func(cmd *cobra.Command, args []string) error {

			// 引数の解析に成功した時点で、エラーが起きてもUsageは表示しない
			cmd.SilenceUsage = true

			cmd.Printf(`Version: %s
Revision: %s
OS: %s
Arch: %s
`, Version, Revision, runtime.GOOS, runtime.GOARCH)

			return nil
		},
	}

	return subCmd
}
