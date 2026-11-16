package cli

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/lancher-dev/lancher/internal/config"
	"github.com/lancher-dev/lancher/internal/storage"
	"github.com/lancher-dev/lancher/internal/ui"
)

type templateInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func newListCmd() *cobra.Command {
	var asJSON bool

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List installed templates",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			templates, err := storage.ListTemplates()
			if err != nil {
				return fmt.Errorf("failed to list templates: %w", err)
			}
			if asJSON {
				return printJSON(cmd.OutOrStdout(), templates)
			}
			printList(cmd.OutOrStdout(), templates)
			return nil
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "print templates as JSON")
	return cmd
}

func printJSON(out io.Writer, templates []string) error {
	infos := make([]templateInfo, 0, len(templates))
	for _, name := range templates {
		path, err := storage.GetTemplatePath(name)
		if err != nil {
			return err
		}
		infos = append(infos, templateInfo{Name: name, Path: path})
	}
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	return enc.Encode(infos)
}

func printList(out io.Writer, templates []string) {
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
		fmt.Fprintf(out, "    %sPath:%s %s\n", ui.ColorGray, ui.ColorReset, path)
		printMetadata(out, result.Config, name)
		if result.Err != nil {
			fmt.Fprintf(out, "    %s⚠ Warning: %v%s\n", ui.ColorYellow, result.Err, ui.ColorReset)
		}
		if len(result.FoundFiles) > 1 {
			fmt.Fprintf(out, "    %s⚠ Warning: Multiple config files found (%v). Using %s%s\n",
				ui.ColorYellow, result.FoundFiles, result.UsedFile, ui.ColorReset)
		}
		if i < len(templates)-1 {
			fmt.Fprintln(out)
		}
	}
}

// printMetadata prints the metadata of cfg, leaving out a name equal to
// skipName.
func printMetadata(out io.Writer, cfg *config.Config, skipName string) {
	for _, f := range cfg.Metadata() {
		if f.Label != "Name" || f.Value != skipName {
			fmt.Fprintf(out, "    %s%s:%s %s\n", ui.ColorGray, f.Label, ui.ColorReset, f.Value)
		}
	}
}
