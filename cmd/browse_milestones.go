package cmd

import (
	"strconv"

	"github.com/pkg/errors"

	"github.com/mpppk/hlb/hlblib"
	"github.com/spf13/cobra"
)

func NewCmdBrowseMilestones(cmdContextFunc func() (*hlblib.CmdContext, error)) *cobra.Command {
	cmd := &cobra.Command{
		Args:  MaximumNumArgs(1),
		Use:   "milestones",
		Short: "browse milestones",
		Long:  ``,
		RunE: NewBrowseCmdFunc(cmdContextFunc, func(cmdContext *hlblib.CmdContext, args []string) (string, error) {
			if len(args) == 0 {
				u, err := cmdContext.Client.GetRepositories().GetMilestonesURL(cmdContext.Remote.Owner, cmdContext.Remote.RepoName)
				return u, errors.Wrap(err, "failed to fetch milestones URL for browse")
			} else {
				id, _ := strconv.Atoi(args[0])
				u, err := cmdContext.Client.GetRepositories().GetMilestoneURL(cmdContext.Remote.Owner, cmdContext.Remote.RepoName, id)
				return u, errors.Wrap(err, "failed to fetch milestone URL for browse")
			}
		}),
	}
	return cmd
}

func init() {
	browseCmd.AddCommand(NewCmdBrowseMilestones(hlblib.NewCmdContext))
}
