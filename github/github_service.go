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

const (
	VALIDATION_FAILED_MSG = "Validation Failed"
	NO_COMMITS_MSG_PREFIX = "No commits between"
	CODE_INVALID          = "invalid"
)

type Client struct {
	rawClient   RawClient
	host        string
	ListOptions *github.ListOptions
}

func (c *Client) GetRepositories() service.RepositoriesService {
	return service.RepositoriesService(&repositoriesService{
		raw: c.rawClient.GetRepositories(),
		host: c.host,
	})
}

func (c *Client) GetIssues() service.IssuesService {
	return service.IssuesService(&issuesService{
		raw: c.rawClient.GetIssues(),
		repositoriesService: c.GetRepositories(),
		ListOptions: c.ListOptions,
	})
}

func (c *Client) GetPullRequests() service.PullRequestsService {
	return service.PullRequestsService(&pullRequestsService{
		raw: c.rawClient.GetPullRequests(),
		repositoriesService: c.GetRepositories(),
		ListOptions: c.ListOptions,
	})
}

func (c *Client) GetAuthorizations() service.AuthorizationsService {
	return service.AuthorizationsService(&authorizationsService{
		raw: c.rawClient.GetAuthorizations(),
	})
}

type ClientBuilder struct {}

func (cb *ClientBuilder) New(ctx context.Context, serviceConfig *etc.ServiceConfig) (service.Client, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: serviceConfig.Token})
	tc := oauth2.NewClient(ctx, ts)
	return newServiceFromClient(serviceConfig, &rawClient{Client: github.NewClient(tc)})
}

func (cb *ClientBuilder) NewViaBasicAuth(ctx context.Context, serviceConfig *etc.ServiceConfig, user, pass string) (service.Client, error) {
	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace(user),
		Password: strings.TrimSpace(pass),
	}
	return newServiceFromClient(serviceConfig, &rawClient{Client: github.NewClient(tp.Client())})
}

func (cb *ClientBuilder) GetType() string {
	return "github"
}

func newServiceFromClient(serviceConfig *etc.ServiceConfig, client RawClient) (service.Client, error) {
	urlStr := serviceConfig.Protocol + "://api." + serviceConfig.Host
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}

	baseUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse base URL based on serviceConfig in github.Client.newServiceFromClient")
	}
	client.SetBaseURL(baseUrl)
	listOpt := &github.ListOptions{PerPage: 100}
	return service.Client(&Client{rawClient: client, host: serviceConfig.Host, ListOptions: listOpt}), nil
}

func (c *Client) GetProjectsURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositories().GetURL(owner, repo)
	return repoUrl + "/projects", errors.Wrap(err, "Error occurred in github.Client.GetProjectsURL")
}

func (c *Client) GetProjectURL(owner, repo string, id int) (string, error) {
	repoUrl, err := c.GetRepositories().GetURL(owner, repo)
	return fmt.Sprintf("%s/projects/%d", repoUrl, id), errors.Wrap(err, "Error occurred in github.Client.GetProjectURL")
}

func (c *Client) GetMilestonesURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositories().GetURL(owner, repo)
	return repoUrl + "/milestones", errors.Wrap(err, "Error occurred in github.Client.GetMilestonesURL")
}

func (c *Client) GetMilestoneURL(owner, repo string, id int) (string, error) {
	repoUrl, err := c.GetRepositories().GetURL(owner, repo)
	return fmt.Sprintf("%s/milestone/%d", repoUrl, id), errors.Wrap(err, "Error occurred in github.Client.GetMilestoneURL")
}

func (c *Client) GetWikisURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositories().GetURL(owner, repo)
	return repoUrl + "/wiki", errors.Wrap(err, "Error occurred in github.Client.GetWikisURL")
}

func (c *Client) GetCommitsURL(owner, repo string) (string, error) {
	repoUrl, err := c.GetRepositories().GetURL(owner, repo)
	return repoUrl + "/commits", errors.Wrap(err, "Error occurred in github.Client.GetCommitsURL")
}

func (c *Client) CreateRelease(ctx context.Context, owner, repo string, newRelease *service.NewRelease) (service.Release, error) {
	newGHRelease := &github.RepositoryRelease{
		TagName: github.String(newRelease.GetTagName()),
		Name:    github.String(newRelease.GetName()),
		Body:    github.String(newRelease.GetBody()),
	}

	createdRelease, _, err := c.rawClient.GetRepositories().CreateRelease(ctx, owner, repo, newGHRelease)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Issues by raw client in github.Client.CreateRelease")
	}
	return createdRelease, nil
}


func hasAuthNote(auths []*github.Authorization, note string) bool {
	for _, a := range auths {
		if a.Note != nil && note == *a.Note {
			return true
		}
	}
	return false
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
