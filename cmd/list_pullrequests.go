package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/mpppk/hlb/hlb"
	"strconv"
	"github.com/mpppk/hlb/etc"
)

// listpullrequestsCmd represents the listpullrequests command
var listpullrequestsCmd = &cobra.Command{
	Use:   "pull-requests",
	Short: "list pull-requests",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		base, err := hlb.NewCmdBase()
		sw := hlb.ServiceWrapper{Base: base}
		etc.PanicIfErrorExist(err)

		pulls, err := sw.GetPullRequests()
		etc.PanicIfErrorExist(err)

		for _, pull := range pulls {
			info := "#" + strconv.Itoa(pull.GetNumber()) + " " + pull.GetTitle()
			fmt.Println(info)
		}
	},
}

func init() {
	listCmd.AddCommand(listpullrequestsCmd)
}
