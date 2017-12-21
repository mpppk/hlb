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

func (c *Client) GetProjects() service.ProjectsService {
	return service.ProjectsService(&projetsService{
		repositoriesService: c.GetRepositories(),
	})
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
