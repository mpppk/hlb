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
		ctx := context.Background()

		var config etc.Config
		err := viper.Unmarshal(&config)
		etc.PanicIfErrorExist(err)

		remote, err := git.GetDefaultRemote(".")
		etc.PanicIfErrorExist(err)

		host, ok := config.FindHost(remote.ServiceHostName)
		if !ok {
			panic("host not found: " + remote.ServiceHostName)
		}

		pj, err := hlb.GetService(ctx, host)
		etc.PanicIfErrorExist(err)

		pulls, err := pj.GetPullRequests(ctx, remote.Owner, remote.RepoName)
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
