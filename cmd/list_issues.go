package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/mpppk/hlb/hlblib"
	"github.com/mpppk/hlb/etc"
	"strconv"
)

// listissuesCmd represents the listissues command
var listissuesCmd = &cobra.Command{
	Use:   "issues",
	Short: "List issues",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		base, err := hlblib.NewCmdBase()
		etc.PanicIfErrorExist(err)
		issues, _, err := base.Client.GetIssues().ListByRepo(base.Context, base.Remote.Owner, base.Remote.RepoName, nil)
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
