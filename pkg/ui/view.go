package ui

import (
	"fmt"
	"strings"

	"github.com/benoit/saga-demonspawn/internal/help"
	"github.com/benoit/saga-demonspawn/pkg/ui/theme"
)

// View renders the current model state to a string for terminal display.
// This is the "V" in Bubble Tea's Model-View-Update pattern.
func (m Model) View() string {
	// Initialize theme if needed
	if m.Config != nil {
		scheme := theme.ColorSchemeDark
		if m.Config.Theme == "light" {
			scheme = theme.ColorSchemeLight
		}
		theme.Init(scheme, m.Config.UseUnicode)
	}

	// Render main content
	var content string
	switch m.CurrentScreen {
	case ScreenMainMenu:
		content = m.viewMainMenu()
	case ScreenCharacterCreation:
		content = m.viewCharacterCreation()
	case ScreenLoadCharacter:
		content = m.viewLoadCharacter()
	case ScreenGameSession:
		content = m.viewGameSession()
	case ScreenCharacterView:
		content = m.viewCharacterView()
	case ScreenCharacterEdit:
		content = m.viewCharacterEdit()
	case ScreenCombatSetup:
		content = m.CombatSetup.View()
	case ScreenCombat:
		content = m.CombatView.View()
	case ScreenInventory:
		content = renderInventoryView(m)
	case ScreenMagic:
		content = m.SpellCasting.Render()
	case ScreenSettings:
		content = m.viewSettings()
	default:
		content = "Unknown screen"
	}

	// Overlay help modal if showing
	if m.ShowingHelp {
		content = m.renderHelpOverlay(content)
	}

	return content
}

// viewMainMenu renders the main menu.
func (m Model) viewMainMenu() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(theme.RenderTitle("SAGAS OF THE DEMONSPAWN - COMPANION"))
	b.WriteString("\n\n")

	choices := m.MainMenu.GetChoices()
	cursor := m.MainMenu.GetCursor()

	for i, choice := range choices {
		b.WriteString(theme.RenderMenuItem(choice, i == cursor) + "\n")
	}

	b.WriteString("\n")
	b.WriteString(theme.RenderKeyHelp(
		"↑/↓ Navigate",
		"Enter Select",
		"q Quit",
		"? Help",
	))

	return b.String()
}

// viewLoadCharacter renders the load character screen.
func (m Model) viewLoadCharacter() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(theme.RenderTitle("LOAD CHARACTER"))
	b.WriteString("\n\n")

	if m.LoadChar.GetError() != nil {
		b.WriteString(theme.RenderError(
			"Load Error",
			fmt.Sprintf("%v", m.LoadChar.GetError()),
			"Press Esc to return to main menu",
		) + "\n\n")
	}

	if !m.LoadChar.HasFiles() {
		b.WriteString(theme.RenderWarning(
			"No Saved Characters",
			"Create a new character to get started.",
		) + "\n\n")
		b.WriteString(theme.RenderKeyHelp("Esc Return to menu"))
		return b.String()
	}

	files := m.LoadChar.GetFiles()
	cursor := m.LoadChar.GetCursor()

	b.WriteString(theme.Current().Body.Render("  Select a character to load:") + "\n\n")

	for i, file := range files {
		selected := i == cursor
		fileInfo := GetFileInfo(file)
		text := fmt.Sprintf("%s (%s)", file, fileInfo)
		b.WriteString("  " + theme.RenderMenuItem(text, selected) + "\n")
	}

	b.WriteString("\n")
	b.WriteString(theme.RenderKeyHelp("↑/↓ Select", "Enter Load", "Esc Cancel", "? Help"))

	return b.String()
}

