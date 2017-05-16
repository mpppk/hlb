package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/mpppk/hlb/hlb"
	"github.com/skratchdot/open-golang/open"
	"github.com/mpppk/hlb/etc"
	"strconv"
	"github.com/mpppk/hlb/git"
	"context"
)

// browsepullrequestsCmd represents the browsepullrequests command
var browsepullrequestsCmd = &cobra.Command{
	Use:   "pull-requests",
	Short: "browse pull-requests",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Println("Too many issue IDs")
		}

		base, err := hlb.NewCmdBase()
		etc.PanicIfErrorExist(err)

		if len(args) == 0 {
			url, err := base.Service.GetPullRequestsURL(base.Remote.Owner, base.Remote.RepoName)
			etc.PanicIfErrorExist(err)
			open.Run(url)
			return
		}else {
			id, err := strconv.Atoi(args[0])
			etc.PanicIfErrorExist(err)

			url, err := base.Service.GetPullRequestURL(base.Remote.Owner, base.Remote.RepoName, id)
			etc.PanicIfErrorExist(err)
			open.Run(url)
			return
		}
	},
}

func init() {
	browseCmd.AddCommand(browsepullrequestsCmd)
}
