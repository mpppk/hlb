package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/mpppk/hlb/hlb"
	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/git"
	"context"
	"github.com/skratchdot/open-golang/open"
)

// browseissuesCmd represents the browseissues command
var browseissuesCmd = &cobra.Command{
	Use:   "issues",
	Short: "browse issues",
	Long: `browse issues`,
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

		url, err := pj.GetIssuesURL(remote.Owner, remote.RepoName)
		etc.PanicIfErrorExist(err)
		open.Run(url)
	},
}

func init() {
	browseCmd.AddCommand(browseissuesCmd)
}
