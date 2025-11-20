package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const ConfigFileName = ".lancher.yaml"

// Config represents the .lancher.yaml configuration file
type Config struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Author      string   `yaml:"author"`
	Version     string   `yaml:"version"`
	Hooks       []string `yaml:"hooks"`
	Ignore      []string `yaml:"ignore"`
}

// LoadConfig loads .lancher.yaml from the template directory
func LoadConfig(templatePath string) (*Config, error) {
	configPath := filepath.Join(templatePath, ConfigFileName)

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, nil // No config file, return nil without error
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// ShouldIgnore checks if a file should be ignored during template creation
func (c *Config) ShouldIgnore(relativePath string) bool {
	if c == nil {
		return false
	}

	for _, pattern := range c.Ignore {
		if matched, _ := filepath.Match(pattern, relativePath); matched {
			return true
		}
		// Also check if pattern matches any part of the path
		if matched, _ := filepath.Match(pattern, filepath.Base(relativePath)); matched {
			return true
		}
	}
	return false
}

// HasHooks returns true if config has hooks defined
func (c *Config) HasHooks() bool {
	return c != nil && len(c.Hooks) > 0
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
