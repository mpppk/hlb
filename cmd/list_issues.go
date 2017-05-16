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
		base, err := hlb.NewCmdBase()
		etc.PanicIfErrorExist(err)

		issues, err := base.Service.GetIssues(base.Context, base.Remote.Owner, base.Remote.RepoName)
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
