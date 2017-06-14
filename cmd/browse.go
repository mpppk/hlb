package cmd

import (
	"fmt"

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
		serviceConfig, ok := config.FindServiceConfig(remote.ServiceHost)

		if !ok {
			fmt.Println(remote.ServiceHost, " is unknown host. Please add the service configuration to config file(~/.hlb.yaml)")
			os.Exit(1)
		}

		if serviceConfig.Token == "" {
			if !hlblib.CanCreateToken(serviceConfig.Type) {
				fmt.Println("The token of ", serviceConfig.Host, " can not create via hlb.")
				fmt.Println(serviceConfig.Host, "Please add token to config file(~/.hlb.yaml) manually.")
				os.Exit(1)
			}
			serviceUrl := serviceConfig.Protocol + "://" + serviceConfig.Host
			fmt.Println(serviceUrl)
			addServiceCmd.Run(cmd, []string{serviceConfig.Type, serviceUrl})
		}
	},
}

func init() {
	RootCmd.AddCommand(browseCmd)
}
