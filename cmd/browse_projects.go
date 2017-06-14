package cmd

import (
	"strconv"

	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/hlblib"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

var browseprojectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "browse projects",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		base, err := hlblib.NewCmdBase()
		etc.PanicIfErrorExist(err)
		sw := hlblib.ClientWrapper{Base: base}
		if len(args) == 0 {
			url, err := sw.GetProjectsURL()
			etc.PanicIfErrorExist(err)
			open.Run(url)
			return
		} else {
			id, err := strconv.Atoi(args[0])
			etc.PanicIfErrorExist(err)

			url, err := sw.GetProjectURL(id)
			etc.PanicIfErrorExist(err)
			open.Run(url)
			return
		}
	},
}

func init() {
	browseCmd.AddCommand(browseprojectsCmd)
}
