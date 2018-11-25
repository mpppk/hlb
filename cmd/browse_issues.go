package cmd

import (
	"strconv"

	"github.com/pkg/errors"

	"github.com/mpppk/hlb/hlblib"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

func NewCmdBrowseIssues(cmdContextFunc func() (*hlblib.CmdContext, error)) *cobra.Command {
	cmd := &cobra.Command{
		Args:  MaximumNumArgs(1),
		Use:   "issues",
		Short: "browse issues",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdContext, err := cmdContextFunc()
			if err != nil {
				return errors.Wrap(err, "failed to get command context")
			}
			var url string
			if len(args) == 0 {
				u, err := cmdContext.Client.GetIssues().GetIssuesURL(cmdContext.Remote.Owner, cmdContext.Remote.RepoName)
				if err != nil {
					return errors.Wrap(err, "failed to get issues URL for browse from: "+url)
				}
				url = u
			} else {
				id, _ := strconv.Atoi(args[0]) // never return err because it already checked by args validator

				u, err := cmdContext.Client.GetIssues().GetURL(cmdContext.Remote.Owner, cmdContext.Remote.RepoName, id)
				if err != nil {
					return errors.Wrap(err, "failed to get issue URL for browse from: "+url)
				}
				url = u
			}

			if urlFlag {
				cmd.Println(url)
			} else {
				if err := open.Run(url); err != nil {
					return errors.Wrap(err, "failed to open repository URL: "+url)
				}
			}
			return nil
		},
	}
	return cmd
}

func init() {
	browseCmd.AddCommand(NewCmdBrowseIssues(hlblib.NewCmdContext))
}
