package cmd

import (
	"github.com/pkg/errors"

	"github.com/mpppk/hlb/hlblib"
	"github.com/skratchdot/open-golang/open"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdContext, err := cmdContextFunc()
			if err != nil {
				return errors.Wrap(err, "failed to get command context")
			}

			url, err := cmdContext.Client.GetRepositories().GetURL(cmdContext.Remote.Owner, cmdContext.Remote.RepoName)
			if err != nil {
				return errors.Wrap(err, "failed to get repository for browse from: "+url)
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

	cmd.PersistentFlags().BoolVarP(&urlFlag, "url", "u", false,
		"outputs the URL rather than opening the browser")
	return cmd
}

func init() {
	RootCmd.AddCommand(browseCmd)
}
