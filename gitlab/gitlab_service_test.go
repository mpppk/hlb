package gitlab

import (
	"testing"

	"fmt"

	"context"

	"github.com/mpppk/hlb/etc"
	"github.com/xanzy/go-gitlab"
)

type MockProjectsService struct {
}

func (m *MockProjectsService) GetProject(pid interface{}, options ...gitlab.OptionFunc) (*gitlab.Project, *gitlab.Response, error) {
	return &gitlab.Project{WebURL: "https://gitlab.com/user/samplerepo"}, nil, nil
}

func (m *MockProjectsService) CreateProject(opt *gitlab.CreateProjectOptions, options ...gitlab.OptionFunc) (*gitlab.Project, *gitlab.Response, error) {
	return &gitlab.Project{WebURL: "https://gitlab.com/user/newrepo"}, nil, nil
}

type MockIssuesService struct{}

func (m *MockIssuesService) ListProjectIssues(pid interface{}, opt *gitlab.ListProjectIssuesOptions, options ...gitlab.OptionFunc) ([]*gitlab.Issue, *gitlab.Response, error) {
	return []*gitlab.Issue{
		{
			IID:    1,
			Title:  "Test Issue",
			WebURL: "https://gitlab.com/issues/1",
		},
		{
			IID:    2,
			Title:  "Test Pull Request",
			WebURL: "https://gitlab.com/pulls/2",
		},
	}, nil, nil
}

type MockMergeRequestsService struct{}

func (m *MockMergeRequestsService) ListMergeRequests(pid interface{}, opt *gitlab.ListMergeRequestsOptions, options ...gitlab.OptionFunc) ([]*gitlab.MergeRequest, *gitlab.Response, error) {
	return []*gitlab.MergeRequest{
		{
			IID:    1,
			Title:  "Test Pull Request",
			WebURL: "https://gitlab.com/pulls/1",
		},
		{
			IID:    2,
			Title:  "Other Pull Request",
			WebURL: "https://sample.com/pulls/2",
		},
	}, nil, nil
}

func (m *MockMergeRequestsService) CreateMergeRequest(pid interface{}, opt *gitlab.CreateMergeRequestOptions, options ...gitlab.OptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error) {
	mrIID := 2
	return &gitlab.MergeRequest{
		IID:    mrIID,
		Title:  *opt.Title,
		WebURL: fmt.Sprintf("https://gitlab.com/%v/%v/merge_requests/%v", "testuser", "newrepo", mrIID),
	}, nil, nil
}

type MockRawClient struct {
	Projects      *MockProjectsService
	Issues        *MockIssuesService
	MergeRequests *MockMergeRequestsService
	BaseURL       string
}

func (m *MockRawClient) GetProjects() projectsService {
	return projectsService(m.Projects)
}

func (m *MockRawClient) GetIssues() issuesService {
	return issuesService(m.Issues)
}

func (m *MockRawClient) GetMergeRequests() mergeRequestsService {
	return mergeRequestsService(m.MergeRequests)
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
		BaseURL:       "https://gitlab.com",
	}
}

type Client_GetRepositoryURLTest struct {
	serviceConfig                     *etc.ServiceConfig
	rawClient                         rawClient
	willBeError                       bool
	user                              string
	repo                              string
	createRepo                        string
	createPRTitle                     string
	createPRMessage                   string
	issueID                           int
	pullRequestID                     int
	expectedRepositoryURL             string
	expectedIssuesURL                 string
	expectedIssueURL                  string
	expectedMergeRequestsURL          string
	expectedPullRequestURL            string
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
			serviceConfig:                     serviceConfig,
			rawClient:                         mockRawClient,
			willBeError:                       false,
			user:                              "testuser",
			repo:                              "testrepo",
			createRepo:                        "newrepo",
			createPRTitle:                     "New Pull Request",
			createPRMessage:                   "New Message",
			issueID:                           1,
			pullRequestID:                     1,
			expectedRepositoryURL:             "https://gitlab.com/testuser/testrepo",
			expectedIssuesURL:                 "https://gitlab.com/testuser/testrepo/issues",
			expectedIssueURL:                  "https://gitlab.com/testuser/testrepo/issues/1",
			expectedMergeRequestsURL:          "https://gitlab.com/testuser/testrepo/merge_requests",
			expectedPullRequestURL:            "https://gitlab.com/testuser/testrepo/merge_requests/1",
			expectedCreatedPullRequestURL:     "https://gitlab.com/testuser/newrepo/merge_requests/2",
			expectedCreatedPullRequestTitle:   "New Pull Request",
			expectedCreatedPullRequestMessage: "New Message",
		},
	}

	for i, test := range client_GetRepositoryTests {
		util := Util{t: t, i: i, test: test}

		client := newClientFromRawClient(test.serviceConfig, test.rawClient)

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

		title = "MergeRequests URL"
		pullRequestsURL, err := client.GetPullRequestsURL(test.user, test.repo)
		if ok := util.printErrorIfUnexpected(err, title); ok && err != nil {
			continue
		}
		util.assertString(pullRequestsURL, test.expectedMergeRequestsURL, title)

		title = "MergeRequest URL"
		pullRequestURL, err := client.GetPullRequestURL(test.user, test.repo, test.pullRequestID)
		if ok := util.printErrorIfUnexpected(err, title); ok && err != nil {
			continue
		}
		util.assertString(pullRequestURL, test.expectedPullRequestURL, title)

		title = "Created MergeRequest"
		createdPullRequest, err := client.CreatePullRequest(context.Background(), test.user, "master", test.user, "feature", test.createRepo, test.createPRTitle, test.createPRMessage)
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
