package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/mpppk/hlb/hlblib"
	"strconv"
	"github.com/mpppk/hlb/etc"
)

// listpullrequestsCmd represents the listpullrequests command
var listpullrequestsCmd = &cobra.Command{
	Use:   "pull-requests",
	Short: "list pull-requests",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		base, err := hlblib.NewCmdBase()
		etc.PanicIfErrorExist(err)

		//pulls, err := sw.GetPullRequests()
		pulls, err := base.Client.GetPullRequests().List(base.Context, base.Remote.Owner, base.Remote.RepoName)

		etc.PanicIfErrorExist(err)

		for _, pull := range pulls {
			info := "#" + strconv.Itoa(pull.GetNumber()) + " " + pull.GetTitle()
			fmt.Println(info)
		}
	},
}

func init() {
	listCmd.AddCommand(listpullrequestsCmd)
}
