package gitlab

import (
	"context"
	"fmt"

	"github.com/mpppk/hlb/service"
	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

type issuesService struct {
	projectsService service.RepositoriesService
	raw             IssuesService
	ListOptions   *gitlab.ListOptions
}

func (i *issuesService) ListByRepo(ctx context.Context, owner, repo string) (serviceIssues []service.Issue, err error) {
	opt := &gitlab.ListProjectIssuesOptions{ListOptions: *i.ListOptions}
	issues, _, err := i.raw.ListProjectIssues(owner+"/"+repo, opt)

	for _, issue := range issues {
		serviceIssues = append(serviceIssues, &Issue{Issue: issue})
	}

	return serviceIssues, errors.Wrap(err, "Failed to get Issues by raw client in gitlab.Client.GetIssues")
}

func (i *issuesService) GetIssuesURL(owner, repo string) (string, error) {
	repoUrl, err := i.projectsService.GetURL(owner, repo)
	return repoUrl + "/issues", errors.Wrap(err, "Error occurred in gitlab.Client.GetIssuesURL")
}

func (i *issuesService) GetURL(owner, repo string, id int) (string, error) {
	url, err := i.GetIssuesURL(owner, repo)
	return fmt.Sprintf("%s/%d", url, id), errors.Wrap(err, "Error occurred in gitlab.Client.GetIssueURL")
}
