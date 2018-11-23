package cmd

import (
	"os"

	"github.com/pkg/errors"

	"github.com/mpppk/hlb/hlblib"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

var urlFlag bool
var browseCmd = NewCmdBrowse()

func NewCmdBrowse() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "browse",
		Short: "browse repo",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			base, err := hlblib.NewCmdBase()
			hlblib.PanicIfErrorExist(err)
			url, err := base.Client.GetRepositories().GetURL(base.Remote.Owner, base.Remote.RepoName)
			if err != nil {
				return errors.Wrap(err, "failed to get repository for browse from: "+url)
			}

			if urlFlag {
				cmd.SetOutput(os.Stdout)
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
