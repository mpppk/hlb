package service

type Repository interface {
	GetHTMLURL() string
	GetGitURL() string
}
