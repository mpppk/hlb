package service

type Release interface {
	GetID() int64
	GetTagName() string
	GetName() string
	GetBody() string
	GetHTMLURL() string
}

type NewRelease struct {
	ID      int
	TagName string
	Name    string
	Body    string
	HTMLURL string
}

func (nr *NewRelease) GetID() int {
	return nr.ID
}

func (nr *NewRelease) GetTagName() string {
	return nr.TagName
}

func (nr *NewRelease) GetName() string {
	return nr.Name
}

func (nr *NewRelease) GetBody() string {
	return nr.Body
}

func (nr *NewRelease) GetHTMLURL() string {
	return nr.HTMLURL
}
