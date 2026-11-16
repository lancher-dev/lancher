package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/lancher-dev/lancher/internal/source"
	"github.com/lancher-dev/lancher/internal/storage"
	"github.com/lancher-dev/lancher/internal/ui"
)

// fetchTemplate installs src as the template name, replacing it only if the
// fetch succeeds, and returns the manifest entry for it without Added.
func fetchTemplate(msg, name string, src source.Source, verbose bool) (storage.Entry, error) {
	var log bytes.Buffer
	var out io.Writer = &log
	if verbose {
		out = os.Stdout
	} else {
		// git fails instead of asking for credentials.
		os.Setenv("GIT_TERMINAL_PROMPT", "0")
	}

	err := ui.Task(msg, verbose, func() error {
		return storage.Replace(name, func(dest string) error {
			return src.Fetch(dest, out)
		})
	})
	if err != nil {
		if tail := lastLines(log.String(), 5); tail != "" {
			return storage.Entry{}, fmt.Errorf("%w\n%s", err, tail)
		}
		return storage.Entry{}, err
	}

	entry := storage.Entry{Kind: src.Kind, Source: src.Display(), Synced: now()}
	if src.Kind == source.Git {
		path, err := storage.GetTemplatePath(name)
		if err != nil {
			return storage.Entry{}, err
		}
		entry.Ref = storage.GitRef(path)
	}
	return entry, nil
}

func now() time.Time {
	return time.Now().UTC().Truncate(time.Second)
}

func fetchMessage(kind source.Kind) string {
	switch kind {
	case source.Git:
		return "Cloning repository..."
	case source.Zip:
		return "Extracting ZIP file..."
	}
	return "Copying template..."
}

// lastLines returns the last n lines of s, without git's "Cloning into" line.
func lastLines(s string, n int) string {
	var lines []string
	for _, line := range strings.Split(strings.TrimSpace(s), "\n") {
		if !strings.HasPrefix(line, "Cloning into") {
			lines = append(lines, line)
		}
	}
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}
