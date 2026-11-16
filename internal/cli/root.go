// Package cli defines the lancher command line.
package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/lancher-dev/lancher/internal/ui"
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
		Version:       version.Version,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	root.SetVersionTemplate("lancher {{.Version}}\n")
	root.PersistentFlags().BoolVar(&opts.verbose, "verbose", false, "show the output of git instead of a spinner")

	root.AddCommand(
		newNewCmd(),
		newAddCmd(opts),
		newListCmd(),
		newRemoveCmd(),
	)
	return root
}

// Execute runs the root command and returns the exit code: 1 on error and
// 130 when a prompt is cancelled.
func Execute() int {
	err := NewRootCmd().Execute()
	switch {
	case err == nil:
		return 0
	case errors.Is(err, ui.ErrCancelled):
		fmt.Printf("%sCancelled.%s\n", ui.ColorYellow, ui.ColorReset)
		return 130
	}
	ui.PrintError(err)
	return 1
}
