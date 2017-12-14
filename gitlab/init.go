package gitlab

import "github.com/mpppk/hlb/hlblib"

func init() {
	hlblib.RegisterClientGenerator(&ClientBuilder{})
}
