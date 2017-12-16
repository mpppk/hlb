package github

import (
	"context"

	"net/url"

	"github.com/google/go-github/github"
)

type RawClient interface {
	GetRepositories() RepositoriesService
	GetPullRequests() PullRequestsService
	GetIssues() IssuesService
	GetAuthorizations() AuthorizationsService
	SetBaseURL(baseUrl *url.URL)
}

type RepositoriesService interface {
	Get(ctx context.Context, owner, repo string) (*github.Repository, *github.Response, error)
	Create(ctx context.Context, org string, repo *github.Repository) (*github.Repository, *github.Response, error)
	CreateRelease(ctx context.Context, owner, repo string, release *github.RepositoryRelease) (*github.RepositoryRelease, *github.Response, error)
}

type IssuesService interface {
	ListByRepo(ctx context.Context, owner, repo string, opt *github.IssueListByRepoOptions) ([]*github.Issue, *github.Response, error)
}

type PullRequestsService interface {
	List(ctx context.Context, owner, repo string, opt *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error)
	Create(ctx context.Context, owner string, repo string, pull *github.NewPullRequest) (*github.PullRequest, *github.Response, error)
}

type AuthorizationsService interface {
	Create(ctx context.Context, authReq *github.AuthorizationRequest) (*github.Authorization, *github.Response, error)
	List(ctx context.Context, options *github.ListOptions) ([]*github.Authorization, *github.Response, error)
}

type rawClient struct {
	*github.Client
}

func (r *rawClient) GetRepositories() RepositoriesService {
	return RepositoriesService(r.Repositories)
}

func (r *rawClient) GetIssues() IssuesService {
	return IssuesService(r.Issues)
}

func (r *rawClient) GetPullRequests() PullRequestsService {
	return PullRequestsService(r.PullRequests)
}

func (r *rawClient) GetAuthorizations() AuthorizationsService {
	return AuthorizationsService(r.Authorizations)
}

func (r *rawClient) SetBaseURL(baseUrl *url.URL) {
	r.BaseURL = baseUrl
}
