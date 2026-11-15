// Package cli defines the lancher command line.
package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/lancher-dev/lancher/internal/cli/shared"
	"github.com/lancher-dev/lancher/internal/version"
)

type globalOptions struct {
	verbose bool
}

// NewRootCmd returns the lancher command with all subcommands.
func NewRootCmd() *cobra.Command {
	opts := &globalOptions{}

	root := &cobra.Command{
		Use:           "lancher",
		Short:         "A minimal local project template manager",
		Version:       version.Get(),
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	root.SetVersionTemplate("lancher {{.Version}}\n")
	root.PersistentFlags().BoolVar(&opts.verbose, "verbose", false, "show command output instead of a spinner")

	root.AddCommand(
		newNewCmd(opts),
		newAddCmd(opts),
		newListCmd(),
		newRemoveCmd(),
	)
	return root
}

// Execute runs the root command and reports a cancelled prompt as success.
func Execute() error {
	err := NewRootCmd().Execute()
	if errors.Is(err, shared.ErrCancelled) {
		fmt.Printf("%sCancelled.%s\n", shared.ColorYellow, shared.ColorReset)
		return nil
	}
	return err
}
