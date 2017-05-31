package service

type PullRequest interface {
	GetNumber() int
	GetTitle() string
	GetHTMLURL() string
}

