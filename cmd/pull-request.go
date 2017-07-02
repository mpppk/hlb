package cmd

import (
	"github.com/mpppk/hlb/etc"
	"github.com/spf13/cobra"
)

var pullrequestCmd = &cobra.Command{
	Use:   "pull-request",
	Short: "Create pull-request",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cprCmd, _, err := cmd.Root().Find([]string{"create", "pull-request"})
		etc.PanicIfErrorExist(err)
		cprCmd.Run(cprCmd, args)
	},
}

func init() {
	RootCmd.AddCommand(pullrequestCmd)
}
