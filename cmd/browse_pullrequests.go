package cmd

import (
	"strconv"

	"github.com/pkg/errors"

	"github.com/mpppk/hlb/hlblib"
	"github.com/spf13/cobra"
)

func NewCmdBrowsePullRequests(cmdContextFunc func() (*hlblib.CmdContext, error)) *cobra.Command {
	cmd := &cobra.Command{
		Args:  MaximumNumArgs(1),
		Use:   "pull-requests",
		Short: "browse pull-requests",
		Long:  ``,
		RunE: NewBrowseCmdFunc(cmdContextFunc, func(cmdContext *hlblib.CmdContext, args []string) (string, error) {
			if len(args) == 0 {
				u, err := cmdContext.Client.GetPullRequests().GetPullRequestsURL(cmdContext.Remote.Owner, cmdContext.Remote.RepoName)
				return u, errors.Wrap(err, "failed to fetch pull requests URL for browse")
			} else {
				id, _ := strconv.Atoi(args[0])
				u, err := cmdContext.Client.GetPullRequests().GetURL(cmdContext.Remote.Owner, cmdContext.Remote.RepoName, id)
				return u, errors.Wrap(err, "failed to fetch pull request URL for browse")
			}
		}),
	}
	return cmd
}

func init() {
	browseCmd.AddCommand(NewCmdBrowsePullRequests(hlblib.NewCmdContext))
}
