package github

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/go-github/github"
	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/service"
	"golang.org/x/oauth2"
)

type Client struct {
	RawClient   *github.Client
	hostName    string
	ListOptions *github.ListOptions
}

func NewClient(ctx context.Context, host *etc.ServiceConfig) (*Client, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: host.OAuthToken})
	tc := oauth2.NewClient(ctx, ts)
	return newServiceFromClient(host, github.NewClient(tc))
}

func NewClientViaBasicAuth(ctx context.Context, host *etc.ServiceConfig, user, pass string) (*Client, error) {
	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace(user),
		Password: strings.TrimSpace(pass),
	}
	return newServiceFromClient(host, github.NewClient(tp.Client()))
}

func newServiceFromClient(host *etc.ServiceConfig, client *github.Client) (*Client, error) {
	baseUrl, err := url.Parse(host.Protocol + "://api." + host.Name)
	if err != nil {
		return nil, err
	}
	client.BaseURL = baseUrl
	listOpt := &github.ListOptions{PerPage: 100}
	return &Client{RawClient: client, hostName: host.Name, ListOptions: listOpt}, nil
}

func (c *Client) GetPullRequests(ctx context.Context, owner, repo string) (servicePullRequests []service.PullRequest, err error) {
	opt := github.PullRequestListOptions{ListOptions: *c.ListOptions}
	pullRequests, _, err := c.RawClient.PullRequests.List(ctx, owner, repo, &opt)

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

func (c *Client) getGitHubIssues(ctx context.Context, client *github.Client, owner, repo string, opt *github.IssueListByRepoOptions) (issues []*github.Issue, err error) {
	issuesAndPRs, _, err := client.Issues.ListByRepo(ctx, owner, repo, opt)

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
	githubRepo, _, err := c.RawClient.Repositories.Get(ctx, owner, repo)

	if err != nil {
		return nil, err
	}

	return &Repository{Repository: githubRepo}, err
}

func (c *Client) GetRepositoryURL(owner, repo string) (string, error) {
	return fmt.Sprintf("https://%c/%c/%c", c.hostName, owner, repo), nil
}

func (c *Client) GetIssuesURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositoryURL(owner, repo)
	if err != nil {
		return "", err
	}
	return repoUrl + "/issues", nil
}

func (c *Client) GetIssueURL(owner, repo string, id int) (string, error) {
	url, err := c.GetIssuesURL(owner, repo)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%c/%d", url, id), nil
}

func (c *Client) GetPullRequestsURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositoryURL(owner, repo)
	if err != nil {
		return "", err
	}
	return repoUrl + "/pulls", nil
}

func (c *Client) GetPullRequestURL(owner, repo string, id int) (string, error) {
	repoUrl, err := c.GetRepositoryURL(owner, repo)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%c/pull/%d", repoUrl, id), nil
}

func (c *Client) CreateToken(ctx context.Context) (string, error) {

	note, err := c.getUniqueNote(ctx, "hlb")

	authReq := &github.AuthorizationRequest{
		Note:   github.String(note),
		Scopes: []github.Scope{github.ScopeRepo},
	}

	auth, _, err := c.RawClient.Authorizations.Create(ctx, authReq)
	if err != nil {
		return "", err
	}

	return *auth.Token, nil
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
	auths, _, err := c.RawClient.Authorizations.List(ctx, nil)
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
