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
	rawClient     RawClient
	host          string
	serviceConfig *etc.ServiceConfig
	ListOptions   *gitlab.ListOptions
}

type ClientBuilder struct {}

func (c *Client) GetRepositories() service.RepositoriesService {
	return service.RepositoriesService(&repositoriesService{
		raw: c.rawClient.GetProjects(),
		host: c.host,
		serviceConfig: c.serviceConfig,
		})
}

func (c *Client) GetIssues() service.IssuesService {
	return service.IssuesService(&issuesService{
		raw: c.rawClient.GetIssues(),
		projectsService: c.GetRepositories(),
		ListOptions: c.ListOptions,
	})
}

func (c *Client) GetPullRequests() service.PullRequestsService {
	return service.PullRequestsService(&pullRequestsService{
		raw:                 c.rawClient.GetMergeRequests(),
		repositoriesService: c.GetRepositories(),
		ListOptions:         c.ListOptions,
	})
}

func (c *Client) GetAuthorizations() service.AuthorizationsService {
	return service.AuthorizationsService(&authorizationsService{})
}

func (cb *ClientBuilder) New(ctx context.Context, serviceConfig *etc.ServiceConfig) (service.Client, error) {
	rawClient := newGitLabRawClient(serviceConfig)
	return newClientFromRawClient(serviceConfig, rawClient), nil
}

func (cb *ClientBuilder) NewViaBasicAuth(ctx context.Context, serviceConfig *etc.ServiceConfig, user, pass string) (service.Client, error) {
	panic("gitlab.ClientBuilder.NewViaBasicAuth is not implemented yet")
}

func (cb *ClientBuilder) GetType() string {
	return "gitlab"
}

func newGitLabRawClient(serviceConfig *etc.ServiceConfig) *rawClient {
	client := gitlab.NewClient(nil, serviceConfig.Token)
	client.SetBaseURL(serviceConfig.Protocol + "://" + serviceConfig.Host + "/api/v3")
	return &rawClient{Client: client}
}

func newClientFromRawClient(serviceConfig *etc.ServiceConfig, rawClient RawClient) service.Client {
	listOpt := &gitlab.ListOptions{PerPage: 100}
	return &Client{rawClient: rawClient, serviceConfig: serviceConfig, host: serviceConfig.Host, ListOptions: listOpt}
}

func (c *Client) GetProjectsURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositories().GetURL(owner, repo)
	return repoUrl + "/boards", errors.Wrap(err, "Error occurred in gitlab.Client.GetProjectsURL")
}

func (c *Client) GetProjectURL(owner, repo string, id int) (string, error) {
	// GitLab can not have multi boards
	return c.GetProjectsURL(owner, repo)
}

func (c *Client) GetMilestonesURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositories().GetURL(owner, repo)
	return repoUrl + "/milestones", errors.Wrap(err, "Error occurred in gitlab.Client.GetMilestonesURL")
}

func (c *Client) GetMilestoneURL(owner, repo string, id int) (string, error) {
	url, err := c.GetMilestonesURL(owner, repo)
	return fmt.Sprintf("%s/%d", url, id), errors.Wrap(err, "Error occurred in gitlab.Client.GetMilestoneURL")
}

func (c *Client) GetWikisURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositories().GetURL(owner, repo)
	return repoUrl + "/wikis", errors.Wrap(err, "Error occurred in gitlab.Client.GetWikisURL")
}

func (c *Client) GetCommitsURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositories().GetURL(owner, repo)
	return repoUrl + "/commits/master", errors.Wrap(err, "Error occurred in gitlab.Client.GetCommitsURL")
}

func (c *Client) CreateRelease(ctx context.Context, owner, repo string, newRelease *service.NewRelease) (service.Release, error) {
	panic("Not Implemented Yet")
	//opt := &gitlab.CreateTagOptions{}
	//tag, _, err := c.rawClient.GetTags().CreateTag(owner+"/"+repo, opt)
	//return tag, err
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
