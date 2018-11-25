package cmd

import (
	"strconv"

	"github.com/pkg/errors"

	"github.com/mpppk/hlb/hlblib"
	"github.com/spf13/cobra"
)

func NewCmdBrowseProjects(cmdContextFunc func() (*hlblib.CmdContext, error)) *cobra.Command {
	cmd := &cobra.Command{
		Args:  MaximumNumArgs(1),
		Use:   "projects",
		Short: "browse projects",
		Long:  ``,
		RunE: NewBrowseCmdFunc(cmdContextFunc, func(cmdContext *hlblib.CmdContext, args []string) (string, error) {
			if len(args) == 0 {
				u, err := cmdContext.Client.GetProjects().GetProjectsURL(cmdContext.Remote.Owner, cmdContext.Remote.RepoName)
				return u, errors.Wrap(err, "failed to fetch projects URL for browse")
			} else {
				id, _ := strconv.Atoi(args[0])
				u, err := cmdContext.Client.GetProjects().GetURL(cmdContext.Remote.Owner, cmdContext.Remote.RepoName, id)
				return u, errors.Wrap(err, "failed to fetch project URL for browse")
			}
		}),
	}
	return cmd
}

func init() {
	browseCmd.AddCommand(NewCmdBrowseProjects(hlblib.NewCmdContext))
}
