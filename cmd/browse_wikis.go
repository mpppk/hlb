package cmd

import (
	"fmt"

	"github.com/mpppk/hlb/etc"
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
		etc.PanicIfErrorExist(err)
		sw := hlblib.ClientWrapper{Base: base}

		url, err := sw.GetWikisURL()
		etc.PanicIfErrorExist(err)
		open.Run(url)
	},
}

func init() {
	browseCmd.AddCommand(browsewikisCmd)
}