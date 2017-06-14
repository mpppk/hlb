package git

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
)

type RawRemoteConfig interface {
}

type RawRemote interface {
	Config() *RawRemoteConfig
}

type Remote struct {
	URL         string
	ServiceHost string
	Owner       string
	RepoName    string
}

func NewRemote(remoteUrl string) (*Remote, error) {
	var assigned *regexp.Regexp
	if strings.HasPrefix(remoteUrl, "http") {
		assigned = regexp.MustCompile(`https?://[.+@]?(.+)/(.+)/(.+)$`)
	} else if strings.HasPrefix(remoteUrl, "git") {
		assigned = regexp.MustCompile(`git@(.+):(.+)/(.+).git`)
	} else {
		return nil, errors.New("unknown remote: " + remoteUrl)
	}

	result := assigned.FindStringSubmatch(remoteUrl)
	if result == nil {
		return nil, errors.New("unknown remoteUrl pattern: " + remoteUrl)
	}
	return &Remote{
		URL:         remoteUrl,
		ServiceHost: result[1],
		Owner:       result[2],
		RepoName:    strings.Replace(result[3], ".git", "", -1),
	}, nil
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
	return NewRemote(remote.Config().URL)

}
