package cmd

import (
	"fmt"

	"strings"

	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/git"
	"github.com/mpppk/hlb/hlb"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// browseCmd represents the browse command
var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "browse repo",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		base, err := hlb.NewCmdBase()
		etc.PanicIfErrorExist(err)
		sw := hlb.ClientWrapper{Base: base}
		url, err := sw.GetRepositoryURL()
		etc.PanicIfErrorExist(err)
		fmt.Println(url)
		open.Run(url)
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var config etc.Config
		err := viper.Unmarshal(&config)
		fmt.Println("parse config")
		etc.PanicIfErrorExist(err)
		remote, err := git.GetDefaultRemote(".")
		fmt.Println("get default remote")
		etc.PanicIfErrorExist(err)
		host, ok := config.FindHost(remote.ServiceHostName)
		fmt.Println("search host")
		if !ok {
			if remote.ServiceHostName == "github" {
				serviceUrl := remote.Remote.Config().URL
				if !strings.Contains(serviceUrl, "http") {
					serviceUrl = "https://" + remote.ServiceHostName
				}

				addServiceCmd.Run(cmd, []string{"github", serviceUrl})
			}
			return
		}

		if host.OAuthToken == "" && host.Type == "github" {
			serviceUrl := host.Protocol + "://" + host.Name
			fmt.Println(serviceUrl)
			addServiceCmd.Run(cmd, []string{"github", serviceUrl})
		}
	},
}

func init() {
	RootCmd.AddCommand(browseCmd)
}
