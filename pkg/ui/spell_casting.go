package ui

import (
	"fmt"
	"strings"

	"github.com/benoit/saga-demonspawn/internal/character"
	"github.com/benoit/saga-demonspawn/internal/dice"
	"github.com/benoit/saga-demonspawn/internal/magic"
	"github.com/benoit/saga-demonspawn/pkg/ui/theme"
)

// SpellCastingModel represents the spell casting screen state.
type SpellCastingModel struct {
	character       *character.Character
	roller          dice.Roller
	cursor          int
	spells          []magic.Spell
	inCombat        bool
	message         string
	awaitingConfirm bool   // Whether awaiting sacrifice confirmation
	confirmSpell    string // Spell awaiting confirmation
	sacrificeAmount int    // Amount of LP to sacrifice
	naturalCheckMsg string // Result of natural inclination check
	returnToCombat  bool   // Whether to return to combat screen on exit
}

// NewSpellCastingModel creates a new spell casting model.
func NewSpellCastingModel(char *character.Character, roller dice.Roller, inCombat bool) SpellCastingModel {
	isDead := char.CurrentLP <= 0
	availableSpells := magic.GetAvailableSpells(inCombat, isDead)

	return SpellCastingModel{
		character:       char,
		roller:          roller,
		cursor:          0,
		spells:          availableSpells,
		inCombat:        inCombat,
		message:         "",
		awaitingConfirm: false,
		naturalCheckMsg: "",
		returnToCombat:  inCombat, // Set to true if casting from combat
	}
}

// SetCharacter updates the character reference.
func (m *SpellCastingModel) SetCharacter(char *character.Character) {
	m.character = char
	isDead := char.CurrentLP <= 0
	m.spells = magic.GetAvailableSpells(m.inCombat, isDead)
}

// MoveUp moves the cursor up.
func (m *SpellCastingModel) MoveUp() {
	if m.cursor > 0 {
		m.cursor--
	}
}

// MoveDown moves the cursor down.
func (m *SpellCastingModel) MoveDown() {
	if m.cursor < len(m.spells) {
		m.cursor++
	}
}

// GetSelectedSpell returns the currently selected spell, or nil if "Natural Inclination Check" is selected.
func (m *SpellCastingModel) GetSelectedSpell() *magic.Spell {
	if m.cursor < len(m.spells) {
		return &m.spells[m.cursor]
	}
	return nil
}

// IsNaturalCheckSelected returns true if "Natural Inclination Check" option is selected.
func (m *SpellCastingModel) IsNaturalCheckSelected() bool {
	return m.cursor == len(m.spells)
}

// PerformNaturalCheck performs the natural inclination check.
func (m *SpellCastingModel) PerformNaturalCheck() {
	success, roll := magic.NaturalInclinationCheck(m.roller)
	if success {
		m.naturalCheckMsg = fmt.Sprintf("Natural Inclination Check: Rolled %d - Fire*Wolf overcomes his aversion to magic!", roll)
	} else {
		m.naturalCheckMsg = fmt.Sprintf("Natural Inclination Check: Rolled %d - Fire*Wolf refuses to use sorcery this section.", roll)
	}
}

// AttemptCast attempts to cast the currently selected spell.
// Returns true if casting should proceed (either immediate or after confirmation).
func (m *SpellCastingModel) AttemptCast() bool {
	spell := m.GetSelectedSpell()
	if spell == nil {
		return false
	}

	// Validate the cast
	isDead := m.character.CurrentLP <= 0
	result := magic.ValidateCast(spell, m.character.CurrentPOW, m.character.CurrentLP, m.inCombat, isDead)

	if result.RequiresSacrifice {
		// Need confirmation for LP sacrifice
		m.awaitingConfirm = true
		m.confirmSpell = spell.Name
		m.sacrificeAmount = result.SacrificeAmount
		m.message = result.Message
		return false
	}

	if !result.Success {
		// Cast failed validation
		m.message = result.Message
		return false
	}

	// Proceed with cast
	return true
}

// ConfirmSacrifice confirms LP sacrifice and proceeds with cast.
func (m *SpellCastingModel) ConfirmSacrifice() bool {
	m.awaitingConfirm = false
	// Sacrifice LP for POW
	m.character.ModifyLP(-m.sacrificeAmount)
	m.character.ModifyPOW(m.sacrificeAmount)
	m.message = fmt.Sprintf("Sacrificed %d LP for %d POW", m.sacrificeAmount, m.sacrificeAmount)
	return true
}

