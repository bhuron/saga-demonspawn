package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/benoit/saga-demonspawn/internal/character"
	"github.com/benoit/saga-demonspawn/internal/combat"
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
	
	// Pass custom messages to combat view when in combat
	default:
		if m.CurrentScreen == ScreenCombat {
			var cmd tea.Cmd
			m.CombatView, cmd = m.CombatView.Update(msg)
			
			// Handle combat end message
			if cmd != nil {
				returnedMsg := cmd()
				if endMsg, ok := returnedMsg.(CombatEndMsg); ok {
					if endMsg.Victory {
						m.CurrentScreen = ScreenGameSession
						m.CombatState = nil
					} else {
						m.CurrentScreen = ScreenGameSession
						m.CombatState = nil
					}
					return m, nil
				}
			}
			
			return m, cmd
		}
	}

	return m, nil
}

// handleKeyPress routes key presses to the appropriate screen handler.
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global quit keys
	if msg.String() == "ctrl+c" {
		return m, tea.Quit
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
			m.LoadChar.Refresh()
			m.CurrentScreen = ScreenLoadCharacter
		case "Exit":
			return m, tea.Quit
		}
	case "q", "esc":
		return m, tea.Quit
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
			// TODO: Phase 4 - Implement magic
		case "Manage Inventory":
			// TODO: Phase 3 - Implement inventory
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
	case "esc", "q":
		// Back to character view
		m.CurrentScreen = ScreenCharacterView
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
	var cmd tea.Cmd
	m.CombatView, cmd = m.CombatView.Update(msg)
	
	// Handle combat end message
	if cmd != nil {
		returnedMsg := cmd()
		if endMsg, ok := returnedMsg.(CombatEndMsg); ok {
			if endMsg.Victory {
				// Victory - character is already updated by combat system
				m.CurrentScreen = ScreenGameSession
				m.CombatState = nil
			} else {
				// Defeat or fled - return to game session
				m.CurrentScreen = ScreenGameSession
				m.CombatState = nil
			}
			return m, nil
		}
	}
	
	switch msg.String() {
	case "esc":
		// Only allow escape back to menu during player turn when waiting for input
		if m.CombatView.waitingForInput && m.CombatState != nil && m.CombatState.PlayerTurn {
			m.CurrentScreen = ScreenGameSession
			m.CombatState = nil
		}
	}
	
	return m, cmd
}
