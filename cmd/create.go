package cmd

import (
	"context"
	"github.com/mpppk/gitany"
	"os"
	"path/filepath"

	"path"

	"time"

	"github.com/AlecAivazis/survey"
	"github.com/briandowns/spinner"
	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/git"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func chooseService(host string, config *etc.Config) (*gitany.ServiceConfig, error) {
	subConfig := config
	if host != "" {
		subConfig = config.FindServiceConfigs(host)
	}

	hosts := subConfig.ListServiceConfigHost()

	var qs = []*survey.Question{
		{
			Name: "serviceHost",
			Prompt: &survey.Select{
				Message: "Choose target service:",
				Options: hosts,
			},
		},
	}

	answers := struct {
		ServiceHost string `survey:"serviceHost"` // survey will match the question and field names
	}{}

	// perform the questions
	err := survey.Ask(qs, &answers)
	if err != nil {
		return nil, errors.Wrap(err, "Error occurred while the user was selecting the git service in create command")
	}
	serviceConfig, ok := config.FindServiceConfig(answers.ServiceHost)

	if !ok {
		return nil, errors.New("host name not found in config file: " + host)
	}

	return serviceConfig, nil
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new public repository",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		var config etc.Config
		err := viper.Unmarshal(&config)
		etc.PanicIfErrorExist(errors.Wrap(err, "Error occurred when unmarshal viper config"))

		host := ""
		if len(args) > 0 {
			host = args[0]
		}

		subConfig := config.FindServiceConfigs(host)
		interactiveFlag := true
		if len(subConfig.Services) == 1 {
			interactiveFlag = false
		}

		var serviceConfig *gitany.ServiceConfig
		if interactiveFlag {
			serviceConfig, err = chooseService(host, &config)
			etc.PanicIfErrorExist(errors.Wrap(err, "Error occurred while selecting the git service in create command"))
		} else {
			serviceConfig = subConfig.Services[0]
		}

		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond) // Build our new spinner
		if interactiveFlag {
			s.Start()
		}

		client, err := gitany.NewClient(ctx, serviceConfig)
		etc.PanicIfErrorExist(errors.Wrap(err, "Error occurred when client creating in create command"))

		currentDirPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
		etc.PanicIfErrorExist(errors.Wrap(err, "Retrieve current directory path is failed in create command"))
		currentDirName := path.Base(currentDirPath)

		newRepo := gitany.NewRepository(currentDirName)
		repo, _, err := client.GetRepositories().Create(ctx, "", newRepo)
		etc.PanicIfErrorExist(errors.Wrap(err, "Repository creating is failed in create command"))

		_, err = git.SetRemote(".", "origin", repo.GetCloneURL())
		etc.PanicIfErrorExist(errors.Wrap(err, "Remote URL setting is failed in create command"))
		s.Stop()
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
}
