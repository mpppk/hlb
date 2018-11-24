package cmd

import (
	"strconv"

	"fmt"

	"github.com/mpppk/hlb/hlblib"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

var browseprojectsCmd = &cobra.Command{
	Aliases: []string{"boards"},
	Use:     "projects",
	Short:   "browse projects",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		base, err := hlblib.NewCmdContext()
		hlblib.PanicIfErrorExist(err)

		var url string
		if len(args) == 0 {
			u, err := base.Client.GetProjects().GetProjectsURL(base.Remote.Owner, base.Remote.RepoName)
			hlblib.PanicIfErrorExist(err)
			url = u
		} else {
			id, err := strconv.Atoi(args[0])
			hlblib.PanicIfErrorExist(err)

			u, err := base.Client.GetProjects().GetURL(base.Remote.Owner, base.Remote.RepoName, id)
			hlblib.PanicIfErrorExist(err)
			url = u
		}

		if urlFlag {
			fmt.Println(url)
		} else {
			open.Run(url)
		}
	},
}

func init() {
	browseCmd.AddCommand(browseprojectsCmd)
}
