package cmd

import (
	"fmt"
	"github.com/mpppk/gitany"

	"os"

	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/hlblib"
	"github.com/spf13/cobra"
)

const (
	DEFAULT_RELEASE_FILE_NAME = "RELEASE_EDITMSG"
)

var createReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Create release page and upload files",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		base, err := hlblib.NewCmdBase()
		etc.PanicIfErrorExist(err)

		if len(args) < 1 {
			fmt.Println("Missed argument TAG")
			os.Exit(1)
		}

		title, message, err := editTitleAndMessage(DEFAULT_RELEASE_FILE_NAME, "", DEFAULT_CS)

		newRelease := &gitany.NewRelease{
			TagName: args[0],
			Name:    title,
			Body:    message,
		}

		release, _, err :=  base.Client.GetRepositories().CreateRelease(base.Context, base.Remote.Owner, base.Remote.RepoName, newRelease)

		etc.PanicIfErrorExist(err)
		fmt.Println(release.GetHTMLURL())
	},
}

func init() {
	createCmd.AddCommand(createReleaseCmd)
}
