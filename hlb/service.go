package hlb

import (
	"errors"
	"github.com/mpppk/hlb/project"
	"github.com/mpppk/hlb/gitlab"
	"context"
	"github.com/mpppk/hlb/github"
	"github.com/mpppk/hlb/etc"
)

func GetService(ctx context.Context, host *etc.Host) (project.Service, error){

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
