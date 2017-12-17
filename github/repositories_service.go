package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"github.com/mpppk/hlb/service"
	"github.com/pkg/errors"
)

type repositoriesService struct {
	raw RepositoriesService
	host string
}

func (r *repositoriesService) Get(ctx context.Context, owner, repo string) (service.Repository, error) {
	githubRepo, _, err := r.raw.Get(ctx, owner, repo)
	return &Repository{Repository: githubRepo}, errors.Wrap(err, "Failed to get Repositories by raw client")
}

func (r *repositoriesService) GetURL(owner, repo string) (string, error) {
	return fmt.Sprintf("https://%s/%s/%s", r.host, owner, repo), checkOwnerAndRepo(owner, repo)
}

func (r *repositoriesService) Create(ctx context.Context, repo string) (service.Repository, error) {
	repository := &github.Repository{Name: github.String(repo)}
	retRepository, _, err := r.raw.Create(ctx, "", repository)
	return &Repository{retRepository}, err
}
