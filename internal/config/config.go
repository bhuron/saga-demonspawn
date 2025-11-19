// Package config manages application configuration and user preferences.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds all user preferences and settings.
type Config struct {
	// Appearance settings
	Theme          string `json:"theme"`           // "dark", "light", or "custom"
	UseUnicode     bool   `json:"use_unicode"`     // Enable Unicode characters vs ASCII
	ShowAnimations bool   `json:"show_animations"` // Enable visual transitions

	// Gameplay settings
	ConfirmActions bool `json:"confirm_actions"` // Require confirmation for risky actions
	AutoSave       bool `json:"auto_save"`       // Save character on application exit
	ShowRollDetails bool `json:"show_roll_details"` // Display dice roll breakdowns

	// Accessibility settings
	HighContrast  bool `json:"high_contrast"`  // Accessibility: enhanced contrast
	ReducedMotion bool `json:"reduced_motion"` // Reduce visual effects

	// File settings
	SaveDirectory   string `json:"save_directory"`    // Character save location
	CharacterBackup bool   `json:"character_backup"`  // Create timestamped backups
}

// Default returns a configuration with default values.
func Default() *Config {
	homeDir, err := os.UserHomeDir()
	saveDir := "."
	if err == nil {
		saveDir = filepath.Join(homeDir, ".saga-demonspawn")
	}

	return &Config{
		Theme:           "dark",
		UseUnicode:      true,
		ShowAnimations:  true,
		ConfirmActions:  true,
		AutoSave:        true,
		ShowRollDetails: true,
		HighContrast:    false,
		ReducedMotion:   false,
		SaveDirectory:   saveDir,
		CharacterBackup: true,
	}
}

// Load loads configuration from the specified path.
// If the file doesn't exist, returns default configuration.
func Load(path string) (*Config, error) {
	// If file doesn't exist, return default
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return Default(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Default(), fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Default(), fmt.Errorf("failed to parse config: %w", err)
	}

	// Validate and apply defaults for any missing fields
	cfg.ApplyDefaults()
	
	if err := cfg.Validate(); err != nil {
		// Return default on validation error
		return Default(), fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

// Save saves the configuration to the specified path.
func (c *Config) Save(path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// Validate checks if the configuration values are valid.
func (c *Config) Validate() error {
	// Validate theme
	validThemes := map[string]bool{
		"dark":   true,
		"light":  true,
		"custom": true,
	}
	if !validThemes[c.Theme] {
		return fmt.Errorf("invalid theme: %s (must be dark, light, or custom)", c.Theme)
	}

	// Validate save directory exists or can be created
	if c.SaveDirectory != "" {
		if err := os.MkdirAll(c.SaveDirectory, 0755); err != nil {
			// Try to use current directory as fallback
			c.SaveDirectory = "."
		}
	}

	return nil
}

// ApplyDefaults fills in any missing fields with default values.
func (c *Config) ApplyDefaults() {
	defaults := Default()

	if c.Theme == "" {
		c.Theme = defaults.Theme
	}
	if c.SaveDirectory == "" {
		c.SaveDirectory = defaults.SaveDirectory
	}
}

// GetConfigPath returns the default configuration file path.
func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "config.json"
	}
	return filepath.Join(homeDir, ".saga-demonspawn", "config.json")
}

// LoadDefault loads configuration from the default location.
func LoadDefault() (*Config, error) {
	return Load(GetConfigPath())
}

// SaveDefault saves configuration to the default location.
func (c *Config) SaveDefault() error {
	return c.Save(GetConfigPath())
}
