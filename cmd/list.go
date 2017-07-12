package cmd

import (
	"fmt"

	"github.com/mpppk/hlb/hlblib"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("list called")
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		base, _ := hlblib.NewCmdBase()
		if base.ServiceConfig.Token == "" && base.ServiceConfig.Type == "github" {
			addServiceCmd.Run(nil, nil)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
