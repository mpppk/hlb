package cmd

import (
	"fmt"

	"strings"

	"os"

	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/git"
	"github.com/mpppk/hlb/hlblib"
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
		if len(args) > 0 {
			fmt.Printf(`Unknown subcommand "%v" for "browse"`, args[0])
			fmt.Println("\nDid you mean this?")

			suggestedCmdNames := cmd.SuggestionsFor(args[0])
			for _, s := range suggestedCmdNames {
				fmt.Println("\t", s)
			}
			os.Exit(1)
		}

		base, err := hlblib.NewCmdBase()
		etc.PanicIfErrorExist(err)
		sw := hlblib.ClientWrapper{Base: base}
		url, err := sw.GetRepositoryURL()
		etc.PanicIfErrorExist(err)
		fmt.Println(url)
		open.Run(url)
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var config etc.Config
		err := viper.Unmarshal(&config)
		etc.PanicIfErrorExist(err)
		remote, err := git.GetDefaultRemote(".")
		etc.PanicIfErrorExist(err)
		host, ok := config.FindHost(remote.ServiceHostName)
		if !ok {
			if remote.ServiceHostName == "github" {
				serviceUrl := remote.URL
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
