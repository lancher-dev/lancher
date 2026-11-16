package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/lancher-dev/lancher/internal/config"
	"github.com/lancher-dev/lancher/internal/storage"
	"github.com/lancher-dev/lancher/internal/ui"
)

type templateInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
	storage.Entry
}

func newListCmd() *cobra.Command {
	var asJSON bool

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List installed templates",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := storage.LoadManifest()
			if err != nil {
				return err
			}
			if asJSON {
				return printJSON(cmd.OutOrStdout(), m)
			}
			printList(cmd.OutOrStdout(), m)
			return nil
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "print templates as JSON")
	return cmd
}

func printJSON(out io.Writer, m *storage.Manifest) error {
	names := slices.Sorted(maps.Keys(m.Templates))
	infos := make([]templateInfo, 0, len(names))
	for _, name := range names {
		path, err := storage.GetTemplatePath(name)
		if err != nil {
			return err
		}
		infos = append(infos, templateInfo{Name: name, Path: path, Entry: m.Templates[name]})
	}
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	return enc.Encode(infos)
}

func printList(out io.Writer, m *storage.Manifest) {
	templates := slices.Sorted(maps.Keys(m.Templates))
	if len(templates) == 0 {
		fmt.Fprintf(out, "%sNo templates found.%s\n", ui.ColorYellow, ui.ColorReset)
		fmt.Fprintf(out, "Add a template with: %slancher add <source>%s\n", ui.ColorCyan, ui.ColorReset)
		return
	}

	fmt.Fprintf(out, "%sAvailable Templates:%s\n\n", ui.ColorBold, ui.ColorReset)
	for i, name := range templates {
		path, _ := storage.GetTemplatePath(name)
		result := config.LoadConfigWithDetails(path)

		fmt.Fprintf(out, "  %s•%s %s%s%s\n", ui.ColorGreen, ui.ColorReset, ui.ColorBold, name, ui.ColorReset)
		printField(out, "Path", path)
		printField(out, "Source", sourceLabel(m.Templates[name]))
		printMetadata(out, result.Config, name)
		printConfigWarnings(out, result)
		if i < len(templates)-1 {
			fmt.Fprintln(out)
		}
	}
}

// printConfigWarnings prints the problems found while loading a template
// config.
func printConfigWarnings(out io.Writer, result *config.LoadResult) {
	if result.Err != nil {
		fmt.Fprintf(out, "    %s⚠ Warning: %v%s\n", ui.ColorYellow, result.Err, ui.ColorReset)
	}
	if len(result.Unknown) > 0 {
		fmt.Fprintf(out, "    %s⚠ Warning: unknown keys in %s: %s%s\n",
			ui.ColorYellow, result.UsedFile, ui.Clean(strings.Join(result.Unknown, ", ")), ui.ColorReset)
	}
	if len(result.FoundFiles) > 1 {
		fmt.Fprintf(out, "    %s⚠ Warning: Multiple config files found (%v). Using %s%s\n",
			ui.ColorYellow, result.FoundFiles, result.UsedFile, ui.ColorReset)
	}
}

// printMetadata prints the metadata of cfg, leaving out a name equal to
// skipName.
func printMetadata(out io.Writer, cfg *config.Config, skipName string) {
	for _, f := range cfg.Metadata() {
		if f.Label != "Name" || f.Value != skipName {
			fmt.Fprintf(out, "    %s%s:%s %s\n", ui.ColorGray, f.Label, ui.ColorReset, ui.Clean(f.Value))
		}
	}
}

func sourceLabel(entry storage.Entry) string {
	if !entry.Known() {
		return string(storage.UnknownKind)
	}
	return fmt.Sprintf("%s (%s)", ui.Clean(entry.Source), entry.Kind)
}

func printField(out io.Writer, label, value string) {
	fmt.Fprintf(out, "    %s%s:%s %s\n", ui.ColorGray, label, ui.ColorReset, value)
}
