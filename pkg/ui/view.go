package ui

import (
	"fmt"
	"strings"
)

// View renders the current model state to a string for terminal display.
// This is the "V" in Bubble Tea's Model-View-Update pattern.
func (m Model) View() string {
	// Route to screen-specific views
	switch m.CurrentScreen {
	case ScreenMainMenu:
		return m.viewMainMenu()
	case ScreenCharacterCreation:
		return m.viewCharacterCreation()
	case ScreenLoadCharacter:
		return m.viewLoadCharacter()
	case ScreenGameSession:
		return m.viewGameSession()
	case ScreenCharacterView:
		return m.viewCharacterView()
	case ScreenCharacterEdit:
		return m.viewCharacterEdit()
	case ScreenCombatSetup:
		return m.CombatSetup.View()
	case ScreenCombat:
		return m.CombatView.View()
	case ScreenInventory:
		return renderInventoryView(m)
	default:
		return "Unknown screen"
	}
}

// viewMainMenu renders the main menu.
func (m Model) viewMainMenu() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString("  ╔════════════════════════════════════════╗\n")
	b.WriteString("  ║  SAGAS OF THE DEMONSPAWN - COMPANION  ║\n")
	b.WriteString("  ╚════════════════════════════════════════╝\n")
	b.WriteString("\n")

	choices := m.MainMenu.GetChoices()
	cursor := m.MainMenu.GetCursor()

	for i, choice := range choices {
		prefix := "  "
		if i == cursor {
			prefix = "> "
		}
		b.WriteString(fmt.Sprintf("%s%s\n", prefix, choice))
	}

	b.WriteString("\n")
	b.WriteString("  Navigation: ↑/↓ to move, Enter to select, q to quit\n")

	return b.String()
}

// viewLoadCharacter renders the load character screen.
func (m Model) viewLoadCharacter() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString("  ╔════════════════════════════════════════╗\n")
	b.WriteString("  ║          LOAD CHARACTER               ║\n")
	b.WriteString("  ╚════════════════════════════════════════╝\n")
	b.WriteString("\n")

	if m.LoadChar.GetError() != nil {
		b.WriteString(fmt.Sprintf("  Error: %v\n\n", m.LoadChar.GetError()))
	}

	if !m.LoadChar.HasFiles() {
		b.WriteString("  No saved characters found.\n\n")
		b.WriteString("  Press Esc to return to main menu\n")
		return b.String()
	}

	files := m.LoadChar.GetFiles()
	cursor := m.LoadChar.GetCursor()

	b.WriteString("  Select a character to load:\n\n")

	for i, file := range files {
		prefix := "  "
		if i == cursor {
			prefix = "> "
		}
		fileInfo := GetFileInfo(file)
		b.WriteString(fmt.Sprintf("%s%s (%s)\n", prefix, file, fileInfo))
	}

	b.WriteString("\n")
	b.WriteString("  Navigation: ↑/↓ to select, Enter to load, Esc to cancel\n")

	return b.String()
}

// viewGameSession renders the game session menu.
func (m Model) viewGameSession() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString("  ╔════════════════════════════════════════╗\n")
	b.WriteString("  ║          GAME SESSION MENU            ║\n")
	b.WriteString("  ╚════════════════════════════════════════╝\n")
	b.WriteString("\n")

	if m.Character != nil {
		b.WriteString(fmt.Sprintf("  Fire*Wolf  |  LP: %d/%d  |  SKL: %d",
			m.Character.CurrentLP, m.Character.MaximumLP, m.Character.Skill))
		if m.Character.MagicUnlocked {
			b.WriteString(fmt.Sprintf("  |  POW: %d/%d",
				m.Character.CurrentPOW, m.Character.MaximumPOW))
		}
		b.WriteString("\n\n")
	}

	choices := m.GameSession.GetChoices()
	cursor := m.GameSession.GetCursor()

	for i, choice := range choices {
		prefix := "  "
		if i == cursor {
			prefix = "> "
		}
		b.WriteString(fmt.Sprintf("%s%s\n", prefix, choice))
	}

	b.WriteString("\n")
	b.WriteString("  Navigation: ↑/↓ to move, Enter to select, q to save and quit\n")

	return b.String()
}

