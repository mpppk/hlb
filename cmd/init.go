package cmd

import (
	"fmt"
	"github.com/mpppk/gitany"
	"io/ioutil"
	"os"

	"github.com/mpppk/hlb/etc"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate setting file to ~/.config/hlb",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		configFilePath, err := etc.GetConfigFilePath()
		etc.PanicIfErrorExist(err)
		if _, err := os.Stat(configFilePath); err != nil {
			hosts := []*gitany.ServiceConfig{
				{
					Host:     "github.com",
					Type:     "github",
					Token:    "",
					Protocol: "https",
				},
				{
					Host:     "gitlab.com",
					Type:     "gitlab",
					Token:    "",
					Protocol: "https",
				},
			}

			config := etc.Config{Services: hosts}
			f, err := yaml.Marshal(config)
			etc.PanicIfErrorExist(err)
			configFileDirPath, err := etc.GetConfigDirPath()
			err = os.MkdirAll(configFileDirPath, 0777)
			etc.PanicIfErrorExist(err)
			err = ioutil.WriteFile(configFilePath, f, 0666)
			etc.PanicIfErrorExist(err)
		} else {
			fmt.Println("config file already exist:", configFilePath)
		}
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