// viewGameSession renders the game session menu.
func (m Model) viewGameSession() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(theme.RenderTitle("GAME SESSION MENU"))
	b.WriteString("\n\n")

	if m.Character != nil {
		// Character status line with health bar
		b.WriteString("  " + theme.Current().Heading.Render("Fire*Wolf") + "\n")
		b.WriteString("  " + theme.RenderHealthBar(m.Character.CurrentLP, m.Character.MaximumLP, 30) + "\n")
		b.WriteString(fmt.Sprintf("  " + theme.RenderLabel("Skill", fmt.Sprintf("%d", m.Character.Skill))))
		if m.Character.MagicUnlocked {
			b.WriteString("  |  " + theme.RenderPOWMeter(m.Character.CurrentPOW, m.Character.MaximumPOW, 20))
		}
		b.WriteString("\n\n")
	}

	choices := m.GameSession.GetChoices()
	cursor := m.GameSession.GetCursor()

	for i, choice := range choices {
		b.WriteString("  " + theme.RenderMenuItem(choice, i == cursor) + "\n")
	}

	b.WriteString("\n")
	b.WriteString(theme.RenderKeyHelp("↑/↓ Navigate", "Enter Select", "q Save & quit", "? Help"))

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
	b.WriteString(theme.RenderTitle("CHARACTER CREATION - STEP 1"))
	b.WriteString("\n")
	b.WriteString(theme.Current().Heading.Render("  Roll Characteristics") + "\n\n")

	b.WriteString(theme.Current().Body.Render("  Roll 2d6 × 8 for each characteristic:") + "\n\n")
	
	// Render stats with color coding
	b.WriteString("    " + theme.RenderLabel("Strength (STR)  ", formatRollColored(str)) + "\n")
	b.WriteString("    " + theme.RenderLabel("Speed (SPD)     ", formatRollColored(spd)) + "\n")
	b.WriteString("    " + theme.RenderLabel("Stamina (STA)   ", formatRollColored(sta)) + "\n")
	b.WriteString("    " + theme.RenderLabel("Courage (CRG)   ", formatRollColored(crg)) + "\n")
	b.WriteString("    " + theme.RenderLabel("Luck (LCK)      ", formatRollColored(lck)) + "\n")
	b.WriteString("    " + theme.RenderLabel("Charm (CHM)     ", formatRollColored(chm)) + "\n")
	b.WriteString("    " + theme.RenderLabel("Attraction (ATT)", formatRollColored(att)) + "\n")
	b.WriteString("\n")

	if m.CharCreation.AreAllRolled() {
		lp := m.CharCreation.GetCalculatedLP()
		b.WriteString("  " + theme.RenderLabel("Life Points (LP)", fmt.Sprintf("%d", lp)) + "\n")
		b.WriteString("  " + theme.RenderLabel("Skill (SKL)     ", "0") + "\n\n")
		b.WriteString(theme.Current().Emphasis.Render("  Press Enter to continue to equipment selection") + "\n")
	} else {
		b.WriteString(theme.Current().Emphasis.Render("  Press 'r' to roll all characteristics") + "\n")
	}

	b.WriteString("\n")
	b.WriteString(theme.RenderKeyHelp("r Roll", "Enter Continue", "Esc Cancel", "? Help"))

	return b.String()
}

// formatRoll formats a roll value, showing "--" if not yet rolled.
func formatRoll(value int) string {
	if value == 0 {
		return "--"
	}
	return fmt.Sprintf("%d", value)
}

// formatRollColored formats a roll value with color coding based on quality.
func formatRollColored(value int) string {
	if value == 0 {
		return theme.Current().MutedText.Render("--")
	}
	// Color code based on characteristic range (16-96)
	// 72+ = green (top 25%), 48-71 = normal, 32-47 = yellow, <32 = red
	return theme.RenderStatValue(value, 16, 96)
}

// viewSelectEquipment renders the equipment selection screen.
func (m Model) viewSelectEquipment() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(theme.RenderTitle("CHARACTER CREATION - STEP 2"))
	b.WriteString("\n")
	b.WriteString(theme.Current().Heading.Render("  Select Starting Equipment") + "\n\n")

	b.WriteString(theme.Current().Body.Render("  Choose your starting weapon:") + "\n")
	weapons := m.CharCreation.GetWeaponOptions()
	weaponCursor := m.CharCreation.GetWeaponCursor()
	for i, weapon := range weapons {
		selected := i == weaponCursor
		text := fmt.Sprintf("%s (+%d damage)", weapon.Name, weapon.DamageBonus)
		b.WriteString("  " + theme.RenderMenuItem(text, selected) + "\n")
	}

	b.WriteString("\n" + theme.Current().Body.Render("  Choose your starting armor:") + "\n")
	armors := m.CharCreation.GetArmorOptions()
	armorCursor := m.CharCreation.GetArmorCursor()
	for i, armor := range armors {
		selected := i == armorCursor
		protection := ""
		if armor.Protection > 0 {
			protection = fmt.Sprintf(" (-%d damage)", armor.Protection)
		}
		text := armor.Name + protection
		b.WriteString("  " + theme.RenderMenuItem(text, selected) + "\n")
	}

	b.WriteString("\n")
	b.WriteString(theme.RenderKeyHelp("↑/↓ Weapons", "←/→ Armor", "Enter Continue", "Esc Back", "? Help"))

	return b.String()
}

