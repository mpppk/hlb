package cmd

import (
	"context"
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
		base, err := hlblib.NewCmdContext()
		hlblib.PanicIfErrorExist(err)

		//pulls, err := sw.GetPullRequests()
		ctx := context.Background()
		pulls, err := base.Client.GetPullRequests().List(ctx, base.Remote.Owner, base.Remote.RepoName)

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
