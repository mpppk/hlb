package git

import (
	"regexp"
	"strings"

	"gopkg.in/src-d/go-git.v4"
)

type Remote struct {
	Remote          *git.Remote
	ServiceHostName string
	Owner           string
	RepoName        string
}

func NewRemote(remote *git.Remote) *Remote {
	remoteConfig := remote.Config()
	url := remoteConfig.URL

	var assigned *regexp.Regexp
	if strings.HasPrefix(url, "http") {
		assigned = regexp.MustCompile(`https?://[.+@]?(.+)/(.+)/(.+)$`)
	} else if strings.HasPrefix(url, "git") {
		assigned = regexp.MustCompile(`git@(.+):(.+)/(.+).git`)
	} else {
		panic("unknown remote: " + url)
	}

	result := assigned.FindStringSubmatch(url)
	if result == nil {
		panic("unknown url pattern: " + url)
	}
	return &Remote{
		Remote:          remote,
		ServiceHostName: result[1],
		Owner:           result[2],
		RepoName:        strings.Replace(result[3], ".git", "", -1),
	}
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
