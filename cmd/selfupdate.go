package cmd

import (
	"bufio"
	"fmt"
	"github.com/blang/semver"
	"github.com/mpppk/hlb/hlblib"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
	"os"
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

		fmt.Printf("New hlb version is available. (current: v%s latest: v%s)\n" +
			"Do you want to update to latest version? (y/n)", version, latest.Version)
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil || (input != "y\n" && input != "n\n") {
			fmt.Println("Invalid input")
			return
		}
		if input == "n\n" {
			return
		}

		exe, err := os.Executable()
		if err != nil {
			fmt.Println("Could not locate executable path")
			return
		}
		if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
			fmt.Println("Error occurred while updating binary:", err)
			return
		}
		fmt.Println("Successfully updated to version", latest.Version)
	},
}

func init() {
	RootCmd.AddCommand(selfupdateCmd)
}
