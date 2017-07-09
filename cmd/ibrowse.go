package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/mpppk/hlb/finder"
	"github.com/mpppk/hlb/hlblib"
	"github.com/mpppk/hlb/etc"
	"github.com/skratchdot/open-golang/open"
)

func toFilterStringerFromIssues(strs []*finder.FilterableIssue) (fstrs []finder.FilterStringer) {
	for _, fissue := range strs {
		fstrs = append(fstrs, finder.FilterStringer(fissue))
	}
	return fstrs
}

var ibrowseCmd = &cobra.Command{
	Use:   "ibrowse",
	Short: "Browse interactive",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		base, err := hlblib.NewCmdBase()
		etc.PanicIfErrorExist(err)
		sw := hlblib.ClientWrapper{Base: base}

		issues, err := sw.GetIssues()
		etc.PanicIfErrorExist(err)

		fissues := finder.ToFilterableIssues(issues)
		fstrs := toFilterStringerFromIssues(fissues)
		selectedFstrs, err := finder.Find(fstrs)
		etc.PanicIfErrorExist(err)

		ss := selectedFstrs[0]
		open.Run(ss.(finder.Linker).GetURL())
	},
}

func init() {
	RootCmd.AddCommand(ibrowseCmd)

}
