package gitlab

import "github.com/xanzy/go-gitlab"

type Repository struct {
	*gitlab.Project
}

func (repo *Repository) GetHTMLURL() string {
	return repo.WebURL
}

func (repo *Repository) GetGitURL() string {
	return repo.HTTPURLToRepo
}

func (repo *Repository) GetCloneURL() string {
	return repo.HTTPURLToRepo
}
