package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/mpppk/hlb/hlb"
	"strconv"
	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/git"
	"context"
)

// listpullrequestsCmd represents the listpullrequests command
var listpullrequestsCmd = &cobra.Command{
	Use:   "pull-requests",
	Short: "list pull-requests",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		base, err := hlb.NewCmdBase()
		etc.PanicIfErrorExist(err)

		pulls, err := base.Service.GetPullRequests(base.Context, base.Remote.Owner, base.Remote.RepoName)
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
