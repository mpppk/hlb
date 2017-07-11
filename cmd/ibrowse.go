package cmd

import (
	"github.com/spf13/cobra"
	"github.com/mpppk/hlb/finder"
	"github.com/mpppk/hlb/hlblib"
	"github.com/mpppk/hlb/etc"
	"github.com/skratchdot/open-golang/open"
	"github.com/mpppk/hlb/service"
)

func toFilterStringerFromIssues(issues []service.Issue) (fss []finder.FilterStringer) {
	for _, fi := range finder.ToFilterableIssues(issues) {
		fss = append(fss, finder.FilterStringer(fi))
	}
	return fss
}

func toFilterStringerFromPullRequests(pulls []service.PullRequest) (fss []finder.FilterStringer) {
	for _, fp := range finder.ToFilterablePullRequests(pulls) {
		fss = append(fss, finder.FilterStringer(fp))
	}
	return fss
}

var ibrowseCmd = &cobra.Command{
	Use:   "ibrowse",
	Short: "Browse interactive",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		base, err := hlblib.NewCmdBase()
		etc.PanicIfErrorExist(err)
		sw := hlblib.ClientWrapper{Base: base}

		fstrs := []finder.FilterStringer{}

		//var links []finder.Linker
		repoUrl, err := sw.GetRepositoryURL()
		etc.PanicIfErrorExist(err)
		issuesUrl, err := sw.GetIssuesURL()
		etc.PanicIfErrorExist(err)
		pullsUrl, err := sw.GetPullRequestsURL()
		etc.PanicIfErrorExist(err)
		commitsUrl, err := sw.GetCommitsURL()
		etc.PanicIfErrorExist(err)
		projectsUrl, err := sw.GetProjectsURL()
		etc.PanicIfErrorExist(err)
		milestonesUrl, err := sw.GetMilestonesURL()
		etc.PanicIfErrorExist(err)
		wikisUrl, err := sw.GetWikisURL()
		etc.PanicIfErrorExist(err)

		fstrs = append(fstrs,
			&finder.FilterableURL{URL: repoUrl, String: "*repo"},
			&finder.FilterableURL{URL: issuesUrl, String: "#issues"},
			&finder.FilterableURL{URL: pullsUrl, String: "!pullrequests"},
			&finder.FilterableURL{URL: projectsUrl, String: "projects"},
			&finder.FilterableURL{URL: milestonesUrl, String: "%milestones"},
			&finder.FilterableURL{URL: commitsUrl, String: "commits"},
			&finder.FilterableURL{URL: wikisUrl, String: "wikis"},
		)

		issues, err := sw.GetIssues()
		etc.PanicIfErrorExist(err)
		fstrs = append(fstrs, toFilterStringerFromIssues(issues)...)

		pulls, err := sw.GetPullRequests()
		etc.PanicIfErrorExist(err)
		fstrs = append(fstrs, toFilterStringerFromPullRequests(pulls)...)

		selectedFstrs, err := finder.Find(fstrs)
		etc.PanicIfErrorExist(err)

		for _, ss := range selectedFstrs {
			open.Run(ss.(finder.Linker).GetURL())
		}
	},
}

func init() {
	RootCmd.AddCommand(ibrowseCmd)

}
