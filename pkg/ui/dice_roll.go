package ui

import (
	"fmt"
	"strings"

	"github.com/benoit/saga-demonspawn/internal/dice"
	"github.com/benoit/saga-demonspawn/pkg/ui/theme"
)

// DiceRollModel handles the dice rolling screen state.
type DiceRollModel struct {
	dice   dice.Roller
	result int
	rolled bool
	msg    string
}

// NewDiceRollModel creates a new dice roll model.
func NewDiceRollModel(roller dice.Roller) DiceRollModel {
	return DiceRollModel{
		dice:   roller,
		result: 0,
		rolled: false,
		msg:    "Press '1' for 1d6, '2' for 2d6, or 'ESC' to return.",
	}
}

// Reset resets the model state.
func (m *DiceRollModel) Reset() {
	m.result = 0
	m.rolled = false
	m.msg = "Press '1' for 1d6, '2' for 2d6, or 'ESC' to return."
}

// Roll1D6 rolls a single 6-sided die.
func (m *DiceRollModel) Roll1D6() {
	m.result = m.dice.Roll1D6()
	m.rolled = true
	m.msg = fmt.Sprintf("You rolled 1d6: %d", m.result)
}

// Roll2D6 rolls two 6-sided dice.
func (m *DiceRollModel) Roll2D6() {
	m.result = m.dice.Roll2D6()
	m.rolled = true
	m.msg = fmt.Sprintf("You rolled 2d6: %d", m.result)
}

// View renders the dice roll screen.
func (m DiceRollModel) View() string {
	var s strings.Builder

	s.WriteString(theme.RenderTitle("Dice Roller"))
	s.WriteString("\n\n")

	if m.rolled {
		s.WriteString(fmt.Sprintf("Result: %d\n\n", m.result))
	}

	s.WriteString(m.msg)
	s.WriteString("\n\n")
	s.WriteString(theme.Current().MutedText.Render("Press '1' or '2' to roll again, or 'ESC' to exit."))

	return s.String()
}
