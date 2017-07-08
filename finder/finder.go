package finder

import (
	"os/exec"
	"io"
	"strings"
)

type FilterStringer interface {
	FilterString() string
}

type FilterableString string

func (f FilterableString) FilterString() string{
	return string(f)
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

func FindFromFilterableStrings(fstrs []FilterableString) ([]FilterableString, error) {
	filterStringers := []FilterStringer{}
	for _, s := range fstrs {
		filterStringers = append(filterStringers, s)
	}
	selectedFStringers, err := Find(filterStringers)
	if err != nil {
		return nil, err
	}

	retFStrs := []FilterableString{}
	for _, sfs := range selectedFStringers {
		for _, s := range fstrs {
			if s.FilterString() == sfs.FilterString() {
				retFStrs = append(retFStrs, s)
			}
		}
	}
	return retFStrs, nil
}
