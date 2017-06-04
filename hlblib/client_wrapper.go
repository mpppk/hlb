package hlblib

import (
	"context"

	"github.com/mpppk/hlb/service"
)

type ClientWrapper struct {
	Base *CmdBase
}

func (c *ClientWrapper) GetRepositoryURL() (string, error) {
	return c.Base.Client.GetRepositoryURL(c.Base.Remote.Owner, c.Base.Remote.RepoName)
}

func (c *ClientWrapper) GetIssuesURL() (string, error) {
	return c.Base.Client.GetIssuesURL(c.Base.Remote.Owner, c.Base.Remote.RepoName)
}

func (c *ClientWrapper) GetIssueURL(no int) (string, error) {
	return c.Base.Client.GetIssueURL(c.Base.Remote.Owner, c.Base.Remote.RepoName, no)
}

func (c *ClientWrapper) GetIssues() ([]service.Issue, error) {
	return c.Base.Client.GetIssues(c.Base.Context, c.Base.Remote.Owner, c.Base.Remote.RepoName)
}

func (c *ClientWrapper) GetPullRequestsURL() (string, error) {
	return c.Base.Client.GetPullRequestsURL(c.Base.Remote.Owner, c.Base.Remote.RepoName)
}

func (c *ClientWrapper) GetPullRequestURL(no int) (string, error) {
	return c.Base.Client.GetPullRequestURL(c.Base.Remote.Owner, c.Base.Remote.RepoName, no)
}

func (c *ClientWrapper) GetPullRequests() ([]service.PullRequest, error) {
	return c.Base.Client.GetPullRequests(c.Base.Context, c.Base.Remote.Owner, c.Base.Remote.RepoName)
}

func (c *ClientWrapper) CreateToken(ctx context.Context) (string, error) {
	return c.Base.Client.CreateToken(ctx)
}
