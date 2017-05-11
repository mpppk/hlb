package gitlab

import "github.com/xanzy/go-gitlab"

type PullRequest struct {
	*gitlab.MergeRequest
}

func (pullRequest *PullRequest) GetNumber() int {
	return pullRequest.IID
}

func (pullRequest *PullRequest) GetTitle() string {
	return pullRequest.Title
}

func (pullRequest *PullRequest) GetHTMLURL() string {
	return pullRequest.WebURL
}

