package cmd

import (
	"github.com/mpppk/hlb/hlblib"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewCmdBrowseCommits(cmdContextFunc func() (*hlblib.CmdContext, error)) *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "commits",
		Short: "browse commits",
		Long:  ``,
		RunE: NewBrowseCmdFunc(cmdContextFunc, func(cmdContext *hlblib.CmdContext, args []string) (string, error) {
			url, err := cmdContext.Client.GetRepositories().GetCommitsURL(cmdContext.Remote.Owner, cmdContext.Remote.RepoName)
			return url, errors.Wrap(err, "failed to get repository commits URL for browse from: "+url)
		}),
	}
	return cmd
}

func init() {
	browseCmd.AddCommand(NewCmdBrowseCommits(hlblib.NewCmdContext))
}
