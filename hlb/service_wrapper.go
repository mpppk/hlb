package hlb

import (
	"context"

	"github.com/mpppk/hlb/service"
)

type ServiceWrapper struct {
	Base *CmdBase
}

func (s *ServiceWrapper) GetRepositoryURL() (string, error) {
	return s.Base.Service.GetRepositoryURL(s.Base.Remote.Owner, s.Base.Remote.RepoName)
}

func (s *ServiceWrapper) GetIssuesURL() (string, error) {
	return s.Base.Service.GetIssuesURL(s.Base.Remote.Owner, s.Base.Remote.RepoName)
}

func (s *ServiceWrapper) GetIssueURL(no int) (string, error) {
	return s.Base.Service.GetIssueURL(s.Base.Remote.Owner, s.Base.Remote.RepoName, no)
}

func (s *ServiceWrapper) GetIssues() ([]service.Issue, error) {
	return s.Base.Service.GetIssues(s.Base.Context, s.Base.Remote.Owner, s.Base.Remote.RepoName)
}

func (s *ServiceWrapper) GetPullRequestsURL() (string, error) {
	return s.Base.Service.GetPullRequestsURL(s.Base.Remote.Owner, s.Base.Remote.RepoName)
}

func (s *ServiceWrapper) GetPullRequestURL(no int) (string, error) {
	return s.Base.Service.GetPullRequestURL(s.Base.Remote.Owner, s.Base.Remote.RepoName, no)
}

func (s *ServiceWrapper) GetPullRequests() ([]service.PullRequest, error) {
	return s.Base.Service.GetPullRequests(s.Base.Context, s.Base.Remote.Owner, s.Base.Remote.RepoName)
}

func (s *ServiceWrapper) CreateToken(ctx context.Context) (string, error) {
	return s.Base.Service.CreateToken(ctx)
}
