package github

import (
	"context"

	"net/url"
	"testing"

	"fmt"

	"github.com/google/go-github/github"
	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/service"
)

const (
	DEFAULT_BASE_URL           = "https://github.com"
	DEFAULT_OWNER_NAME         = "testuser"
	DEFAULT_REPO_NAME          = "testrepo"
	DEFAULT_CREATED_REPO_NAME  = "newrepo"
	DEFAULT_CREATED_PR_TITLE   = "New PullRequest"
	DEFAULT_CREATED_PR_MESSAGE = "New Message"
	DEFAULT_BASE_BRANCH        = "master"
	DEFAULT_HEAD_BRANCH        = "feature"
)

type MockRepositoriesService struct{}

func (m *MockRepositoriesService) Get(ctx context.Context, owner, repo string) (*github.Repository, *github.Response, error) {
	return &github.Repository{HTMLURL: github.String(fmt.Sprintf("%v/%v/%v", DEFAULT_BASE_URL, owner, repo))}, nil, nil
}

func (m *MockRepositoriesService) Create(ctx context.Context, org string, repo *github.Repository) (*github.Repository, *github.Response, error) {
	return &github.Repository{HTMLURL: github.String(fmt.Sprintf("%v/%v/%v", DEFAULT_BASE_URL, DEFAULT_OWNER_NAME, repo))}, nil, nil
}

func (m *MockRepositoriesService) CreateRelease(ctx context.Context, owner, repo string, release *github.RepositoryRelease) (*github.RepositoryRelease, *github.Response, error) {
	return &github.RepositoryRelease{}, nil, nil
}

type MockIssuesService struct{}

func (m *MockIssuesService) ListByRepo(ctx context.Context, owner, repo string, opt *github.IssueListByRepoOptions) ([]*github.Issue, *github.Response, error) {
	return []*github.Issue{
		{
			Number:           github.Int(1),
			Title:            github.String("Test Issue"),
			HTMLURL:          github.String(fmt.Sprintf("%v/%v/%v/issues/1", DEFAULT_BASE_URL, owner, repo)),
			PullRequestLinks: nil,
		},
		{
			Number:           github.Int(2),
			Title:            github.String("Test Pull Request"),
			HTMLURL:          github.String(fmt.Sprintf("%v/%v/%v/pulls/1", DEFAULT_BASE_URL, owner, repo)),
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
			HTMLURL: github.String(fmt.Sprintf("%v/%v/%v/pulls/1", DEFAULT_BASE_URL, owner, repo)),
		},
		{
			Number:  github.Int(2),
			Title:   github.String("Other Pull Request"),
			HTMLURL: github.String(fmt.Sprintf("%v/%v/%v/pulls/2", DEFAULT_BASE_URL, owner, repo)),
		},
	}, nil, nil
}

