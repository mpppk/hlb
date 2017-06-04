package hlblib

import (
	"context"
	"errors"

	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/github"
	"github.com/mpppk/hlb/gitlab"
	"github.com/mpppk/hlb/service"
)

func GetClient(ctx context.Context, host *etc.ServiceConfig) (service.Client, error) {
	switch host.Type {
	case etc.HOST_TYPE_GITHUB.String():
		client, err := github.NewClient(ctx, host)
		if err != nil {
			return nil, err
		}
		return service.Client(client), nil
	}
	switch host.Type {
	case etc.HOST_TYPE_GITLAB.String():
		client, err := gitlab.NewClient(host)
		if err != nil {
			return nil, err
		}

		return service.Client(client), nil
	}
	return nil, errors.New("unknown host type: " + host.Type)
}

func CreateToken(ctx context.Context, host *etc.ServiceConfig, username, pass string) (string, error) {
	//user, pass := project.PromptUserAndPassword(host.Name)

	var s service.Client
	switch host.Type {
	case etc.HOST_TYPE_GITHUB.String():
		client, err := github.NewClientViaBasicAuth(ctx, host, username, pass)
		if err != nil {
			return "", err
		}
		s = service.Client(client)
	}
	switch host.Type {
	case etc.HOST_TYPE_GITLAB.String():
		//service, err := gitlab.NewClientViaBasicAuth(host, user, name)
		//if err != nil {
		//	return nil, err
		//}
		//
		//return project.ServiceConfig(service), nil
	}
	return s.CreateToken(ctx)
}
