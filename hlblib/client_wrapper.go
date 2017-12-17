package hlblib

import (
	"context"

	"github.com/mpppk/hlb/service"
)

type ClientWrapper struct {
	Base *CmdBase
}

func (c *ClientWrapper) GetRepositoryURL() (string, error) {
	return c.Base.Client.GetRepositories().GetURL(c.Base.Remote.Owner, c.Base.Remote.RepoName)
}

func (c *ClientWrapper) GetIssuesURL() (string, error) {
	return c.Base.Client.GetIssues().GetIssuesURL(c.Base.Remote.Owner, c.Base.Remote.RepoName)
}

func (c *ClientWrapper) GetIssueURL(no int) (string, error) {
	return c.Base.Client.GetIssues().GetURL(c.Base.Remote.Owner, c.Base.Remote.RepoName, no)
}

func (c *ClientWrapper) GetIssues() ([]service.Issue, error) {
	return c.Base.Client.GetIssues().ListByRepo(c.Base.Context, c.Base.Remote.Owner, c.Base.Remote.RepoName)
}

func (c *ClientWrapper) GetPullRequestsURL() (string, error) {
	return c.Base.Client.GetPullRequests().GetPullRequestsURL(c.Base.Remote.Owner, c.Base.Remote.RepoName)
}

func (c *ClientWrapper) GetPullRequestURL(no int) (string, error) {
	return c.Base.Client.GetPullRequests().GetURL(c.Base.Remote.Owner, c.Base.Remote.RepoName, no)
}

func (c *ClientWrapper) GetPullRequests() ([]service.PullRequest, error) {
	return c.Base.Client.GetPullRequests().List(c.Base.Context, c.Base.Remote.Owner, c.Base.Remote.RepoName)
}

func (c *ClientWrapper) GetProjectsURL() (string, error) {
	return c.Base.Client.GetProjectsURL(c.Base.Remote.Owner, c.Base.Remote.RepoName)
}

func (c *ClientWrapper) GetProjectURL(no int) (string, error) {
	return c.Base.Client.GetProjectURL(c.Base.Remote.Owner, c.Base.Remote.RepoName, no)
}

func (c *ClientWrapper) GetMilestonesURL() (string, error) {
	return c.Base.Client.GetMilestonesURL(c.Base.Remote.Owner, c.Base.Remote.RepoName)
}

func (c *ClientWrapper) GetMilestoneURL(no int) (string, error) {
	return c.Base.Client.GetMilestoneURL(c.Base.Remote.Owner, c.Base.Remote.RepoName, no)
}

func (c *ClientWrapper) GetWikisURL() (string, error) {
	return c.Base.Client.GetWikisURL(c.Base.Remote.Owner, c.Base.Remote.RepoName)
}

func (c *ClientWrapper) GetCommitsURL() (string, error) {
	return c.Base.Client.GetCommitsURL(c.Base.Remote.Owner, c.Base.Remote.RepoName)
}

func (c *ClientWrapper) CreateRepository(repo string) (service.Repository, error) {
	return c.Base.Client.GetRepositories().Create(c.Base.Context, repo)
}

func (c *ClientWrapper) CreatePullRequest(newPR *service.NewPullRequest) (service.PullRequest, error) {
	return c.Base.Client.GetPullRequests().Create(c.Base.Context, c.Base.Remote.RepoName, newPR)
}

func (c *ClientWrapper) CreateRelease(newRelease *service.NewRelease) (service.Release, error) {
	return c.Base.Client.CreateRelease(c.Base.Context, c.Base.Remote.Owner, c.Base.Remote.RepoName, newRelease)
}

func (c *ClientWrapper) CreateToken(ctx context.Context) (string, error) {
	return c.Base.Client.GetAuthorizations().CreateToken(ctx)
}
