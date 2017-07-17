package cmd

import (
	"github.com/mpppk/hlb/etc"
	"github.com/spf13/cobra"
)

// forkCmd represents the fork command
var forkCmd = &cobra.Command{
	Use:   "fork",
	Short: "Create Fork Repository",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfCmd, _, err := cmd.Root().Find([]string{"create", "fork"})
		etc.PanicIfErrorExist(err)
		cfCmd.Run(cfCmd, args)
	},
}

func init() {
	RootCmd.AddCommand(forkCmd)
}
