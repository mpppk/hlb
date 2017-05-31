package github

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/go-github/github"
	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/project"
	"golang.org/x/oauth2"
)

type Service struct {
	Client      *github.Client
	hostName    string
	ListOptions *github.ListOptions
}

func NewService(ctx context.Context, host *etc.ServiceConfig) (*Service, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: host.OAuthToken})
	tc := oauth2.NewClient(ctx, ts)
	return newServiceFromClient(host, github.NewClient(tc))
}

func NewServiceViaBasicAuth(ctx context.Context, host *etc.ServiceConfig, user, pass string) (*Service, error) {
	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace(user),
		Password: strings.TrimSpace(pass),
	}
	return newServiceFromClient(host, github.NewClient(tp.Client()))
}

func newServiceFromClient(host *etc.ServiceConfig, client *github.Client) (*Service, error) {
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

func (s *Service) GetPullRequestsURL(owner, repo string) (string, error) {
	repoUrl, err := s.GetRepositoryURL(owner, repo)
	if err != nil {
		return "", err
	}
	return repoUrl + "/pulls", nil
}

func (s *Service) GetPullRequestURL(owner, repo string, id int) (string, error) {
	repoUrl, err := s.GetRepositoryURL(owner, repo)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/pull/%d", repoUrl, id), nil
}

func (s *Service) CreateToken(ctx context.Context) (string, error) {

	note, err := s.getUniqueNote(ctx, "hlb")

	authReq := &github.AuthorizationRequest{
		Note:   github.String(note),
		Scopes: []github.Scope{github.ScopeRepo},
	}

	auth, _, err := s.Client.Authorizations.Create(ctx, authReq)
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

func (s *Service) getUniqueNote(ctx context.Context, orgNote string) (string, error) {
	auths, _, err := s.Client.Authorizations.List(ctx, nil)
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
