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
			hosts := []*etc.Host{
				{
					Name:       "github.com",
					Type:       "github",
					OAuthToken: "",
					Protocol:   "https",
				},
				{
					Name:       "gitlab.com",
					Type:       "gitlab",
					OAuthToken: "",
					Protocol:   "https",
				},
			}

			config := etc.Config{Hosts: hosts}
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
