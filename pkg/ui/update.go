package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/benoit/saga-demonspawn/internal/character"
	"github.com/benoit/saga-demonspawn/internal/combat"
	"github.com/benoit/saga-demonspawn/internal/help"
	"github.com/benoit/saga-demonspawn/internal/magic"
	"github.com/benoit/saga-demonspawn/pkg/ui/theme"
)

// Init initializes the Bubble Tea application.
// This is called once when the program starts.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles incoming messages and updates the model state.
// This is the core of the Elm Architecture pattern.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	
	// Handle custom messages from combat view
	case CastSpellMsg:
		if m.CurrentScreen == ScreenCombat {
			// Switch to spell casting screen in combat mode
			m.SpellCasting = NewSpellCastingModel(m.Character, m.Dice, true)
			m.CurrentScreen = ScreenMagic
			return m, nil
		}
	
	case CombatEndMsg:
		if msg.Victory {
			m.CurrentScreen = ScreenGameSession
			m.CombatState = nil
		} else {
			m.CurrentScreen = ScreenGameSession
			m.CombatState = nil
		}
		return m, nil
	
	// Pass other messages to combat view when in combat
	default:
		if m.CurrentScreen == ScreenCombat {
			var cmd tea.Cmd
			m.CombatView, cmd = m.CombatView.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

// handleKeyPress routes key presses to the appropriate screen handler.
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global quit keys
	if msg.String() == "ctrl+c" {
		// Auto-save if enabled
		if m.Config != nil && m.Config.AutoSave {
			_ = m.SaveCharacter()
		}
		return m, tea.Quit
	}

	// Handle help modal if showing
	if m.ShowingHelp {
		return m.handleHelpModalKeys(msg)
	}

	// Global help key
	if msg.String() == "?" {
		// Determine which help screen to show based on current screen
		var helpScreen help.Screen
		switch m.CurrentScreen {
		case ScreenMainMenu:
			helpScreen = help.ScreenMainMenu
		case ScreenCharacterCreation:
			helpScreen = help.ScreenCharacterCreation
		case ScreenCharacterEdit:
			helpScreen = help.ScreenCharacterEdit
		case ScreenCombat, ScreenCombatSetup:
			helpScreen = help.ScreenCombat
		case ScreenMagic:
			helpScreen = help.ScreenMagic
		default:
			helpScreen = help.ScreenGlobal
		}
		m.ShowHelp(helpScreen)
		return m, nil
	}

	// Route to screen-specific handlers
	switch m.CurrentScreen {
	case ScreenMainMenu:
		return m.handleMainMenuKeys(msg)
	case ScreenCharacterCreation:
		return m.handleCharacterCreationKeys(msg)
	case ScreenLoadCharacter:
		return m.handleLoadCharacterKeys(msg)
	case ScreenGameSession:
		return m.handleGameSessionKeys(msg)
	case ScreenCharacterView:
		return m.handleCharacterViewKeys(msg)
	case ScreenCharacterEdit:
		return m.handleCharacterEditKeys(msg)
	case ScreenCombatSetup:
		return m.handleCombatSetupKeys(msg)
	case ScreenCombat:
		return m.handleCombatKeys(msg)
	case ScreenInventory:
		return m.handleInventoryKeys(msg)
	case ScreenMagic:
		return m.handleMagicKeys(msg)
	case ScreenSettings:
		return m.handleSettingsKeys(msg)
	case ScreenDiceRoll:
		return m.handleDiceRollKeys(msg)
	default:
		return m, nil
	}
}

// handleMainMenuKeys processes key presses on the main menu.
func (m Model) handleMainMenuKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.MainMenu.MoveUp()
	case "down", "j":
		m.MainMenu.MoveDown()
	case "enter":
		selected := m.MainMenu.GetSelected()
		switch selected {
		case "New Character":
			m.CharCreation.Reset()
			m.CurrentScreen = ScreenCharacterCreation
		case "Load Character":
			// Use configured save directory
			saveDir := "."
			if m.Config != nil && m.Config.SaveDirectory != "" {
				saveDir = m.Config.SaveDirectory
			}
			m.LoadChar.RefreshFromDirectory(saveDir)
			m.CurrentScreen = ScreenLoadCharacter
		case "Settings":
			m.Settings.Reset(m.Config)
			m.CurrentScreen = ScreenSettings
		case "Help":
			m.ShowHelp(help.ScreenGlobal)
		case "Exit":
			// Auto-save if enabled
			if m.Config != nil && m.Config.AutoSave {
				_ = m.SaveCharacter()
			}
			return m, tea.Quit
		}
	case "q", "esc":
		// Auto-save if enabled
		if m.Config != nil && m.Config.AutoSave {
			_ = m.SaveCharacter()
		}
		return m, tea.Quit
	case "h":
		// Capital H for comprehensive help
		m.ShowHelp(help.ScreenGlobal)
	}
	return m, nil
}

