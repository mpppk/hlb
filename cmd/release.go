package cmd

import (
	"github.com/mpppk/hlb/etc"
	"github.com/spf13/cobra"
)

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Create release page and upload files",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		crCmd, _, err := cmd.Root().Find([]string{"create", "release"})
		etc.PanicIfErrorExist(err)
		crCmd.Run(crCmd, args)
	},
}

func init() {
	RootCmd.AddCommand(releaseCmd)
}
