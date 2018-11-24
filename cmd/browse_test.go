package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/mpppk/gitany/mock"

	"github.com/mpppk/gitany"
	"github.com/mpppk/hlb/git"

	"github.com/mpppk/hlb/hlblib"
)

func newMockClient(serviceConfig *gitany.ServiceConfig) gitany.Client {
	client := mock.NewClient()
	client.Repositories.URL = fmt.Sprintf("%s://%s", serviceConfig.Protocol, serviceConfig.Host)
	return client
}

func newMockCmdContext() (*hlblib.CmdContext, error) {
	serviceConfig := &gitany.ServiceConfig{
		Host:     "example.com",
		Type:     "gitlab",
		Token:    "xxxx",
		Protocol: "http",
	}
	client := newMockClient(serviceConfig)
	return &hlblib.CmdContext{
		Remote:        &git.Remote{},
		ServiceConfig: serviceConfig,
		Client:        client,
	}, nil
}

func TestBrowse(t *testing.T) {
	cases := []struct {
		cmdArgs        []string
		cmdContext     func() (hlblib.CmdContext, error)
		expectedOutput string
	}{
		{
			cmdArgs:        []string{"-u"},
			expectedOutput: "http://example.com\n",
		},
	}

	for _, c := range cases {
		buf := new(bytes.Buffer)
		cmd := NewCmdBrowse(newMockCmdContext)
		cmd.SetOutput(buf)
		cmd.SetArgs(c.cmdArgs)
		if err := cmd.Execute(); err != nil {
			t.Errorf("failed to execute browse cmd: %s\n", err)
			continue
		}
		output := buf.String()
		if c.expectedOutput != output {
			t.Errorf("unexpected response: expected:%s, actual:%s", c.expectedOutput, output)
		}
	}
}
