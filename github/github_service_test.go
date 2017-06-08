package github

import (
	"context"

	"net/url"
	"testing"

	"github.com/google/go-github/github"
	"github.com/mpppk/hlb/etc"
)

type MockRepositoriesService struct {
}

func (m *MockRepositoriesService) Get(ctx context.Context, owner, repo string) (*github.Repository, *github.Response, error) {
	return &github.Repository{HTMLURL: github.String("https://github.com/samplerepo")}, nil, nil
}

type MockIssuesService struct{}

func (m *MockIssuesService) ListByRepo(ctx context.Context, owner, repo string, opt *github.IssueListByRepoOptions) ([]*github.Issue, *github.Response, error) {
	return []*github.Issue{
		{
			Number:           github.Int(1),
			Title:            github.String("Test Issue"),
			HTMLURL:          github.String("https://github.com/issues/1"),
			PullRequestLinks: nil,
		},
		{
			Number:           github.Int(2),
			Title:            github.String("Test Pull Request"),
			HTMLURL:          github.String("https://github.com/pulls/2"),
			PullRequestLinks: &github.PullRequestLinks{},
		},
	}, nil, nil
}

type MockPullRequestsService struct{}

func (m *MockPullRequestsService) List(ctx context.Context, owner, repo string, opt *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error) {
	return []*github.PullRequest{
		{
			Number:  github.Int(1),
			Title:   github.String("Test Pull Request"),
			HTMLURL: github.String("https://github.com/pulls/1"),
		},
		{
			Number:  github.Int(2),
			Title:   github.String("Other Pull Request"),
			HTMLURL: github.String("https://sample.com/pulls/2"),
		},
	}, nil, nil
}

type MockAuthorizationsService struct{}

func (m *MockAuthorizationsService) Create(ctx context.Context, authReq *github.AuthorizationRequest) (*github.Authorization, *github.Response, error) {
	return &github.Authorization{
		Token: github.String("sample token"),
		Note:  github.String("sample note"),
	}, nil, nil
}

func (m *MockAuthorizationsService) List(ctx context.Context, options *github.ListOptions) ([]*github.Authorization, *github.Response, error) {
	return []*github.Authorization{
		{
			Token: github.String("sampletoken"),
			Note:  github.String("sample note"),
		},
	}, nil, nil
}

type MockRawClient struct {
	Repositories   *MockRepositoriesService
	Issues         *MockIssuesService
	PullRequests   *MockPullRequestsService
	Authorizations *MockAuthorizationsService
	BaseURL        *url.URL
}

func (m *MockRawClient) GetRepositories() repositoriesService {
	return m.Repositories
}

func (m *MockRawClient) GetIssues() issuesService {
	return issuesService(m.Issues)
}

func (m *MockRawClient) GetPullRequests() pullRequestsService {
	return pullRequestsService(m.PullRequests)
}

func (m *MockRawClient) GetAuthorizations() authorizationsService {
	return authorizationsService(m.Authorizations)
}

func (m *MockRawClient) SetBaseURL(baseUrl *url.URL) {
	m.BaseURL = baseUrl
}

func newMockRawClient() *MockRawClient {
	baseURL, _ := url.Parse("https://github.com")
	return &MockRawClient{
		Repositories:   &MockRepositoriesService{},
		Issues:         &MockIssuesService{},
		PullRequests:   &MockPullRequestsService{},
		Authorizations: &MockAuthorizationsService{},
		BaseURL:        baseURL,
	}
}

func TestNewServiceFromClient(t *testing.T) {
	serviceConfig := &etc.ServiceConfig{
		Name:       "github.com",
		Type:       "github",
		OAuthToken: "testtoken",
		Protocol:   "https",
	}

	client, err := newServiceFromClient(serviceConfig, newMockRawClient())

	if err != nil {
		t.Errorf("client creating failed: %v", err)
	}

	repoURL, err := client.GetRepositoryURL("testuser", "testrepo")
	if err != nil {
		t.Errorf("GetRepositoryURL return error: %v", err)
	}

	if repoURL != "https://github.com/testuser/testrepo" {
		t.Errorf("invalid repository URL: %v", repoURL)
	}

}
