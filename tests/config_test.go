package tests

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/lancher-dev/lancher/internal/config"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		yamlContent string
		wantName    string
		wantDesc    string
		wantAuthor  string
		wantVersion string
		wantHooks   int
		wantErr     bool
	}{
		{
			name: "valid config with all fields",
			yamlContent: `name: Test Template
description: A test template
author: Test Author
version: 1.0.0
hooks:
  - npm install
  - git init`,
			wantName:    "Test Template",
			wantDesc:    "A test template",
			wantAuthor:  "Test Author",
			wantVersion: "1.0.0",
			wantHooks:   2,
			wantErr:     false,
		},
		{
			name: "minimal config",
			yamlContent: `name: Minimal
description: Minimal template`,
			wantName:  "Minimal",
			wantDesc:  "Minimal template",
			wantHooks: 0,
			wantErr:   false,
		},
		{
			name:        "invalid yaml",
			yamlContent: "name: [unclosed\nhooks:\n  - npm install",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, config.ConfigFileNames[0])

			// Write config file
			if err := os.WriteFile(configPath, []byte(tt.yamlContent), 0644); err != nil {
				t.Fatalf("failed to write test config: %v", err)
			}

			// Load config
			cfg, err := config.LoadConfig(tmpDir)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if cfg == nil {
				t.Fatal("expected config but got nil")
			}

			if cfg.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", cfg.Name, tt.wantName)
			}
			if cfg.Description != tt.wantDesc {
				t.Errorf("Description = %q, want %q", cfg.Description, tt.wantDesc)
			}
			if cfg.Author != tt.wantAuthor {
				t.Errorf("Author = %q, want %q", cfg.Author, tt.wantAuthor)
			}
			if cfg.Version != tt.wantVersion {
				t.Errorf("Version = %q, want %q", cfg.Version, tt.wantVersion)
			}
			if len(cfg.Hooks) != tt.wantHooks {
				t.Errorf("Hooks count = %d, want %d", len(cfg.Hooks), tt.wantHooks)
			}
		})
	}
}

func TestLoadConfigNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	cfg, err := config.LoadConfig(tmpDir)

	if err != nil {
		t.Errorf("expected no error when config not found, got: %v", err)
	}
	if cfg != nil {
		t.Errorf("expected nil config when file not found, got: %+v", cfg)
	}
}

func TestHasHooks(t *testing.T) {
	tests := []struct {
		name string
		cfg  *config.Config
		want bool
	}{
		{
			name: "config with hooks",
			cfg: &config.Config{
				Hooks: []string{"npm install", "git init"},
			},
			want: true,
		},
		{
			name: "config without hooks",
			cfg: &config.Config{
				Hooks: []string{},
			},
			want: false,
		},
		{
			name: "nil hooks",
			cfg:  &config.Config{},
			want: false,
		},
		{
			name: "nil config",
			cfg:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.HasHooks()
			if got != tt.want {
				t.Errorf("HasHooks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetadata(t *testing.T) {
	cfg := &config.Config{Name: "Minimal", Version: "0.1.0"}
	want := []config.Field{{Label: "Name", Value: "Minimal"}, {Label: "Version", Value: "0.1.0"}}
	if got := cfg.Metadata(); !slices.Equal(got, want) {
		t.Errorf("Metadata() = %v, want %v", got, want)
	}

	var nilCfg *config.Config
	if got := nilCfg.Metadata(); len(got) != 0 {
		t.Errorf("Metadata() on nil config = %v, want none", got)
	}
}

func TestLoadConfigWithDetailsInvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, ".lancher.yaml"), []byte("name: [unclosed"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "lancher.yml"), []byte("name: fallback"), 0644); err != nil {
		t.Fatal(err)
	}

	result := config.LoadConfigWithDetails(tmpDir)

	if result.Err == nil {
		t.Fatal("expected an error for invalid YAML")
	}
	if result.Config != nil {
		t.Errorf("expected no config, got %+v (lower-priority files must not be used as fallback)", result.Config)
	}
	if result.UsedFile != ".lancher.yaml" {
		t.Errorf("UsedFile = %q, want .lancher.yaml", result.UsedFile)
	}
	if len(result.FoundFiles) != 2 {
		t.Errorf("FoundFiles = %v, want 2 entries", result.FoundFiles)
	}
}

func TestUnknownKeys(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "lancher.yml"), "name: x\nignore:\n  - dist\nhook:\n  - ls\n")

	result := config.LoadConfigWithDetails(dir)
	if result.Err != nil {
		t.Fatal(result.Err)
	}
	if want := []string{"hook", "ignore"}; !slices.Equal(result.Unknown, want) {
		t.Errorf("Unknown = %v, want %v", result.Unknown, want)
	}
}
