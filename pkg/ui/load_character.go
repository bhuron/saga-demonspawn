package ui

import (
	"os"
	"path/filepath"
	"strings"
)

// LoadCharacterModel represents the load character screen state.
type LoadCharacterModel struct {
	files  []string
	cursor int
	err    error
}

// NewLoadCharacterModel creates a new load character model.
func NewLoadCharacterModel() LoadCharacterModel {
	return LoadCharacterModel{
		files:  []string{},
		cursor: 0,
		err:    nil,
	}
}

// Refresh scans the current directory for character save files.
func (m *LoadCharacterModel) Refresh() {
	m.files = []string{}
	m.cursor = 0
	m.err = nil

	// Look for JSON files in current directory
	files, err := filepath.Glob("character_*.json")
	if err != nil {
		m.err = err
		return
	}

	m.files = files
}

// MoveUp moves the cursor up.
func (m *LoadCharacterModel) MoveUp() {
	if m.cursor > 0 {
		m.cursor--
	}
}

// MoveDown moves the cursor down.
func (m *LoadCharacterModel) MoveDown() {
	if m.cursor < len(m.files)-1 {
		m.cursor++
	}
}

// GetCursor returns the current cursor position.
func (m *LoadCharacterModel) GetCursor() int {
	return m.cursor
}

// GetFiles returns all available save files.
func (m *LoadCharacterModel) GetFiles() []string {
	return m.files
}

// GetSelectedFile returns the currently selected file.
func (m *LoadCharacterModel) GetSelectedFile() string {
	if m.cursor < len(m.files) {
		return m.files[m.cursor]
	}
	return ""
}

// HasFiles returns true if there are save files available.
func (m *LoadCharacterModel) HasFiles() bool {
	return len(m.files) > 0
}

// GetError returns any error encountered during refresh.
func (m *LoadCharacterModel) GetError() error {
	return m.err
}

// GetFileInfo returns formatted information about a save file.
func GetFileInfo(filepath string) string {
	info, err := os.Stat(filepath)
	if err != nil {
		return ""
	}

	// Extract character name from filename if possible
	basename := strings.TrimPrefix(filepath, "character_")
	basename = strings.TrimSuffix(basename, ".json")
	
	// Format: filename (size, modified date)
	return info.ModTime().Format("2006-01-02 15:04")
}