// handleLoadCharacterKeys processes key presses on the load character screen.
func (m Model) handleLoadCharacterKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.LoadChar.MoveUp()
	case "down", "j":
		m.LoadChar.MoveDown()
	case "enter":
		// Load the selected character
		if m.LoadChar.HasFiles() {
			filename := m.LoadChar.GetSelectedFile()
			char, err := character.Load(filename)
			if err != nil {
				m.Err = err
				return m, nil
			}
			m.LoadCharacter(char)
			m.GameSession.UpdateMagicVisibility(char.MagicUnlocked)
		}
	case "esc", "q":
		// Return to main menu
		m.CurrentScreen = ScreenMainMenu
	}
	return m, nil
}

// handleCharacterCreationKeys processes key presses during character creation.
func (m Model) handleCharacterCreationKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.CharCreation.GetStep() {
	case StepRollCharacteristics:
		return m.handleRollCharacteristicsKeys(msg)
	case StepSelectEquipment:
		return m.handleSelectEquipmentKeys(msg)
	case StepReviewCharacter:
		return m.handleReviewCharacterKeys(msg)
	}
	return m, nil
}

// handleRollCharacteristicsKeys handles keys during stat rolling.
func (m Model) handleRollCharacteristicsKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "r":
		// Roll all characteristics
		m.CharCreation.RollAll()
	case "enter":
		// Proceed to equipment selection if all rolled
		if m.CharCreation.AreAllRolled() {
			m.CharCreation.NextStep()
		}
	case "esc", "q":
		// Return to main menu
		m.CurrentScreen = ScreenMainMenu
	}
	return m, nil
}

// handleSelectEquipmentKeys handles keys during equipment selection.
func (m Model) handleSelectEquipmentKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.CharCreation.MoveWeaponCursorUp()
	case "down", "j":
		m.CharCreation.MoveWeaponCursorDown()
	case "left", "h":
		m.CharCreation.MoveArmorCursorUp()
	case "right", "l":
		m.CharCreation.MoveArmorCursorDown()
	case "enter":
		// Proceed to review
		m.CharCreation.NextStep()
	case "esc":
		// Go back to rolling
		m.CharCreation.PreviousStep()
	}
	return m, nil
}

// handleReviewCharacterKeys handles keys on the review screen.
func (m Model) handleReviewCharacterKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Finalize character creation
		char, err := m.CharCreation.CreateCharacter()
		if err != nil {
			m.Err = err
			return m, nil
		}
		m.LoadCharacter(char)
		m.GameSession.UpdateMagicVisibility(char.MagicUnlocked)
	case "esc":
		// Go back to equipment selection
		m.CharCreation.PreviousStep()
	case "q":
		// Cancel and return to main menu
		m.CurrentScreen = ScreenMainMenu
	}
	return m, nil
}

// handleGameSessionKeys processes key presses on the game session menu.
func (m Model) handleGameSessionKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.GameSession.MoveUp()
	case "down", "j":
		m.GameSession.MoveDown()
	case "enter":
		selected := m.GameSession.GetSelected()
		switch selected {
		case "View Character":
			m.CurrentScreen = ScreenCharacterView
		case "Edit Character Stats":
			m.CurrentScreen = ScreenCharacterEdit
		case "Combat":
			// Start combat setup
			m.CombatSetup.Reset()
			m.CurrentScreen = ScreenCombatSetup
		case "Cast Spell":
			// Initialize spell casting screen
			m.SpellCasting = NewSpellCastingModel(m.Character, m.Dice, false)
			m.CurrentScreen = ScreenMagic
		case "Manage Inventory":
			// Initialize inventory with current character
			m.Inventory = NewInventoryManagementModel(m.Character, false)
			m.CurrentScreen = ScreenInventory
		case "Roll Dice":
			m.DiceRoll.Reset()
			m.CurrentScreen = ScreenDiceRoll
		case "Save & Exit":
			if err := m.SaveCharacter(); err != nil {
				m.Err = err
				return m, nil
			}
			m.CurrentScreen = ScreenMainMenu
		}
	case "q", "esc":
		// Save and return to main menu
		if err := m.SaveCharacter(); err != nil {
			m.Err = err
			return m, nil
		}
		m.CurrentScreen = ScreenMainMenu
	}
	return m, nil
}

