package gitlab

import "github.com/xanzy/go-gitlab"

type Issue struct {
	*gitlab.Issue
}

func (issue *Issue) GetNumber() int {
	return issue.IID
}

func (issue *Issue) GetTitle() string {
	return issue.Title
}

func (issue *Issue) GetHTMLURL() string {
	return issue.WebURL
}

