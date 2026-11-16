package cli

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/spf13/cobra"

	"github.com/lancher-dev/lancher/internal/storage"
	"github.com/lancher-dev/lancher/internal/ui"
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
			fmt.Printf("%sNo templates found.%s\n", ui.ColorYellow, ui.ColorReset)
			return nil
		}
		if names, err = ui.MultiSelect("Select templates to remove:", templates); err != nil {
			return err
		}
		if len(names) == 0 {
			fmt.Printf("%sNo templates selected.%s\n", ui.ColorYellow, ui.ColorReset)
			return nil
		}
	}

	slices.Sort(names)
	names = slices.Compact(names)
	for _, name := range names {
		if err := ui.SanitizeTemplateName(name); err != nil {
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
			err = os.RemoveAll(path)
		}
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to remove template '%s': %w", name, err))
			continue
		}
		fmt.Printf("%s✓ Template '%s' removed%s\n", ui.ColorGreen, name, ui.ColorReset)
		removed++
	}
	if removed > 1 {
		fmt.Printf("\n%s✓ Removed %d templates%s\n", ui.ColorGreen, removed, ui.ColorReset)
	}
	return errors.Join(errs...)
}
