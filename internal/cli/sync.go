package cli

import (
	"errors"
	"fmt"
	"maps"
	"slices"

	"github.com/spf13/cobra"

	"github.com/lancher-dev/lancher/internal/source"
	"github.com/lancher-dev/lancher/internal/storage"
	"github.com/lancher-dev/lancher/internal/ui"
)

func newSyncCmd(global *globalOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "sync [template...]",
		Short: "Fetch templates again from their sources",
		Long: `Fetch templates again from their sources and replace them. A template is
replaced only if the fetch succeeds. Without arguments, lancher syncs every
template.`,
		ValidArgsFunction: completeTemplates(false),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSync(args, global.verbose)
		},
	}
}

func runSync(names []string, verbose bool) error {
	m, err := storage.LoadManifest()
	if err != nil {
		return err
	}
	if len(names) == 0 {
		names = slices.Sorted(maps.Keys(m.Templates))
		if len(names) == 0 {
			fmt.Printf("%sNo templates found.%s\n", ui.ColorYellow, ui.ColorReset)
			return nil
		}
	} else {
		slices.Sort(names)
		names = slices.Compact(names)
		for _, name := range names {
			if _, ok := m.Templates[name]; !ok {
				return fmt.Errorf("template '%s' not found", name)
			}
		}
	}

	var errs []error
	for _, name := range names {
		entry := m.Templates[name]
		src, skip := syncSource(name, entry)
		if skip != "" {
			fmt.Printf("%s⚠ '%s' skipped: %s%s\n", ui.ColorYellow, name, skip, ui.ColorReset)
			continue
		}

		synced, err := fetchTemplate(fmt.Sprintf("Syncing '%s'...", name), name, src, verbose)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to sync '%s': %w", name, err))
			continue
		}
		synced.Added = entry.Added
		m.Templates[name] = synced
		if err := m.Save(); err != nil {
			return errors.Join(append(errs, err)...)
		}

		ref := ""
		if synced.Ref != "" {
			ref = " at " + shortRef(synced.Ref, 7)
		}
		fmt.Printf("%s✓ Template '%s' synced%s%s\n", ui.ColorGreen, name, ref, ui.ColorReset)
	}
	return errors.Join(errs...)
}

// syncSource returns the source of entry, or why the template is skipped.
func syncSource(name string, entry storage.Entry) (source.Source, string) {
	if !entry.Known() {
		return source.Source{}, fmt.Sprintf("unknown source (re-add it with: lancher add --force <source> %s)", name)
	}
	src, err := source.Parse(entry.Source)
	if err != nil {
		return source.Source{}, err.Error()
	}
	return src, ""
}

func shortRef(ref string, n int) string {
	if len(ref) > n {
		return ref[:n]
	}
	return ref
}
