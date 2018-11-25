package cmd

import (
	"context"
	"fmt"

	"github.com/mpppk/gitany"

	"os"

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
		base, err := hlblib.NewCmdContext()
		hlblib.PanicIfErrorExist(err)

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

		ctx := context.Background()
		release, _, err := base.Client.GetRepositories().CreateRelease(ctx, base.Remote.Owner, base.Remote.RepoName, newRelease)

		hlblib.PanicIfErrorExist(err)
		fmt.Println(release.GetHTMLURL())
	},
}

func init() {
	createCmd.AddCommand(createReleaseCmd)
}
