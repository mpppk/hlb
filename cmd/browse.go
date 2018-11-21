package cmd

import (
	"fmt"

	"github.com/mpppk/hlb/hlblib"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

var urlFlag bool

// browseCmd represents the browse command
var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "browse repo",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		base, err := hlblib.NewCmdBase()
		hlblib.PanicIfErrorExist(err)
		url, err := base.Client.GetRepositories().GetURL(base.Remote.Owner, base.Remote.RepoName)

		hlblib.PanicIfErrorExist(err)

		if urlFlag {
			fmt.Println(url)
		} else {
			if err := open.Run(url); err != nil {
				hlblib.PanicIfErrorExist(err)
			}
		}
	},
}

func init() {
	browseCmd.PersistentFlags().BoolVarP(&urlFlag, "url", "u", false,
		"outputs the URL rather than opening the browser")
	RootCmd.AddCommand(browseCmd)
}
