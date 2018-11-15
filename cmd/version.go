package cmd

import (
	"fmt"

	"github.com/mpppk/hlb/hlblib"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show hlb version",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v"+hlblib.Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
