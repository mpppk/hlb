package service

import (
	"context"
)

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
	CreateToken(ctx context.Context) (string, error)
}
