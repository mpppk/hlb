package github

import (
	"context"

	"net/url"
	"testing"

	"fmt"

	"github.com/google/go-github/github"
	"github.com/mpppk/hlb/etc"
)

type MockRepositoriesService struct {
}

func (m *MockRepositoriesService) Get(ctx context.Context, owner, repo string) (*github.Repository, *github.Response, error) {
	return &github.Repository{HTMLURL: github.String("https://github.com/samplerepo")}, nil, nil
}

func (m *MockRepositoriesService) Create(ctx context.Context, org string, repo *github.Repository) (*github.Repository, *github.Response, error) {
	return &github.Repository{HTMLURL: github.String("https://github.com/newrepo")}, nil, nil
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

func (m *MockPullRequestsService) Create(ctx context.Context, owner string, repo string, pull *github.NewPullRequest) (*github.PullRequest, *github.Response, error) {
	prNumber := 2
	return &github.PullRequest{
		Number:  github.Int(prNumber),
		Title:   pull.Title,
		HTMLURL: github.String(fmt.Sprintf("https://github.com/%v/%v/pull/%v", owner, repo, prNumber)),
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

type Client_GetRepositoryURLTest struct {
	serviceConfig                   *etc.ServiceConfig
	rawClient                       rawClient
	willBeError                     bool
	user                            string
	repo                            string
	createRepo                      string
	issueID                         int
	pullRequestID                   int
	expectedRepositoryURL           string
	expectedIssuesURL               string
	expectedIssueURL                string
	expectedPullRequestsURL         string
	expectedPullRequestURL          string
	expectedCreatedPullRequestURL   string
	expectedCreatedPullRequestTitle string
}

type Util struct {
	i    int
	test *Client_GetRepositoryURLTest
	t    *testing.T
}

func (u *Util) printErrorIfUnexpected(err error, msg string) bool {
	ok := err == nil || u.test.willBeError
	if !ok {
		u.t.Errorf("%v: %v return error: %v", u.i, msg, err)
	}
	return ok
}

func (u *Util) assertString(actual, expected string, msg string) bool {
	ok := actual == expected
	if !ok {
		u.t.Errorf("%v: Expected %v: %v, Actual: %v",
			u.i, msg, expected, actual)
	}
	return ok
}

func TestClient_GetRepositoryURL(t *testing.T) {

	serviceConfig := &etc.ServiceConfig{
		Host:     "github.com",
		Type:     "github",
		Token:    "testtoken",
		Protocol: "https",
	}

	mockRawClient := newMockRawClient()

	client_GetRepositoryTests := []*Client_GetRepositoryURLTest{
		{
			serviceConfig: serviceConfig,
			rawClient:     mockRawClient,
			willBeError:   true,
			user:          "",
			repo:          "testrepo",
		},
		{
			serviceConfig: serviceConfig,
			rawClient:     mockRawClient,
			willBeError:   true,
			user:          "testuser",
			repo:          "",
		},
		{
			serviceConfig:                   serviceConfig,
			rawClient:                       mockRawClient,
			willBeError:                     false,
			user:                            "testuser",
			repo:                            "testrepo",
			createRepo:                      "newrepo",
			issueID:                         1,
			pullRequestID:                   1,
			expectedRepositoryURL:           "https://github.com/testuser/testrepo",
			expectedIssuesURL:               "https://github.com/testuser/testrepo/issues",
			expectedIssueURL:                "https://github.com/testuser/testrepo/issues/1",
			expectedPullRequestsURL:         "https://github.com/testuser/testrepo/pulls",
			expectedPullRequestURL:          "https://github.com/testuser/testrepo/pull/1",
			expectedCreatedPullRequestURL:   "https://github.com/testuser/newrepo/pull/2",
			expectedCreatedPullRequestTitle: "Created PullRequest",
		},
	}

	for i, test := range client_GetRepositoryTests {
		util := Util{t: t, i: i, test: test}

		client, err := newServiceFromClient(test.serviceConfig, test.rawClient)
		if ok := util.printErrorIfUnexpected(err, "client create fail"); !ok && test.willBeError {
			continue
		}

		title := "Repository URL"
		repoURL, err := client.GetRepositoryURL(test.user, test.repo)
		if ok := util.printErrorIfUnexpected(err, title); ok && err != nil {
			continue
		}
		util.assertString(repoURL, test.expectedRepositoryURL, title)

		title = "Issues URL"
		issuesURL, err := client.GetIssuesURL(test.user, test.repo)
		if ok := util.printErrorIfUnexpected(err, title); ok && err != nil {
			continue
		}
		util.assertString(issuesURL, test.expectedIssuesURL, title)

		title = "Issue URL"
		issueURL, err := client.GetIssueURL(test.user, test.repo, test.issueID)
		if ok := util.printErrorIfUnexpected(err, title); ok && err != nil {
			continue
		}
		util.assertString(issueURL, test.expectedIssueURL, title)

		title = "PullRequests URL"
		pullRequestsURL, err := client.GetPullRequestsURL(test.user, test.repo)
		if ok := util.printErrorIfUnexpected(err, title); ok && err != nil {
			continue
		}
		util.assertString(pullRequestsURL, test.expectedPullRequestsURL, title)

		title = "PullRequest URL"
		pullRequestURL, err := client.GetPullRequestURL(test.user, test.repo, test.pullRequestID)
		if ok := util.printErrorIfUnexpected(err, title); ok && err != nil {
			continue
		}
		util.assertString(pullRequestURL, test.expectedPullRequestURL, title)

		title = "Create PullRequest"
		createdPullRequest, err := client.CreatePullRequest(context.Background(), test.user, "master", test.user, "feature", test.createRepo, "New PullRequest", "Test Message")
		if ok := util.printErrorIfUnexpected(err, title); ok && err != nil {
			continue
		}
		util.assertString(createdPullRequest.GetHTMLURL(), test.expectedCreatedPullRequestURL, title)

		if test.willBeError {
			t.Errorf("%v: Error expected, params: %#v",
				i, test)
		}
	}
}
