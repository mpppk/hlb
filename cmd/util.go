package cmd

import (
	"fmt"
	"strconv"

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
