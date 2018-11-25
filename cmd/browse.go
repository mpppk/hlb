package cmd

import (
	"github.com/mpppk/hlb/hlblib"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var urlFlag bool
var browseCmd = NewCmdBrowse(hlblib.NewCmdContext)

func NewCmdBrowse(cmdContextFunc func() (*hlblib.CmdContext, error)) *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "browse",
		Short: "browse repo",
		Long:  ``,
		RunE: NewBrowseCmdFunc(cmdContextFunc, func(cmdContext *hlblib.CmdContext, args []string) (string, error) {
			url, err := cmdContext.Client.GetRepositories().GetURL(cmdContext.Remote.Owner, cmdContext.Remote.RepoName)
			return url, errors.Wrap(err, "failed to get repository commits URL for browse from: "+url)
		}),
	}

	cmd.PersistentFlags().BoolVarP(&urlFlag, "url", "u", false,
		"outputs the URL rather than opening the browser")
	return cmd
}

func init() {
	RootCmd.AddCommand(browseCmd)
}
