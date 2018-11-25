package cmd

import (
	"github.com/pkg/errors"

	"github.com/mpppk/hlb/hlblib"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

func NewCmdBrowseCommits(cmdContextFunc func() (*hlblib.CmdContext, error)) *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "commits",
		Short: "browse commits",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdContext, err := cmdContextFunc()
			if err != nil {
				return errors.Wrap(err, "failed to get command context")
			}
			url, err := cmdContext.Client.GetRepositories().GetCommitsURL(cmdContext.Remote.Owner, cmdContext.Remote.RepoName)
			if err != nil {
				return errors.Wrap(err, "failed to get repository commits URL for browse from: "+url)
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
	browseCmd.AddCommand(NewCmdBrowseCommits(hlblib.NewCmdContext))
}
