package cmd

import (
	"context"
	"io/ioutil"
	"net/url"

	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/hlb"
	"github.com/mpppk/hlb/project"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// add-serviceCmd represents the add-service command
var addServiceCmd = &cobra.Command{
	Use:   "add-service",
	Short: "Add service",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		var config etc.Config
		err := viper.Unmarshal(&config)
		etc.PanicIfErrorExist(err)

		if len(args) < 2 {
			panic("invalid args")
		}

		serviceType := args[0]
		serviceUrl := args[1]

		parsedUrl, err := url.Parse(serviceUrl)
		etc.PanicIfErrorExist(err)

		username, password := project.PromptUserAndPassword(serviceType)

		host := &etc.Host{
			Name:       parsedUrl.Host,
			Type:       serviceType,
			OAuthToken: "",
			Protocol:   parsedUrl.Scheme,
		}
		token, err := hlb.CreateToken(ctx, host, username, password)
		etc.PanicIfErrorExist(err)
		host.OAuthToken = token
		config.Hosts = append(config.Hosts, host)
		f, err := yaml.Marshal(config)
		ioutil.WriteFile("test.yaml", f, 0666)
	},
}

func init() {
	RootCmd.AddCommand(addServiceCmd)
}
