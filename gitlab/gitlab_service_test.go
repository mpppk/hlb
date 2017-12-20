package gitlab

import (
	"testing"

	"fmt"

	"context"

	"strings"

	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/service"
	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

const (
	DEFAULT_BASE_URL           = "https://gitlab.com"
	DEFAULT_OWNER_NAME         = "testuser"
	DEFAULT_REPO_NAME          = "testrepo"
	DEFAULT_CREATED_REPO_NAME  = "newrepo"
	DEFAULT_CREATED_PR_TITLE   = "New PullRequest"
	DEFAULT_CREATED_PR_MESSAGE = "New Message"
	DEFAULT_BASE_BRANCH        = "master"
	DEFAULT_HEAD_BRANCH        = "feature"
)

func getOwnerAndRepoFromPid(pid interface{}) (string, string, error) {
	pidStr, ok := pid.(string)
	if !ok {
		return "", "", errors.New("pid is not string")
	}

	ownerAndRepo := strings.Split(pidStr, "/")
	return ownerAndRepo[0], ownerAndRepo[1], nil
}

type MockProjectsService struct {
}

func (m *MockProjectsService) GetProject(pid interface{}, options ...gitlab.OptionFunc) (*gitlab.Project, *gitlab.Response, error) {
	owner, repo, _ := getOwnerAndRepoFromPid(pid)

	return &gitlab.Project{WebURL: fmt.Sprintf("%v/%v/%v", DEFAULT_BASE_URL, owner, repo)}, nil, nil
}

func (m *MockProjectsService) CreateProject(opt *gitlab.CreateProjectOptions, options ...gitlab.OptionFunc) (*gitlab.Project, *gitlab.Response, error) {
	return &gitlab.Project{WebURL: fmt.Sprintf("%v/%v/%v", DEFAULT_BASE_URL, DEFAULT_OWNER_NAME, opt.Name)}, nil, nil
}

type MockIssuesService struct{}

func (m *MockIssuesService) ListProjectIssues(pid interface{}, opt *gitlab.ListProjectIssuesOptions, options ...gitlab.OptionFunc) ([]*gitlab.Issue, *gitlab.Response, error) {
	owner, repo, _ := getOwnerAndRepoFromPid(pid)

	return []*gitlab.Issue{
		{
			IID:    1,
			Title:  "Test Issue",
			WebURL: fmt.Sprintf("%v/%v/%v/issues/1", DEFAULT_BASE_URL, owner, repo),
		},
		{
			IID:    2,
			Title:  "Test Pull Request",
			WebURL: fmt.Sprintf("%v/%v/%v/issues/2", DEFAULT_BASE_URL, owner, repo),
		},
	}, nil, nil
}

type MockMergeRequestsService struct{}

func (m *MockMergeRequestsService) ListProjectMergeRequests(pid interface{}, opt *gitlab.ListProjectMergeRequestsOptions, options ...gitlab.OptionFunc) ([]*gitlab.MergeRequest, *gitlab.Response, error) {
	return []*gitlab.MergeRequest{
		{
			IID:    1,
			Title:  "Test Pull Request",
			WebURL: fmt.Sprintf("%v/%v/%v/merge_requests/1", DEFAULT_BASE_URL, "testowner", "testrepo"),
		},
		{
			IID:    2,
			Title:  "Other Pull Request",
			WebURL: fmt.Sprintf("%v/%v/%v/merge_requests/2", DEFAULT_BASE_URL, "testowner", "testrepo"),
		},
	}, nil, nil
}

func (m *MockMergeRequestsService) CreateMergeRequest(pid interface{}, opt *gitlab.CreateMergeRequestOptions, options ...gitlab.OptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error) {
	mrIID := 2
	owner, repo, _ := getOwnerAndRepoFromPid(pid)

	return &gitlab.MergeRequest{
		IID:    mrIID,
		Title:  *opt.Title,
		WebURL: fmt.Sprintf("%v/%v/%v/merge_requests/%v", DEFAULT_BASE_URL, owner, repo, mrIID),
	}, nil, nil
}

type MockRawClient struct {
	Projects      *MockProjectsService
	Issues        *MockIssuesService
	MergeRequests *MockMergeRequestsService
	BaseURL       string
}

func (m *MockRawClient) GetProjects() ProjectsService {
	return ProjectsService(m.Projects)
}

func (m *MockRawClient) GetIssues() IssuesService {
	return IssuesService(m.Issues)
}

func (m *MockRawClient) GetMergeRequests() MergeRequestsService {
	return MergeRequestsService(m.MergeRequests)
}

func (m *MockRawClient) SetBaseURL(baseUrl string) error {
	m.BaseURL = baseUrl
	return nil
}

func newMockRawClient() *MockRawClient {
	return &MockRawClient{
		Projects:      &MockProjectsService{},
		Issues:        &MockIssuesService{},
		MergeRequests: &MockMergeRequestsService{},
		BaseURL:       "%v",
	}
}

type Client_GetRepositoryURLTest struct {
	serviceConfig            *etc.ServiceConfig
	rawClient                RawClient
	willBeError              bool
	user                     string
	repo                     string
	createRepo               string
	createPRTitle            string
	createPRMessage          string
	issueID                  int
	pullRequestID            int
	expectedRepositoryURL    string
	expectedIssuesURL        string
	expectedIssueURL         string
	expectedMergeRequestsURL string
	expectedPullRequestURL   string
	expectedCreatedPullRequestURL     string
	expectedCreatedPullRequestTitle   string
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
		Host:     "gitlab.com",
		Type:     "gitlab",
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
			createPRTitle:                     DEFAULT_CREATED_PR_TITLE,
			createPRMessage:                   DEFAULT_CREATED_PR_MESSAGE,
			issueID:                           1,
			pullRequestID:                     1,
			expectedRepositoryURL:             fmt.Sprintf("%v/%v/%v", DEFAULT_BASE_URL, DEFAULT_OWNER_NAME, DEFAULT_REPO_NAME),
			expectedIssuesURL:                 fmt.Sprintf("%v/%v/%v/issues", DEFAULT_BASE_URL, DEFAULT_OWNER_NAME, DEFAULT_REPO_NAME),
			expectedIssueURL:                  fmt.Sprintf("%v/%v/%v/issues/1", DEFAULT_BASE_URL, DEFAULT_OWNER_NAME, DEFAULT_REPO_NAME),
			expectedMergeRequestsURL:          fmt.Sprintf("%v/%v/%v/merge_requests", DEFAULT_BASE_URL, DEFAULT_OWNER_NAME, DEFAULT_REPO_NAME),
			expectedPullRequestURL:            fmt.Sprintf("%v/%v/%v/merge_requests/1", DEFAULT_BASE_URL, DEFAULT_OWNER_NAME, DEFAULT_REPO_NAME),
			expectedCreatedPullRequestURL:     fmt.Sprintf("%v/%v/%v/merge_requests/2", DEFAULT_BASE_URL, DEFAULT_OWNER_NAME, DEFAULT_CREATED_REPO_NAME),
			expectedCreatedPullRequestTitle:   DEFAULT_CREATED_PR_TITLE,
			expectedCreatedPullRequestMessage: DEFAULT_CREATED_PR_MESSAGE,
		},
	}

	for i, test := range client_GetRepositoryTests {
		util := Util{t: t, i: i, test: test}

		client := newClientFromRawClient(test.serviceConfig, test.rawClient)

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

		title = "MergeRequests URL"
		pullRequestsURL, err := client.GetPullRequests().GetPullRequestsURL(test.user, test.repo)
		if ok := util.printErrorIfUnexpected(err, title); ok && err != nil {
			continue
		}
		util.assertString(pullRequestsURL, test.expectedMergeRequestsURL, title)

		title = "MergeRequest URL"
		pullRequestURL, err := client.GetPullRequests().GetURL(test.user, test.repo, test.pullRequestID)
		if ok := util.printErrorIfUnexpected(err, title); ok && err != nil {
			continue
		}
		util.assertString(pullRequestURL, test.expectedPullRequestURL, title)

		title = "Created MergeRequest"
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
		util.assertString(createdPullRequest.GetHTMLURL(), test.expectedCreatedPullRequestURL, title+" URL")
		util.assertString(createdPullRequest.GetTitle(), test.expectedCreatedPullRequestTitle, title+" Title")

		if test.willBeError {
			t.Errorf("%v: Error expected, params: %#v",
				i, test)
		}
	}
}
