package cli

import (
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"

	"github.com/lancher-dev/lancher/internal/config"
	"github.com/lancher-dev/lancher/internal/storage"
	"github.com/lancher-dev/lancher/internal/ui"
)

func newInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "info <template>",
		Short:             "Show the source, metadata and hooks of a template",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completeTemplates(true),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInfo(cmd.OutOrStdout(), args[0])
		},
	}
}

func runInfo(out io.Writer, name string) error {
	m, err := storage.LoadManifest()
	if err != nil {
		return err
	}
	entry, ok := m.Templates[name]
	if !ok {
		return fmt.Errorf("template '%s' not found", name)
	}
	path, _ := storage.GetTemplatePath(name)
	result := config.LoadConfigWithDetails(path)

	fmt.Fprintf(out, "%s%s%s\n", ui.ColorBold, name, ui.ColorReset)
	printField(out, "Path", path)
	printField(out, "Source", sourceLabel(entry))
	if entry.Ref != "" {
		printField(out, "Ref", shortRef(entry.Ref, 12))
	}
	printTime(out, "Added", entry.Added)
	printTime(out, "Synced", entry.Synced)
	printMetadata(out, result.Config, name)

	if !result.Config.HasHooks() {
		printField(out, "Hooks", "none")
	} else {
		fmt.Fprintf(out, "    %sHooks:%s\n", ui.ColorGray, ui.ColorReset)
		for i, hook := range result.Config.Hooks {
			fmt.Fprintf(out, "      %d. %s\n", i+1, ui.Clean(hook))
		}
	}
	printConfigWarnings(out, result)
	return nil
}

func printTime(out io.Writer, label string, t time.Time) {
	if !t.IsZero() {
		printField(out, label, t.Local().Format("2006-01-02 15:04"))
	}
}
