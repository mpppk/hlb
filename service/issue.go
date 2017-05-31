package service

type Issue interface {
	GetNumber() int
	GetTitle() string
	GetHTMLURL() string
}

