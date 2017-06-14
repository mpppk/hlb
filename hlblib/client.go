package hlblib

import (
	"context"
	"errors"

	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/github"
	"github.com/mpppk/hlb/gitlab"
	"github.com/mpppk/hlb/service"
)

func GetClient(ctx context.Context, serviceConfig *etc.ServiceConfig) (service.Client, error) {
	switch serviceConfig.Type {
	case etc.HOST_TYPE_GITHUB.String():
		client, err := github.NewClient(ctx, serviceConfig)
		if err != nil {
			return nil, err
		}
		return service.Client(client), nil
	}
	switch serviceConfig.Type {
	case etc.HOST_TYPE_GITLAB.String():
		client, err := gitlab.NewClient(serviceConfig)
		if err != nil {
			return nil, err
		}

		return service.Client(client), nil
	}
	return nil, errors.New("unknown serviceConfig type: " + serviceConfig.Type)
}

func CreateToken(ctx context.Context, serviceConfig *etc.ServiceConfig, username, pass string) (string, error) {
	//user, pass := project.PromptUserAndPassword(serviceConfig.Host)

	var s service.Client
	switch serviceConfig.Type {
	case etc.HOST_TYPE_GITHUB.String():
		client, err := github.NewClientViaBasicAuth(ctx, serviceConfig, username, pass)
		if err != nil {
			return "", err
		}
		s = service.Client(client)
	}
	switch serviceConfig.Type {
	case etc.HOST_TYPE_GITLAB.String():
		//service, err := gitlab.NewClientViaBasicAuth(serviceConfig, user, name)
		//if err != nil {
		//	return nil, err
		//}
		//
		//return project.ServiceConfig(service), nil
	}
	return s.CreateToken(ctx)
}
