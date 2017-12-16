package github

import "github.com/mpppk/hlb/hlblib"

func init() {
	hlblib.RegisterClientGenerator(&ClientBuilder{})
}
