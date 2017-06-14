package etc

import "strings"

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

type ServiceConfig struct {
	Host     string
	Type     string
	Token    string `mapstructure:"oauth_token" yaml:"oauth_token"`
	Protocol string
}

type Config struct {
	Services []*ServiceConfig
}

func (c *Config) FindServiceConfig(host string) (*ServiceConfig, bool) {
	for _, h := range c.Services {
		if strings.Contains(host, h.Host) {
			return h, true
		}
	}
	return nil, false
}
