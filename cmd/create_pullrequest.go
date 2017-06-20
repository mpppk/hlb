package cmd

import (
	"fmt"

	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/hlblib"
	"github.com/spf13/cobra"
)

var createpullrequestCmd = &cobra.Command{
	Use:   "pull-request",
	Short: "Create pull request",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		title := "Sample title"
		message := "Sample message"
		baseOwner := "mpppk"
		baseBranch := "master"
		headBranch := "some_feature"

		base, err := hlblib.NewCmdBase()
		etc.PanicIfErrorExist(err)
		sw := hlblib.ClientWrapper{Base: base}
		pr, err := sw.CreatePullRequest(baseOwner, baseBranch, headBranch, title, message)
		etc.PanicIfErrorExist(err)
		fmt.Println(pr)
	},
}

func init() {
	createCmd.AddCommand(createpullrequestCmd)
}