// handleCharacterViewKeys processes key presses on the character view screen.
func (m Model) handleCharacterViewKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "e":
		// Enter edit mode
		m.CurrentScreen = ScreenCharacterEdit
	case "b", "esc", "q":
		// Back to game session
		m.CurrentScreen = ScreenGameSession
	}
	return m, nil
}

// handleCharacterEditKeys processes key presses on the character edit screen.
func (m Model) handleCharacterEditKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.CharEdit.IsUnlockMode() {
		return m.handleUnlockMagicKeys(msg)
	}

	if m.CharEdit.IsInputMode() {
		return m.handleCharacterEditInputKeys(msg)
	}

	switch msg.String() {
	case "up", "k":
		m.CharEdit.MoveUp()
	case "down", "j":
		m.CharEdit.MoveDown()
	case "enter":
		// Start editing the selected field
		m.CharEdit.StartInput()
	case "u", "U":
		// Unlock magic (only if not already unlocked)
		if m.Character != nil && !m.Character.MagicUnlocked {
			m.CharEdit.StartUnlockMode()
		}
	case "esc", "q":
		// Back to character view
		m.CurrentScreen = ScreenCharacterView
	}
	return m, nil
}

// handleUnlockMagicKeys handles keys during magic unlock mode.
func (m Model) handleUnlockMagicKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Confirm unlock
		if m.CharEdit.ConfirmUnlock() {
			// Update game session menu to show Cast Spell
			m.GameSession.UpdateMagicVisibility(m.Character.MagicUnlocked)
		}
	case "esc":
		// Cancel unlock
		m.CharEdit.CancelUnlockMode()
	case "backspace":
		m.CharEdit.Backspace()
	default:
		// Append numeric input
		if len(msg.String()) == 1 && msg.String()[0] >= '0' && msg.String()[0] <= '9' {
			m.CharEdit.AppendInput(msg.String())
		}
	}
	return m, nil
}

// handleCharacterEditInputKeys handles keys when actively editing a value.
func (m Model) handleCharacterEditInputKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Apply the edit
		m.applyCharacterEdit()
		m.CharEdit.CancelInput()
	case "esc":
		// Cancel the edit
		m.CharEdit.CancelInput()
	case "backspace":
		m.CharEdit.Backspace()
	default:
		// Append numeric input
		if len(msg.String()) == 1 && msg.String()[0] >= '0' && msg.String()[0] <= '9' {
			m.CharEdit.AppendInput(msg.String())
		}
		// Allow minus sign for negative values
		if msg.String() == "-" && len(m.CharEdit.GetInputBuffer()) == 0 {
			m.CharEdit.AppendInput("-")
		}
	}
	return m, nil
}

// applyCharacterEdit applies the input buffer to the selected character field.
func (m *Model) applyCharacterEdit() {
	if m.Character == nil {
		return
	}

	buffer := m.CharEdit.GetInputBuffer()
	if buffer == "" {
		return
	}

	// Parse the input
	var value int
	_, err := fmt.Sscanf(buffer, "%d", &value)
	if err != nil {
		m.Err = err
		return
	}

	// Apply to the selected field
	cursor := m.CharEdit.GetCursor()
	switch EditField(cursor) {
	case EditFieldStrength:
		// For characteristics, we set directly (book might say "your STR is now 75")
		m.Character.Strength = value
	case EditFieldSpeed:
		m.Character.Speed = value
	case EditFieldStamina:
		m.Character.Stamina = value
	case EditFieldCourage:
		m.Character.Courage = value
	case EditFieldLuck:
		m.Character.Luck = value
	case EditFieldCharm:
		m.Character.Charm = value
	case EditFieldAttraction:
		m.Character.Attraction = value
	case EditFieldCurrentLP:
		m.Character.SetLP(value)
	case EditFieldMaxLP:
		m.Character.SetMaxLP(value)
	case EditFieldSkill:
		m.Character.SetSkill(value)
	case EditFieldCurrentPOW:
		m.Character.SetPOW(value)
	case EditFieldMaxPOW:
		m.Character.SetMaxPOW(value)
	}
}

