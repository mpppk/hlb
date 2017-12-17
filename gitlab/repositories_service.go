package gitlab

import (
	"context"
	"fmt"

	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/service"
	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

type repositoriesService struct {
	raw ProjectsService
	host string
	serviceConfig *etc.ServiceConfig
}

func (p *repositoriesService) GetURL(owner, repo string) (string, error) {
	err := checkOwnerAndRepo(owner, repo)
	return fmt.Sprintf("%s://%s/%s/%s", p.serviceConfig.Protocol, p.host, owner, repo), errors.Wrap(err, "Error occurred in gitlab.Client.GetRepositoryURL")
}

func (p *repositoriesService) Get(ctx context.Context, owner, repo string) (service.Repository, error) {
	project, _, err := p.raw.GetProject(owner + "/" + repo)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Repository by raw client in gitlab.Client.GetRepository")
	}

	return &Repository{Project: project}, err
}

func (p *repositoriesService) Create(ctx context.Context, repo string) (service.Repository, error) {
	opt := &gitlab.CreateProjectOptions{Name: &repo}
	retRepository, _, err := p.raw.CreateProject(opt)
	return &Repository{retRepository}, err
}
