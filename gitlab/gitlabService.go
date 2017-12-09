package gitlab

import (
	"context"
	"fmt"

	"strings"

	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/service"
	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

type Client struct {
	RawClient     rawClient
	host          string
	serviceConfig *etc.ServiceConfig
	ListOptions   *gitlab.ListOptions
}

func NewClient(serviceConfig *etc.ServiceConfig) (service.Client, error) {
	rawClient := newGitLabRawClient(serviceConfig)
	return newClientFromRawClient(serviceConfig, rawClient), nil
}

func newGitLabRawClient(serviceConfig *etc.ServiceConfig) *RawClient {
	client := gitlab.NewClient(nil, serviceConfig.Token)
	client.SetBaseURL(serviceConfig.Protocol + "://" + serviceConfig.Host + "/api/v3")
	return &RawClient{Client: client}
}

func newClientFromRawClient(serviceConfig *etc.ServiceConfig, rawClient rawClient) service.Client {
	listOpt := &gitlab.ListOptions{PerPage: 100}
	return service.Client(&Client{RawClient: rawClient, serviceConfig: serviceConfig, host: serviceConfig.Host, ListOptions: listOpt})
}

func (c *Client) GetIssues(ctx context.Context, owner, repo string) (serviceIssues []service.Issue, err error) {
	opt := &gitlab.ListProjectIssuesOptions{ListOptions: *c.ListOptions}
	issues, _, err := c.RawClient.GetIssues().ListProjectIssues(owner+"/"+repo, opt)

	for _, issue := range issues {
		serviceIssues = append(serviceIssues, &Issue{Issue: issue})
	}

	return serviceIssues, errors.Wrap(err, "Failed to get Issues by raw client in gitlab.Client.GetIssues")
}

func (c *Client) GetPullRequests(ctx context.Context, owner, repo string) (servicePullRequests []service.PullRequest, err error) {
	opt := &gitlab.ListMergeRequestsOptions{ListOptions: *c.ListOptions}
	mergeRequests, _, err := c.RawClient.GetMergeRequests().ListMergeRequests(opt)

	for _, mergeRequest := range mergeRequests {
		servicePullRequests = append(servicePullRequests, &PullRequest{MergeRequest: mergeRequest})
	}

	return servicePullRequests, errors.Wrap(err, "Failed to get Pull Requests by raw client in gitlab.Client.GetPullRequests")
}

func (c *Client) GetRepository(ctx context.Context, owner, repo string) (service.Repository, error) {
	project, _, err := c.RawClient.GetProjects().GetProject(owner + "/" + repo)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Repository by raw client in gitlab.Client.GetRepository")
	}

	return &Repository{Project: project}, err
}

func (c *Client) GetRepositoryURL(owner, repo string) (string, error) {
	err := checkOwnerAndRepo(owner, repo)
	return fmt.Sprintf("%s://%s/%s/%s", c.serviceConfig.Protocol, c.host, owner, repo), errors.Wrap(err, "Error occurred in gitlab.Client.GetRepositoryURL")
}

func (c *Client) GetIssuesURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositoryURL(owner, repo)
	return repoUrl + "/issues", errors.Wrap(err, "Error occurred in gitlab.Client.GetIssuesURL")
}

func (c *Client) GetIssueURL(owner, repo string, id int) (string, error) {
	url, err := c.GetIssuesURL(owner, repo)
	return fmt.Sprintf("%s/%d", url, id), errors.Wrap(err, "Error occurred in gitlab.Client.GetIssueURL")
}
func (c *Client) GetPullRequestsURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositoryURL(owner, repo)
	return repoUrl + "/merge_requests", errors.Wrap(err, "Error occurred in gitlab.Client.GetPullRequestsURL")
}

func (c *Client) GetPullRequestURL(owner, repo string, id int) (string, error) {
	url, err := c.GetPullRequestsURL(owner, repo)
	return fmt.Sprintf("%s/%d", url, id), errors.Wrap(err, "Error occurred in gitlab.Client.GetPUllRequestURL")
}

func (c *Client) GetProjectsURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositoryURL(owner, repo)
	return repoUrl + "/boards", errors.Wrap(err, "Error occurred in gitlab.Client.GetProjectsURL")
}

func (c *Client) GetProjectURL(owner, repo string, id int) (string, error) {
	// GitLab can not have multi boards
	return c.GetProjectsURL(owner, repo)
}

func (c *Client) GetMilestonesURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositoryURL(owner, repo)
	return repoUrl + "/milestones", errors.Wrap(err, "Error occurred in gitlab.Client.GetMilestonesURL")
}

func (c *Client) GetMilestoneURL(owner, repo string, id int) (string, error) {
	url, err := c.GetMilestonesURL(owner, repo)
	return fmt.Sprintf("%s/%d", url, id), errors.Wrap(err, "Error occurred in gitlab.Client.GetMilestoneURL")
}

func (c *Client) GetWikisURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositoryURL(owner, repo)
	return repoUrl + "/wikis", errors.Wrap(err, "Error occurred in gitlab.Client.GetWikisURL")
}

func (c *Client) GetCommitsURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositoryURL(owner, repo)
	return repoUrl + "/commits/master", errors.Wrap(err, "Error occurred in gitlab.Client.GetCommitsURL")
}

func (c *Client) CreateRepository(ctx context.Context, repo string) (service.Repository, error) {
	opt := &gitlab.CreateProjectOptions{Name: &repo}
	retRepository, _, err := c.RawClient.GetProjects().CreateProject(opt)
	return &Repository{retRepository}, err
}

func (c *Client) CreatePullRequest(ctx context.Context, repo string, newPR *service.NewPullRequest) (service.PullRequest, error) {
	opt := &gitlab.CreateMergeRequestOptions{
		Title:        &newPR.Title,
		Description:  &newPR.Body,
		SourceBranch: &newPR.HeadBranch,
		TargetBranch: &newPR.BaseBranch,
	}
	newMergeRequest, _, err := c.RawClient.GetMergeRequests().CreateMergeRequest(newPR.BaseOwner+"/"+repo, opt)
	return &PullRequest{MergeRequest: newMergeRequest}, err
}

func (c *Client) CreateRelease(ctx context.Context, owner, repo string, newRelease *service.NewRelease) (service.Release, error) {
	panic("Not Implemented Yet")
	//opt := &gitlab.CreateTagOptions{}
	//tag, _, err := c.RawClient.GetTags().CreateTag(owner+"/"+repo, opt)
	//return tag, err
}

func (c *Client) CreateToken(ctx context.Context) (string, error) {
	return "", errors.New("Not Implemented Yet")
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
