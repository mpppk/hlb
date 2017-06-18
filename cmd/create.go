package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"path"

	"github.com/AlecAivazis/survey"
	"github.com/mpppk/hlb/etc"
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
		fmt.Println("in create")
		ctx := context.Background()
		var config etc.Config
		err := viper.Unmarshal(&config)
		etc.PanicIfErrorExist(errors.Wrap(err, "Error occurred when unmarshal viper config"))

		subConfig := config.FindServiceConfigs(args[0])
		fmt.Println("sub config:", subConfig)
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
		etc.PanicIfErrorExist(err)

		fmt.Println("answers:", answers)

		serviceConfig, _ := config.FindServiceConfig(answers.ServiceHost)

		client, err := hlblib.GetClient(ctx, serviceConfig)
		etc.PanicIfErrorExist(err)

		currentDirPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
		etc.PanicIfErrorExist(err)
		currentDirName := path.Base(currentDirPath)
		fmt.Println(currentDirName)

		repo, err := client.CreateRepository(ctx, currentDirName)
		etc.PanicIfErrorExist(err)
		fmt.Println(repo)
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
}