// viewReviewCharacter renders the final character review screen.
func (m Model) viewReviewCharacter() string {
	var b strings.Builder

	str, spd, sta, crg, lck, chm, att := m.CharCreation.GetCharacteristics()

	b.WriteString("\n")
	b.WriteString(theme.RenderTitle("CHARACTER CREATION - STEP 3"))
	b.WriteString("\n")
	b.WriteString(theme.Current().Heading.Render("  Review Fire*Wolf") + "\n\n")

	b.WriteString(theme.Current().Body.Render("  Characteristics:") + "\n")
	b.WriteString(fmt.Sprintf("    STR: %s  SPD: %s  STA: %s  CRG: %s\n",
		formatRollColored(str), formatRollColored(spd), formatRollColored(sta), formatRollColored(crg)))
	b.WriteString(fmt.Sprintf("    LCK: %s  CHM: %s  ATT: %s\n",
		formatRollColored(lck), formatRollColored(chm), formatRollColored(att)))
	b.WriteString("\n")

	lp := m.CharCreation.GetCalculatedLP()
	b.WriteString("  " + theme.RenderLabel("Life Points", fmt.Sprintf("%d", lp)) + "\n")
	b.WriteString("  " + theme.RenderLabel("Skill", "0") + "\n\n")

	weapon := m.CharCreation.GetSelectedWeapon()
	armor := m.CharCreation.GetSelectedArmor()

	b.WriteString(theme.Current().Body.Render("  Equipment:") + "\n")
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
	b.WriteString(theme.Current().Emphasis.Render("  Press Enter to begin your adventure!") + "\n")
	b.WriteString("\n")
	b.WriteString(theme.RenderKeyHelp("Enter Begin", "Esc Change equipment", "q Cancel", "? Help"))

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

	// Show unlock magic dialog if in unlock mode
	if m.CharEdit.IsUnlockMode() {
		b.WriteString("  ───────────────────────────────────────────\n")
		b.WriteString("  UNLOCK MAGIC SYSTEM\n\n")
		b.WriteString("  Enter initial POWER value: [" + m.CharEdit.GetInputBuffer() + "_]\n\n")
		if m.CharEdit.GetUnlockMessage() != "" {
			b.WriteString("  " + m.CharEdit.GetUnlockMessage() + "\n\n")
		}
		b.WriteString("  Enter to confirm, Esc to cancel\n")
		b.WriteString("  ───────────────────────────────────────────\n")
		return b.String()
	}

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
	
	// Show unlock message if present
	if m.CharEdit.GetUnlockMessage() != "" {
		b.WriteString("  " + m.CharEdit.GetUnlockMessage() + "\n\n")
	}
	
	if m.CharEdit.IsInputMode() {
		b.WriteString("  Type new value, Enter to confirm, Esc to cancel\n")
	} else {
		b.WriteString("  Navigation: ↑/↓ to select field, Enter to edit, Esc to return\n")
		if m.Character != nil && !m.Character.MagicUnlocked {
			b.WriteString("  Press 'U' to unlock magic when gamebook allows\n")
		}
	}

	return b.String()
}

