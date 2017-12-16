package main

import "github.com/mpppk/hlb/cmd"
import _ "github.com/mpppk/hlb/gitlab"
import _ "github.com/mpppk/hlb/github"

func main() {
	cmd.Execute()
}
