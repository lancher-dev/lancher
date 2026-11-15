// Package source resolves template sources and fetches them into a directory.
package source

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lancher-dev/lancher/internal/fileutil"
)

// Kind is the type of a template source.
type Kind string

const (
	Git   Kind = "git"
	Zip   Kind = "zip"
	Local Kind = "local"
)

type alias struct {
	prefix  string
	cli     string
	baseURL string
}

var aliases = []alias{
	{"gh:", "gh", "https://github.com/"},
	{"gl:", "glab", "https://gitlab.com/"},
}

// Source is a parsed template source.
type Source struct {
	// Raw is the source as given by the user.
	Raw  string
	Kind Kind
	// Location is the clone URL for git sources and the absolute path otherwise.
	Location string
	// CLI and Repo are set for gh: and gl: aliases.
	CLI  string
	Repo string
}

// Parse resolves raw into a Source. Local paths and ZIP files must exist.
func Parse(raw string) (Source, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return Source{}, fmt.Errorf("source cannot be empty")
	}

	for _, a := range aliases {
		if repo, ok := strings.CutPrefix(raw, a.prefix); ok {
			repo = strings.Trim(repo, "/")
			if repo == "" {
				return Source{}, fmt.Errorf("missing repository in '%s'", raw)
			}
			return Source{
				Raw:      raw,
				Kind:     Git,
				Location: a.baseURL + strings.TrimSuffix(repo, ".git") + ".git",
				CLI:      a.cli,
				Repo:     repo,
			}, nil
		}
	}

	if isGitURL(raw) {
		return Source{Raw: raw, Kind: Git, Location: raw}, nil
	}

	abs, err := filepath.Abs(raw)
	if err != nil {
		return Source{}, fmt.Errorf("invalid path '%s': %w", raw, err)
	}
	info, err := os.Stat(abs)

	if strings.HasSuffix(strings.ToLower(raw), ".zip") {
		if err != nil || info.IsDir() {
			return Source{}, fmt.Errorf("ZIP file not found: '%s'", abs)
		}
		return Source{Raw: raw, Kind: Zip, Location: abs}, nil
	}

	if err != nil {
		return Source{}, fmt.Errorf("directory not found: '%s'", abs)
	}
	if !info.IsDir() {
		return Source{}, fmt.Errorf("not a directory or ZIP file: '%s'", abs)
	}
	return Source{Raw: raw, Kind: Local, Location: abs}, nil
}

func isGitURL(s string) bool {
	return strings.HasPrefix(s, "http://") ||
		strings.HasPrefix(s, "https://") ||
		strings.HasPrefix(s, "git@") ||
		strings.HasSuffix(s, ".git")
}

// Name derives a template name from the source: the repository, directory or
// ZIP file name.
func (s Source) Name() string {
	loc := strings.TrimRight(s.Location, "/")
	var base string
	switch s.Kind {
	case Git:
		// Covers URLs, scp-like git@host:owner/repo and local paths.
		if i := strings.LastIndexAny(loc, "/:"); i >= 0 {
			base = loc[i+1:]
		} else {
			base = loc
		}
		base = strings.TrimSuffix(base, ".git")
	case Zip:
		base = filepath.Base(loc)
		base = strings.TrimSuffix(base, filepath.Ext(base))
	default:
		base = filepath.Base(loc)
	}
	return base
}

// Display returns the source as shown to the user.
func (s Source) Display() string {
	if s.CLI != "" {
		return s.Raw
	}
	return s.Location
}

// Fetch writes the template into dest, which must not exist, and applies the
// template's ignore patterns. Command output goes to out. On failure dest is
// removed.
func (s Source) Fetch(dest string, out io.Writer) error {
	if _, err := os.Lstat(dest); err == nil {
		return fmt.Errorf("destination already exists: '%s'", dest)
	}

	if err := s.fetch(dest, out); err != nil {
		os.RemoveAll(dest)
		return err
	}
	return nil
}

func (s Source) fetch(dest string, out io.Writer) error {
	switch s.Kind {
	case Git:
		var cmd *exec.Cmd
		if s.CLI != "" && commandExists(s.CLI) {
			cmd = exec.Command(s.CLI, "repo", "clone", s.Repo, dest, "--", "--depth", "1")
		} else {
			cmd = exec.Command("git", "clone", "--depth", "1", s.Location, dest)
		}
		cmd.Stdout = out
		cmd.Stderr = out
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to clone repository: %w", err)
		}
	case Zip:
		if err := fileutil.UnzipToDir(s.Location, dest); err != nil {
			return fmt.Errorf("failed to extract ZIP: %w", err)
		}
	case Local:
		if err := fileutil.CopyDir(s.Location, dest); err != nil {
			return fmt.Errorf("failed to copy template: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("unknown source kind '%s'", s.Kind)
	}

	if err := fileutil.ApplyIgnorePatterns(dest); err != nil {
		return fmt.Errorf("failed to apply ignore patterns: %w", err)
	}
	return nil
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
