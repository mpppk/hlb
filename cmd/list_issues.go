package cmd

import (
	"context"
	"fmt"

	"strconv"

	"github.com/mpppk/hlb/hlblib"
	"github.com/spf13/cobra"
)

// listissuesCmd represents the listissues command
var listissuesCmd = &cobra.Command{
	Use:   "issues",
	Short: "List issues",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		base, err := hlblib.NewCmdContext()
		hlblib.PanicIfErrorExist(err)

		ctx := context.Background()
		issues, _, err := base.Client.GetIssues().ListByRepo(ctx, base.Remote.Owner, base.Remote.RepoName, nil)
		hlblib.PanicIfErrorExist(err)

		for _, issue := range issues {
			info := "#" + strconv.Itoa(issue.GetNumber()) + " " + issue.GetTitle()
			fmt.Println(info)
		}
	},
}

func init() {
	listCmd.AddCommand(listissuesCmd)
}
