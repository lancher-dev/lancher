package config

import (
	"fmt"
	"os"
	"path/filepath"

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
	Ignore      []string `yaml:"ignore"`
}

// LoadResult contains the loaded config and metadata about the loading process
type LoadResult struct {
	Config     *Config
	FoundFiles []string // All config files found
	UsedFile   string   // The config file that was loaded
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
	}

	return result
}

// HasHooks returns true if config has hooks defined
func (c *Config) HasHooks() bool {
	return c != nil && len(c.Hooks) > 0
}

// HasIgnore returns true if config has deprecated ignore patterns defined in YAML
func (c *Config) HasIgnore() bool {
	return c != nil && len(c.Ignore) > 0
}

// GetMetadata returns formatted metadata string for display
func (c *Config) GetMetadata() string {
	if c == nil {
		return ""
	}

	var metadata string
	if c.Name != "" {
		metadata += fmt.Sprintf("Name: %s\n", c.Name)
	}
	if c.Description != "" {
		metadata += fmt.Sprintf("Description: %s\n", c.Description)
	}
	if c.Author != "" {
		metadata += fmt.Sprintf("Author: %s\n", c.Author)
	}
	if c.Version != "" {
		metadata += fmt.Sprintf("Version: %s\n", c.Version)
	}
	return metadata
}
