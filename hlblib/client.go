package hlblib

import (
	"context"
	"errors"

	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/service"
)

var clientGenerators []service.ClientGenerator

func RegisterClientGenerator(clientGenerator service.ClientGenerator) {
	clientGenerators = append(clientGenerators, clientGenerator)
}

func GetClient(ctx context.Context, serviceConfig *etc.ServiceConfig) (service.Client, error) {
	for _, clientGenerator := range clientGenerators {
		if clientGenerator.GetType() == serviceConfig.Type {
			return clientGenerator.New(ctx, serviceConfig)
		}
	}
	return nil, errors.New("unknown serviceConfig type: " + serviceConfig.Type)
}

func CanCreateToken(serviceType string) bool {
	switch serviceType {
	case etc.HOST_TYPE_GITHUB.String():
		return true
	default:
		return false
	}
}

func CreateToken(ctx context.Context, serviceConfig *etc.ServiceConfig, username, pass string) (string, error) {
	for _, clientGenerator := range clientGenerators {
		if clientGenerator.GetType() == serviceConfig.Type {
			client, err := clientGenerator.NewViaBasicAuth(ctx, serviceConfig, username, pass)
			if err != nil {
				return "", err
			}
			return client.GetAuthorizations().CreateToken(ctx)
		}
	}
	return "", errors.New("token creating failed because unknown serviceConfig type is provided: " + serviceConfig.Type)
}
