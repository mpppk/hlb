package service

import "context"

type NewPullRequest struct {
	Title      string
	Body       string
	BaseBranch string
	HeadBranch string
	BaseOwner  string
	HeadOwner  string
}

type Client interface {
	GetPullRequests(ctx context.Context, owner, repo string) ([]PullRequest, error)
	GetIssues(ctx context.Context, owner, repo string) ([]Issue, error)
	GetRepository(ctx context.Context, owner, repo string) (Repository, error)
	GetRepositoryURL(owner, repo string) (string, error)
	GetIssuesURL(owner, repo string) (string, error)
	GetIssueURL(owner, repo string, no int) (string, error)
	GetPullRequestsURL(owner, repo string) (string, error)
	GetPullRequestURL(owner, repo string, no int) (string, error)
	GetProjectsURL(owner, repo string) (string, error)
	GetProjectURL(owner, repo string, no int) (string, error)
	GetMilestonesURL(owner, repo string) (string, error)
	GetMilestoneURL(owner, repo string, no int) (string, error)
	GetWikisURL(owner, repo string) (string, error)
	GetCommitsURL(owner, repo string) (string, error)
	CreateRepository(ctx context.Context, repo string) (Repository, error)
	CreatePullRequest(ctx context.Context, repo string, opt *NewPullRequest) (PullRequest, error)
	CreateFork(ctx context.Context, owner, repo string) (Repository, error)
	CreateToken(ctx context.Context) (string, error)
}
