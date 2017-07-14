package cmd

import (
	"io"
	"os/exec"

	"os"

	"strings"

	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/finder"
	"github.com/mpppk/hlb/hlblib"
	"github.com/mpppk/hlb/service"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
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
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		base, err := hlblib.NewCmdBase()
		etc.PanicIfErrorExist(err)
		sw := hlblib.ClientWrapper{Base: base}

		list := []finder.FilterStringer{}

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

		list = append(list,
			&finder.FilterableURL{URL: repoUrl, String: "*repo"},
			&finder.FilterableURL{URL: issuesUrl, String: "#issues"},
			&finder.FilterableURL{URL: pullsUrl, String: "!pullrequests"},
			&finder.FilterableURL{URL: projectsUrl, String: "projects"},
			&finder.FilterableURL{URL: milestonesUrl, String: "%milestones"},
			&finder.FilterableURL{URL: commitsUrl, String: "commits"},
			&finder.FilterableURL{URL: wikisUrl, String: "wikis"},
		)

		mycmd := exec.Command("peco")
		stdin, _ := mycmd.StdinPipe()

		for _, fstr := range list {
			io.WriteString(stdin, fstr.FilterString()+"\n")
		}

		issuesChan := make(chan []finder.FilterStringer)
		pullsChan := make(chan []finder.FilterStringer)

		go func() {
			issues, err := sw.GetIssues()
			etc.PanicIfErrorExist(err)
			fstrs := toFilterStringerFromIssues(issues)
			for _, fstr := range fstrs {
				io.WriteString(stdin, fstr.FilterString()+"\n")
			}
			issuesChan <- fstrs
		}()

		go func() {
			pulls, err := sw.GetPullRequests()
			etc.PanicIfErrorExist(err)
			fstrs := toFilterStringerFromPullRequests(pulls)
			for _, fstr := range fstrs {
				io.WriteString(stdin, fstr.FilterString()+"\n")
			}
			pullsChan <- fstrs
		}()

		out, err := mycmd.Output()

		select {
		case issues := <-issuesChan:
			list = append(list, issues...)
		default:
			close(issuesChan)
		}

		select {
		case pulls := <-pullsChan:
			list = append(list, pulls...)
		default:
			close(issuesChan)
		}

		stdin.Close()
		etc.PanicIfErrorExist(err)

		selectedStr := strings.TrimSpace(string(out))
		selectedStr = strings.Trim(selectedStr, "\n")

		for _, l := range list {
			if l.FilterString() == selectedStr {
				open.Run(l.(finder.Linker).GetURL())
				os.Exit(0)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(ibrowseCmd)

}
