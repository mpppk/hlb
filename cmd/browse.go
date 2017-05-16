package cmd

import (
	"github.com/spf13/cobra"
	"github.com/mpppk/hlb/etc"
	"github.com/skratchdot/open-golang/open"
	"github.com/mpppk/hlb/hlb"
)

// browseCmd represents the browse command
var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "browse repo",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		base, err := hlb.NewCmdBase()
		etc.PanicIfErrorExist(err)
		url, err := base.Service.GetRepositoryURL(base.Remote.Owner, base.Remote.RepoName)
		etc.PanicIfErrorExist(err)
		open.Run(url)
	},
}

func init() {
	RootCmd.AddCommand(browseCmd)
}
