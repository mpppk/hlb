package cmd

import (
	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/git"
	"github.com/mpppk/hlb/hlblib"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var createForkCmd = &cobra.Command{
	Use:   "fork",
	Short: "Create Fork Repository",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		base, err := hlblib.NewCmdBase()
		etc.PanicIfErrorExist(err)
		sw := hlblib.ClientWrapper{Base: base}
		repo, err := sw.CreateFork()
		etc.PanicIfErrorExist(errors.Wrap(err, "Repository forking is failed in create fork command"))

		_, err = git.SetRemote(".", base.ServiceConfig.User, repo.GetGitURL())
		etc.PanicIfErrorExist(errors.Wrap(err, "Remote URL setting is failed in create fork command"))
	},
}

func init() {
	createCmd.AddCommand(createForkCmd)
}
