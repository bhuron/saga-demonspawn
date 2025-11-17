package ui

import "github.com/benoit/saga-demonspawn/internal/character"

// EditField represents which field is currently being edited.
type EditField int

const (
	EditFieldStrength EditField = iota
	EditFieldSpeed
	EditFieldStamina
	EditFieldCourage
	EditFieldLuck
	EditFieldCharm
	EditFieldAttraction
	EditFieldCurrentLP
	EditFieldMaxLP
	EditFieldSkill
	EditFieldCurrentPOW
	EditFieldMaxPOW
)

// CharacterEditModel represents the character editing screen state.
type CharacterEditModel struct {
	character    *character.Character
	originalChar *character.Character // Backup for canceling
	cursor       int
	fields       []string
	inputMode    bool   // Whether we're actively editing a value
	inputBuffer  string // Current input being typed
}

// NewCharacterEditModel creates a new character edit model.
func NewCharacterEditModel() CharacterEditModel {
	return CharacterEditModel{
		character:    nil,
		originalChar: nil,
		cursor:       0,
		fields: []string{
			"Strength",
			"Speed",
			"Stamina",
			"Courage",
			"Luck",
			"Charm",
			"Attraction",
			"Current LP",
			"Maximum LP",
			"Skill",
			"Current POW",
			"Maximum POW",
		},
		inputMode:   false,
		inputBuffer: "",
	}
}

// SetCharacter sets the character to edit and creates a backup.
func (m *CharacterEditModel) SetCharacter(char *character.Character) {
	m.character = char
	// Create a deep copy for backup (we'll implement proper cloning if needed)
	m.originalChar = char
}

// GetCharacter returns the current character.
func (m *CharacterEditModel) GetCharacter() *character.Character {
	return m.character
}

// MoveUp moves the cursor up.
func (m *CharacterEditModel) MoveUp() {
	if !m.inputMode && m.cursor > 0 {
		m.cursor--
	}
}

// MoveDown moves the cursor down.
func (m *CharacterEditModel) MoveDown() {
	maxCursor := len(m.fields) - 1
	// Don't show POW fields if magic is not unlocked
	if m.character != nil && !m.character.MagicUnlocked {
		maxCursor -= 2
	}
	
	if !m.inputMode && m.cursor < maxCursor {
		m.cursor++
	}
}

// GetCursor returns the current cursor position.
func (m *CharacterEditModel) GetCursor() int {
	return m.cursor
}

// GetFields returns all editable fields.
func (m *CharacterEditModel) GetFields() []string {
	if m.character != nil && !m.character.MagicUnlocked {
		// Exclude POW fields
		return m.fields[:len(m.fields)-2]
	}
	return m.fields
}

// IsInputMode returns whether we're in input mode.
func (m *CharacterEditModel) IsInputMode() bool {
	return m.inputMode
}

// StartInput begins editing the current field.
func (m *CharacterEditModel) StartInput() {
	m.inputMode = true
	m.inputBuffer = ""
}

// CancelInput cancels the current input.
func (m *CharacterEditModel) CancelInput() {
	m.inputMode = false
	m.inputBuffer = ""
}

// GetInputBuffer returns the current input.
func (m *CharacterEditModel) GetInputBuffer() string {
	return m.inputBuffer
}

// AppendInput adds a character to the input buffer.
func (m *CharacterEditModel) AppendInput(char string) {
	m.inputBuffer += char
}

// Backspace removes the last character from input.
func (m *CharacterEditModel) Backspace() {
	if len(m.inputBuffer) > 0 {
		m.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
	}
}

// GetCurrentValue returns the current value of the selected field.
func (m *CharacterEditModel) GetCurrentValue() int {
	if m.character == nil {
		return 0
	}

	switch EditField(m.cursor) {
	case EditFieldStrength:
		return m.character.Strength
	case EditFieldSpeed:
		return m.character.Speed
	case EditFieldStamina:
		return m.character.Stamina
	case EditFieldCourage:
		return m.character.Courage
	case EditFieldLuck:
		return m.character.Luck
	case EditFieldCharm:
		return m.character.Charm
	case EditFieldAttraction:
		return m.character.Attraction
	case EditFieldCurrentLP:
		return m.character.CurrentLP
	case EditFieldMaxLP:
		return m.character.MaximumLP
	case EditFieldSkill:
		return m.character.Skill
	case EditFieldCurrentPOW:
		return m.character.CurrentPOW
	case EditFieldMaxPOW:
		return m.character.MaximumPOW
	default:
		return 0
	}
}
