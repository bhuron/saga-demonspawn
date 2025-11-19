package ui

// GameSessionModel represents the game session menu state.
type GameSessionModel struct {
	cursor  int
	choices []string
	showMagic bool // Whether to show Cast Spell option
}

// NewGameSessionModel creates a new game session model.
func NewGameSessionModel() GameSessionModel {
	return GameSessionModel{
		cursor: 0,
		choices: []string{
			"View Character",
			"Edit Character Stats",
			"Combat",
			"Manage Inventory",
			"Roll Dice",
			"Save & Exit",
		},
		showMagic: false,
	}
}

// UpdateMagicVisibility sets whether the magic option should be shown.
func (m *GameSessionModel) UpdateMagicVisibility(magicUnlocked bool) {
	m.showMagic = magicUnlocked
	// Rebuild choices based on magic availability
	if magicUnlocked {
		m.choices = []string{
			"View Character",
			"Edit Character Stats",
			"Combat",
			"Cast Spell",
			"Manage Inventory",
			"Roll Dice",
			"Save & Exit",
		}
	} else {
		m.choices = []string{
			"View Character",
			"Edit Character Stats",
			"Combat",
			"Manage Inventory",
			"Roll Dice",
			"Save & Exit",
		}
	}
}

// MoveUp moves the cursor up.
func (m *GameSessionModel) MoveUp() {
	if m.cursor > 0 {
		m.cursor--
	}
}

// MoveDown moves the cursor down.
func (m *GameSessionModel) MoveDown() {
	if m.cursor < len(m.choices)-1 {
		m.cursor++
	}
}

// GetSelected returns the currently selected menu item.
func (m *GameSessionModel) GetSelected() string {
	return m.choices[m.cursor]
}

// GetCursor returns the current cursor position.
func (m *GameSessionModel) GetCursor() int {
	return m.cursor
}

// GetChoices returns all menu choices.
func (m *GameSessionModel) GetChoices() []string {
	return m.choices
}
