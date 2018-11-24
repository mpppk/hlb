package cmd

import (
	"fmt"

	"github.com/mpppk/hlb/hlblib"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

var browsecommitsCmd = &cobra.Command{
	Use:   "commits",
	Short: "browse commits",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Println("warning: `browse commits` does not accept any args. They are ignored.")
		}

		base, err := hlblib.NewCmdContext()
		hlblib.PanicIfErrorExist(err)
		url, err := base.Client.GetRepositories().GetCommitsURL(base.Remote.Owner, base.Remote.RepoName)
		hlblib.PanicIfErrorExist(err)

		if urlFlag {
			fmt.Println(url)
		} else {
			open.Run(url)
		}
	},
}

func init() {
	browseCmd.AddCommand(browsecommitsCmd)
}
