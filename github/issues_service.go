package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"github.com/mpppk/hlb/service"
	"github.com/pkg/errors"
)

type issuesService struct {
	repositoriesService service.RepositoriesService
	raw             IssuesService
	ListOptions   *github.ListOptions
}

func (i *issuesService) GetIssuesURL(owner, repo string) (string, error) {
	if err := checkOwnerAndRepo(owner, repo); err != nil {
		return "", errors.Wrap(err, "Invalid owner or repo was passed to GetIssuesURL")
	}

	repoUrl, err := i.repositoriesService.GetURL(owner, repo)
	return repoUrl + "/issues", errors.Wrap(err, "Error occurred in github.Client.GetIssuesURL")
}

func (i *issuesService) GetURL(owner, repo string, id int) (string, error) {
	url, err := i.GetIssuesURL(owner, repo)
	return fmt.Sprintf("%s/%d", url, id), errors.Wrap(err, "Error occurred in github.Client.GetIssueURL")
}

func (i *issuesService) ListByRepo(ctx context.Context, owner, repo string) (serviceIssues []service.Issue, err error) {
	opt := &github.IssueListByRepoOptions{ListOptions: *i.ListOptions}
	issues, err := i.getGitHubIssues(ctx, owner, repo, opt)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Issues by raw client in github.Client.GetIssues")
	}

	for _, issue := range issues {
		serviceIssues = append(serviceIssues, &Issue{Issue: issue})
	}

	return serviceIssues, errors.Wrap(err, "Error occurred in github.Client.GetIssues")
}

func (i *issuesService) getGitHubIssues(ctx context.Context, owner, repo string, opt *github.IssueListByRepoOptions) (issues []*github.Issue, err error) {
	issuesAndPRs, _, err := i.raw.ListByRepo(ctx, owner, repo, opt)

	if err != nil {
		return nil, errors.Wrap(err, "Error occurred in github.Client.getGitHubIssues")
	}

	for _, issueOrPR := range issuesAndPRs {
		if issueOrPR.PullRequestLinks == nil {
			issues = append(issues, issueOrPR)
		}
	}
	return issues, nil
}
