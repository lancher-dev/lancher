package cli

import (
	"slices"

	"github.com/spf13/cobra"

	"github.com/lancher-dev/lancher/internal/storage"
)

// completeTemplates completes template names. With single set, only the first
// argument is completed; otherwise names already given are left out.
func completeTemplates(single bool) cobra.CompletionFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		if single && len(args) > 0 {
			return nil, cobra.ShellCompDirectiveDefault
		}
		templates, err := storage.ListTemplates()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		var names []cobra.Completion
		for _, name := range templates {
			if !slices.Contains(args, name) {
				names = append(names, name)
			}
		}
		return names, cobra.ShellCompDirectiveNoFileComp
	}
}
