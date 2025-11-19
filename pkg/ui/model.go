// Package ui provides Bubble Tea UI components for the Saga application.
package ui

import (
	"github.com/benoit/saga-demonspawn/internal/character"
	"github.com/benoit/saga-demonspawn/internal/combat"
	"github.com/benoit/saga-demonspawn/internal/config"
	"github.com/benoit/saga-demonspawn/internal/dice"
	"github.com/benoit/saga-demonspawn/internal/help"
)

// Screen represents the different screens in the application.
type Screen int

const (
	// ScreenMainMenu is the initial screen with New/Load/Exit options
	ScreenMainMenu Screen = iota
	// ScreenCharacterCreation is the character creation flow
	ScreenCharacterCreation
	// ScreenLoadCharacter is the character loading screen
	ScreenLoadCharacter
	// ScreenGameSession is the main game menu after loading a character
	ScreenGameSession
	// ScreenCharacterView displays full character sheet
	ScreenCharacterView
	// ScreenCharacterEdit allows editing character stats
	ScreenCharacterEdit
	// ScreenCombatSetup is the enemy entry screen
	ScreenCombatSetup
	// ScreenCombat is the combat interface (Phase 2)
	ScreenCombat
	// ScreenInventory is the inventory management screen (Phase 3)
	ScreenInventory
	// ScreenMagic is the spell casting interface (Phase 4)
	ScreenMagic
	// ScreenSettings is the settings/configuration screen (Phase 5)
	ScreenSettings
	// ScreenDiceRoll is the dice rolling interface
	ScreenDiceRoll
)

// Model is the root Bubble Tea model containing all application state.
type Model struct {
	// CurrentScreen tracks which screen is active
	CurrentScreen Screen

	// Character is the currently loaded character (nil if none)
	Character *character.Character

	// Dice roller for all random events
	Dice dice.Roller

	// Configuration
	Config *config.Config

	// Screen-specific models
	MainMenu        MainMenuModel
	CharCreation    CharacterCreationModel
	LoadChar        LoadCharacterModel
	GameSession     GameSessionModel
	CharView        CharacterViewModel
	CharEdit        CharacterEditModel
	CombatSetup     CombatSetupModel
	CombatView      CombatViewModel
	CombatState     *combat.CombatState
	Inventory       InventoryManagementModel
	SpellCasting    SpellCastingModel
	Settings        SettingsModel
	DiceRoll        DiceRollModel

	// Help modal state
	ShowingHelp    bool
	HelpScreen     help.Screen
	HelpScroll     int
	HelpMaxScroll  int

	// Application state
	Width  int // Terminal width
	Height int // Terminal height
	Err    error // Last error encountered
}

// NewModel creates a new root model with initial state.
func NewModel() Model {
	roller := dice.NewStandardRoller()

	// Load configuration
	cfg, err := config.LoadDefault()
	if err != nil {
		// Use default config if load fails
		cfg = config.Default()
	}

	return Model{
		CurrentScreen: ScreenMainMenu,
		Character:     nil,
		Dice:          roller,
		Config:        cfg,
		MainMenu:      NewMainMenuModel(),
		CharCreation:  NewCharacterCreationModel(roller),
		LoadChar:      NewLoadCharacterModel(),
		GameSession:   NewGameSessionModel(),
		CharView:      NewCharacterViewModel(),
		CharEdit:      NewCharacterEditModel(),
		CombatSetup:   NewCombatSetupModel(),
		CombatView:    CombatViewModel{}, // Will be initialized when combat starts
		CombatState:   nil,
		Settings:      NewSettingsModel(cfg),
		DiceRoll:      NewDiceRollModel(roller),
		ShowingHelp:   false,
		HelpScreen:    help.ScreenGlobal,
		HelpScroll:    0,
		HelpMaxScroll: 0,
		Width:         80,
		Height:        24,
		Err:           nil,
	}
}

// LoadCharacter loads a character and transitions to the game session.
func (m *Model) LoadCharacter(char *character.Character) {
	m.Character = char
	m.CurrentScreen = ScreenGameSession
	m.CharView.SetCharacter(char)
	m.CharEdit.SetCharacter(char)
}

// SaveCharacter saves the current character to the configured location.
func (m *Model) SaveCharacter() error {
	if m.Character == nil {
		return nil
	}
	// Use configured save directory
	saveDir := m.Config.SaveDirectory
	if saveDir == "" {
		saveDir = "."
	}
	return m.Character.Save(saveDir)
}

// ShowHelp displays the help modal for the specified screen.
func (m *Model) ShowHelp(screen help.Screen) {
	m.ShowingHelp = true
	m.HelpScreen = screen
	m.HelpScroll = 0
	// Calculate max scroll based on content and terminal height
	lines := help.GetLines(screen)
	visibleLines := m.Height - 6 // Account for header and footer
	if visibleLines < 10 {
		visibleLines = 10
	}
	maxScroll := len(lines) - visibleLines
	if maxScroll < 0 {
		maxScroll = 0
	}
	m.HelpMaxScroll = maxScroll
}

// HideHelp closes the help modal.
func (m *Model) HideHelp() {
	m.ShowingHelp = false
	m.HelpScroll = 0
}

// ScrollHelpUp scrolls the help content up.
func (m *Model) ScrollHelpUp() {
	if m.HelpScroll > 0 {
		m.HelpScroll--
	}
}

// ScrollHelpDown scrolls the help content down.
func (m *Model) ScrollHelpDown() {
	if m.HelpScroll < m.HelpMaxScroll {
		m.HelpScroll++
	}
}
