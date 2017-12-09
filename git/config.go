package git

import (
	"github.com/pkg/errors"
	"github.com/tcnksm/go-gitconfig"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
)

const EDITOR_KEY = "core.editor" // FIXME
const DEFAULT_EDITOR_NAME = "vim" // FIXME

func GetConfig(repoPath string) (*config.Config, error) {
	errMsg := "Error occurred in git.GetConfig: "
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}

	c, err := r.Config()
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}

	return c, nil
}

func GetEditorName() (string, error) {
	localEditorName, err := gitconfig.Local(EDITOR_KEY)

	if err == nil && localEditorName != "" {
		return localEditorName, nil
	}

	globalEditorName, err := gitconfig.Global(EDITOR_KEY)
	if err == nil && globalEditorName != "" {
		return globalEditorName, nil
	}

	return DEFAULT_EDITOR_NAME, err
}
