package ui

import "github.com/benoit/saga-demonspawn/internal/character"

// CharacterViewModel represents the character view screen state.
type CharacterViewModel struct {
	character *character.Character
}

// NewCharacterViewModel creates a new character view model.
func NewCharacterViewModel() CharacterViewModel {
	return CharacterViewModel{
		character: nil,
	}
}

// SetCharacter sets the character to display.
func (m *CharacterViewModel) SetCharacter(char *character.Character) {
	m.character = char
}

// GetCharacter returns the current character.
func (m *CharacterViewModel) GetCharacter() *character.Character {
	return m.character
}
