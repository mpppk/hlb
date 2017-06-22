package git

import (
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func GetCurrentBranch(path string) (string, error) {
	errMsg := "Error occured in git.GetCurrentBranch: "
	r, err := git.PlainOpen(path)
	if err != nil {
		return "", errors.Wrap(err, errMsg)
	}

	head, err := r.Head()

	if !head.IsBranch() {
		return "", errors.New(errMsg + "You are in detached branch")
	}

	if head.Type() == plumbing.InvalidReference {
		return "", errors.New(errMsg + "HEAD is invalid reference")
	}

	return head.Name().Short(), nil
}
