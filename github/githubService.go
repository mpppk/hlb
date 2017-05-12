package github

import (
	"context"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"net/url"
	"github.com/mpppk/hlb/project"
	"fmt"
	"github.com/mpppk/hlb/etc"
)

type Service struct {
	Client *github.Client
	hostName string
	ListOptions *github.ListOptions
}

func NewService(ctx context.Context, host *etc.Host) (*Service, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: host.OAuthToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	baseUrl, err := url.Parse(host.Protocol + "://api." + host.Name)
	if err != nil {
		return nil, err
	}

	client.BaseURL = baseUrl
	listOpt := &github.ListOptions{PerPage: 100}
	return &Service{Client: client, hostName: host.Name, ListOptions: listOpt}, nil
}

func (s *Service) GetPullRequests(ctx context.Context, owner, repo string) (pullRequests []project.PullRequest, err error) {
	opt := github.PullRequestListOptions{ListOptions: *s.ListOptions}
	gitHubPullRequests, _, err := s.Client.PullRequests.List(ctx, owner, repo, &opt)

	if err != nil {
		return nil, err
	}

	for _, gitHubPullRequest := range gitHubPullRequests {
		pullRequests = append(pullRequests, &PullRequest{PullRequest: gitHubPullRequest})
	}

	return pullRequests, err
}

func (s *Service) GetIssues(ctx context.Context, owner, repo string) (issues []project.Issue, err error) {
	opt := &github.IssueListByRepoOptions{ListOptions: *s.ListOptions}
	gitHubIssues, err := s.getGitHubIssues(ctx, s.Client, owner, repo, opt)

	if err != nil {
		return nil, err
	}

	for _, gitHubIssue := range gitHubIssues {
		issues = append(issues, &Issue{Issue: gitHubIssue})
	}

	return issues, err
}

func (s *Service) getGitHubIssues(ctx context.Context, client *github.Client, owner, repo string, opt *github.IssueListByRepoOptions) (issues []*github.Issue, err error) {
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

func (s *Service) GetRepository(ctx context.Context, owner, repo string) (project.Repository, error) {
	githubRepo, _, err := s.Client.Repositories.Get(ctx, owner, repo)

	if err != nil {
		return nil, err
	}

	return &Repository{Repository: githubRepo}, err
}

func (s *Service) GetRepositoryURL(owner, repo string) (string, error) {
	return fmt.Sprintf("https://%s/%s/%s", s.hostName, owner, repo), nil
}

func (s *Service) GetIssuesURL(owner, repo string) (string, error) {
	repoUrl, err := s.GetRepositoryURL(owner, repo)
	if err != nil {
		return "", err
	}
	return repoUrl + "/issues", nil
}

func (s *Service) GetIssueURL(owner, repo string, id int) (string, error) {
	url, err := s.GetIssuesURL(owner, repo)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%d", url, id), nil
}