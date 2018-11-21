package cmd

import (
	"fmt"

	"github.com/mpppk/hlb/hlblib"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

var browsewikisCmd = &cobra.Command{
	Use:   "wikis",
	Short: "browse wikis",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Println("warning: `browse wikis` does not accept any args. They are ignored.")
		}

		base, err := hlblib.NewCmdBase()
		hlblib.PanicIfErrorExist(err)
		url, err := base.Client.GetRepositories().GetWikisURL(base.Remote.Owner, base.Remote.RepoName)

		hlblib.PanicIfErrorExist(err)

		if urlFlag {
			fmt.Println(url)
		} else {
			open.Run(url)
		}
	},
}

func init() {
	browseCmd.AddCommand(browsewikisCmd)
}