// CancelSacrifice cancels the sacrifice and returns to spell selection.
func (m *SpellCastingModel) CancelSacrifice() {
	m.awaitingConfirm = false
	m.message = "Sacrifice cancelled"
}

// PerformCast performs the actual spell cast (FFR check + effect).
// Returns the spell effect result.
func (m *SpellCastingModel) PerformCast() (magic.SpellEffect, bool) {
	spell := m.GetSelectedSpell()
	if spell == nil {
		return magic.SpellEffect{}, false
	}

	// Deduct power cost
	m.character.ModifyPOW(-spell.PowerCost)

	// Perform FFR check
	castResult := magic.PerformCast(spell, m.roller)
	if castResult.FFRFailed {
		m.message = castResult.Message
		return magic.SpellEffect{Success: false, Message: castResult.Message}, false
	}

	// Apply spell effect
	var effect magic.SpellEffect
	switch spell.Name {
	case "ARMOUR":
		effect = magic.ApplyARMOUR()
		m.character.AddSpellEffect("ARMOUR", 10)
	case "CRYPT":
		effect = magic.ApplyCRYPT()
		// Restore POW to maximum
		m.character.SetPOW(m.character.MaximumPOW)
	case "FIREBALL":
		effect = magic.ApplyFIREBALL()
	case "INVISIBILITY":
		effect = magic.ApplyINVISIBILITY(m.inCombat)
	case "PARALYSIS":
		effect = magic.ApplyPARALYSIS()
	case "POISON NEEDLE":
		effect = magic.ApplyPOISONNEEDLE(m.roller)
	case "RESURRECTION":
		effect = magic.ApplyRESURRECTION()
	case "RETRACE":
		// For now, just show message (actual navigation handled by UI layer)
		effect = magic.ApplyRETRACE("Previous Section")
	case "TIMEWARP":
		effect = magic.ApplyTIMEWARP()
		// Restore character LP to max (simplified - actual implementation would track section entry LP)
		m.character.SetLP(m.character.MaximumLP)
	case "XENOPHOBIA":
		effect = magic.ApplyXENOPHOBIA()
		// Effect is handled in combat damage calculation
	default:
		effect = magic.SpellEffect{Success: false, Message: "Unknown spell"}
	}

	m.message = fmt.Sprintf("%s\n%s\n\nSpent %d POW (Remaining: %d/%d)", 
		castResult.Message, effect.Message, spell.PowerCost, m.character.CurrentPOW, m.character.MaximumPOW)

	return effect, true
}

// GetMessage returns the current message.
func (m *SpellCastingModel) GetMessage() string {
	return m.message
}

// IsAwaitingConfirmation returns true if awaiting sacrifice confirmation.
func (m *SpellCastingModel) IsAwaitingConfirmation() bool {
	return m.awaitingConfirm
}

// GetNaturalCheckMessage returns the natural inclination check message.
func (m *SpellCastingModel) GetNaturalCheckMessage() string {
	return m.naturalCheckMsg
}

