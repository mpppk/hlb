package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

type authorizationsService struct {
	raw AuthorizationsService
}

func (a *authorizationsService) CreateToken(ctx context.Context) (string, error) {
	note, err := a.getUniqueNote(ctx, "hlb")
	if err != nil {
		return "", err
	}

	println(note)

	authReq := &github.AuthorizationRequest{
		Note:   github.String(note),
		Scopes: []github.Scope{github.ScopeRepo},
	}

	auth, _, err := a.raw.Create(ctx, authReq)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get authorizations by raw client in github.Client.CreateToken")
	}

	return *auth.Token, nil
}

func (a *authorizationsService) getAllAuthorizations(ctx context.Context, orgNote string) (allAuths []*github.Authorization, err error) {
	opt := &github.ListOptions{
		PerPage: 5,
	}

	for {
		auths, resp, err := a.raw.List(ctx, opt)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get authorizations by raw client in github.Client.GetPullRequestsURL")
		}

		allAuths = append(allAuths, auths...)

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allAuths, err
}


func (a *authorizationsService) getUniqueNote(ctx context.Context, orgNote string) (string, error) {
	allAuths, err := a.getAllAuthorizations(ctx, orgNote)
	if err != nil {
		return "", err
	}

	cnt := 1
	note := orgNote
	for {
		if !hasAuthNote(allAuths, note) {
			return note, nil
		}
		cnt++
		note = fmt.Sprint(orgNote, cnt)
	}
}