// renderHelpOverlay renders the help modal overlay on top of the current screen.
func (m Model) renderHelpOverlay(background string) string {
	var b strings.Builder

	// Show background (dimmed effect would be nice but complex in terminal)
	b.WriteString(background)
	b.WriteString("\n\n")

	// Render help box
	title := help.GetTitle(m.HelpScreen)
	b.WriteString(theme.RenderTitle(title))
	b.WriteString("\n\n")

	// Get help content lines
	lines := help.GetLines(m.HelpScreen)
	visibleLines := m.Height - 6
	if visibleLines < 10 {
		visibleLines = 10
	}

	// Calculate visible range
	start := m.HelpScroll
	end := start + visibleLines
	if end > len(lines) {
		end = len(lines)
	}

	// Render visible lines
	for i := start; i < end; i++ {
		b.WriteString("  " + lines[i] + "\n")
	}

	// Footer
	b.WriteString("\n")
	if m.HelpMaxScroll > 0 {
		b.WriteString(theme.RenderKeyHelp(
			"↑/↓ Scroll",
			fmt.Sprintf("Line %d/%d", m.HelpScroll+1, len(lines)),
			"Esc/? to close",
		))
	} else {
		b.WriteString(theme.RenderKeyHelp("Esc/? to close"))
	}

	return b.String()
}

// viewSettings renders the settings screen.
func (m Model) viewSettings() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(theme.RenderTitle("SETTINGS"))
	b.WriteString("\n\n")

	cfg := m.Settings.GetConfig()
	cursor := m.Settings.GetCursor()

	// Appearance section
	b.WriteString(theme.Current().Heading.Render("  Appearance") + "\n")
	renderSetting(&b, 0, cursor, "Color Scheme", cfg.Theme)
	renderSetting(&b, 1, cursor, "Use Unicode", boolToString(cfg.UseUnicode))
	renderSetting(&b, 2, cursor, "Show Animations", boolToString(cfg.ShowAnimations))
	b.WriteString("\n")

	// Gameplay section
	b.WriteString(theme.Current().Heading.Render("  Gameplay") + "\n")
	renderSetting(&b, 3, cursor, "Confirm Actions", boolToString(cfg.ConfirmActions))
	renderSetting(&b, 4, cursor, "Auto-save on Exit", boolToString(cfg.AutoSave))
	renderSetting(&b, 5, cursor, "Show Roll Details", boolToString(cfg.ShowRollDetails))
	b.WriteString("\n")

	// Accessibility section
	b.WriteString(theme.Current().Heading.Render("  Accessibility") + "\n")
	renderSetting(&b, 6, cursor, "High Contrast", boolToString(cfg.HighContrast))
	renderSetting(&b, 7, cursor, "Reduced Motion", boolToString(cfg.ReducedMotion))
	b.WriteString("\n")

	// Actions
	b.WriteString(theme.Current().Heading.Render("  Actions") + "\n")
	renderAction(&b, 8, cursor, "[Save]")
	renderAction(&b, 9, cursor, "[Cancel]")
	renderAction(&b, 10, cursor, "[Reset to Defaults]")
	b.WriteString("\n")

	// Status message
	if msg := m.Settings.GetMessage(); msg != "" {
		if m.Settings.IsSaved() {
			b.WriteString(theme.RenderSuccess(msg) + "\n\n")
		} else {
			b.WriteString(theme.Current().WarningMsg.Render(msg) + "\n\n")
		}
	}

	// Help text
	b.WriteString(theme.RenderKeyHelp(
		"↑/↓ Navigate",
		"Enter Toggle/Select",
		"Esc Return to menu",
		"? Help",
	))

	return b.String()
}

// renderSetting renders a settings field.
func renderSetting(b *strings.Builder, index, cursor int, label, value string) {
	prefix := "  "
	if index == cursor {
		prefix = "> "
		label = theme.Current().MenuItemSel.Render(label)
		value = theme.Current().Emphasis.Render(value)
	} else {
		label = theme.Current().MenuItem.Render(label)
		value = theme.Current().Value.Render(value)
	}
	b.WriteString(fmt.Sprintf("%s%-25s: %s\n", prefix, label, value))
}

// renderAction renders an action button.
func renderAction(b *strings.Builder, index, cursor int, label string) {
	prefix := "  "
	if index == cursor {
		prefix = "> "
		label = theme.Current().ButtonFocus.Render(label)
	} else {
		label = theme.Current().Button.Render(label)
	}
	b.WriteString(fmt.Sprintf("%s%s\n", prefix, label))
}

// boolToString converts a boolean to a user-friendly string.
func boolToString(value bool) string {
	if value {
		return "Enabled"
	}
	return "Disabled"
}
