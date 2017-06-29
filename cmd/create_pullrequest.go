package cmd

import (
	"os"
	"os/exec"

	"io/ioutil"

	"fmt"

	"bufio"
	"bytes"
	"io"
	"regexp"
	"strings"

	"github.com/mpppk/hlb/etc"
	"github.com/mpppk/hlb/git"
	"github.com/mpppk/hlb/github"
	"github.com/mpppk/hlb/hlblib"
	"github.com/spf13/cobra"
	gogit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

const (
	DEFAULT_BRANCH_NAME = "master"
)

func readTitleAndMessage(reader io.Reader, cs string) (title, body string, err error) {
	var titleParts, bodyParts []string

	r := regexp.MustCompile("\\S")
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, cs) {
			continue
		}

		if len(bodyParts) == 0 && r.MatchString(line) {
			titleParts = append(titleParts, line)
		} else {
			bodyParts = append(bodyParts, line)
		}
	}

	if err = scanner.Err(); err != nil {
		return
	}

	title = strings.Join(titleParts, " ")
	title = strings.TrimSpace(title)

	body = strings.Join(bodyParts, "\n")
	body = strings.TrimSpace(body)

	return
}

var createpullrequestCmd = &cobra.Command{
	Use:   "pull-request",
	Short: "Create pull request",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		base, err := hlblib.NewCmdBase()
		etc.PanicIfErrorExist(err)
		sw := hlblib.ClientWrapper{Base: base}

		baseOwner := base.Remote.Owner
		baseBranch := DEFAULT_BRANCH_NAME
		headBranch, err := git.GetCurrentBranch(".")

		r, err := gogit.PlainOpen(".")
		logs, err := r.Log(&gogit.LogOptions{})
		etc.PanicIfErrorExist(err)

		branches, err := r.Branches()
		etc.PanicIfErrorExist(err)
		var masterHash plumbing.Hash
		branches.ForEach(func(ref *plumbing.Reference) error {
			if ref.Name().Short() == "master" {
				masterHash = ref.Hash()
			}
			fmt.Println(ref.Name(), ref.Hash())
			return nil
		})

		initMsg := ""

		co, err := logs.Next()
		etc.PanicIfErrorExist(err)
		co2, err2 := logs.Next()
		etc.PanicIfErrorExist(err2)

		if err == nil && err2 == nil && co2.Hash == masterHash {
			initMsg = co.Message
		}

		logs.Close()

		comments, err := github.RenderPullRequestTpl(initMsg, "#", baseBranch, headBranch, "")
		etc.PanicIfErrorExist(err)

		pullreqFileName := "PULLREQ_EDITMSG"
		ioutil.WriteFile(pullreqFileName, []byte(comments), 0777)

		editorName, err := git.GetEditorName()
		etc.PanicIfErrorExist(err)

		c := exec.Command(editorName, pullreqFileName)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Run()

		contents, err := ioutil.ReadFile(pullreqFileName)
		etc.PanicIfErrorExist(err)

		err = os.Remove(pullreqFileName)
		etc.PanicIfErrorExist(err)

		title, message, err := readTitleAndMessage(bytes.NewReader(contents), "#")
		etc.PanicIfErrorExist(err)

		pr, err := sw.CreatePullRequest(baseOwner, baseBranch, headBranch, title, message)
		etc.PanicIfErrorExist(err)
		fmt.Println(pr.GetHTMLURL())
	},
}

func init() {
	createCmd.AddCommand(createpullrequestCmd)
}
