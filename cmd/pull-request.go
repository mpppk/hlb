package cmd

import (
	"github.com/mpppk/hlb/hlblib"
	"github.com/spf13/cobra"
)

var pullrequestCmd = &cobra.Command{
	Use:   "pull-request",
	Short: "Create pull-request (experimental)",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cprCmd, _, err := cmd.Root().Find([]string{"create", "pull-request"})
		hlblib.PanicIfErrorExist(err)
		cprCmd.Run(cprCmd, args)
	},
}

func init() {
	RootCmd.AddCommand(pullrequestCmd)
}
