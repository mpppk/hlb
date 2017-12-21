package github

import (
	"fmt"

	"github.com/mpppk/hlb/service"
	"github.com/pkg/errors"
)

type projectsService struct {
	repositoriesService service.RepositoriesService
}

func (p *projectsService) GetProjectsURL(owner, repo string) (string, error) {
	repoUrl, err := p.repositoriesService.GetURL(owner, repo)
	return repoUrl + "/projects", errors.Wrap(err, "Error occurred in github.Client.GetProjectsURL")
}

func (p *projectsService) GetURL(owner, repo string, id int) (string, error) {
	repoUrl, err := p.repositoriesService.GetURL(owner, repo)
	return fmt.Sprintf("%s/projects/%d", repoUrl, id), errors.Wrap(err, "Error occurred in github.Client.GetProjectURL")
}
