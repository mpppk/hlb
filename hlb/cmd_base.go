package hlb

import (
	"context"
	"errors"

	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/git"
	"github.com/mpppk/hlb/service"
	"github.com/spf13/viper"
)

type CmdBase struct {
	Context context.Context
	Config  *etc.Config
	Remote  *git.Remote
	Host    *etc.ServiceConfig
	Service service.ServiceClient
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

	host, ok := config.FindHost(remote.ServiceHostName)
	if !ok {
		errors.New("host not found" + remote.ServiceHostName)
	}

	service, err := GetService(ctx, host)
	if err != nil {
		return nil, err
	}

	return &CmdBase{
		Context: ctx,
		Config:  &config,
		Remote:  remote,
		Host:    host,
		Service: service,
	}, nil
}
