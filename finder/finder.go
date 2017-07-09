package finder

import (
	"os/exec"
	"io"
	"strings"
	"github.com/mpppk/hlb/service"
	"strconv"
)

type FilterStringer interface {
	FilterString() string
}

type Linker interface {
	GetURL() string
}

type FilterableIssue struct {
	service.Issue
}

func (f *FilterableIssue) GetURL() string{
	return f.GetHTMLURL()
}

func (f *FilterableIssue) FilterString() string{
	return "#" + strconv.Itoa(f.GetNumber()) + " " + f.GetTitle()
}

func ToFilterableIssues(issues []service.Issue) (fis []*FilterableIssue) {
	for _, issue := range issues {
		fi := &FilterableIssue{Issue: issue}
		fis = append(fis, fi)
	}
	return fis
}


func PipeToPeco(texts []string) (string, error) {
	cmd := exec.Command("peco")
	stdin, _ := cmd.StdinPipe()
	io.WriteString(stdin, strings.Join(texts, "\n"))
	stdin.Close()
	out, err := cmd.Output()

	if err != nil {
		return "", err
	}

	return strings.Trim(string(out), " \n"), nil
}


func Find(stringers []FilterStringer) ([]FilterStringer, error){
	strs := []string{}

	for _, stringer := range stringers {
		strs = append(strs, stringer.FilterString())
	}

	selectedStr, err := PipeToPeco(strs)
	if err != nil {
		return nil, err
	}

	selectedStringers := []FilterStringer{}
	for _, stringer := range stringers {
		if stringer.FilterString() == selectedStr {
			selectedStringers = append(selectedStringers, stringer)
		}
	}
	return selectedStringers, nil
}
