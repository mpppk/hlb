package git

import "testing"

type newRemoteTest struct {
	url              string
	willBeError      bool
	expectedHostName string
	expectedOwner    string
	expectedRepoName string
}

func TestNewRemote(t *testing.T) {
	newRemoteTests := []*newRemoteTest{
		{
			url:              "git@github.com:mpppk/hlb.git",
			willBeError:      false,
			expectedHostName: "github.com",
			expectedOwner:    "mpppk",
			expectedRepoName: "hlb",
		},
		{
			url:         "git@github.com:mpppk.git",
			willBeError: true,
		},
		{
			url:         "git@github.com:/hlb.git",
			willBeError: true,
		},
		{
			url:         "git@github.commpppk/hlb.git",
			willBeError: true,
		},
		{
			url:         "github.com/mpppk/hlb.git",
			willBeError: true,
		},
		{
			url:              "https://github.com/mpppk/hlb",
			willBeError:      false,
			expectedHostName: "github.com",
			expectedOwner:    "mpppk",
			expectedRepoName: "hlb",
		},
		{
			url:              "http://github.com/mpppk/hlb",
			willBeError:      false,
			expectedHostName: "github.com",
			expectedOwner:    "mpppk",
			expectedRepoName: "hlb",
		},
		{
			url:         "http://github.com/mpppk",
			willBeError: true,
		},
	}

	for i, nr := range newRemoteTests {
		remote, err := NewRemote(nr.url)

		if err != nil {
			if !nr.willBeError {
				t.Errorf("%v: Unexpected error ocured: %v, params: %v",
					i, err, nr)
			} else {
				continue
			}
		} else if nr.willBeError {
			t.Errorf("%v: Error expected, params: %v",
				i, nr)
		}

		if remote.ServiceHostName != nr.expectedHostName {
			t.Errorf("ServiceHostName field must be host name "+
				"that extracted from provided URL "+
				"%v: expected: %v, actual(extract from %v): %v",
				i, nr.expectedHostName, nr.url, remote.ServiceHostName)
		}

		if remote.Owner != nr.expectedOwner {
			t.Errorf("Owner field must be owner name that extracted from provided URL, "+
				"%v: expected: %v, actual(extract from %v): %v",
				i, nr.expectedOwner, nr.url, remote.Owner)
		}

		if remote.RepoName != nr.expectedRepoName {
			t.Errorf("RepoName field must be repository name that extracted from provided URL, "+
				"%v: expected: %v, actual(extract from %v): %v",
				i, nr.expectedRepoName, nr.url, remote.RepoName)
		}
	}
}
