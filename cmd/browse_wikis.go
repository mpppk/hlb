package cmd

import (
	"github.com/pkg/errors"

	"github.com/mpppk/hlb/hlblib"
	"github.com/spf13/cobra"
)

func NewCmdBrowseWikis(cmdContextFunc func() (*hlblib.CmdContext, error)) *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "wikis",
		Short: "browse wikis",
		Long:  ``,
		RunE: NewBrowseCmdFunc(cmdContextFunc, func(cmdContext *hlblib.CmdContext, args []string) (string, error) {
			url, err := cmdContext.Client.GetRepositories().GetWikisURL(cmdContext.Remote.Owner, cmdContext.Remote.RepoName)
			return url, errors.Wrap(err, "failed to get repository wikis URL for browse from: "+url)
		}),
	}
	return cmd
}

func init() {
	browseCmd.AddCommand(NewCmdBrowseWikis(hlblib.NewCmdContext))
}
