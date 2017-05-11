package git

import (
	"gopkg.in/src-d/go-git.v4"
	"regexp"
)

type Remote struct {
	Remote *git.Remote
	ServiceHostName string
	Owner    string
	RepoName string
}

func NewRemote(remote *git.Remote) *Remote {
	remoteConfig := remote.Config()
	url := remoteConfig.URL

	assigned := regexp.MustCompile(`git@(.+):(.+)/(.+).git`)

	newRemote := &Remote{}
	if assigned != nil {
		result := assigned.FindStringSubmatch(url)
		newRemote.ServiceHostName = result[1]
		newRemote.Owner = result[2]
		newRemote.RepoName = result[3]
	}
	return newRemote
}

func GetDefaultRemote(path string) (*Remote, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	remote, err := r.Remote(git.DefaultRemoteName)
	if err != nil {
		return nil, err
	}
	return NewRemote(remote), nil

}