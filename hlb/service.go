package hlb

import (
	"context"
	"errors"

	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/github"
	"github.com/mpppk/hlb/gitlab"
	"github.com/mpppk/hlb/project"
)

func GetService(ctx context.Context, host *etc.Host) (project.Service, error) {
	switch host.Type {
	case etc.HOST_TYPE_GITHUB.String():
		service, err := github.NewService(ctx, host)
		if err != nil {
			return nil, err
		}
		return project.Service(service), nil
	}
	switch host.Type {
	case etc.HOST_TYPE_GITLAB.String():
		service, err := gitlab.NewService(host)
		if err != nil {
			return nil, err
		}

		return project.Service(service), nil
	}
	return nil, errors.New("unknown host type: " + host.Type)
}

func CreateToken(ctx context.Context, host *etc.Host, username, pass string) (string, error) {
	//user, pass := project.PromptUserAndPassword(host.Name)

	var s project.Service
	switch host.Type {
	case etc.HOST_TYPE_GITHUB.String():
		service, err := github.NewServiceViaBasicAuth(ctx, host, username, pass)
		if err != nil {
			return "", err
		}
		s = project.Service(service)
	}
	switch host.Type {
	case etc.HOST_TYPE_GITLAB.String():
		//service, err := gitlab.NewServiceViaBasicAuth(host, user, name)
		//if err != nil {
		//	return nil, err
		//}
		//
		//return project.Service(service), nil
	}
	return s.CreateToken(ctx)
}