// Render returns the spell casting screen view.
func (m SpellCastingModel) Render() string {
	var b strings.Builder
	t := theme.Current()

	b.WriteString("\n")
	b.WriteString(theme.RenderTitle("SPELL CASTING"))
	b.WriteString("\n\n")

	// Show POW status
	b.WriteString("  " + theme.RenderPOWMeter(m.character.CurrentPOW, m.character.MaximumPOW, 30) + "\n")
	if m.inCombat {
		b.WriteString("  " + t.WarningMsg.Render("⚔ IN COMBAT") + "\n")
	}
	b.WriteString("\n")

	// Show natural check message if present
	if m.naturalCheckMsg != "" {
		b.WriteString("  " + t.Emphasis.Render(m.naturalCheckMsg) + "\n\n")
	}

	// Show awaiting confirmation dialog
	if m.awaitingConfirm {
		b.WriteString(theme.RenderSeparator(60) + "\n")
		b.WriteString(theme.RenderWarning("LP Sacrifice Required", m.message) + "\n\n")
		b.WriteString(theme.RenderKeyHelp("Y Sacrifice LP", "N Cancel") + "\n")
		b.WriteString(theme.RenderSeparator(60) + "\n")
		return b.String()
	}

	// Show spell list with scrolling viewport
	b.WriteString(t.Heading.Render("  Available Spells") + "\n\n")
	
	// Total items = spells + Natural Inclination Check option
	totalItems := len(m.spells) + 1
	
	// Adjust visible items based on whether message/effects are shown
	maxVisibleItems := 8
	if m.message != "" {
		// Reduce visible items when message is displayed (messages can be multi-line)
		maxVisibleItems = 5
	} else if len(m.character.ActiveSpellEffects) > 0 {
		// Reduce slightly if showing active effects
		maxVisibleItems = 7
	}
	
	startIdx := 0
	endIdx := totalItems
	
	if totalItems > maxVisibleItems {
		// Calculate viewport window centered on cursor
		viewportMid := maxVisibleItems / 2
		startIdx = m.cursor - viewportMid
		endIdx = m.cursor + viewportMid
		
		// Adjust if at start of list
		if startIdx < 0 {
			startIdx = 0
			endIdx = maxVisibleItems
		}
		
		// Adjust if at end of list
		if endIdx > totalItems {
			endIdx = totalItems
			startIdx = endIdx - maxVisibleItems
			if startIdx < 0 {
				startIdx = 0
			}
		}
	}
	
	// Show scroll indicator if there are items above
	if startIdx > 0 {
		b.WriteString(t.MutedText.Render("  ↑ More spells above...") + "\n\n")
	}
	
	// Render visible spells only
	for i := startIdx; i < endIdx; i++ {
		if i < len(m.spells) {
			// Render spell
			spell := m.spells[i]
			selected := i == m.cursor

			// Check if player can afford
			var affordIcon string
			if spell.PowerCost > m.character.CurrentPOW {
				affordIcon = t.Error.Render("✗")
			} else {
				affordIcon = t.SuccessMsg.Render("✓")
			}

			context := ""
			if spell.CombatOnly {
				context = t.WarningMsg.Render(" (Combat only)")
			} else if spell.Name == "CRYPT" || spell.Name == "RETRACE" {
				context = t.MutedText.Render(" (Non-combat)")
			}

			spellLine := fmt.Sprintf("%-15s [%d POW] %s%s", spell.Name, spell.PowerCost, affordIcon, context)
			
			if selected {
				b.WriteString("  " + theme.RenderMenuItem(spellLine, true) + "\n")
				b.WriteString("    " + t.MutedText.Render(spell.Description) + "\n\n")
			} else {
				b.WriteString("  " + t.MenuItem.Render(spellLine) + "\n")
				b.WriteString("    " + t.MutedText.Render(spell.Description) + "\n\n")
			}
		} else {
			// Render Natural Inclination Check option
			selected := m.cursor == len(m.spells)
			checkText := "Natural Inclination Check (Roll 2d6, need 4+)"
			
			if selected {
				b.WriteString("  " + theme.RenderMenuItem(checkText, true) + "\n\n")
			} else {
				b.WriteString("  " + t.MenuItem.Render(checkText) + "\n\n")
			}
		}
	}
	
	// Show scroll indicator if there are items below
	if endIdx < totalItems {
		b.WriteString(t.MutedText.Render("  ↓ More spells below...") + "\n\n")
	}

	// Show message if present
	if m.message != "" {
		b.WriteString(theme.RenderSeparator(60) + "\n")
		b.WriteString(t.Emphasis.Render("  "+m.message) + "\n")
		b.WriteString(theme.RenderSeparator(60) + "\n\n")
	}

	// Show active spell effects
	if len(m.character.ActiveSpellEffects) > 0 {
		b.WriteString(t.Heading.Render("  Active Effects") + "\n")
		for effect, value := range m.character.ActiveSpellEffects {
			b.WriteString(fmt.Sprintf("  %s %s (%d)\n", t.SuccessMsg.Render("•"), t.Value.Render(effect), value))
		}
		b.WriteString("\n")
	}

	b.WriteString(theme.RenderKeyHelp("↑/↓ Navigate", "Enter Cast/Check", "Esc Back", "? Help") + "\n")

	return b.String()
}
