package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/mpppk/hlb/finder"
	"github.com/mpppk/hlb/etc"
)

var ibrowseCmd = &cobra.Command{
	Use:   "ibrowse",
	Short: "Browse interactive",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {

		strs := []finder.FilterableString{
			"hoge",
			"fuga",
			"piyo",
		}

		str, err := finder.FindFromFilterableStrings(strs)
		etc.PanicIfErrorExist(err)

		fmt.Println(str)
	},
}

func init() {
	RootCmd.AddCommand(ibrowseCmd)

}
