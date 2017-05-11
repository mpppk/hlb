package project

type Issue interface {
	GetNumber() int
	GetTitle() string
	GetHTMLURL() string
}

