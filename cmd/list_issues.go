package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/mpppk/hlb/hlb"
	"github.com/mpppk/hlb/etc"
	"strconv"
)

// listissuesCmd represents the listissues command
var listissuesCmd = &cobra.Command{
	Use:   "issues",
	Short: "list issuses",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		base, err := hlb.NewCmdBase()
		etc.PanicIfErrorExist(err)
		sw := hlb.ServiceWrapper{Base: base}

		issues, err := sw.GetIssues()
		etc.PanicIfErrorExist(err)

		for _, issue := range issues {
			info := "#" + strconv.Itoa(issue.GetNumber()) + " " + issue.GetTitle()
			fmt.Println(info)
		}
	},
}

func init() {
	listCmd.AddCommand(listissuesCmd)
}
