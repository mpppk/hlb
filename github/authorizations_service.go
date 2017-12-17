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

	authReq := &github.AuthorizationRequest{
		Note:   github.String(note),
		Scopes: []github.Scope{github.ScopeRepo},
	}

	auth, _, err := a.raw.Create(ctx, authReq)
	return *auth.Token, errors.Wrap(err, "Failed to get authorizations by raw client in github.Client.CreateToken")
}

func (a *authorizationsService) getUniqueNote(ctx context.Context, orgNote string) (string, error) {
	auths, _, err := a.raw.List(ctx, nil)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get authorizations by raw client in github.Client.GetPullRequestsURL")
	}

	cnt := 1
	note := orgNote
	for {
		if !hasAuthNote(auths, note) {
			return note, nil
		}
		cnt++
		note = fmt.Sprint(orgNote, cnt)
	}
}