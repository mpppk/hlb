package github

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/go-github/github"
	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/service"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type Client struct {
	RawClient   rawClient
	hostName    string
	ListOptions *github.ListOptions
}

func NewClient(ctx context.Context, host *etc.ServiceConfig) (*Client, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: host.OAuthToken})
	tc := oauth2.NewClient(ctx, ts)
	return newServiceFromClient(host, &RawClient{Client: github.NewClient(tc)})
}

func NewClientViaBasicAuth(ctx context.Context, host *etc.ServiceConfig, user, pass string) (*Client, error) {
	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace(user),
		Password: strings.TrimSpace(pass),
	}
	return newServiceFromClient(host, &RawClient{Client: github.NewClient(tp.Client())})
}

func newServiceFromClient(host *etc.ServiceConfig, client rawClient) (*Client, error) {
	baseUrl, err := url.Parse(host.Protocol + "://api." + host.Name)
	if err != nil {
		return nil, err
	}
	client.SetBaseURL(baseUrl)
	listOpt := &github.ListOptions{PerPage: 100}
	return &Client{RawClient: client, hostName: host.Name, ListOptions: listOpt}, nil
}

func (c *Client) GetPullRequests(ctx context.Context, owner, repo string) (servicePullRequests []service.PullRequest, err error) {
	opt := github.PullRequestListOptions{ListOptions: *c.ListOptions}
	pullRequests, _, err := c.RawClient.GetPullRequests().List(ctx, owner, repo, &opt)

	if err != nil {
		return nil, err
	}

	for _, gitHubPullRequest := range pullRequests {
		servicePullRequests = append(servicePullRequests, &PullRequest{PullRequest: gitHubPullRequest})
	}

	return servicePullRequests, err
}

func (c *Client) GetIssues(ctx context.Context, owner, repo string) (serviceIssues []service.Issue, err error) {
	opt := &github.IssueListByRepoOptions{ListOptions: *c.ListOptions}
	issues, err := c.getGitHubIssues(ctx, c.RawClient, owner, repo, opt)

	if err != nil {
		return nil, err
	}

	for _, issue := range issues {
		serviceIssues = append(serviceIssues, &Issue{Issue: issue})
	}

	return serviceIssues, err
}

func (c *Client) getGitHubIssues(ctx context.Context, client rawClient, owner, repo string, opt *github.IssueListByRepoOptions) (issues []*github.Issue, err error) {
	issuesAndPRs, _, err := client.GetIssues().ListByRepo(ctx, owner, repo, opt)

	if err != nil {
		return nil, err
	}

	for _, issueOrPR := range issuesAndPRs {
		if issueOrPR.PullRequestLinks == nil {
			issues = append(issues, issueOrPR)
		}
	}
	return issues, nil
}

func (c *Client) GetRepository(ctx context.Context, owner, repo string) (service.Repository, error) {
	githubRepo, _, err := c.RawClient.GetRepositories().Get(ctx, owner, repo)
	return &Repository{Repository: githubRepo}, err
}

func (c *Client) GetRepositoryURL(owner, repo string) (string, error) {
	return fmt.Sprintf("https://%s/%s/%s", c.hostName, owner, repo), checkOwnerAndRepo(owner, repo)
}

func (c *Client) GetIssuesURL(owner, repo string) (string, error) {
	if err := checkOwnerAndRepo(owner, repo); err != nil {
		return "", err
	}

	repoUrl, err := c.GetRepositoryURL(owner, repo)
	return repoUrl + "/issues", err
}

func (c *Client) GetIssueURL(owner, repo string, id int) (string, error) {
	url, err := c.GetIssuesURL(owner, repo)
	return fmt.Sprintf("%s/%d", url, id), err
}

func (c *Client) GetPullRequestsURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositoryURL(owner, repo)
	return repoUrl + "/pulls", err
}

func (c *Client) GetPullRequestURL(owner, repo string, id int) (string, error) {
	repoUrl, err := c.GetRepositoryURL(owner, repo)
	return fmt.Sprintf("%s/pull/%d", repoUrl, id), err
}

func (c *Client) GetProjectsURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositoryURL(owner, repo)
	return repoUrl + "/projects", err
}

func (c *Client) GetProjectURL(owner, repo string, id int) (string, error) {
	repoUrl, err := c.GetRepositoryURL(owner, repo)
	return fmt.Sprintf("%s/projects/%d", repoUrl, id), err
}

func (c *Client) GetMilestonesURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositoryURL(owner, repo)
	return repoUrl + "/milestones", err
}

func (c *Client) GetMilestoneURL(owner, repo string, id int) (string, error) {
	repoUrl, err := c.GetRepositoryURL(owner, repo)
	return fmt.Sprintf("%s/milestone/%d", repoUrl, id), err
}

func (c *Client) GetWikisURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositoryURL(owner, repo)
	return repoUrl + "/wiki", err
}

func (c *Client) GetCommitsURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositoryURL(owner, repo)
	return repoUrl + "/commits", err
}

func (c *Client) CreateToken(ctx context.Context) (string, error) {

	note, err := c.getUniqueNote(ctx, "hlb")

	authReq := &github.AuthorizationRequest{
		Note:   github.String(note),
		Scopes: []github.Scope{github.ScopeRepo},
	}

	auth, _, err := c.RawClient.GetAuthorizations().Create(ctx, authReq)
	return *auth.Token, err
}

func hasAuthNote(auths []*github.Authorization, note string) bool {
	for _, a := range auths {
		if a.Note != nil && note == *a.Note {
			return true
		}
	}
	return false
}

func (c *Client) getUniqueNote(ctx context.Context, orgNote string) (string, error) {
	auths, _, err := c.RawClient.GetAuthorizations().List(ctx, nil)
	if err != nil {
		return "", err
	}

	cnt := 1
	note := orgNote
	for {
		if !hasAuthNote(auths, note) {
			return note, nil
		}
		cnt++
		note = fmt.Sprint(orgNote, cnt)
	}
}

func checkOwnerAndRepo(owner, repo string) error {
	m := map[string]string{"owner": owner, "repo": repo}

	for strType, name := range m {
		if err := validateOwnerOrRepo(strType, name); err != nil {
			return err
		}
	}
	return nil
}

func validateOwnerOrRepo(strType, name string) error {
	if name == "" {
		return errors.New(strType + " is empty")
	}
	if strings.Contains(name, "/") {
		return errors.New(fmt.Sprintf("invalid %v: %v", strType, name))
	}
	return nil
}
