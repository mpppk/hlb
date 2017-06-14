package cmd

import (
	"github.com/mpppk/hlb/etc"
	"github.com/spf13/cobra"
)

var browseboardsCmd = &cobra.Command{
	Use:   "boards",
	Short: "browse boards",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		pjCmd, _, err := cmd.Root().Find([]string{"browse", "projects"})
		etc.PanicIfErrorExist(err)
		pjCmd.Run(pjCmd, args)
	},
}

func init() {
	browseCmd.AddCommand(browseboardsCmd)
}
