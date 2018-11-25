package cmd

import (
	"strconv"

	"github.com/mpppk/hlb/hlblib"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewCmdBrowseIssues(cmdContextFunc func() (*hlblib.CmdContext, error)) *cobra.Command {
	cmd := &cobra.Command{
		Args:  MaximumNumArgs(1),
		Use:   "issues",
		Short: "browse issues",
		Long:  ``,
		RunE: NewBrowseCmdFunc(cmdContextFunc, func(cmdContext *hlblib.CmdContext, args []string) (string, error) {
			var url string
			if len(args) == 0 {
				u, err := cmdContext.Client.GetIssues().GetIssuesURL(cmdContext.Remote.Owner, cmdContext.Remote.RepoName)
				return u, errors.Wrap(err, "failed to get issues URL for browse from: "+url)
			} else {
				id, _ := strconv.Atoi(args[0]) // never return err because it already checked by args validator
				u, err := cmdContext.Client.GetIssues().GetURL(cmdContext.Remote.Owner, cmdContext.Remote.RepoName, id)
				return u, errors.Wrap(err, "failed to get issue URL for browse from: "+url)
			}
		}),
	}
	return cmd
}

func init() {
	browseCmd.AddCommand(NewCmdBrowseIssues(hlblib.NewCmdContext))
}
