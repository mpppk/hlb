package cmd

import (
	"context"
	"fmt"
	"github.com/mpppk/gitany"
	"github.com/mpppk/hlb/hlblib"
	"io/ioutil"
	"net/url"
	"os"

	"time"

	"github.com/AlecAivazis/survey"
	"github.com/briandowns/spinner"
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

		var config hlblib.Config
		err := viper.Unmarshal(&config)
		hlblib.PanicIfErrorExist(err)

		if len(args) < 2 {
			panic("invalid args")
		}

		serviceType := args[0]
		serviceUrl := args[1]

		parsedUrl, err := url.Parse(serviceUrl)
		hlblib.PanicIfErrorExist(err)

		serviceConfig, ok := config.FindServiceConfig(parsedUrl.Host)
		if ok {
			if serviceConfig.Token != "" {
				msg := "token for " + parsedUrl.Host + " is already exist.\n"
				msg += "Are you sure to over write token?"

				replaceOAuthToken := false
				prompt := &survey.Confirm{
					Message: msg,
				}
				survey.AskOne(prompt, &replaceOAuthToken, nil)

				if !replaceOAuthToken {
					os.Exit(0)
				}

			}
		} else {
			serviceConfig = &gitany.ServiceConfig{
				Host:     parsedUrl.Host,
				Type:     serviceType,
				Token:    "",
				Protocol: parsedUrl.Scheme,
			}
		}

		username, password := gitany.PromptUserAndPassword(serviceType)

		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond) // Build our new spinner
		s.Start()                                                    // Start the spinner
		token, err := gitany.CreateToken(ctx, serviceConfig, username, password)
		hlblib.PanicIfErrorExist(err)
		serviceConfig.Token = token
		s.Stop()
		if !ok {
			fmt.Println("Add new service:", parsedUrl.Host)
			config.Services = append(config.Services, serviceConfig)
		} else {
			fmt.Println("Update service:", parsedUrl.Host)
		}

		f, err := yaml.Marshal(config)
		configFilePath, err := hlblib.GetConfigFilePath()
		hlblib.PanicIfErrorExist(err)
		ioutil.WriteFile(configFilePath, f, 0666)
	},
}

func init() {
	RootCmd.AddCommand(addServiceCmd)
}
