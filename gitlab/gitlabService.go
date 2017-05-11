package gitlab

import (
	"context"
	"github.com/xanzy/go-gitlab"
	"errors"
	"github.com/mpppk/hlb/project"
)

type Service struct {
	Client *gitlab.Client
	ListOptions *gitlab.ListOptions
}

func NewService(token string, baseUrlStrs ...string) (project.Service, error) {
	if len(baseUrlStrs) > 1 {
		return nil, errors.New("too many base urls")
	}

	client := gitlab.NewClient(nil, token)

	if len(baseUrlStrs) == 1 {
		client.SetBaseURL(baseUrlStrs[0])
	}

	listOpt := &gitlab.ListOptions{PerPage: 100}
	return project.Service(&Service{Client: client, ListOptions: listOpt}), nil
}

func (s *Service) GetIssues(ctx context.Context, owner, repo string) (issues []project.Issue, err error) {
	opt := &gitlab.ListProjectIssuesOptions{ListOptions: *s.ListOptions}
	gitLabIssues, _, err := s.Client.Issues.ListProjectIssues(owner + "/" + repo, opt)

	for _, gitLabIssue := range gitLabIssues {
		issues = append(issues, &Issue{Issue: gitLabIssue})
	}

	return issues, err
}

func (s *Service) GetPullRequests(ctx context.Context, owner, repo string) (pullRequests []project.PullRequest, err error) {
	opt := &gitlab.ListMergeRequestsOptions{ListOptions: *s.ListOptions}
	gitLabMergeRequests, _, err := s.Client.MergeRequests.ListMergeRequests(owner + "/" + repo, opt)

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
