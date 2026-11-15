package cli

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/lancher-dev/lancher/internal/cli/shared"
	"github.com/lancher-dev/lancher/internal/config"
	"github.com/lancher-dev/lancher/internal/storage"
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
		fmt.Fprintf(out, "%sNo templates found.%s\n", shared.ColorYellow, shared.ColorReset)
		fmt.Fprintf(out, "Add a template with: %slancher add <source>%s\n", shared.ColorCyan, shared.ColorReset)
		return
	}

	fmt.Fprintf(out, "%sAvailable Templates:%s\n\n", shared.ColorBold, shared.ColorReset)
	for i, name := range templates {
		path, _ := storage.GetTemplatePath(name)
		result := config.LoadConfigWithDetails(path)

		fmt.Fprintf(out, "  %s•%s %s%s%s\n", shared.ColorGreen, shared.ColorReset, shared.ColorBold, name, shared.ColorReset)
		field(out, "Path", path)
		if cfg := result.Config; cfg != nil {
			if cfg.Name != name {
				field(out, "Name", cfg.Name)
			}
			field(out, "Description", cfg.Description)
			field(out, "Author", cfg.Author)
			field(out, "Version", cfg.Version)
		}
		if result.Err != nil {
			fmt.Fprintf(out, "    %s⚠ Warning: %v%s\n", shared.ColorYellow, result.Err, shared.ColorReset)
		}
		if len(result.FoundFiles) > 1 {
			fmt.Fprintf(out, "    %s⚠ Warning: Multiple config files found (%v). Using %s%s\n",
				shared.ColorYellow, result.FoundFiles, result.UsedFile, shared.ColorReset)
		}
		if i < len(templates)-1 {
			fmt.Fprintln(out)
		}
	}
}

func field(out io.Writer, label, value string) {
	if value != "" {
		fmt.Fprintf(out, "    %s%s:%s %s\n", shared.ColorGray, label, shared.ColorReset, value)
	}
}
