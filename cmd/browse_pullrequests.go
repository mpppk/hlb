package cmd

import (
	"fmt"

	"strconv"

	"github.com/mpppk/hlb/hlblib"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

var browsepullrequestsCmd = &cobra.Command{
	Use:   "pull-requests",
	Short: "browse pull-requests",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Println("Too many issue IDs")
		}

		base, err := hlblib.NewCmdBase()
		hlblib.PanicIfErrorExist(err)

		var url string
		if len(args) == 0 {
			u, err := base.Client.GetPullRequests().GetPullRequestsURL(base.Remote.Owner, base.Remote.RepoName)
			hlblib.PanicIfErrorExist(err)
			url = u
		} else {
			id, err := strconv.Atoi(args[0])
			hlblib.PanicIfErrorExist(err)

			u, err := base.Client.GetPullRequests().GetURL(base.Remote.Owner, base.Remote.RepoName, id)
			hlblib.PanicIfErrorExist(err)
			url = u
		}

		if urlFlag {
			fmt.Println(url)
		} else {
			open.Run(url)
		}
	},
}

func init() {
	browseCmd.AddCommand(browsepullrequestsCmd)
}
