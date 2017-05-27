package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	"time"

	"github.com/AlecAivazis/survey"
	"github.com/briandowns/spinner"
	"github.com/mitchellh/go-homedir"
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

		host, ok := config.FindHost(parsedUrl.Host)
		if ok {
			if host.OAuthToken != "" {
				msg := "oauth token for " + parsedUrl.Host + " is already exist.\n"
				msg += "Are you sure to over write oauth token?"

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
			host = &etc.Host{
				Name:       parsedUrl.Host,
				Type:       serviceType,
				OAuthToken: "",
				Protocol:   parsedUrl.Scheme,
			}
		}

		username, password := project.PromptUserAndPassword(serviceType)

		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond) // Build our new spinner
		s.Start()                                                    // Start the spinner
		token, err := hlb.CreateToken(ctx, host, username, password)
		etc.PanicIfErrorExist(err)
		host.OAuthToken = token
		s.Stop()
		if !ok {
			fmt.Println("Add new service:", parsedUrl.Host)
			config.Hosts = append(config.Hosts, host)
		} else {
			fmt.Println("Update service:", parsedUrl.Host)
		}

		f, err := yaml.Marshal(config)
		homeDir, err := homedir.Dir()
		etc.PanicIfErrorExist(err)
		configFilePath := filepath.Join(homeDir, ".hlb.yaml")
		ioutil.WriteFile(configFilePath, f, 0666)
	},
}

func init() {
	RootCmd.AddCommand(addServiceCmd)
}
