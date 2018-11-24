package cmd

import (
	"context"
	"io"
	"os/exec"

	"github.com/mpppk/gitany"

	"os"

	"strings"

	"github.com/mpppk/hlb/finder"
	"github.com/mpppk/hlb/hlblib"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

func toFilterStringerFromIssues(issues []gitany.Issue) (fss []finder.FilterStringer) {
	for _, fi := range finder.ToFilterableIssues(issues) {
		fss = append(fss, finder.FilterStringer(fi))
	}
	return fss
}

func toFilterStringerFromPullRequests(pulls []gitany.PullRequest) (fss []finder.FilterStringer) {
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
		base, err := hlblib.NewCmdContext()
		hlblib.PanicIfErrorExist(err)

		var list []finder.FilterStringer

		repoUrl, err := base.Client.GetRepositories().GetURL(base.Remote.Owner, base.Remote.RepoName)
		hlblib.PanicIfErrorExist(err)
		issuesUrl, err := base.Client.GetIssues().GetIssuesURL(base.Remote.Owner, base.Remote.RepoName)
		hlblib.PanicIfErrorExist(err)
		pullsUrl, err := base.Client.GetPullRequests().GetPullRequestsURL(base.Remote.Owner, base.Remote.RepoName)
		hlblib.PanicIfErrorExist(err)
		commitsUrl, err := base.Client.GetRepositories().GetCommitsURL(base.Remote.Owner, base.Remote.RepoName)
		hlblib.PanicIfErrorExist(err)
		projectsUrl, err := base.Client.GetProjects().GetProjectsURL(base.Remote.Owner, base.Remote.RepoName)
		hlblib.PanicIfErrorExist(err)
		milestonesUrl, err := base.Client.GetRepositories().GetMilestonesURL(base.Remote.Owner, base.Remote.RepoName)
		hlblib.PanicIfErrorExist(err)
		wikisUrl, err := base.Client.GetRepositories().GetWikisURL(base.Remote.Owner, base.Remote.RepoName)
		hlblib.PanicIfErrorExist(err)

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

		ctx := context.Background()
		go func() {
			issues, _, err := base.Client.GetIssues().ListByRepo(
				ctx,
				base.Remote.Owner,
				base.Remote.RepoName,
				nil,
			)
			hlblib.PanicIfErrorExist(err)
			fstrs := toFilterStringerFromIssues(issues)
			for _, fstr := range fstrs {
				io.WriteString(stdin, fstr.FilterString()+"\n")
			}
			issuesChan <- fstrs
		}()

		go func() {
			pulls, err := base.Client.GetPullRequests().List(ctx, base.Remote.Owner, base.Remote.RepoName)
			hlblib.PanicIfErrorExist(err)
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
			close(pullsChan)
		}

		stdin.Close()
		hlblib.PanicIfErrorExist(err)

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
