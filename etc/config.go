package etc

import (
	"strings"
)

type HostType int

const (
	HOST_TYPE_GITHUB HostType = iota + 1
	HOST_TYPE_GITLAB
	HOST_TYPE_BITBUCKET
	HOST_TYPE_GITBUCKET
)

func (s HostType) String() string {
	switch s {
	case HOST_TYPE_GITHUB:
		return "github"
	case HOST_TYPE_GITLAB:
		return "gitlab"
	case HOST_TYPE_BITBUCKET:
		return "bitbucket"
	case HOST_TYPE_GITBUCKET:
		return "gitbucket"
	default:
		return "Unknown"
	}
}

type Host struct {
	Name       string
	Type       string
	OAuthToken string `mapstructure:"oauth_token"`
	Protocol   string
}

type Config struct {
	Hosts []*Host
}

func (c *Config) FindHost(name string) (*Host, bool) {
	for _, h := range c.Hosts {
		if strings.Contains(name, h.Name) {
			return h, true
		}
	}
	return nil, false
}