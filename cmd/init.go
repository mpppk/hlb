package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/mpppk/hlb/etc"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// ファイルが無ければ作成
		dirName, err := homedir.Dir()
		if err != nil {
			panic(err)
		}
		configFilePath := path.Join(dirName, ".hlb.yaml")
		if _, err := os.Stat(configFilePath); err != nil {
			hosts := []*etc.ServiceConfig{
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
			ioutil.WriteFile(configFilePath, f, 0666)
		} else {
			fmt.Println("config file already exist:", configFilePath)
		}
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
