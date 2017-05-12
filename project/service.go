package project

import (
	"context"
)

type Service interface{
	GetPullRequests(ctx context.Context, owner, repo string) ([]PullRequest, error)
	GetIssues(ctx context.Context, owner, repo string) ([]Issue, error)
	GetRepository(ctx context.Context, owner, repo string) (Repository, error)
	GetRepositoryURL(owner, reop string) (string, error)
}


