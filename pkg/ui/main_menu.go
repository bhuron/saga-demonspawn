package ui

// MainMenuModel represents the state of the main menu screen.
type MainMenuModel struct {
	cursor int
	choices []string
}

// NewMainMenuModel creates a new main menu model.
func NewMainMenuModel() MainMenuModel {
	return MainMenuModel{
		cursor: 0,
		choices: []string{
			"New Character",
			"Load Character",
			"Settings",
			"Help",
			"Exit",
		},
	}
}

// MoveUp moves the cursor up in the menu.
func (m *MainMenuModel) MoveUp() {
	if m.cursor > 0 {
		m.cursor--
	}
}

// MoveDown moves the cursor down in the menu.
func (m *MainMenuModel) MoveDown() {
	if m.cursor < len(m.choices)-1 {
		m.cursor++
	}
}

// GetSelected returns the currently selected menu item.
func (m *MainMenuModel) GetSelected() string {
	return m.choices[m.cursor]
}

// GetCursor returns the current cursor position.
func (m *MainMenuModel) GetCursor() int {
	return m.cursor
}

// GetChoices returns all menu choices.
func (m *MainMenuModel) GetChoices() []string {
	return m.choices
}
