package cmd

import (
	"fmt"
	"strconv"

	"github.com/mpppk/hlb/hlblib"

	"github.com/pkg/errors"
	"github.com/skratchdot/open-golang/open"

	"github.com/spf13/cobra"
)

func MaximumNumArgs(n int) cobra.PositionalArgs {
	maximumArgs := cobra.MaximumNArgs(n)
	return func(cmd *cobra.Command, args []string) error {
		if err := maximumArgs(cmd, args); err != nil {
			return err
		}

		if len(args) == 0 {
			return nil
		}

		if _, err := strconv.Atoi(args[0]); err != nil {
			return fmt.Errorf("accepts only number arg, received %s", args[0])
		}
		return nil
	}
}

func NewBrowseCmdFunc(cmdContextFunc hlblib.CmdContextFunc, fetchURLFunc func(cmdContext *hlblib.CmdContext, args []string) (string, error)) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cmdContext, err := cmdContextFunc()
		if err != nil {
			return errors.Wrap(err, "failed to get command context")
		}

		url, err := fetchURLFunc(cmdContext, args)
		if err != nil {
			return errors.Wrap(err, "failed to fetch URL for browse")
		}

		if urlFlag {
			cmd.Println(url)
		} else {
			if err := open.Run(url); err != nil {
				return errors.Wrap(err, "failed to open repository URL: "+url)
			}
		}
		return nil
	}
}