// handleCombatSetupKeys processes key presses on the combat setup screen.
func (m Model) handleCombatSetupKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.CombatSetup, cmd = m.CombatSetup.Update(msg)
	
	// Handle combat start message
	if cmd != nil {
		returnedMsg := cmd()
		if _, ok := returnedMsg.(CombatStartMsg); ok {
			// Start combat - create enemy and initialize combat state
			name, str, spd, sta, crg, lck, skill, currentLP, maxLP, weaponBonus, armorProtection, isDemonspawn := m.CombatSetup.GetEnemyData()
			enemy, err := combat.NewEnemy(name, str, spd, sta, crg, lck, skill, currentLP, maxLP, weaponBonus, armorProtection, isDemonspawn)
			if err != nil {
				m.Err = err
				return m, nil
			}
			
			// Initialize combat
			m.CombatState = combat.StartCombat(m.Character, enemy, m.Dice)
			m.CombatState.AddLogEntry(fmt.Sprintf("Combat begins against %s!", enemy.Name))
			m.CombatState.AddLogEntry(fmt.Sprintf("[Initiative] Player: %d, Enemy: %d", m.CombatState.PlayerInitiative, m.CombatState.EnemyInitiative))
			
			if m.CombatState.PlayerFirstStrike {
				m.CombatState.AddLogEntry("You strike first!")
			} else {
				m.CombatState.AddLogEntry("Enemy strikes first!")
			}
			
			m.CombatView = NewCombatViewModel(m.Character, m.CombatState, m.Dice)
			m.CurrentScreen = ScreenCombat
			return m, nil
		}
	}
	
	switch msg.String() {
	case "esc":
		if !m.CombatSetup.inputMode {
			m.CurrentScreen = ScreenGameSession
		}
	}
	
	return m, nil
}

// handleCombatKeys processes key presses during combat.
func (m Model) handleCombatKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Only allow escape back to menu during player turn when waiting for input
		if m.CombatView.waitingForInput && m.CombatState != nil && m.CombatState.PlayerTurn {
			m.CurrentScreen = ScreenGameSession
			m.CombatState = nil
			return m, nil
		}
	}
	
	// Update combat view with key message - this may produce CastSpellMsg or CombatEndMsg
	var cmd tea.Cmd
	m.CombatView, cmd = m.CombatView.Update(msg)
	return m, cmd
}

// handleInventoryKeys processes key presses on the inventory management screen.
func (m Model) handleInventoryKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.Inventory.MoveUp()
	case "down", "j":
		m.Inventory.MoveDown()
	case "enter":
		m.Inventory.HandleEnter()
		// Rebuild to reflect changes
		m.CharView.SetCharacter(m.Character)
	case "u":
		m.Inventory.HandleUse()
	case "a":
		// Acquire special items (for testing/when finding items)
		m.Inventory.HandleAcquire()
		// Rebuild to reflect changes
		m.CharView.SetCharacter(m.Character)
	case "r":
		// Request recharge confirmation
		if m.Inventory.HandleRecharge() {
			// Show confirmation or confirm directly
			m.Inventory.ConfirmRecharge()
		}
	case "i":
		// Show item info - for now just show in message
		item := m.Inventory.GetCurrentItem()
		if item != nil {
			if item.SpecialItem != "" {
				switch item.SpecialItem {
				case "healing_stone":
					m.Inventory.message = "Healing Stone: Use during combat to restore 1d6×10 LP. Recharge with 'R' when gamebook allows."
				case "doombringer":
					m.Inventory.message = "Doombringer: +20 damage, -10 LP per attack, heal LP equal to damage dealt on hit."
				case "orb":
					m.Inventory.message = "The Orb: Hold to double damage vs Demonspawn, or throw for instant kill (4+ to hit)."
				}
			} else {
				// Show general help for regular items
				m.Inventory.message = "Tip: Special items acquired during adventure can be activated with 'A' key."
			}
		}
	case "esc", "q":
		// Back to game session
		m.CurrentScreen = ScreenGameSession
	}
	return m, nil
}

// handleMagicKeys processes key presses on the spell casting screen.
func (m Model) handleMagicKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle confirmation dialog
	if m.SpellCasting.IsAwaitingConfirmation() {
		switch msg.String() {
		case "y", "Y":
			// Confirm sacrifice
			if m.SpellCasting.ConfirmSacrifice() {
				// Now proceed with cast
				effect, success := m.SpellCasting.PerformCast()
				if success {
					m.handleSpellEffect(effect)
				}
				m.SpellCasting.SetCharacter(m.Character)
			}
		case "n", "N", "esc":
			// Cancel sacrifice
			m.SpellCasting.CancelSacrifice()
		}
		return m, nil
	}

	// Normal spell selection
	switch msg.String() {
	case "up", "k":
		m.SpellCasting.MoveUp()
	case "down", "j":
		m.SpellCasting.MoveDown()
	case "enter":
		// Check if natural inclination check is selected
		if m.SpellCasting.IsNaturalCheckSelected() {
			m.SpellCasting.PerformNaturalCheck()
		} else {
			// Attempt to cast spell
			if m.SpellCasting.AttemptCast() {
				// Cast validated, perform the cast
				effect, success := m.SpellCasting.PerformCast()
				if success {
					m.handleSpellEffect(effect)
				}
				m.SpellCasting.SetCharacter(m.Character)
			}
		}
	case "esc", "q":
		// Back to game session or combat
		if m.SpellCasting.returnToCombat {
			m.CurrentScreen = ScreenCombat
		} else {
			m.CurrentScreen = ScreenGameSession
		}
	}
	return m, nil
}

