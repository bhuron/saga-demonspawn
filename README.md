# Sagas of the Demonspawn - Rules Engine

A command-line companion application for the "Sagas of the Demonspawn" gamebook, built with Go and Bubble Tea.

## Overview

This application serves as a rules engine and character management tool for players of the gamebook. It handles:
- Character creation and stat management
- Combat resolution with automated calculations
- Magic system (unlocked during gameplay)
- Inventory and equipment management

## Features

### Phase 1: Foundation & Character Management ✓ (Complete)
- [x] Project structure and Go module setup
- [x] Character creation with dice rolling (7 characteristics)
- [x] Character stat editing
- [x] Save/load characters to JSON
- [x] Main menu navigation

### Phase 2: Combat System ✓ (Complete)
- [x] Combat engine with initiative calculation
- [x] To-hit rolls with skill and luck modifiers
- [x] Damage calculation with strength and equipment bonuses
- [x] Stamina-based endurance system
- [x] Death save mechanism
- [x] Enemy data entry and combat UI
- [x] Turn-based combat with combat log

### Phase 3: Items & Inventory ✓ (Complete)
- [x] Equipment system (weapons, armor, shields)
- [x] Inventory management UI with scrolling
- [x] Special items (Healing Stone, Doombringer, The Orb)
- [x] Item acquisition and equipment switching
- [x] Healing Stone charge tracking and recharge
- [x] Shield/Orb mutual exclusion rules
- [x] Combat equipment lock

### Phase 4: Magic System ✓ (Complete)
- [x] Magic engine with spell catalog
- [x] All 10 spells implemented (ARMOUR, CRYPT, FIREBALL, INVISIBILITY, PARALYSIS, POISON NEEDLE, RESURRECTION, RETRACE, TIMEWARP, XENOPHOBIA)
- [x] Natural inclination check (optional)
- [x] Fundamental Failure Rate (FFR) check
- [x] POW cost validation with LP sacrifice
- [x] Spell casting UI with clear feedback
- [x] Magic unlock via Character Edit (Press 'U')
- [x] Combat integration (ARMOUR, XENOPHOBIA, FIREBALL, etc.)
- [x] ActiveSpellEffects tracking
- [x] Unit tests for magic engine

### Phase 5: Polish ✓ (Complete)
- [x] Lipgloss theming system with dark/light schemes
- [x] Comprehensive help system with context-sensitive content
- [x] Configuration management with Settings screen
- [x] User preferences persistence (theme, auto-save, Unicode, etc.)
- [x] Help modal overlay with scrolling
- [x] Enhanced main menu with Settings and Help options
- [x] Auto-save on exit (configurable)
- [ ] Complete screen theming (foundation ready, partial implementation)
- [ ] Enhanced error messages with visual styling
- [ ] Additional quality of life improvements

## Installation

```bash
# Clone the repository
git clone https://github.com/benoit/saga-demonspawn.git
cd saga-demonspawn

# Build the application
go build -o saga ./cmd/saga

# Run the application
./saga
```

## Usage

Run the application and follow the on-screen menu:

```bash
./saga
```

Navigation is keyboard-only:
- Arrow keys: Navigate menus
- Enter: Select option
- Esc/q: Go back/quit
- e: Edit character stats (when viewing character)
- ?: Open context-sensitive help on any screen

### New in Phase 5: Polish & Configuration

**Settings:**
Access from Main Menu → Settings to configure:
- Color scheme (dark/light)
- Unicode/ASCII character display
- Auto-save on exit
- And more appearance and gameplay options

**Help System:**
- Press `?` on any screen for context-specific help
- Select "Help" from Main Menu for comprehensive guide
- Help covers all systems: character, combat, magic, inventory, rules

**Configuration File:**
Settings are saved to `~/.saga-demonspawn/config.json` and persist across sessions.

## Project Structure

```
saga-demonspawn/
├── cmd/saga/           # Main application entry point
├── internal/           # Private application packages
│   ├── character/      # Character state and operations
│   ├── combat/         # Combat resolution engine
│   ├── dice/           # Random number generation
│   ├── items/          # Inventory and equipment
│   ├── magic/          # Spell casting system
│   └── rules/          # Game rules constants
├── pkg/ui/             # Bubble Tea UI components
└── data/               # Game data (enemies, items)
```

## Learning Objectives

This project is designed as a learning tool for Go development, demonstrating:

### Go Fundamentals
- Module initialization and dependency management
- Package organization (internal vs. pkg)
- Struct design with methods
- Interface-based polymorphism
- Error handling patterns

### Testing
- Table-driven tests
- Test fixtures and golden files
- Mock implementations with interfaces
- Code coverage analysis

### Bubble Tea Framework
- The Elm Architecture (Model-View-Update)
- Immutable state management
- Message passing and commands
- Screen navigation

### Best Practices
- Clean code and separation of concerns
- Documentation comments
- Idiomatic Go patterns
- Configuration management

## Game Rules

The complete ruleset is documented in `saga_demonspawn_ruleset.md`. Key mechanics:

**Character Stats (Characteristics):**
- STR (Strength), SPD (Speed), STA (Stamina)
- CRG (Courage), LCK (Luck)
- CHM (Charm), ATT (Attraction)
- Each rolled as 2d6 × 8 (range 16-96, nobody is perfect!)

**Derived Values:**
- Life Points (LP) = Sum of all characteristics
- Skill (SKL) = Starts at 0, +1 per enemy defeated
- Power (POW) = Acquired during adventure (initially 0)

**Combat:**
- Initiative: 2d6 + SPD + CRG + LCK
- To-Hit: Roll 7+ on 2d6 (modified by SKL and LCK)
- Damage: (Roll × 5) + STR bonus + Weapon bonus - Armor

**Special Mechanics:**
- Death saves (once per combat)
- Stamina-based endurance
- Special weapons (Doombringer, The Orb)
- Magic system (unlocked during play)

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/character
```

### Building

```bash
# Build for current platform
go build -o saga ./cmd/saga

# Build with optimizations
go build -ldflags="-s -w" -o saga ./cmd/saga
```

## Contributing

This is primarily a learning project, but suggestions and improvements are welcome!

## License

This project is for educational purposes. The "Sagas of the Demonspawn" gamebook and its rules are property of their respective copyright holders.

## Acknowledgments

- Bubble Tea framework by Charm
- The Sagas of the Demonspawn gamebook
- Go community for excellent documentation and examples
