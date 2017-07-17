package etc

import (
	"path"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
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

type ServiceConfig struct {
	Host     string
	Type     string
	Token    string `mapstructure:"oauth_token" yaml:"oauth_token"`
	Protocol string
	User     string
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

func (c *Config) FindServiceConfigs(host string) *Config {
	var serviceConfigs []*ServiceConfig
	for _, s := range c.Services {
		if strings.Contains(s.Host, host) {
			serviceConfigs = append(serviceConfigs, s)
		}
	}
	return &Config{Services: serviceConfigs}
}

func (c *Config) FindServiceConfigsByType(serviceType string) *Config {
	var serviceConfigs []*ServiceConfig
	for _, s := range c.Services {
		if s.Type == serviceType {
			serviceConfigs = append(serviceConfigs, s)
		}
	}
	return &Config{Services: serviceConfigs}
}

func (c *Config) ListServiceConfigHost() (hosts []string) {
	for _, s := range c.Services {
		hosts = append(hosts, s.Host)
	}
	return hosts
}

func GetConfigDirName() string {
	return path.Join(".config", "hlb")
}

func GetConfigFileName() string {
	return ".hlb.yaml"
}

func GetConfigDirPath() (string, error) {
	dir, err := homedir.Dir()
	return path.Join(dir, GetConfigDirName()), errors.Wrap(err, "Error occurred in etc.GetConfigDirPath")
}

func GetConfigFilePath() (string, error) {
	configDirPath, err := GetConfigDirPath()
	return path.Join(configDirPath, GetConfigFileName()), errors.Wrap(err, "Error occurred in etc.GetConfigFilePath")
}
