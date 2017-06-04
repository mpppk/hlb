package etc

import (
	"reflect"
	"testing"
)

type findHostTest struct {
	services          []*ServiceConfig
	targetServiceName string
	expectedStatus    bool
	expectedService   *ServiceConfig
}

func TestFindHost(t *testing.T) {
	services := []*ServiceConfig{
		{
			Name:       "nameA",
			Type:       "typeA",
			OAuthToken: "OAuthTokenA",
			Protocol:   "protocolA",
		},
		{
			Name:       "nameB",
			Type:       "typeB",
			OAuthToken: "OAuthTokenB",
			Protocol:   "protocolB",
		},
	}

	findHostTests := []*findHostTest{
		{
			services:          services,
			targetServiceName: "nameA",
			expectedStatus:    true,
			expectedService:   services[0],
		},
		{
			services:          services,
			targetServiceName: "nameB",
			expectedStatus:    true,
			expectedService:   services[1],
		},
		{
			services:          services,
			targetServiceName: "nameC",
			expectedStatus:    false,
			expectedService:   nil,
		},
	}

	for i, f := range findHostTests {
		config := Config{f.services}

		s, ok := config.FindHost(f.targetServiceName)

		if ok != f.expectedStatus {
			t.Errorf("%v: expected find status %v, but %v", i, f.expectedStatus, ok)
		}

		if !ok {
			continue
		}

		if !reflect.DeepEqual(s, f.expectedService) {
			t.Errorf("%v: expected service %v, but %v", i, f.expectedService, s)
		}
	}
}
