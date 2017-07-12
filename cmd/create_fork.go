package cmd

import (
	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/hlblib"
	"github.com/spf13/cobra"
)

var createforkCmd = &cobra.Command{
	Use:   "fork",
	Short: "Create Fork Repository",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		base, err := hlblib.NewCmdBase()
		etc.PanicIfErrorExist(err)
		sw := hlblib.ClientWrapper{Base: base}
		_, err = sw.CreateFork()
		etc.PanicIfErrorExist(err)

		// Add remote repository
	},
}

func init() {
	createCmd.AddCommand(createforkCmd)
}
