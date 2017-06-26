package git

import (
	"github.com/pkg/errors"
	"github.com/tcnksm/go-gitconfig"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
)

const EDITOR_KEY = "core.editor"
const DEFAULT_EDITOR_NAME = "vim"

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
	errMsg := "Error occurred in git.GetEditor: "
	localEditorName, err := gitconfig.Local(EDITOR_KEY)

	_, isErrNotFound := err.(*gitconfig.ErrNotFound)

	if err != nil && !isErrNotFound {
		return "", errors.Wrap(err, errMsg+"when retrieve local git config")
	}

	if !isErrNotFound {
		return localEditorName, nil
	}

	globalEditorName, err := gitconfig.Global(EDITOR_KEY)
	_, isErrNotFound = err.(*gitconfig.ErrNotFound)
	if err != nil && !isErrNotFound {
		return "", errors.Wrap(err, errMsg+"when retrieve global git config")
	}

	if !isErrNotFound {
		return globalEditorName, nil
	}

	return DEFAULT_EDITOR_NAME, nil
}
