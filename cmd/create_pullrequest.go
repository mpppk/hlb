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
	"github.com/mpppk/hlb/service"
	"github.com/spf13/cobra"
	gogit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

const (
	DEFAULT_BRANCH_NAME  = "master"
	DEFAULT_PR_FILE_NAME = "PULLREQ_EDITMSG"
	DEFAULT_CS           = "#"
)

var baseBranch string
var headBranch string
var argMessage string

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

func getInitMessage(baseBranch string) string {
	initMsg := ""

	r, err := gogit.PlainOpen(".")
	logs, err := r.Log(&gogit.LogOptions{})
	etc.PanicIfErrorExist(err)

	branches, err := r.Branches()
	etc.PanicIfErrorExist(err)
	var baseBranchHash plumbing.Hash
	branches.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().Short() == baseBranch {
			baseBranchHash = ref.Hash()
		}
		return nil
	})

	co, err := logs.Next()
	co2, err2 := logs.Next()

	if err == nil && err2 == nil && co2.Hash == baseBranchHash {
		initMsg = co.Message
	}

	logs.Close()

	return initMsg
}

func editPRTitleAndMessage(pullreqFileName, initMsg, cs string) (title, body string, err error) {
	// TODO Add commit logs
	comments, err := github.RenderPullRequestTpl(initMsg, cs, baseBranch, headBranch, "")
	etc.PanicIfErrorExist(err)

	ioutil.WriteFile(pullreqFileName, []byte(comments), 0777)

	editorName, err := git.GetEditorName()
	etc.PanicIfErrorExist(err)

	c := exec.Command(editorName, pullreqFileName)
	vimr := regexp.MustCompile("[mg]?vi[m]$")
	if vimr.MatchString(editorName) {
		c.Args = append(c.Args, "--cmd")
		c.Args = append(c.Args, "set ft=gitcommit tw=0 wrap lbr")
	}

	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Run()

	contents, err := ioutil.ReadFile(pullreqFileName)
	etc.PanicIfErrorExist(err)

	err = os.Remove(pullreqFileName)
	etc.PanicIfErrorExist(err)

	return readTitleAndMessage(bytes.NewReader(contents), cs)
}

var createpullrequestCmd = &cobra.Command{
	Use:   "pull-request",
	Short: "Create pull request",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		base, err := hlblib.NewCmdBase()
		etc.PanicIfErrorExist(err)
		sw := hlblib.ClientWrapper{Base: base}

		newPR := &service.NewPullRequest{
			// TODO Set head owner from config file
			HeadOwner:  base.Remote.Owner,
			BaseBranch: baseBranch,
		}
		newPR.BaseOwner = base.Remote.Owner

		if headBranch == "" {
			headBranch, err = git.GetCurrentBranch(".")
			etc.PanicIfErrorExist(err)
		}
		newPR.HeadBranch = headBranch

		initMsg := getInitMessage(baseBranch)

		var title, body string
		if argMessage == "" {
			title, body, err = editPRTitleAndMessage(DEFAULT_PR_FILE_NAME, initMsg, DEFAULT_CS)
		} else {
			title, body, err = readTitleAndMessage(strings.NewReader(argMessage), DEFAULT_CS)
		}
		etc.PanicIfErrorExist(err)
		newPR.Title = title
		newPR.Body = body

		pr, err := sw.CreatePullRequest(newPR)
		etc.PanicIfErrorExist(err)
		fmt.Println(pr.GetHTMLURL())
	},
}

func init() {
	createCmd.AddCommand(createpullrequestCmd)
	createpullrequestCmd.PersistentFlags().StringVarP(&baseBranch, "base", "b", DEFAULT_BRANCH_NAME, "Base branch(Default is master)")
	createpullrequestCmd.PersistentFlags().StringVarP(&headBranch, "head", "H", "", "Head branch(Default is current branch)")
	createpullrequestCmd.PersistentFlags().StringVarP(&argMessage, "message", "m", "", "Pull Request title and body")
}