// handleSpellEffect applies the spell effect to the game state.
func (m *Model) handleSpellEffect(effect magic.SpellEffect) {
	// Handle combat effects
	if effect.CombatEnded && m.CombatState != nil {
		if effect.Victory {
			m.CombatState.AddLogEntry("Combat ended via magic (victory)!")
		} else {
			m.CombatState.AddLogEntry("Combat ended via magic (escape)!")
		}
		m.CurrentScreen = ScreenGameSession
		m.CombatState = nil
	}

	// Handle enemy damage
	if effect.DamageDealt > 0 && m.CombatState != nil {
		m.CombatState.Enemy.CurrentLP -= effect.DamageDealt
		m.CombatState.AddLogEntry(fmt.Sprintf("Spell deals %d damage to %s!", effect.DamageDealt, m.CombatState.Enemy.Name))
		if m.CombatState.Enemy.CurrentLP <= 0 {
			m.CombatState.AddLogEntry(fmt.Sprintf("%s is defeated!", m.CombatState.Enemy.Name))
		}
	}

	// Handle enemy killed
	if effect.EnemyKilled && m.CombatState != nil {
		m.CombatState.Enemy.CurrentLP = 0
		m.CombatState.AddLogEntry(fmt.Sprintf("%s is killed by magic!", m.CombatState.Enemy.Name))
	}

	// Handle navigation
	if effect.NavigateTo != "" {
		// For now, just show message (actual navigation would require section system)
		// CRYPT: restore POW to max
		if effect.NavigateTo == "CRYPT" {
			m.Character.SetPOW(m.Character.MaximumPOW)
		}
	}

	// Handle RESURRECTION (requires stat reroll)
	if effect.RequiresReroll {
		// For now, just restore LP (full implementation would reroll all stats)
		m.Character.SetLP(m.Character.MaximumLP)
	}
}

// handleHelpModalKeys processes key presses when help modal is shown.
func (m Model) handleHelpModalKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.ScrollHelpUp()
	case "down", "j":
		m.ScrollHelpDown()
	case "esc", "?":
		m.HideHelp()
	}
	return m, nil
}

// handleSettingsKeys processes key presses on the settings screen.
func (m Model) handleSettingsKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.Settings.MoveCursorUp()
	case "down", "j":
		m.Settings.MoveCursorDown()
	case "enter", " ":
		cursor := m.Settings.GetCursor()
		field := SettingField(cursor)
		
		switch field {
		case SettingTheme:
			m.Settings.CycleTheme()
			// Apply theme immediately
			cfg := m.Settings.GetConfig()
			scheme := theme.ColorSchemeDark
			if cfg.Theme == "light" {
				scheme = theme.ColorSchemeLight
			}
			theme.Init(scheme, cfg.UseUnicode)
		case SettingSave:
			if err := m.Settings.Save(); err == nil {
				// Update main config
				*m.Config = *m.Settings.GetConfig()
				// Reinitialize theme
				scheme := theme.ColorSchemeDark
				if m.Config.Theme == "light" {
					scheme = theme.ColorSchemeLight
				}
				theme.Init(scheme, m.Config.UseUnicode)
			}
		case SettingCancel:
			m.Settings.Cancel()
			m.CurrentScreen = ScreenMainMenu
		case SettingReset:
			m.Settings.ResetToDefaults()
		default:
			// Toggle boolean settings
			m.Settings.ToggleCurrentSetting()
		}
	case "esc", "q":
		m.CurrentScreen = ScreenMainMenu
	}
	return m, nil
}

// handleDiceRollKeys processes key presses on the dice roll screen.
func (m Model) handleDiceRollKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "1":
		m.DiceRoll.Roll1D6()
	case "2":
		m.DiceRoll.Roll2D6()
	case "esc", "q":
		m.CurrentScreen = ScreenGameSession
	}
	return m, nil
}

