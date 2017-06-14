package cmd

import (
	"github.com/mpppk/hlb/etc"
	"github.com/spf13/cobra"
)

var browsemergerequestsCmd = &cobra.Command{
	Use:   "merge-requests",
	Short: "browse merge-requests",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		prCmd, _, err := cmd.Root().Find([]string{"browse", "pull-requests"})
		etc.PanicIfErrorExist(err)
		prCmd.Run(prCmd, args)
	},
}

func init() {
	browseCmd.AddCommand(browsemergerequestsCmd)
}
