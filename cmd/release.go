package cmd

import (
	"github.com/mpppk/hlb/hlblib"
	"github.com/spf13/cobra"
)

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Create release page and upload files (experimental)",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		crCmd, _, err := cmd.Root().Find([]string{"create", "release"})
		hlblib.PanicIfErrorExist(err)
		crCmd.Run(crCmd, args)
	},
}

func init() {
	RootCmd.AddCommand(releaseCmd)
}
