package cmd

import (
	"github.com/spf13/cobra"
	"github.com/mpppk/hlb/hlblib"
	"github.com/mpppk/hlb/etc"
	"github.com/skratchdot/open-golang/open"
	"fmt"
	"strconv"
)

// browseissuesCmd represents the browseissues command
var browseissuesCmd = &cobra.Command{
	Use:   "issues",
	Short: "browse issues",
	Long: `browse issues`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Println("Too many issue IDs")
		}

		base, err := hlblib.NewCmdBase()
		etc.PanicIfErrorExist(err)
		sw := hlblib.ClientWrapper{Base: base}

		if len(args) == 0 {
			url, err := sw.GetIssuesURL()
			etc.PanicIfErrorExist(err)
			open.Run(url)
			return
		}else {
			id, err := strconv.Atoi(args[0])
			etc.PanicIfErrorExist(err)

			url, err := sw.GetIssueURL(id)
			etc.PanicIfErrorExist(err)
			open.Run(url)
			return
		}
	},
}

func init() {
	browseCmd.AddCommand(browseissuesCmd)
}
