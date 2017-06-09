package gitlab

import (
	"context"
	"fmt"

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
	client := gitlab.NewClient(nil, host.OAuthToken)
	client.SetBaseURL(host.Protocol + "://" + host.Name + "/api/v3")
	listOpt := &gitlab.ListOptions{PerPage: 100}
	rawClient := &RawClient{Client: client}
	return service.Client(&Client{RawClient: rawClient, hostName: host.Name, ListOptions: listOpt}), nil
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
	return fmt.Sprintf("https://%s/%s/%s", c.hostName, owner, repo), nil
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
	return fmt.Sprintf("%s/%d", url, id), nil
}
func (c *Client) GetPullRequestsURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositoryURL(owner, repo)
	if err != nil {
		return "", err
	}
	return repoUrl + "/merge_requests", nil
}

func (c *Client) GetPullRequestURL(owner, repo string, id int) (string, error) {
	url, err := c.GetPullRequestsURL(owner, repo)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%d", url, id), nil
}

func (c *Client) CreateToken(ctx context.Context) (string, error) {
	return "Not Implemented Yet", nil
}
