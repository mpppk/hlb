package cmd

import (
	"fmt"

	"strconv"

	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/hlblib"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

var browsepullrequestsCmd = &cobra.Command{
	Use:   "pull-requests",
	Short: "browse pull-requests",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Println("Too many issue IDs")
		}

		base, err := hlblib.NewCmdBase()
		etc.PanicIfErrorExist(err)
		sw := hlblib.ClientWrapper{Base: base}

		if len(args) == 0 {
			url, err := sw.GetPullRequestsURL()
			etc.PanicIfErrorExist(err)
			open.Run(url)
			return
		} else {
			id, err := strconv.Atoi(args[0])
			etc.PanicIfErrorExist(err)

			url, err := sw.GetPullRequestURL(id)
			etc.PanicIfErrorExist(err)
			open.Run(url)
			return
		}
	},
}

func init() {
	browseCmd.AddCommand(browsepullrequestsCmd)
}
