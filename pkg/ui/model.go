// Package ui provides Bubble Tea UI components for the Saga application.
package ui

import (
	"github.com/benoit/saga-demonspawn/internal/character"
	"github.com/benoit/saga-demonspawn/internal/combat"
	"github.com/benoit/saga-demonspawn/internal/dice"
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
)

// Model is the root Bubble Tea model containing all application state.
type Model struct {
	// CurrentScreen tracks which screen is active
	CurrentScreen Screen

	// Character is the currently loaded character (nil if none)
	Character *character.Character

	// Dice roller for all random events
	Dice dice.Roller

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

	// Application state
	Width  int // Terminal width
	Height int // Terminal height
	Err    error // Last error encountered
}

// NewModel creates a new root model with initial state.
func NewModel() Model {
	roller := dice.NewStandardRoller()

	return Model{
		CurrentScreen: ScreenMainMenu,
		Character:     nil,
		Dice:          roller,
		MainMenu:      NewMainMenuModel(),
		CharCreation:  NewCharacterCreationModel(roller),
		LoadChar:      NewLoadCharacterModel(),
		GameSession:   NewGameSessionModel(),
		CharView:      NewCharacterViewModel(),
		CharEdit:      NewCharacterEditModel(),
		CombatSetup:   NewCombatSetupModel(),
		CombatView:    CombatViewModel{}, // Will be initialized when combat starts
		CombatState:   nil,
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

// SaveCharacter saves the current character to the default location.
func (m *Model) SaveCharacter() error {
	if m.Character == nil {
		return nil
	}
	// Save to user's home directory/.saga-demonspawn/
	// For now, use current directory
	return m.Character.Save(".")
}
