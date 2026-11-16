package config

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"go.yaml.in/yaml/v3"
)

// ConfigFileNames lists all supported configuration file names in order of priority
var ConfigFileNames = []string{
	".lancher.yaml",
	".lancher.yml",
	"lancher.yaml",
	"lancher.yml",
}

// Config represents the lancher configuration file
type Config struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Author      string   `yaml:"author"`
	Version     string   `yaml:"version"`
	Hooks       []string `yaml:"hooks"`
}

// LoadResult contains the loaded config and metadata about the loading process
type LoadResult struct {
	Config     *Config
	FoundFiles []string // All config files found
	UsedFile   string   // The config file that was loaded
	Unknown    []string // Top-level keys of UsedFile that lancher does not use
	Err        error    // Read or parse error of UsedFile
}

// LoadConfig loads the highest-priority config file from the template directory.
// It returns nil and no error when no config file exists.
func LoadConfig(templatePath string) (*Config, error) {
	result := LoadConfigWithDetails(templatePath)
	return result.Config, result.Err
}

// LoadConfigWithDetails loads the highest-priority config file and lists all
// config files found. A file that cannot be read or parsed sets Err.
func LoadConfigWithDetails(templatePath string) *LoadResult {
	result := &LoadResult{
		FoundFiles: []string{},
	}

	for _, fileName := range ConfigFileNames {
		configPath := filepath.Join(templatePath, fileName)
		if _, err := os.Stat(configPath); err != nil {
			continue
		}
		result.FoundFiles = append(result.FoundFiles, fileName)

		if result.UsedFile != "" {
			continue
		}
		result.UsedFile = fileName

		data, err := os.ReadFile(configPath)
		if err != nil {
			result.Err = fmt.Errorf("failed to read %s: %w", fileName, err)
			continue
		}

		var cfg Config
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			result.Err = fmt.Errorf("invalid YAML in %s: %w", fileName, err)
			continue
		}
		result.Config = &cfg
		result.Unknown = unknownKeys(data)
	}

	return result
}

var knownKeys = []string{"name", "description", "author", "version", "hooks"}

func unknownKeys(data []byte) []string {
	var raw map[string]any
	if yaml.Unmarshal(data, &raw) != nil {
		return nil
	}
	var unknown []string
	for key := range raw {
		if !slices.Contains(knownKeys, key) {
			unknown = append(unknown, key)
		}
	}
	slices.Sort(unknown)
	return unknown
}

// HasHooks returns true if config has hooks defined
func (c *Config) HasHooks() bool {
	return c != nil && len(c.Hooks) > 0
}

// Field is a labelled metadata value.
type Field struct {
	Label, Value string
}

// Metadata returns the metadata fields that are set, in display order.
func (c *Config) Metadata() []Field {
	if c == nil {
		return nil
	}
	var fields []Field
	for _, f := range []Field{
		{"Name", c.Name},
		{"Description", c.Description},
		{"Author", c.Author},
		{"Version", c.Version},
	} {
		if f.Value != "" {
			fields = append(fields, f)
		}
	}
	return fields
}
