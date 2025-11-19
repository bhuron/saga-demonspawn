// Package help provides context-sensitive help content for the application.
package help

import (
	_ "embed"
	"strings"
)

// Embedded help content files
//go:embed content/global.txt
var globalHelp string

//go:embed content/main_menu.txt
var mainMenuHelp string

//go:embed content/character_creation.txt
var characterCreationHelp string

//go:embed content/character_edit.txt
var characterEditHelp string

//go:embed content/combat.txt
var combatHelp string

//go:embed content/magic.txt
var magicHelp string

// Screen represents different screens that have context-specific help.
type Screen int

const (
	// ScreenGlobal is the comprehensive help guide
	ScreenGlobal Screen = iota
	// ScreenMainMenu is the main menu help
	ScreenMainMenu
	// ScreenCharacterCreation is character creation help
	ScreenCharacterCreation
	// ScreenCharacterEdit is character editing help
	ScreenCharacterEdit
	// ScreenCombat is combat system help
	ScreenCombat
	// ScreenMagic is magic system help
	ScreenMagic
)

// GetHelp returns the help content for the specified screen.
func GetHelp(screen Screen) string {
	switch screen {
	case ScreenGlobal:
		return globalHelp
	case ScreenMainMenu:
		return mainMenuHelp
	case ScreenCharacterCreation:
		return characterCreationHelp
	case ScreenCharacterEdit:
		return characterEditHelp
	case ScreenCombat:
		return combatHelp
	case ScreenMagic:
		return magicHelp
	default:
		return "No help available for this screen.\n\nPress Esc or ? to close."
	}
}

// GetLines splits help content into lines for scrolling.
func GetLines(screen Screen) []string {
	content := GetHelp(screen)
	return strings.Split(content, "\n")
}

// GetTitle returns a brief title for the help screen.
func GetTitle(screen Screen) string {
	switch screen {
	case ScreenGlobal:
		return "COMPREHENSIVE HELP GUIDE"
	case ScreenMainMenu:
		return "MAIN MENU HELP"
	case ScreenCharacterCreation:
		return "CHARACTER CREATION HELP"
	case ScreenCharacterEdit:
		return "CHARACTER EDIT HELP"
	case ScreenCombat:
		return "COMBAT HELP"
	case ScreenMagic:
		return "MAGIC SYSTEM HELP"
	default:
		return "HELP"
	}
}
