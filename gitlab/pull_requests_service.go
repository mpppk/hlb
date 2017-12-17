package gitlab

import (
	"context"
	"fmt"

	"github.com/mpppk/hlb/service"
	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

type pullRequestsService struct {
	repositoriesService service.RepositoriesService
	raw                 MergeRequestsService
	ListOptions         *gitlab.ListOptions
}

func (p *pullRequestsService) List(ctx context.Context, owner, repo string) (servicePullRequests []service.PullRequest, err error) {
	opt := &gitlab.ListProjectMergeRequestsOptions{ListOptions: *p.ListOptions}
	mergeRequests, _, err := p.raw.ListProjectMergeRequests(owner+"/"+repo, opt)

	for _, mergeRequest := range mergeRequests {
		servicePullRequests = append(servicePullRequests, &PullRequest{MergeRequest: mergeRequest})
	}

	return servicePullRequests, errors.Wrap(err, "Failed to get Pull Requests by raw client in gitlab.Client.GetPullRequests")
}

func (p *pullRequestsService) GetPullRequestsURL(owner, repo string) (string, error) {
	repoUrl, err := p.repositoriesService.GetURL(owner, repo)
	return repoUrl + "/merge_requests", errors.Wrap(err, "Error occurred in gitlab.Client.GetPullRequestsURL")
}

func (p *pullRequestsService) GetURL(owner, repo string, id int) (string, error) {
	url, err := p.GetPullRequestsURL(owner, repo)
	return fmt.Sprintf("%s/%d", url, id), errors.Wrap(err, "Error occurred in gitlab.Client.GetPUllRequestURL")
}

func (p *pullRequestsService) Create(ctx context.Context, repo string, newPR *service.NewPullRequest) (service.PullRequest, error) {
	opt := &gitlab.CreateMergeRequestOptions{
		Title:        &newPR.Title,
		Description:  &newPR.Body,
		SourceBranch: &newPR.HeadBranch,
		TargetBranch: &newPR.BaseBranch,
	}
	newMergeRequest, _, err := p.raw.CreateMergeRequest(newPR.BaseOwner+"/"+repo, opt)
	return &PullRequest{MergeRequest: newMergeRequest}, err
}
