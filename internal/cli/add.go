package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/lancher-dev/lancher/internal/source"
	"github.com/lancher-dev/lancher/internal/storage"
	"github.com/lancher-dev/lancher/internal/ui"
)

func newAddCmd(global *globalOptions) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add [source] [name]",
		Short: "Add a template from a directory, ZIP file or git repository",
		Long: `Add a template from a directory, ZIP file or git repository.

Sources: a local path, a .zip file, a git URL (HTTPS or SSH), gh:owner/repo
or gl:owner/repo. The name defaults to the repository, directory or ZIP name.
Without arguments, lancher asks for the source and the name.`,
		Example: `  lancher add ~/projects/api-base
  lancher add gh:owner/repo my-template
  lancher add ./starter.zip --force`,
		Args: cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdd(args, force, global.verbose)
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "replace an existing template with the same name")
	return cmd
}

func runAdd(args []string, force, verbose bool) error {
	var src source.Source
	var name string
	var err error

	if len(args) == 0 {
		if src, err = promptSource(); err != nil {
			return err
		}
		name, err = ui.PromptStringWithDefault("Template name:", src.Name())
		if err != nil {
			return err
		}
	} else {
		if src, err = source.Parse(args[0]); err != nil {
			return err
		}
		name = src.Name()
		if len(args) == 2 {
			name = args[1]
		}
	}

	exists, err := storage.TemplateExists(name)
	if err != nil {
		return err
	}
	if exists && !force {
		return fmt.Errorf("template '%s' already exists (use --force to replace it)", name)
	}

	m, err := storage.LoadManifest()
	if err != nil {
		return err
	}
	entry, err := fetchTemplate(fetchMessage(src.Kind), name, src, verbose)
	if err != nil {
		return err
	}
	entry.Added = entry.Synced
	m.Templates[name] = entry
	if err := m.Save(); err != nil {
		return err
	}

	done := "added"
	if exists {
		done = "replaced"
	}
	fmt.Printf("%s✓ Template '%s' %s%s\n", ui.ColorGreen, name, done, ui.ColorReset)
	path, _ := storage.GetTemplatePath(name)
	fmt.Printf("  %sSource:%s %s\n", ui.ColorYellow, ui.ColorReset, src.Display())
	fmt.Printf("  %sStored:%s %s\n", ui.ColorYellow, ui.ColorReset, path)
	return nil
}

func promptSource() (source.Source, error) {
	raw, err := ui.PromptStringWithDefault("Source (local path, git URL or ZIP file):", ".")
	if err != nil {
		return source.Source{}, err
	}
	return source.Parse(raw)
}
