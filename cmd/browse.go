package cmd

import (
	"github.com/spf13/cobra"
	"github.com/mpppk/hlb/git"
	"context"
	"github.com/spf13/viper"
	"github.com/mpppk/hlb/etc"
	"github.com/skratchdot/open-golang/open"
	"github.com/mpppk/hlb/hlb"
)

// browseCmd represents the browse command
var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "browse repo",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		var config etc.Config
		err := viper.Unmarshal(&config)
		etc.PanicIfErrorExist(err)

		remote, err := git.GetDefaultRemote(".")
		etc.PanicIfErrorExist(err)

		host, ok := config.FindHost(remote.ServiceHostName)
		if !ok {
			panic("host not found" + remote.ServiceHostName)
		}

		pj, err := hlb.GetService(ctx, host)
		etc.PanicIfErrorExist(err)

		repo, err := pj.GetRepository(ctx, remote.Owner, remote.RepoName)
		etc.PanicIfErrorExist(err)
		open.Run(repo.GetHTMLURL())
	},
}

func init() {
	RootCmd.AddCommand(browseCmd)
}
