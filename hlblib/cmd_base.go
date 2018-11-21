package hlblib

import (
	"context"
	"errors"

	"github.com/mpppk/gitany"

	"github.com/mpppk/hlb/git"
	"github.com/spf13/viper"
)

type CmdBase struct {
	Context       context.Context
	Config        *Config
	Remote        *git.Remote
	ServiceConfig *gitany.ServiceConfig
	Client        gitany.Client
}

func NewCmdBase() (*CmdBase, error) {
	ctx := context.Background()

	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	remote, err := git.GetDefaultRemote(".")
	if err != nil {
		return nil, err
	}

	serviceConfig, ok := config.FindServiceConfig(remote.ServiceHost)
	if !ok {
		return nil, errors.New("serviceConfig not found" + remote.ServiceHost)
	}

	client, err := gitany.NewClient(ctx, serviceConfig)
	if err != nil {
		return nil, err
	}

	return &CmdBase{
		Context:       ctx,
		Config:        &config,
		Remote:        remote,
		ServiceConfig: serviceConfig,
		Client:        client,
	}, nil
}
