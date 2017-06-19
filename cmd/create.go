package cmd

import (
	"context"
	"os"
	"path/filepath"

	"path"

	"github.com/AlecAivazis/survey"
	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/git"
	"github.com/mpppk/hlb/hlblib"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new public repository",
	Long:  `Sample: hlb create github/gitlab`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		var config etc.Config
		err := viper.Unmarshal(&config)
		etc.PanicIfErrorExist(errors.Wrap(err, "Error occurred when unmarshal viper config"))

		subConfig := &config
		if len(args) > 0 {
			subConfig = config.FindServiceConfigs(args[0])
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
		err = survey.Ask(qs, &answers)
		etc.PanicIfErrorExist(errors.Wrap(err, "Error occurred while the user was selecting the git service in create command"))

		serviceConfig, _ := config.FindServiceConfig(answers.ServiceHost)

		client, err := hlblib.GetClient(ctx, serviceConfig)
		etc.PanicIfErrorExist(errors.Wrap(err, "Error occurred when client creating in create command"))

		currentDirPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
		etc.PanicIfErrorExist(errors.Wrap(err, "Retrieve current directory path is failed in create command"))
		currentDirName := path.Base(currentDirPath)

		repo, err := client.CreateRepository(ctx, currentDirName)
		etc.PanicIfErrorExist(errors.Wrap(err, "Repository creating is failed in create command"))

		// service.repositoryにgitのURLを取得するAPIを追加する
		_, err = git.SetRemote(".", "origin", repo.GetGitURL())
		etc.PanicIfErrorExist(errors.Wrap(err, "Remote URL setting is failed in create command"))
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
}
