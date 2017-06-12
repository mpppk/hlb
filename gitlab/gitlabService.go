package gitlab

import (
	"context"
	"errors"
	"fmt"

	"strings"

	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/service"
	"github.com/xanzy/go-gitlab"
)

type Client struct {
	RawClient   rawClient
	hostName    string
	ListOptions *gitlab.ListOptions
}

func NewClient(host *etc.ServiceConfig) (service.Client, error) {
	rawClient := newGitLabRawClient(host)
	return newClientFromRawClient(host, rawClient), nil
}

func newGitLabRawClient(host *etc.ServiceConfig) *RawClient {
	client := gitlab.NewClient(nil, host.OAuthToken)
	client.SetBaseURL(host.Protocol + "://" + host.Name + "/api/v3")
	return &RawClient{Client: client}
}

func newClientFromRawClient(host *etc.ServiceConfig, rawClient rawClient) service.Client {
	listOpt := &gitlab.ListOptions{PerPage: 100}
	return service.Client(&Client{RawClient: rawClient, hostName: host.Name, ListOptions: listOpt})
}

func (c *Client) GetIssues(ctx context.Context, owner, repo string) (serviceIssues []service.Issue, err error) {
	opt := &gitlab.ListProjectIssuesOptions{ListOptions: *c.ListOptions}
	issues, _, err := c.RawClient.GetIssues().ListProjectIssues(owner+"/"+repo, opt)

	for _, issue := range issues {
		serviceIssues = append(serviceIssues, &Issue{Issue: issue})
	}

	return serviceIssues, err
}

func (c *Client) GetPullRequests(ctx context.Context, owner, repo string) (servicePullRequests []service.PullRequest, err error) {
	opt := &gitlab.ListMergeRequestsOptions{ListOptions: *c.ListOptions}
	mergeRequests, _, err := c.RawClient.GetMergeRequests().ListMergeRequests(owner+"/"+repo, opt)

	for _, mergeRequest := range mergeRequests {
		servicePullRequests = append(servicePullRequests, &PullRequest{MergeRequest: mergeRequest})
	}

	return servicePullRequests, err
}

func (c *Client) GetRepository(ctx context.Context, owner, repo string) (service.Repository, error) {
	project, _, err := c.RawClient.GetProjects().GetProject(owner + "/" + repo)

	if err != nil {
		return nil, err
	}

	return &Repository{Project: project}, err
}

func (c *Client) GetRepositoryURL(owner, repo string) (string, error) {
	err := checkOwnerAndRepo(owner, repo)
	return fmt.Sprintf("https://%s/%s/%s", c.hostName, owner, repo), err
}

func (c *Client) GetIssuesURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositoryURL(owner, repo)
	return repoUrl + "/issues", err
}

func (c *Client) GetIssueURL(owner, repo string, id int) (string, error) {
	url, err := c.GetIssuesURL(owner, repo)
	return fmt.Sprintf("%s/%d", url, id), err
}
func (c *Client) GetPullRequestsURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositoryURL(owner, repo)
	return repoUrl + "/merge_requests", err
}

func (c *Client) GetPullRequestURL(owner, repo string, id int) (string, error) {
	url, err := c.GetPullRequestsURL(owner, repo)
	return fmt.Sprintf("%s/%d", url, id), err
}

func (c *Client) GetProjectsURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositoryURL(owner, repo)
	return repoUrl + "/boards", err
}

func (c *Client) GetProjectURL(owner, repo string, id int) (string, error) {
	// GitLab can not have multi boards
	return c.GetProjectsURL(owner, repo)
}

func (c *Client) GetMilestonesURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositoryURL(owner, repo)
	return repoUrl + "/milestones", err
}

func (c *Client) GetMilestoneURL(owner, repo string, id int) (string, error) {
	url, err := c.GetMilestonesURL(owner, repo)
	return fmt.Sprintf("%s/%d", url, id), err
}

func (c *Client) GetWikisURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositoryURL(owner, repo)
	return repoUrl + "/wikis", err
}

func (c *Client) CreateToken(ctx context.Context) (string, error) {
	return "Not Implemented Yet", nil
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
