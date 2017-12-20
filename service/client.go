package service

import (
	"context"

	"github.com/mpppk/hlb/etc"
)

type NewPullRequest struct {
	Title      string
	Body       string
	BaseBranch string
	HeadBranch string
	BaseOwner  string
	HeadOwner  string
}

type RepositoriesService interface {
	Get(ctx context.Context, owner, repo string) (Repository, error)
	GetURL(owner, repo string) (string, error)
	GetWikisURL(owner, repo string) (string, error)
	GetMilestonesURL(owner, repo string) (string, error)
	GetMilestoneURL(owner, repo string, no int) (string, error)
	GetCommitsURL(owner, repo string) (string, error)
	Create(ctx context.Context, repo string) (Repository, error)
	CreateRelease(ctx context.Context, owner, repo string, newRelease *NewRelease) (Release, error)
}

type IssuesService interface {
	ListByRepo(ctx context.Context, owner, repo string) ([]Issue, error)
	GetIssuesURL(owner, repo string) (string, error)
	GetURL(owner, repo string, no int) (string, error)
}

type PullRequestsService interface {
	List(ctx context.Context, owner, repo string) ([]PullRequest, error)
	Create(ctx context.Context, repo string, pull *NewPullRequest) (PullRequest, error)
	GetPullRequestsURL(owner, repo string) (string, error)
	GetURL(owner, repo string, no int) (string, error)
}

type AuthorizationsService interface {
	CreateToken(ctx context.Context) (string, error)
}

type ProjectsService interface {
	GetProjectsURL(owner, repo string) (string, error)
	GetURL(owner, repo string, no int) (string, error)
}

type Client interface {
	GetRepositories() RepositoriesService
	GetPullRequests() PullRequestsService
	GetIssues() IssuesService
	GetProjects() ProjectsService
	GetAuthorizations() AuthorizationsService
}

type ClientGenerator interface {
	New(ctx context.Context, serviceConfig *etc.ServiceConfig) (Client, error)
	NewViaBasicAuth(ctx context.Context, serviceConfig *etc.ServiceConfig, username, pass string) (Client, error)
	GetType() string
}