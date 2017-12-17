package gitlab

import (
	"context"

	"github.com/pkg/errors"
)

type authorizationsService struct {
	//raw AuthorizationsService
}

func (a *authorizationsService) CreateToken(ctx context.Context) (string, error) {
	return "", errors.New("Not Implemented Yet")
}
