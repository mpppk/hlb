package project

type PullRequest interface {
	GetNumber() int
	GetTitle() string
	GetHTMLURL() string
}

