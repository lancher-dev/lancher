package cli

import (
	"errors"
	"fmt"
	"slices"

	"github.com/spf13/cobra"

	"github.com/lancher-dev/lancher/internal/cli/shared"
	"github.com/lancher-dev/lancher/internal/fileutil"
	"github.com/lancher-dev/lancher/internal/storage"
)

func newRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "remove [name...]",
		Aliases:           []string{"rm"},
		Short:             "Remove one or more templates",
		Long:              "Remove one or more templates. Without arguments, lancher asks which ones to remove.",
		ValidArgsFunction: completeTemplates(false),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRemove(args)
		},
	}
}

func runRemove(names []string) error {
	if len(names) == 0 {
		templates, err := storage.ListTemplates()
		if err != nil {
			return fmt.Errorf("failed to list templates: %w", err)
		}
		if len(templates) == 0 {
			fmt.Printf("%sNo templates found.%s\n", shared.ColorYellow, shared.ColorReset)
			return nil
		}
		if names, err = shared.MultiSelect("Select templates to remove:", templates); err != nil {
			return err
		}
		if len(names) == 0 {
			fmt.Printf("%sNo templates selected.%s\n", shared.ColorYellow, shared.ColorReset)
			return nil
		}
	}

	slices.Sort(names)
	names = slices.Compact(names)
	for _, name := range names {
		if err := shared.SanitizeTemplateName(name); err != nil {
			return fmt.Errorf("invalid template name '%s': %w", name, err)
		}
		exists, err := storage.TemplateExists(name)
		if err != nil {
			return fmt.Errorf("failed to check template '%s': %w", name, err)
		}
		if !exists {
			return fmt.Errorf("template '%s' not found", name)
		}
	}

	var errs []error
	removed := 0
	for _, name := range names {
		path, err := storage.GetTemplatePath(name)
		if err == nil {
			err = fileutil.RemoveDir(path)
		}
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to remove template '%s': %w", name, err))
			continue
		}
		fmt.Printf("%s✓ Template '%s' removed%s\n", shared.ColorGreen, name, shared.ColorReset)
		removed++
	}
	if removed > 1 {
		fmt.Printf("\n%s✓ Removed %d templates%s\n", shared.ColorGreen, removed, shared.ColorReset)
	}
	return errors.Join(errs...)
}
