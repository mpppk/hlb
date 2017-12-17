package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"github.com/mpppk/hlb/service"
	"github.com/pkg/errors"
)

type pullRequestsService struct {
	repositoriesService service.RepositoriesService
	raw                 PullRequestsService
	ListOptions         *github.ListOptions
}

func (p *pullRequestsService) GetPullRequestsURL(owner, repo string) (string, error) {
	repoUrl, err := p.repositoriesService.GetURL(owner, repo)
	return repoUrl + "/pulls", errors.Wrap(err, "Error occurred in github.Client.GetPullRequestsURL")
}

func (p *pullRequestsService) GetURL(owner, repo string, id int) (string, error) {
	repoUrl, err := p.repositoriesService.GetURL(owner, repo)
	return fmt.Sprintf("%s/pull/%d", repoUrl, id), errors.Wrap(err, "Error occurred in github.Client.GetPullRequestURL")
}

func (p *pullRequestsService) Create(ctx context.Context, repo string, newPR *service.NewPullRequest) (service.PullRequest, error) {
	head := fmt.Sprintf("%s:%s", newPR.HeadOwner, newPR.HeadBranch)

	fmt.Println("github body:", newPR.Body)

	newPullRequest := &github.NewPullRequest{
		Title: github.String(newPR.Title),
		Body:  github.String(newPR.Body),
		Base:  github.String(newPR.BaseBranch),
		Head:  github.String(head),
	}

	createdPullRequest, _, err := p.raw.Create(ctx, newPR.BaseOwner, repo, newPullRequest)

	if e, ok := err.(*github.ErrorResponse); ok && e.Message == VALIDATION_FAILED_MSG {
		for _, es := range e.Errors {
			if es.Message == NO_COMMITS_MSG_PREFIX {
				return createdPullRequest, errors.Wrap(err, es.Message)
			}
			if es.Field == "head" && es.Code == CODE_INVALID {
				errMsg := fmt.Sprintf("head branch(%v) is invalid.\n", head)
				errMsg += "The branch you are trying to create a pull request may not exist in the remote repository. Please try the following command.\n"
				errMsg += fmt.Sprintf("git push origin %v\n", newPR.HeadBranch)
				return createdPullRequest, errors.Wrap(err, errMsg)
			}
		}
	}
	return createdPullRequest, errors.Wrap(err, "Error occurred in github.CreatePullRequest")
}

func (p *pullRequestsService) List(ctx context.Context, owner, repo string) (servicePullRequests []service.PullRequest, err error) {
	opt := github.PullRequestListOptions{ListOptions: *p.ListOptions}
	pullRequests, _, err := p.raw.List(ctx, owner, repo, &opt)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Pull Requests by raw client in github.Client.GetPullRequests")
	}

	for _, pullRequest := range pullRequests {
		servicePullRequests = append(servicePullRequests, &PullRequest{PullRequest: pullRequest})
	}

	return servicePullRequests, nil
}
