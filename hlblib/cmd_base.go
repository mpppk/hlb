package hlblib

import (
	"context"
	"errors"

	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/git"
	"github.com/mpppk/hlb/service"
	"github.com/spf13/viper"
)

type CmdBase struct {
	Context       context.Context
	Config        *etc.Config
	Remote        *git.Remote
	ServiceConfig *etc.ServiceConfig
	Client        service.Client
}

func NewCmdBase() (*CmdBase, error) {
	ctx := context.Background()

	var config etc.Config
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
		errors.New("serviceConfig not found" + remote.ServiceHost)
	}

	client, err := GetClient(ctx, serviceConfig)
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
