package gitlab

import (
	"github.com/mpppk/hlb/service"
	"github.com/pkg/errors"
)

type projetsService struct {
	repositoriesService service.RepositoriesService
}

func (p *projetsService) GetProjectsURL(owner, repo string) (string, error) {
	repoUrl, err := p.repositoriesService.GetURL(owner, repo)
	return repoUrl + "/boards", errors.Wrap(err, "Error occurred in gitlab.Client.GetProjectsURL")
}

func (p *projetsService) GetURL(owner, repo string, id int) (string, error) {
	// TODO: Support multi issue boards
	return p.GetProjectsURL(owner, repo)
}



