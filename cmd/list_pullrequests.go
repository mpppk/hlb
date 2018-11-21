package cmd

import (
	"fmt"

	"strconv"

	"github.com/mpppk/hlb/hlblib"
	"github.com/spf13/cobra"
)

// listpullrequestsCmd represents the listpullrequests command
var listpullrequestsCmd = &cobra.Command{
	Use:   "pull-requests",
	Short: "list pull-requests",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		base, err := hlblib.NewCmdBase()
		hlblib.PanicIfErrorExist(err)

		//pulls, err := sw.GetPullRequests()
		pulls, err := base.Client.GetPullRequests().List(base.Context, base.Remote.Owner, base.Remote.RepoName)

		hlblib.PanicIfErrorExist(err)

		for _, pull := range pulls {
			info := "#" + strconv.Itoa(pull.GetNumber()) + " " + pull.GetTitle()
			fmt.Println(info)
		}
	},
}

func init() {
	listCmd.AddCommand(listpullrequestsCmd)
}
