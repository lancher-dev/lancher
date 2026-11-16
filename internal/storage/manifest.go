package storage

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"go.yaml.in/yaml/v3"

	"github.com/lancher-dev/lancher/internal/source"
)

// ManifestFile is the name of the manifest, stored next to the templates
// directory.
const ManifestFile = "templates.yaml"

// UnknownKind marks a template whose source is not recorded.
const UnknownKind source.Kind = "unknown"

const manifestVersion = 1

// Entry records where a template comes from.
type Entry struct {
	Kind   source.Kind `yaml:"kind" json:"kind"`
	Source string      `yaml:"source,omitempty" json:"source,omitempty"`
	Ref    string      `yaml:"ref,omitempty" json:"ref,omitempty"`
	Added  time.Time   `yaml:"added,omitempty" json:"added,omitzero"`
	Synced time.Time   `yaml:"synced,omitempty" json:"synced,omitzero"`
}

// Known reports whether the source of the template is recorded.
func (e Entry) Known() bool {
	return e.Source != "" && e.Kind != UnknownKind
}

// Manifest maps template names to their entries.
type Manifest struct {
	Version   int              `yaml:"version"`
	Templates map[string]Entry `yaml:"templates"`

	path string
}

// ManifestPath returns the path of the manifest.
func ManifestPath() (string, error) {
	dir, err := GetTemplatesDir()
	if err != nil {
		return "", err
	}
	return manifestPath(dir), nil
}

func manifestPath(templatesDir string) string {
	return filepath.Join(filepath.Dir(templatesDir), ManifestFile)
}

// LoadManifest reads the manifest and aligns it with the template
// directories: entries without a directory are dropped and directories
// without an entry are registered. The file is rewritten when anything
// changes.
func LoadManifest() (*Manifest, error) {
	templatesDir, err := GetTemplatesDir()
	if err != nil {
		return nil, err
	}
	path := manifestPath(templatesDir)
	m := &Manifest{}
	data, err := os.ReadFile(path)
	switch {
	case err == nil:
		if err := yaml.Unmarshal(data, m); err != nil {
			return nil, fmt.Errorf("invalid manifest %s (fix or delete it): %w", path, err)
		}
		if m.Version > manifestVersion {
			return nil, fmt.Errorf("manifest %s has version %d, newer than this lancher supports", path, m.Version)
		}
	case !os.IsNotExist(err):
		return nil, err
	}
	m.path = path
	m.Version = manifestVersion
	if m.Templates == nil {
		m.Templates = map[string]Entry{}
	}

	names, err := ListTemplates()
	if err != nil {
		return nil, err
	}
	changed := false
	for name := range m.Templates {
		if !slices.Contains(names, name) {
			delete(m.Templates, name)
			changed = true
		}
	}
	for _, name := range names {
		if _, ok := m.Templates[name]; !ok {
			m.Templates[name] = detect(filepath.Join(templatesDir, name))
			changed = true
		}
	}
	if changed {
		if err := m.Save(); err != nil {
			return nil, err
		}
	}
	return m, nil
}

// Save writes the manifest atomically.
func (m *Manifest) Save() error {
	var data bytes.Buffer
	enc := yaml.NewEncoder(&data)
	enc.SetIndent(2)
	if err := enc.Encode(m); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(m.path), ".templates-*.yaml")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(data.Bytes()); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmp.Name(), m.path); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}
	return nil
}

// detect builds the entry of an unregistered template from its git remote.
func detect(dir string) Entry {
	if remote := gitOutput(dir, "config", "--get", "remote.origin.url"); remote != "" {
		return Entry{Kind: source.Git, Source: remote, Ref: GitRef(dir)}
	}
	return Entry{Kind: UnknownKind}
}

// GitRef returns the commit checked out in dir, or "" when dir has no
// repository of its own.
func GitRef(dir string) string {
	return gitOutput(dir, "rev-parse", "HEAD")
}

func gitOutput(dir string, args ...string) string {
	if _, err := os.Stat(filepath.Join(dir, ".git")); err != nil {
		return ""
	}
	out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