// viewCharacterCreation renders the character creation screens.
func (m Model) viewCharacterCreation() string {
	switch m.CharCreation.GetStep() {
	case StepRollCharacteristics:
		return m.viewRollCharacteristics()
	case StepSelectEquipment:
		return m.viewSelectEquipment()
	case StepReviewCharacter:
		return m.viewReviewCharacter()
	default:
		return "Unknown creation step"
	}
}

// viewRollCharacteristics renders the stat rolling screen.
func (m Model) viewRollCharacteristics() string {
	var b strings.Builder

	str, spd, sta, crg, lck, chm, att := m.CharCreation.GetCharacteristics()

	b.WriteString("\n")
	b.WriteString("  ╔════════════════════════════════════════╗\n")
	b.WriteString("  ║      CHARACTER CREATION - STEP 1      ║\n")
	b.WriteString("  ║         Roll Characteristics          ║\n")
	b.WriteString("  ╚════════════════════════════════════════╝\n")
	b.WriteString("\n")

	b.WriteString("  Roll 2d6 × 8 for each characteristic:\n\n")
	b.WriteString(fmt.Sprintf("    Strength (STR)   : %s\n", formatRoll(str)))
	b.WriteString(fmt.Sprintf("    Speed (SPD)      : %s\n", formatRoll(spd)))
	b.WriteString(fmt.Sprintf("    Stamina (STA)    : %s\n", formatRoll(sta)))
	b.WriteString(fmt.Sprintf("    Courage (CRG)    : %s\n", formatRoll(crg)))
	b.WriteString(fmt.Sprintf("    Luck (LCK)       : %s\n", formatRoll(lck)))
	b.WriteString(fmt.Sprintf("    Charm (CHM)      : %s\n", formatRoll(chm)))
	b.WriteString(fmt.Sprintf("    Attraction (ATT) : %s\n", formatRoll(att)))
	b.WriteString("\n")

	if m.CharCreation.AreAllRolled() {
		b.WriteString(fmt.Sprintf("  Life Points (LP): %d\n", m.CharCreation.GetCalculatedLP()))
		b.WriteString("  Skill (SKL): 0\n\n")
		b.WriteString("  Press Enter to continue to equipment selection\n")
	} else {
		b.WriteString("  Press 'r' to roll all characteristics\n")
	}

	b.WriteString("  Press Esc to return to main menu\n")

	return b.String()
}

// formatRoll formats a roll value, showing "--" if not yet rolled.
func formatRoll(value int) string {
	if value == 0 {
		return "--"
	}
	return fmt.Sprintf("%d", value)
}

// viewSelectEquipment renders the equipment selection screen.
func (m Model) viewSelectEquipment() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString("  ╔════════════════════════════════════════╗\n")
	b.WriteString("  ║      CHARACTER CREATION - STEP 2      ║\n")
	b.WriteString("  ║        Select Starting Equipment       ║\n")
	b.WriteString("  ╚════════════════════════════════════════╝\n")
	b.WriteString("\n")

	b.WriteString("  Choose your starting weapon:\n")
	weapons := m.CharCreation.GetWeaponOptions()
	weaponCursor := m.CharCreation.GetWeaponCursor()
	for i, weapon := range weapons {
		prefix := "    "
		if i == weaponCursor {
			prefix = "  > "
		}
		b.WriteString(fmt.Sprintf("%s%s (+%d damage)\n", prefix, weapon.Name, weapon.DamageBonus))
	}

	b.WriteString("\n  Choose your starting armor:\n")
	armors := m.CharCreation.GetArmorOptions()
	armorCursor := m.CharCreation.GetArmorCursor()
	for i, armor := range armors {
		prefix := "    "
		if i == armorCursor {
			prefix = "  > "
		}
		protection := ""
		if armor.Protection > 0 {
			protection = fmt.Sprintf(" (-%d damage)", armor.Protection)
		}
		b.WriteString(fmt.Sprintf("%s%s%s\n", prefix, armor.Name, protection))
	}

	b.WriteString("\n")
	b.WriteString("  Navigation: ↑/↓ weapons, ←/→ armor, Enter to continue, Esc to go back\n")

	return b.String()
}

