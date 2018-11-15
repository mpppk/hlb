package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey"
	"github.com/blang/semver"
	"github.com/briandowns/spinner"
	"github.com/mpppk/hlb/hlblib"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var selfupdateCmd = &cobra.Command{
	Use:   "selfupdate",
	Short: "Update hlb",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		version := hlblib.Version
		latest, found, err := selfupdate.DetectLatest("mpppk/hlb")
		if err != nil {
			fmt.Println("Error occurred while detecting version:", err)
			return
		}

		v := semver.MustParse(version)
		if !found || latest.Version.LTE(v) {
			fmt.Println("Current version is the latest")
			return
		}

		msg := fmt.Sprintf("New hlb version is available. (current: v%s latest: v%s)\n" +
			"Do you want to update to latest version? (y/n)", version, latest.Version)

		confirm := false
		prompt := &survey.Confirm{
			Message: msg,
		}
		err = survey.AskOne(prompt, &confirm, nil)
		if err != nil {
			fmt.Println("Sorry, internal error is occurred.")
			return
		}

		if !confirm {
			return
		}

		exe, err := os.Executable()
		if err != nil {
			fmt.Println("Could not locate executable path")
			return
		}

		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond) // Build our new spinner
		s.Start()

		if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
			s.Stop()
			fmt.Println("Error occurred while updating binary:", err)
			return
		}
		s.Stop()
		fmt.Printf("Successfully updated to v%s\n", latest.Version)
	},
}

func init() {
	RootCmd.AddCommand(selfupdateCmd)
}
