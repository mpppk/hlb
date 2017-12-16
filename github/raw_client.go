package github

import (
	"context"

	"net/url"

	"github.com/google/go-github/github"
)

type rawClient interface {
	GetRepositories() repositoriesService
	GetPullRequests() pullRequestsService
	GetIssues() issuesService
	GetAuthorizations() authorizationsService
	SetBaseURL(baseUrl *url.URL)
}

type repositoriesService interface {
	Get(ctx context.Context, owner, repo string) (*github.Repository, *github.Response, error)
	Create(ctx context.Context, org string, repo *github.Repository) (*github.Repository, *github.Response, error)
	CreateRelease(ctx context.Context, owner, repo string, release *github.RepositoryRelease) (*github.RepositoryRelease, *github.Response, error)
}

type issuesService interface {
	ListByRepo(ctx context.Context, owner, repo string, opt *github.IssueListByRepoOptions) ([]*github.Issue, *github.Response, error)
}

type pullRequestsService interface {
	List(ctx context.Context, owner, repo string, opt *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error)
	Create(ctx context.Context, owner string, repo string, pull *github.NewPullRequest) (*github.PullRequest, *github.Response, error)
}

type authorizationsService interface {
	Create(ctx context.Context, authReq *github.AuthorizationRequest) (*github.Authorization, *github.Response, error)
	List(ctx context.Context, options *github.ListOptions) ([]*github.Authorization, *github.Response, error)
}

type RawClient struct {
	*github.Client
}

func (r *RawClient) GetRepositories() repositoriesService {
	return repositoriesService(r.Repositories)
}

func (r *RawClient) GetIssues() issuesService {
	return issuesService(r.Issues)
}

func (r *RawClient) GetPullRequests() pullRequestsService {
	return pullRequestsService(r.PullRequests)
}

func (r *RawClient) GetAuthorizations() authorizationsService {
	return authorizationsService(r.Authorizations)
}

func (r *RawClient) SetBaseURL(baseUrl *url.URL) {
	r.BaseURL = baseUrl
}
