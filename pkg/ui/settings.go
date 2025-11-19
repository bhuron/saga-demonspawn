package ui

import (
	"github.com/benoit/saga-demonspawn/internal/config"
)

// SettingField represents individual settings that can be edited.
type SettingField int

const (
	SettingTheme SettingField = iota
	SettingUseUnicode
	SettingShowAnimations
	SettingConfirmActions
	SettingAutoSave
	SettingShowRollDetails
	SettingHighContrast
	SettingReducedMotion
	SettingSave
	SettingCancel
	SettingReset
)

// SettingsModel handles the settings screen state.
type SettingsModel struct {
	config       *config.Config  // Working copy of configuration
	originalConfig *config.Config // Original config for cancel
	cursor       int
	saved        bool
	message      string
}

// NewSettingsModel creates a new settings model.
func NewSettingsModel(cfg *config.Config) SettingsModel {
	// Create a copy for editing
	workingCfg := *cfg
	originalCfg := *cfg
	
	return SettingsModel{
		config:       &workingCfg,
		originalConfig: &originalCfg,
		cursor:       0,
		saved:        false,
		message:      "",
	}
}

// Reset resets the model with new config.
func (m *SettingsModel) Reset(cfg *config.Config) {
	workingCfg := *cfg
	originalCfg := *cfg
	m.config = &workingCfg
	m.originalConfig = &originalCfg
	m.cursor = 0
	m.saved = false
	m.message = ""
}

// GetCursor returns the current cursor position.
func (m SettingsModel) GetCursor() int {
	return m.cursor
}

// MoveCursorUp moves cursor up.
func (m *SettingsModel) MoveCursorUp() {
	if m.cursor > 0 {
		m.cursor--
	}
}

// MoveCursorDown moves cursor down.
func (m *SettingsModel) MoveCursorDown() {
	maxCursor := int(SettingReset)
	if m.cursor < maxCursor {
		m.cursor++
	}
}

// GetConfig returns the current working configuration.
func (m *SettingsModel) GetConfig() *config.Config {
	return m.config
}

// ToggleCurrentSetting toggles the currently selected boolean setting.
func (m *SettingsModel) ToggleCurrentSetting() {
	field := SettingField(m.cursor)
	
	switch field {
	case SettingUseUnicode:
		m.config.UseUnicode = !m.config.UseUnicode
	case SettingShowAnimations:
		m.config.ShowAnimations = !m.config.ShowAnimations
	case SettingConfirmActions:
		m.config.ConfirmActions = !m.config.ConfirmActions
	case SettingAutoSave:
		m.config.AutoSave = !m.config.AutoSave
	case SettingShowRollDetails:
		m.config.ShowRollDetails = !m.config.ShowRollDetails
	case SettingHighContrast:
		m.config.HighContrast = !m.config.HighContrast
	case SettingReducedMotion:
		m.config.ReducedMotion = !m.config.ReducedMotion
	}
}

// CycleTheme cycles through available themes.
func (m *SettingsModel) CycleTheme() {
	switch m.config.Theme {
	case "dark":
		m.config.Theme = "light"
	case "light":
		m.config.Theme = "dark"
	default:
		m.config.Theme = "dark"
	}
}

// Save saves the current configuration.
func (m *SettingsModel) Save() error {
	if err := m.config.SaveDefault(); err != nil {
		m.message = "Failed to save configuration: " + err.Error()
		m.saved = false
		return err
	}
	
	// Update original to match saved
	*m.originalConfig = *m.config
	m.message = "Configuration saved successfully!"
	m.saved = true
	return nil
}

// Cancel reverts changes to original configuration.
func (m *SettingsModel) Cancel() {
	*m.config = *m.originalConfig
	m.message = "Changes canceled"
	m.saved = false
}

// ResetToDefaults resets all settings to default values.
func (m *SettingsModel) ResetToDefaults() {
	defaults := config.Default()
	*m.config = *defaults
	m.message = "Reset to defaults"
	m.saved = false
}

// GetMessage returns the current status message.
func (m SettingsModel) GetMessage() string {
	return m.message
}

// ClearMessage clears the status message.
func (m *SettingsModel) ClearMessage() {
	m.message = ""
}

// IsSaved returns whether settings have been saved.
func (m SettingsModel) IsSaved() bool {
	return m.saved
}
