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

func (r *repositoriesService) GetMilestonesURL(owner, repo string) (string, error) {
	repoUrl, err := r.GetURL(owner, repo)
	return repoUrl + "/milestones", errors.Wrap(err, "Error occurred in github.Client.GetMilestonesURL")
}

func (r *repositoriesService) GetMilestoneURL(owner, repo string, id int) (string, error) {
	repoUrl, err := r.GetURL(owner, repo)
	return fmt.Sprintf("%s/milestone/%d", repoUrl, id), errors.Wrap(err, "Error occurred in github.Client.GetMilestoneURL")
}

func (r *repositoriesService) GetWikisURL(owner, repo string) (string, error) {
	repoUrl, err := r.GetURL(owner, repo)
	return repoUrl + "/wiki", errors.Wrap(err, "Error occurred in github.Client.GetWikisURL")
}

func (r *repositoriesService) GetCommitsURL(owner, repo string) (string, error) {
	repoUrl, err := r.GetURL(owner, repo)
	return repoUrl + "/commits", errors.Wrap(err, "Error occurred in github.Client.GetCommitsURL")
}

func (r *repositoriesService) CreateRelease(ctx context.Context, owner, repo string, newRelease *service.NewRelease) (service.Release, error) {
	newGHRelease := &github.RepositoryRelease{
		TagName: github.String(newRelease.GetTagName()),
		Name:    github.String(newRelease.GetName()),
		Body:    github.String(newRelease.GetBody()),
	}

	createdRelease, _, err := r.raw.CreateRelease(ctx, owner, repo, newGHRelease)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Issues by raw client in github.Client.CreateRelease")
	}
	return createdRelease, nil
}
