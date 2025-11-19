# Project Structure

## Directory Layout

```
saga-demonspawn/
├── cmd/saga/              # Application entry point
├── internal/              # Private application packages
│   ├── character/         # Character state and operations
│   ├── combat/            # Combat resolution engine
│   ├── config/            # Configuration management
│   ├── dice/              # Random number generation
│   ├── help/              # Help system with content files
│   ├── items/             # Inventory and equipment
│   ├── magic/             # Spell casting system
│   └── rules/             # Game rules constants
├── pkg/ui/                # Public UI components (Bubble Tea)
│   └── theme/             # Theming system
└── data/                  # Game data (enemies, items)
```

## Package Organization

### `internal/` - Private Packages

**Purpose**: Core business logic that should not be imported by external projects.

- **character**: Character struct with methods for stat modification, equipment, special items, save/load
- **combat**: Combat state, initiative, to-hit calculations, damage formulas, endurance tracking
- **config**: User preferences and configuration persistence
- **dice**: Dice rolling interface and implementations (supports mocking for tests)
- **help**: Context-sensitive help content and rendering
- **items**: Equipment definitions (weapons, armor, shields, special items)
- **magic**: Spell catalog, casting mechanics, POW management, spell effects
- **rules**: Game constants and rule definitions

### `pkg/ui/` - Public UI Package

**Purpose**: Bubble Tea components implementing The Elm Architecture.

Each screen has its own model file:
- `main_menu.go` - Initial menu
- `character_creation.go` - Character creation flow
- `character_view.go` - Character sheet display
- `character_edit.go` - Stat editing interface
- `load_character.go` - Character loading
- `game_session.go` - Main game menu
- `combat_setup.go` - Enemy entry
- `combat_view.go` - Combat interface
- `inventory_management.go` / `inventory_view.go` - Inventory system
- `spell_casting.go` - Magic interface
- `settings.go` - Configuration screen

Root files:
- `model.go` - Root model with screen navigation
- `update.go` - Message routing and state updates
- `view.go` - Rendering and help modal

### `cmd/saga/` - Entry Point

Single `main.go` that creates the root model and starts the Bubble Tea program.

## Code Conventions

### Naming

- **Packages**: Lowercase, single word (e.g., `character`, `combat`, `dice`)
- **Files**: Lowercase with underscores (e.g., `character_creation.go`)
- **Types**: PascalCase (e.g., `Character`, `CombatState`)
- **Functions/Methods**: PascalCase for exported, camelCase for private
- **Constants**: PascalCase or SCREAMING_SNAKE_CASE for enums

### Documentation

- All exported types, functions, and methods have doc comments
- Doc comments start with the name of the item being documented
- Package-level doc comments explain the package purpose

### Error Handling

- Return errors rather than panicking
- Wrap errors with context using `fmt.Errorf("context: %w", err)`
- Validate inputs at package boundaries

### Testing

- Test files named `*_test.go` in the same package
- Table-driven tests preferred
- Use `dice.Roller` interface for deterministic testing
- Test coverage for core logic (character, combat, magic)

### State Management

- Immutable state transitions in Bubble Tea components
- Character modifications through methods, not direct field access
- Combat state encapsulated in `CombatState` struct
- Configuration persisted to `~/.saga-demonspawn/config.json`

### File Persistence

- Characters saved as JSON with timestamps: `character_YYYYMMDD-HHMMSS.json`
- Default save location: current directory (configurable)
- JSON marshaling with indentation for readability
