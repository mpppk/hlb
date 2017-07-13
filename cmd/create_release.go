package cmd

import (
	"fmt"

	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/hlblib"
	"github.com/mpppk/hlb/service"
	"github.com/spf13/cobra"
)

var createReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Create release page and upload files",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		base, err := hlblib.NewCmdBase()
		etc.PanicIfErrorExist(err)
		sw := hlblib.ClientWrapper{Base: base}

		newRelease := &service.NewRelease{
			TagName: "v0.0.1",
			Name:    "Test Release",
			Body:    "This is test release",
		}

		release, err := sw.CreateRelease(newRelease)
		etc.PanicIfErrorExist(err)
		fmt.Println(release.GetHTMLURL())
	},
}

func init() {
	createCmd.AddCommand(createReleaseCmd)
}
