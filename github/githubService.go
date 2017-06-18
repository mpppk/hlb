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
	host        string
	ListOptions *github.ListOptions
}

func NewClient(ctx context.Context, serviceConfig *etc.ServiceConfig) (*Client, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: serviceConfig.Token})
	tc := oauth2.NewClient(ctx, ts)
	return newServiceFromClient(serviceConfig, &RawClient{Client: github.NewClient(tc)})
}

func NewClientViaBasicAuth(ctx context.Context, serviceConfig *etc.ServiceConfig, user, pass string) (*Client, error) {
	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace(user),
		Password: strings.TrimSpace(pass),
	}
	return newServiceFromClient(serviceConfig, &RawClient{Client: github.NewClient(tp.Client())})
}

func newServiceFromClient(serviceConfig *etc.ServiceConfig, client rawClient) (*Client, error) {
	baseUrl, err := url.Parse(serviceConfig.Protocol + "://api." + serviceConfig.Host)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse base URL based on serviceConfig in github.Client.newServiceFromClient")
	}
	client.SetBaseURL(baseUrl)
	listOpt := &github.ListOptions{PerPage: 100}
	return &Client{RawClient: client, host: serviceConfig.Host, ListOptions: listOpt}, nil
}

func (c *Client) GetPullRequests(ctx context.Context, owner, repo string) (servicePullRequests []service.PullRequest, err error) {
	opt := github.PullRequestListOptions{ListOptions: *c.ListOptions}
	pullRequests, _, err := c.RawClient.GetPullRequests().List(ctx, owner, repo, &opt)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Pull Requests by raw client in github.Client.GetPullRequests")
	}

	for _, pullRequest := range pullRequests {
		servicePullRequests = append(servicePullRequests, &PullRequest{PullRequest: pullRequest})
	}

	return servicePullRequests, nil
}

func (c *Client) GetIssues(ctx context.Context, owner, repo string) (serviceIssues []service.Issue, err error) {
	opt := &github.IssueListByRepoOptions{ListOptions: *c.ListOptions}
	issues, err := c.getGitHubIssues(ctx, c.RawClient, owner, repo, opt)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Issues by raw client in github.Client.GetIssues")
	}

	for _, issue := range issues {
		serviceIssues = append(serviceIssues, &Issue{Issue: issue})
	}

	return serviceIssues, errors.Wrap(err, "Error occurred in github.Client.GetIssues")
}

func (c *Client) getGitHubIssues(ctx context.Context, client rawClient, owner, repo string, opt *github.IssueListByRepoOptions) (issues []*github.Issue, err error) {
	issuesAndPRs, _, err := client.GetIssues().ListByRepo(ctx, owner, repo, opt)

	if err != nil {
		return nil, errors.Wrap(err, "Error occurred in github.Client.getGitHubIssues")
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
	return &Repository{Repository: githubRepo}, errors.Wrap(err, "Failed to get Repositories by raw client")
}

func (c *Client) GetRepositoryURL(owner, repo string) (string, error) {
	return fmt.Sprintf("https://%s/%s/%s", c.host, owner, repo), checkOwnerAndRepo(owner, repo)
}

func (c *Client) GetIssuesURL(owner, repo string) (string, error) {
	if err := checkOwnerAndRepo(owner, repo); err != nil {
		return "", errors.Wrap(err, "Invalid owner or repo was passed to GetIssuesURL")
	}

	repoUrl, err := c.GetRepositoryURL(owner, repo)
	return repoUrl + "/issues", errors.Wrap(err, "Error occurred in github.Client.GetIssuesURL")
}

func (c *Client) GetIssueURL(owner, repo string, id int) (string, error) {
	url, err := c.GetIssuesURL(owner, repo)
	return fmt.Sprintf("%s/%d", url, id), errors.Wrap(err, "Error occurred in github.Client.GetIssueURL")
}

func (c *Client) GetPullRequestsURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositoryURL(owner, repo)
	return repoUrl + "/pulls", errors.Wrap(err, "Error occurred in github.Client.GetPullRequestsURL")
}

func (c *Client) GetPullRequestURL(owner, repo string, id int) (string, error) {
	repoUrl, err := c.GetRepositoryURL(owner, repo)
	return fmt.Sprintf("%s/pull/%d", repoUrl, id), errors.Wrap(err, "Error occurred in github.Client.GetPullRequestURL")
}

func (c *Client) GetProjectsURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositoryURL(owner, repo)
	return repoUrl + "/projects", errors.Wrap(err, "Error occurred in github.Client.GetProjectsURL")
}

func (c *Client) GetProjectURL(owner, repo string, id int) (string, error) {
	repoUrl, err := c.GetRepositoryURL(owner, repo)
	return fmt.Sprintf("%s/projects/%d", repoUrl, id), errors.Wrap(err, "Error occurred in github.Client.GetProjectURL")
}

func (c *Client) GetMilestonesURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositoryURL(owner, repo)
	return repoUrl + "/milestones", errors.Wrap(err, "Error occurred in github.Client.GetMilestonesURL")
}

func (c *Client) GetMilestoneURL(owner, repo string, id int) (string, error) {
	repoUrl, err := c.GetRepositoryURL(owner, repo)
	return fmt.Sprintf("%s/milestone/%d", repoUrl, id), errors.Wrap(err, "Error occurred in github.Client.GetMilestoneURL")
}

func (c *Client) GetWikisURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositoryURL(owner, repo)
	return repoUrl + "/wiki", errors.Wrap(err, "Error occurred in github.Client.GetWikisURL")
}

func (c *Client) GetCommitsURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositoryURL(owner, repo)
	return repoUrl + "/commits", errors.Wrap(err, "Error occurred in github.Client.GetCommitsURL")
}

func (c *Client) CreateRepository(ctx context.Context, repo string) (service.Repository, error) {
	repository := &github.Repository{Name: github.String(repo)}
	retRepository, _, err := c.RawClient.GetRepositories().Create(ctx, "", repository)
	return &Repository{retRepository}, err
}

func (c *Client) CreateToken(ctx context.Context) (string, error) {

	note, err := c.getUniqueNote(ctx, "hlb")

	authReq := &github.AuthorizationRequest{
		Note:   github.String(note),
		Scopes: []github.Scope{github.ScopeRepo},
	}

	auth, _, err := c.RawClient.GetAuthorizations().Create(ctx, authReq)
	return *auth.Token, errors.Wrap(err, "Failed to get authorizations by raw client in github.Client.CreateToken")
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
		return "", errors.Wrap(err, "Failed to get authorizations by raw client in github.Client.GetPullRequestsURL")
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
