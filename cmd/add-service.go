package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	"time"

	"github.com/AlecAivazis/survey"
	"github.com/briandowns/spinner"
	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/hlblib"
	"github.com/mpppk/hlb/service"
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

		if !hlblib.CanCreateToken(serviceType) {
			fmt.Println("Unsupported service type: ", serviceType)
			filePath, _ := etc.GetConfigFilePath()

			fmt.Println("Please add the service configuration to config file(" + filePath + ") manually")
			os.Exit(1)
		}

		parsedUrl, err := url.Parse(serviceUrl)
		etc.PanicIfErrorExist(err)

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
			serviceConfig = &etc.ServiceConfig{
				Host:     parsedUrl.Host,
				Type:     serviceType,
				Token:    "",
				Protocol: parsedUrl.Scheme,
			}
		}

		username, password := service.PromptUserAndPassword(serviceType)
		serviceConfig.User = username

		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond) // Build our new spinner
		s.Start()                                                    // Start the spinner
		token, err := hlblib.CreateToken(ctx, serviceConfig, username, password)
		etc.PanicIfErrorExist(err)
		serviceConfig.Token = token
		s.Stop()
		if !ok {
			fmt.Println("Add new service:", parsedUrl.Host)
			config.Services = append(config.Services, serviceConfig)
		} else {
			fmt.Println("Update service:", parsedUrl.Host)
		}

		f, err := yaml.Marshal(config)
		configFilePath, err := etc.GetConfigFilePath()
		etc.PanicIfErrorExist(err)
		ioutil.WriteFile(configFilePath, f, 0666)
	},
}

func init() {
	RootCmd.AddCommand(addServiceCmd)
}
