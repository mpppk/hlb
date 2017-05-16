package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/mpppk/hlb/hlb"
	"context"
	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/git"
	"strconv"
)

// listissuesCmd represents the listissues command
var listissuesCmd = &cobra.Command{
	Use:   "issues",
	Short: "list issuses",
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

		issues, err := pj.GetIssues(ctx, remote.Owner, remote.RepoName)
		etc.PanicIfErrorExist(err)

		for _, issue := range issues {
			info := "#" + strconv.Itoa(issue.GetNumber()) + " " + issue.GetTitle()
			fmt.Println(info)
		}
	},
}

func init() {
	listCmd.AddCommand(listissuesCmd)
}