// viewReviewCharacter renders the final character review screen.
func (m Model) viewReviewCharacter() string {
	var b strings.Builder

	str, spd, sta, crg, lck, chm, att := m.CharCreation.GetCharacteristics()

	b.WriteString("\n")
	b.WriteString("  ╔════════════════════════════════════════╗\n")
	b.WriteString("  ║      CHARACTER CREATION - STEP 3      ║\n")
	b.WriteString("  ║          Review Fire*Wolf             ║\n")
	b.WriteString("  ╚════════════════════════════════════════╝\n")
	b.WriteString("\n")

	b.WriteString("  Characteristics:\n")
	b.WriteString(fmt.Sprintf("    STR: %d  SPD: %d  STA: %d  CRG: %d\n", str, spd, sta, crg))
	b.WriteString(fmt.Sprintf("    LCK: %d  CHM: %d  ATT: %d\n", lck, chm, att))
	b.WriteString("\n")

	b.WriteString(fmt.Sprintf("  Life Points: %d\n", m.CharCreation.GetCalculatedLP()))
	b.WriteString("  Skill: 0\n\n")

	weapon := m.CharCreation.GetSelectedWeapon()
	armor := m.CharCreation.GetSelectedArmor()

	b.WriteString("  Equipment:\n")
	if weapon != nil {
		b.WriteString(fmt.Sprintf("    Weapon: %s (+%d)\n", weapon.Name, weapon.DamageBonus))
	}
	if armor != nil {
		b.WriteString(fmt.Sprintf("    Armor: %s", armor.Name))
		if armor.Protection > 0 {
			b.WriteString(fmt.Sprintf(" (-%d)", armor.Protection))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString("  Press Enter to begin your adventure!\n")
	b.WriteString("  Press Esc to change equipment, q to cancel\n")

	return b.String()
}

// viewCharacterView renders the character sheet display.
func (m Model) viewCharacterView() string {
	var b strings.Builder

	if m.Character == nil {
		return "No character loaded"
	}

	char := m.Character

	b.WriteString("\n")
	b.WriteString("  ╔════════════════════════════════════════╗\n")
	b.WriteString("  ║          CHARACTER SHEET              ║\n")
	b.WriteString("  ╚════════════════════════════════════════╝\n")
	b.WriteString("\n")

	b.WriteString("  ┌─ Characteristics ─────────────────────┐\n")
	b.WriteString(fmt.Sprintf("  │ STR: %-3d  SPD: %-3d  STA: %-3d  CRG: %-3d │\n",
		char.Strength, char.Speed, char.Stamina, char.Courage))
	b.WriteString(fmt.Sprintf("  │ LCK: %-3d  CHM: %-3d  ATT: %-3d          │\n",
		char.Luck, char.Charm, char.Attraction))
	b.WriteString("  └───────────────────────────────────────┘\n\n")

	b.WriteString("  ┌─ Resources ───────────────────────────┐\n")
	b.WriteString(fmt.Sprintf("  │ Life Points: %d / %d\n", char.CurrentLP, char.MaximumLP))
	b.WriteString(fmt.Sprintf("  │ Skill: %d\n", char.Skill))
	if char.MagicUnlocked {
		b.WriteString(fmt.Sprintf("  │ Power: %d / %d\n", char.CurrentPOW, char.MaximumPOW))
	}
	b.WriteString("  └───────────────────────────────────────┘\n\n")

	b.WriteString("  ┌─ Equipment ───────────────────────────┐\n")
	if char.EquippedWeapon != nil {
		b.WriteString(fmt.Sprintf("  │ Weapon: %s (+%d)\n",
			char.EquippedWeapon.Name, char.EquippedWeapon.DamageBonus))
	}
	if char.EquippedArmor != nil {
		b.WriteString(fmt.Sprintf("  │ Armor: %s", char.EquippedArmor.Name))
		if char.EquippedArmor.Protection > 0 {
			b.WriteString(fmt.Sprintf(" (-%d)", char.EquippedArmor.Protection))
		}
		b.WriteString("\n")
	}
	if char.HasShield {
		b.WriteString("  │ Shield: Equipped\n")
	}
	totalProtection := char.GetArmorProtection()
	b.WriteString(fmt.Sprintf("  │ Total Protection: -%d damage\n", totalProtection))
	b.WriteString("  └───────────────────────────────────────┘\n\n")

	b.WriteString("  ┌─ Progress ────────────────────────────┐\n")
	b.WriteString(fmt.Sprintf("  │ Enemies Defeated: %d\n", char.EnemiesDefeated))
	b.WriteString("  └───────────────────────────────────────┘\n\n")
	
	// Special Items section (Phase 3)
	hasSpecialItems := char.HealingStoneCharges > 0 || char.DoombringerPossessed || char.OrbPossessed
	if hasSpecialItems {
		b.WriteString("  ┌─ Special Items ───────────────────────┐\n")
		
		// Healing Stone
		if char.HealingStoneCharges > 0 {
			b.WriteString(fmt.Sprintf("  │ Healing Stone    [AVAILABLE] %d/50 charges\n", char.HealingStoneCharges))
		}
		
		// Doombringer
		if char.DoombringerPossessed {
			if char.EquippedWeapon != nil && char.EquippedWeapon.Name == "Doombringer" {
				b.WriteString("  │ Doombringer      [EQUIPPED] +20 damage\n")
			} else {
				b.WriteString("  │ Doombringer      [POSSESSED] +20 damage\n")
			}
		}
		
		// The Orb
		if char.OrbPossessed {
			if char.OrbDestroyed {
				b.WriteString("  │ The Orb          [DESTROYED]\n")
			} else if char.OrbEquipped {
				b.WriteString("  │ The Orb          [EQUIPPED] Left hand\n")
			} else {
				b.WriteString("  │ The Orb          [POSSESSED] Not equipped\n")
			}
		}
		
		b.WriteString("  └───────────────────────────────────────┘\n\n")
	}

	b.WriteString("  Press 'e' to edit stats, 'b' to return to menu\n")

	return b.String()
}

// viewCharacterEdit renders the character editing screen.
func (m Model) viewCharacterEdit() string {
	var b strings.Builder

	if m.Character == nil {
		return "No character loaded"
	}

	b.WriteString("\n")
	b.WriteString("  ╔════════════════════════════════════════╗\n")
	b.WriteString("  ║         EDIT CHARACTER STATS          ║\n")
	b.WriteString("  ╚════════════════════════════════════════╝\n")
	b.WriteString("\n")

	fields := m.CharEdit.GetFields()
	cursor := m.CharEdit.GetCursor()

	for i, field := range fields {
		prefix := "  "
		if i == cursor {
			prefix = "> "
		}

		if i == cursor && m.CharEdit.IsInputMode() {
			// Show input buffer when editing
			b.WriteString(fmt.Sprintf("%s%-15s: [%s_]\n", prefix, field, m.CharEdit.GetInputBuffer()))
		} else {
			// Get the actual value for each field
			var value int
			switch EditField(i) {
			case EditFieldStrength:
				value = m.Character.Strength
			case EditFieldSpeed:
				value = m.Character.Speed
			case EditFieldStamina:
				value = m.Character.Stamina
			case EditFieldCourage:
				value = m.Character.Courage
			case EditFieldLuck:
				value = m.Character.Luck
			case EditFieldCharm:
				value = m.Character.Charm
			case EditFieldAttraction:
				value = m.Character.Attraction
			case EditFieldCurrentLP:
				value = m.Character.CurrentLP
			case EditFieldMaxLP:
				value = m.Character.MaximumLP
			case EditFieldSkill:
				value = m.Character.Skill
			case EditFieldCurrentPOW:
				value = m.Character.CurrentPOW
			case EditFieldMaxPOW:
				value = m.Character.MaximumPOW
			}
			b.WriteString(fmt.Sprintf("%s%-15s: %d\n", prefix, field, value))
		}
	}

	b.WriteString("\n")
	if m.CharEdit.IsInputMode() {
		b.WriteString("  Type new value, Enter to confirm, Esc to cancel\n")
	} else {
		b.WriteString("  Navigation: ↑/↓ to select field, Enter to edit, Esc to return\n")
	}

	return b.String()
}
