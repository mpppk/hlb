package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mpppk/gitany"
	"github.com/mpppk/hlb/hlblib"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate setting file to ~/.config/hlb",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		configFilePath, err := hlblib.GetConfigFilePath()
		hlblib.PanicIfErrorExist(err)
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

			config := hlblib.Config{Services: hosts}
			f, err := yaml.Marshal(config)
			hlblib.PanicIfErrorExist(err)
			configFileDirPath, err := hlblib.GetConfigDirPath()
			err = os.MkdirAll(configFileDirPath, 0777)
			hlblib.PanicIfErrorExist(err)
			err = ioutil.WriteFile(configFilePath, f, 0666)
			hlblib.PanicIfErrorExist(err)
		} else {
			fmt.Println("config file already exist:", configFilePath)
		}
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
