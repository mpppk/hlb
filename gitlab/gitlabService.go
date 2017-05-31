package gitlab

import (
	"context"
	"fmt"

	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/project"
	"github.com/xanzy/go-gitlab"
)

type Service struct {
	Client      *gitlab.Client
	hostName    string
	ListOptions *gitlab.ListOptions
}

func NewService(host *etc.ServiceConfig) (project.Service, error) {
	client := gitlab.NewClient(nil, host.OAuthToken)
	client.SetBaseURL(host.Protocol + "://" + host.Name + "/api/v3")
	listOpt := &gitlab.ListOptions{PerPage: 100}
	return project.Service(&Service{Client: client, hostName: host.Name, ListOptions: listOpt}), nil
}

func (s *Service) GetIssues(ctx context.Context, owner, repo string) (issues []project.Issue, err error) {
	opt := &gitlab.ListProjectIssuesOptions{ListOptions: *s.ListOptions}
	gitLabIssues, _, err := s.Client.Issues.ListProjectIssues(owner+"/"+repo, opt)

	for _, gitLabIssue := range gitLabIssues {
		issues = append(issues, &Issue{Issue: gitLabIssue})
	}

	return issues, err
}

func (s *Service) GetPullRequests(ctx context.Context, owner, repo string) (pullRequests []project.PullRequest, err error) {
	opt := &gitlab.ListMergeRequestsOptions{ListOptions: *s.ListOptions}
	gitLabMergeRequests, _, err := s.Client.MergeRequests.ListMergeRequests(owner+"/"+repo, opt)

	for _, gitLabMergeRequest := range gitLabMergeRequests {
		pullRequests = append(pullRequests, &PullRequest{MergeRequest: gitLabMergeRequest})
	}

	return pullRequests, err
}

func (s *Service) GetRepository(ctx context.Context, owner, repo string) (project.Repository, error) {
	gitLabProject, _, err := s.Client.Projects.GetProject(owner + "/" + repo)

	if err != nil {
		return nil, err
	}

	return &Repository{Project: gitLabProject}, err
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
	return repoUrl + "/merge_requests", nil
}

func (s *Service) GetPullRequestURL(owner, repo string, id int) (string, error) {
	url, err := s.GetPullRequestsURL(owner, repo)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%d", url, id), nil
}

func (s *Service) CreateToken(ctx context.Context) (string, error) {
	return "Not Implemented Yet", nil
}
