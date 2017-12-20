package cmd

import (
	"fmt"
	"strconv"

	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/hlblib"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

var browsemilestonesCmd = &cobra.Command{
	Use:   "milestones",
	Short: "browse milestones",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Println("Too many milestone IDs")
		}

		base, err := hlblib.NewCmdBase()
		etc.PanicIfErrorExist(err)

		if len(args) == 0 {
			url, err :=  base.Client.GetRepositories().GetMilestonesURL(base.Remote.Owner, base.Remote.RepoName)

			etc.PanicIfErrorExist(err)
			open.Run(url)
			return
		} else {
			id, err := strconv.Atoi(args[0])
			etc.PanicIfErrorExist(err)

			url, err :=  base.Client.GetRepositories().GetMilestoneURL(base.Remote.Owner, base.Remote.RepoName, id)

			etc.PanicIfErrorExist(err)
			open.Run(url)
			return
		}
	},
}

func init() {
	browseCmd.AddCommand(browsemilestonesCmd)
}