func (m *MockPullRequestsService) Create(ctx context.Context, owner string, repo string, pull *github.NewPullRequest) (*github.PullRequest, *github.Response, error) {
	prNumber := 2
	return &github.PullRequest{
		Number:  github.Int(prNumber),
		Title:   pull.Title,
		HTMLURL: github.String(fmt.Sprintf("%v/%v/%v/pull/%v", DEFAULT_BASE_URL, owner, repo, prNumber)),
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

func (m *MockRawClient) GetRepositories() RepositoriesService {
	return m.Repositories
}

func (m *MockRawClient) GetIssues() IssuesService {
	return IssuesService(m.Issues)
}

func (m *MockRawClient) GetPullRequests() PullRequestsService {
	return PullRequestsService(m.PullRequests)
}

func (m *MockRawClient) GetAuthorizations() AuthorizationsService {
	return AuthorizationsService(m.Authorizations)
}

func (m *MockRawClient) SetBaseURL(baseUrl *url.URL) {
	m.BaseURL = baseUrl
}

func newMockRawClient() *MockRawClient {
	baseURL, _ := url.Parse(DEFAULT_BASE_URL)
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
	rawClient                       RawClient
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
	expectedCreatedPullRequestMessage string
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
			repo:          DEFAULT_REPO_NAME,
		},
		{
			serviceConfig: serviceConfig,
			rawClient:     mockRawClient,
			willBeError:   true,
			user:          DEFAULT_OWNER_NAME,
			repo:          "",
		},
		{
			serviceConfig:                     serviceConfig,
			rawClient:                         mockRawClient,
			willBeError:                       false,
			user:                              DEFAULT_OWNER_NAME,
			repo:                              DEFAULT_REPO_NAME,
			createRepo:                        DEFAULT_CREATED_REPO_NAME,
			issueID:                           1,
			pullRequestID:                     1,
			expectedRepositoryURL:             fmt.Sprintf("%v/%v/%v", DEFAULT_BASE_URL, DEFAULT_OWNER_NAME, DEFAULT_REPO_NAME),
			expectedIssuesURL:                 fmt.Sprintf("%v/%v/%v/issues", DEFAULT_BASE_URL, DEFAULT_OWNER_NAME, DEFAULT_REPO_NAME),
			expectedIssueURL:                  fmt.Sprintf("%v/%v/%v/issues/1", DEFAULT_BASE_URL, DEFAULT_OWNER_NAME, DEFAULT_REPO_NAME),
			expectedPullRequestsURL:           fmt.Sprintf("%v/%v/%v/pulls", DEFAULT_BASE_URL, DEFAULT_OWNER_NAME, DEFAULT_REPO_NAME),
			expectedPullRequestURL:            fmt.Sprintf("%v/%v/%v/pull/1", DEFAULT_BASE_URL, DEFAULT_OWNER_NAME, DEFAULT_REPO_NAME),
			expectedCreatedPullRequestURL:     fmt.Sprintf("%v/%v/%v/pull/2", DEFAULT_BASE_URL, DEFAULT_OWNER_NAME, DEFAULT_CREATED_REPO_NAME),
			expectedCreatedPullRequestTitle:   DEFAULT_CREATED_PR_TITLE,
			expectedCreatedPullRequestMessage: DEFAULT_CREATED_PR_MESSAGE,
		},
	}

	for i, test := range client_GetRepositoryTests {
		util := Util{t: t, i: i, test: test}

		client, err := newServiceFromClient(test.serviceConfig, test.rawClient)
		if ok := util.printErrorIfUnexpected(err, "client create fail"); !ok && test.willBeError {
			continue
		}

		title := "Repository URL"
		repoURL, err := client.GetRepositories().GetURL(test.user, test.repo)
		if ok := util.printErrorIfUnexpected(err, title); ok && err != nil {
			continue
		}
		util.assertString(repoURL, test.expectedRepositoryURL, title)

		title = "Issues URL"
		issuesURL, err := client.GetIssues().GetIssuesURL(test.user, test.repo)
		if ok := util.printErrorIfUnexpected(err, title); ok && err != nil {
			continue
		}
		util.assertString(issuesURL, test.expectedIssuesURL, title)

		title = "Issue URL"
		issueURL, err := client.GetIssues().GetURL(test.user, test.repo, test.issueID)
		if ok := util.printErrorIfUnexpected(err, title); ok && err != nil {
			continue
		}
		util.assertString(issueURL, test.expectedIssueURL, title)

		title = "PullRequests URL"
		pullRequestsURL, err := client.GetPullRequests().GetPullRequestsURL(test.user, test.repo)
		if ok := util.printErrorIfUnexpected(err, title); ok && err != nil {
			continue
		}
		util.assertString(pullRequestsURL, test.expectedPullRequestsURL, title)

		title = "PullRequest URL"
		pullRequestURL, err := client.GetPullRequests().GetURL(test.user, test.repo, test.pullRequestID)
		if ok := util.printErrorIfUnexpected(err, title); ok && err != nil {
			continue
		}
		util.assertString(pullRequestURL, test.expectedPullRequestURL, title)

		title = "Create PullRequest"
		newPR := &service.NewPullRequest{
			BaseOwner:  test.user,
			BaseBranch: DEFAULT_BASE_BRANCH,
			HeadOwner:  test.user,
			HeadBranch: DEFAULT_HEAD_BRANCH,
			Title:      DEFAULT_CREATED_PR_TITLE,
			Body:       DEFAULT_CREATED_PR_MESSAGE,
		}
		createdPullRequest, err := client.GetPullRequests().Create(context.Background(), test.createRepo, newPR)
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
