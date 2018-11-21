package hlblib

import (
	"github.com/mpppk/gitany"
	"path"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

type HostType int

type Config struct {
	Services []*gitany.ServiceConfig
}

func (c *Config) FindServiceConfig(host string) (*gitany.ServiceConfig, bool) {
	var matchedHosts []*gitany.ServiceConfig
	for _, h := range c.Services {
		if strings.Contains(host, h.Host) {
			matchedHosts = append(matchedHosts, h)
		}
	}

	matchedHostsLen := len(matchedHosts)

	if matchedHostsLen <= 0 {
		return nil, false
	}

	if matchedHostsLen == 1 {
		return matchedHosts[0], true
	}

	// priority
	// query: example.com, matchedHosts: [example.com, example.com:81]
	// => return example.com
	// query: example.com:81, matchedHosts: [example.com, example.com:81]
	// => return example.com:81
	queryHasPortNumber := strings.Contains(host, ":") // Check that query host has port number
	for _, h := range matchedHosts {
		hostHasPortNumber := strings.Contains(h.Host, ":")
		if queryHasPortNumber && hostHasPortNumber {
			return h, true
		}

		if !queryHasPortNumber && !hostHasPortNumber {
			return h, true
		}
	}

	return matchedHosts[0], true
}

func (c *Config) FindServiceConfigs(host string) *Config {
	var serviceConfigs []*gitany.ServiceConfig
	for _, s := range c.Services {
		if strings.Contains(s.Host, host) {
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
	return path.Join(dir, GetConfigDirName()), errors.Wrap(err, "Error occurred in hlblib.GetConfigDirPath")
}

func GetConfigFilePath() (string, error) {
	configDirPath, err := GetConfigDirPath()
	return path.Join(configDirPath, GetConfigFileName()), errors.Wrap(err, "Error occurred in hlblib.GetConfigFilePath")
}
