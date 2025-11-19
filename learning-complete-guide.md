# Complete Learning Guide: Sagas of the Demonspawn

A comprehensive guide to understanding this Go + Bubble Tea educational project.

---

## Table of Contents

1. [Project Overview](#project-overview)
2. [Architecture & Design](#architecture--design)
3. [Reading Order](#reading-order)
4. [Core Systems Deep Dive](#core-systems-deep-dive)
5. [Go Language Features](#go-language-features)
6. [Bubble Tea Patterns](#bubble-tea-patterns)
7. [Testing Strategy](#testing-strategy)
8. [Phase-by-Phase Evolution](#phase-by-phase-evolution)

---

## Project Overview

### What This Project Is

A command-line companion application for the "Sagas of the Demonspawn" gamebook that demonstrates:
- **Clean Architecture** in Go
- **The Elm Architecture** via Bubble Tea
- **Test-Driven Development** with table-driven tests
- **Domain-Driven Design** with clear business logic separation
- **Progressive Enhancement** through 5 development phases

### Project Statistics

- **Total Lines**: ~8,000+ lines of Go code
- **Packages**: 8 internal packages, 1 public UI package
- **Test Coverage**: 100% of core business logic
- **Phases**: 5 complete development phases
- **Screens**: 12+ interactive UI screens

### Key Technologies

- **Go 1.24**: Modern Go with generics and improved error handling
- **Bubble Tea**: TUI framework implementing The Elm Architecture
- **Lipgloss**: Terminal styling and theming
- **Standard Library**: Heavy use of `encoding/json`, `os`, `filepath`, `time`

---

## Architecture & Design

### High-Level Structure

```
saga-demonspawn/
├── cmd/saga/              # Application entry point (minimal)
├── internal/              # Private business logic
│   ├── character/         # Character domain model
│   ├── combat/            # Combat rules engine
│   ├── config/            # Configuration management
│   ├── dice/              # RNG abstraction
│   ├── help/              # Help system with embedded content
│   ├── items/             # Equipment and special items
│   ├── magic/             # Spell casting system
│   └── rules/             # Game constants (future)
└── pkg/ui/                # Public UI layer (Bubble Tea)
    ├── theme/             # Theming system
    └── *.go               # Screen models and views
```

### Design Principles

**1. Separation of Concerns**
- Business logic in `internal/` knows nothing about UI
- UI in `pkg/ui/` orchestrates but doesn't implement rules
- Each package has a single, clear responsibility

**2. Dependency Inversion**
- `dice.Roller` interface allows deterministic testing
- Character methods don't depend on UI state
- Combat engine is pure calculation

**3. Immutable State (Bubble Tea)**
- State updates return new state, never mutate in place
- Messages drive all state transitions
- Commands represent side effects

**4. Progressive Complexity**
- Phase 1: Foundation (character, save/load)
- Phase 2: Combat system
- Phase 3: Inventory management
- Phase 4: Magic system
- Phase 5: Polish and UX

---

## Reading Order

### For Complete Beginners

**Start Here** (Foundation):
1. `README.md` - Project overview and features
2. `saga_demonspawn_ruleset.md` - Game rules (understand the domain)
3. `cmd/saga/main.go` - Entry point (5 lines!)
4. `internal/dice/dice.go` - Simple abstraction example

**Core Architecture** (Understanding the flow):
5. `pkg/ui/model.go` - Root application state
6. `pkg/ui/update.go` - Message routing (Elm Architecture)
7. `pkg/ui/view.go` - Rendering logic

**Simple Screens** (Learn patterns):
8. `pkg/ui/main_menu.go` - Simplest screen model
9. `pkg/ui/game_session.go` - Menu with conditional items
10. `pkg/ui/load_character.go` - File system interaction

**Domain Logic** (Business rules):
11. `internal/character/character.go` - Core domain model
12. `internal/items/items.go` - Declarative data
13. `internal/combat/combat.go` - Complex calculations

### For Intermediate Developers

**Focus on Patterns**:
1. Character creation workflow (`character_creation.go`) - Multi-step state machine
2. Character editing (`character_edit.go`) - Input handling and validation
3. Combat system (`combat.go` + `combat_view.go`) - Complex state management
4. Magic system (`magic/` package) - Validation chains and effects
5. Theming system (`theme/theme.go`) - Reusable UI components

**Testing Approach**:
6. `internal/combat/combat_test.go` - Table-driven tests
7. `internal/magic/casting_test.go` - Mock dependencies
8. `internal/character/character_test.go` - Domain model testing

### For Advanced Study

**Architecture Decisions**:
1. Why `internal/` vs `pkg/`?
2. How does Bubble Tea's message passing work?
3. Why use interfaces for dice rolling?
4. How is state synchronized between screens?

**Extension Points**:
5. Adding a new spell
6. Creating a new screen
7. Implementing a new special item
8. Adding enemy presets

---

## Core Systems Deep Dive

### 1. Character System (`internal/character/`)

**Purpose**: Manage Fire*Wolf's complete state and progression.

**Key Concepts**:
```go
type Character struct {
    // Core characteristics (rolled 2d6 × 8)
    Strength, Speed, Stamina, Courage, Luck, Charm, Attraction int
    
    // Derived values
    CurrentLP, MaximumLP int  // Life Points = sum of characteristics
    Skill int                  // +1 per enemy defeated
    
    // Magic (unlocked during adventure)
    CurrentPOW, MaximumPOW int
    MagicUnlocked bool
    ActiveSpellEffects map[string]int
    
    // Equipment
    EquippedWeapon *items.Weapon
    EquippedArmor *items.Armor
    HasShield bool
    
    // Special items
    HealingStoneCharges int
    DoombringerPossessed, OrbPossessed, OrbEquipped, OrbDestroyed bool
    
    // Progress
    EnemiesDefeated int
    CreatedAt, LastSaved time.Time
}
```

**Design Decisions**:
- **Pointer receivers** for all modifying methods (efficiency + correctness)
- **Validation** at package boundaries (NewCharacter, ModifyX methods)
- **JSON tags** for clean serialization
- **Timestamped saves** for easy versioning
- **Derived values** calculated on demand (GetArmorProtection)

**Learning Points**:
- How to design a rich domain model
- Validation patterns in Go
- JSON marshaling with custom types
- File I/O with error wrapping

### 2. Combat System (`internal/combat/`)

**Purpose**: Implement turn-based combat with all game rules.

**State Machine**:
```
Initialize → Roll Initiative → [Player Turn ↔ Enemy Turn] → Check Victory/Defeat
                                      ↓
                              Check Endurance → Rest if needed
                                      ↓
                              Death Save if LP ≤ 0
```

**Key Functions**:
```go
// Pure calculations (no side effects)
CalculateInitiative(player, enemy, roller) (int, int, bool)
CalculateToHitRequirement(skill, luck) int
CalculateDamage(roll, strength, weaponBonus) int

// State mutations
ExecutePlayerAttack(state, player, roller) AttackResult
ExecuteEnemyAttack(state, player, roller) AttackResult
NextTurn(state)  // Advances turn and round

// Combat lifecycle
StartCombat(player, enemy, roller) *CombatState
AttemptDeathSave(player, state, roller) (int, bool)
ResolveCombatVictory(player)
```

**Design Decisions**:
- **Separate calculation from state** (pure functions testable in isolation)
- **CombatState encapsulates everything** (no hidden state)
- **AttackResult struct** for detailed feedback
- **Combat log** for player transparency

**Learning Points**:
- State machine implementation
- Pure vs impure functions
- Complex calculation testing
- Turn-based game logic

### 3. Magic System (`internal/magic/`)

**Purpose**: Spell casting with validation, costs, and effects.

**Validation Chain**:
```
1. Natural Inclination Check (optional, 2d6 ≥ 4)
2. Power Cost Check (can afford or sacrifice LP?)
3. Fundamental Failure Rate (2d6 ≥ 6)
4. Apply Effect
```

**Spell Categories**:
- **Offensive**: FIREBALL, POISON NEEDLE
- **Defensive**: ARMOUR, XENOPHOBIA
- **Tactical**: INVISIBILITY, PARALYSIS
- **Navigation**: CRYPT, RETRACE, TIMEWARP
- **Recovery**: RESURRECTION

**Design Decisions**:
- **Spell catalog** as slice of structs (data-driven)
- **Validation separate from effects** (single responsibility)
- **Context-aware availability** (combat vs non-combat)
- **Effect structs** for complex results

**Learning Points**:
- Validation chains
- Data-driven design
- Effect systems
- Context-sensitive logic

### 4. UI System (`pkg/ui/`)

**Purpose**: Bubble Tea screens implementing The Elm Architecture.

**The Elm Architecture**:
```
Model (State) → View (Render) → User Input → Update (New State) → Model
```

**Screen Pattern**:
```go
// Every screen follows this pattern:

type ScreenModel struct {
    // Screen-specific state
    cursor int
    items []string
    inputBuffer string
    // ...
}

func (m ScreenModel) Update(msg tea.Msg) (ScreenModel, tea.Cmd) {
    // Handle messages, return new state + optional command
}

func (m ScreenModel) View() string {
    // Pure rendering from state
}
```

**Root Model Coordination**:
```go
type Model struct {
    CurrentScreen Screen
    Character *character.Character
    
    // Screen models
    MainMenu MainMenuModel
    CharCreation CharacterCreationModel
    // ... one for each screen
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Route to appropriate screen handler
    switch m.CurrentScreen {
    case ScreenMainMenu:
        return m.handleMainMenuKeys(msg)
    // ...
    }
}
```

**Design Decisions**:
- **One model per screen** (focused responsibility)
- **Root model coordinates** (screen transitions, shared state)
- **Message passing** for async operations
- **Pure view functions** (no side effects in rendering)

**Learning Points**:
- The Elm Architecture
- State management patterns
- Message-driven programming
- Screen navigation

### 5. Theming System (`pkg/ui/theme/`)

**Purpose**: Consistent styling across all screens.

**Theme Structure**:
```go
type Theme struct {
    // Colors
    Primary, Secondary, Success, Warning, Danger lipgloss.Color
    
    // Typography
    Title, Heading, Label, Value, Body lipgloss.Style
    
    // Components
    Box, MenuItem, Button, Panel lipgloss.Style
    
    // Settings
    UseUnicode bool
}
```

**Utility Functions**:
```go
theme.RenderTitle(text)                    // Styled screen title
theme.RenderMenuItem(text, selected)       // Menu item with cursor
theme.RenderHealthBar(current, max, width) // Visual HP bar
theme.RenderPOWMeter(current, max, width)  // Visual POW gauge
theme.RenderError(title, desc, suggestion) // Structured error
```

**Design Decisions**:
- **Centralized theme** (single source of truth)
- **Utility functions** (DRY principle)
- **Dark/light schemes** (user preference)
- **Unicode fallback** (terminal compatibility)

**Learning Points**:
- Lipgloss styling
- Theme abstraction
- Reusable UI components
- Terminal capabilities

---

## Go Language Features

### Features Demonstrated

**1. Modules and Packages**
```go
// go.mod defines module
module github.com/benoit/saga-demonspawn

// internal/ packages are private
import "github.com/benoit/saga-demonspawn/internal/character"

// pkg/ packages are public
import "github.com/benoit/saga-demonspawn/pkg/ui"
```

**2. Interfaces**
```go
// dice.Roller enables testing
type Roller interface {
    Roll1D6() int
    Roll2D6() int
}

// Multiple implementations
type StandardRoller struct { rand *rand.Rand }
type SeededRoller struct { rand *rand.Rand }
```

**3. Struct Tags**
```go
type Character struct {
    Strength int `json:"strength"`  // JSON serialization
    Speed    int `json:"speed"`
}
```

**4. Pointer Receivers**
```go
// Mutating methods use pointers
func (c *Character) ModifyLP(delta int) {
    c.CurrentLP += delta
}

// Read-only methods can use values
func (c Character) IsAlive() bool {
    return c.CurrentLP > 0
}
```

**5. Error Handling**
```go
// Explicit error returns
func New(str, spd int) (*Character, error) {
    if str < 0 {
        return nil, fmt.Errorf("strength cannot be negative: %d", str)
    }
    return &Character{Strength: str}, nil
}

// Error wrapping
if err := os.WriteFile(path, data, 0644); err != nil {
    return fmt.Errorf("failed to write save file: %w", err)
}
```

**6. Embedding**
```go
//go:embed content/*.txt
var helpContent embed.FS
```

**7. Table-Driven Tests**
```go
tests := []struct {
    name string
    input int
    want int
}{
    {"zero", 0, 0},
    {"positive", 5, 25},
}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        got := Calculate(tt.input)
        if got != tt.want {
            t.Errorf("got %d, want %d", got, tt.want)
        }
    })
}
```

**8. Composition**
```go
// Model composes screen models
type Model struct {
    MainMenu MainMenuModel
    CharCreation CharacterCreationModel
    // No inheritance, just composition
}
```

---

## Bubble Tea Patterns

### The Elm Architecture

**Model**: Application state
```go
type Model struct {
    CurrentScreen Screen
    Character *character.Character
    Width, Height int
}
```

**Update**: State transitions
```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        return m.handleKeyPress(msg)
    case tea.WindowSizeMsg:
        m.Width, m.Height = msg.Width, msg.Height
        return m, nil
    }
    return m, nil
}
```

**View**: Pure rendering
```go
func (m Model) View() string {
    switch m.CurrentScreen {
    case ScreenMainMenu:
        return m.viewMainMenu()
    case ScreenCombat:
        return m.viewCombat()
    }
    return ""
}
```

### Message Passing

**Built-in Messages**:
- `tea.KeyMsg` - Keyboard input
- `tea.WindowSizeMsg` - Terminal resize

**Custom Messages**:
```go
type CombatEndMsg struct {
    Victory bool
}

type CastSpellMsg struct{}

// Return command that produces message
return m, func() tea.Msg {
    return CombatEndMsg{Victory: true}
}
```

### Commands

**Commands represent side effects**:
```go
// No side effect
return m, nil

// Quit application
return m, tea.Quit

// Custom async operation
return m, func() tea.Msg {
    // Do work
    return SomeResultMsg{}
}
```

### Screen Navigation

**Pattern**:
```go
// Change screen
m.CurrentScreen = ScreenCombat

// Initialize screen model
m.CombatView = NewCombatViewModel(...)

// Return updated model
return m, nil
```

---

## Testing Strategy

### Unit Tests

**Pure Functions** (easiest to test):
```go
func TestCalculateDamage(t *testing.T) {
    tests := []struct {
        roll, str, weapon, want int
    }{
        {8, 65, 10, 85},  // (8×5) + (6×5) + 10 = 85
        {2, 10, 5, 20},   // (2×5) + (1×5) + 5 = 20
    }
    for _, tt := range tests {
        got := CalculateDamage(tt.roll, tt.str, tt.weapon)
        if got != tt.want {
            t.Errorf("got %d, want %d", got, tt.want)
        }
    }
}
```

**With Dependencies** (use mocks):
```go
type MockRoller struct {
    NextRoll int
}

func (m *MockRoller) Roll2D6() int {
    return m.NextRoll
}

func TestWithDice(t *testing.T) {
    roller := &MockRoller{NextRoll: 7}
    result := SomeFunction(roller)
    // Test with deterministic roll
}
```

**Domain Models** (test invariants):
```go
func TestCharacterValidation(t *testing.T) {
    _, err := character.New(-5, 50, 50, 50, 50, 50, 50)
    if err == nil {
        t.Error("expected error for negative strength")
    }
}
```

### Test Coverage

**What to Test**:
- ✅ All calculation functions
- ✅ Domain model validation
- ✅ State transitions
- ✅ Edge cases (zero, negative, max values)

**What Not to Test**:
- ❌ UI rendering (too brittle)
- ❌ File I/O (use integration tests)
- ❌ Third-party libraries

### Running Tests

```bash
# All tests
go test ./...

# With coverage
go test -cover ./...

# Specific package
go test ./internal/combat/

# Verbose output
go test -v ./internal/magic/
```

---

## Phase-by-Phase Evolution

### Phase 1: Foundation (Character Management)

**Goal**: Basic character creation, viewing, editing, save/load.

**What You Learn**:
- Go project structure
- Bubble Tea basics
- JSON serialization
- File I/O
- Simple state machines

**Key Files**:
- `internal/character/character.go`
- `pkg/ui/character_creation.go`
- `pkg/ui/character_edit.go`

### Phase 2: Combat System

**Goal**: Turn-based combat with all game rules.

**What You Learn**:
- Complex state management
- Pure function design
- Table-driven testing
- Message passing
- State machines

**Key Files**:
- `internal/combat/combat.go`
- `internal/combat/combat_test.go`
- `pkg/ui/combat_view.go`

### Phase 3: Inventory Management

**Goal**: Equipment system with special items.

**What You Learn**:
- Data-driven design
- Mutual exclusions
- Context-sensitive actions
- Scrolling viewports

**Key Files**:
- `internal/items/items.go`
- `pkg/ui/inventory_management.go`

### Phase 4: Magic System

**Goal**: Spell casting with validation and effects.

**What You Learn**:
- Validation chains
- Effect systems
- Context-aware logic
- Integration patterns

**Key Files**:
- `internal/magic/spells.go`
- `internal/magic/casting.go`
- `internal/magic/effects.go`

### Phase 5: Polish & UX

**Goal**: Professional styling and user experience.

**What You Learn**:
- Theming systems
- Configuration management
- Help systems
- Lipgloss styling

**Key Files**:
- `pkg/ui/theme/theme.go`
- `internal/config/config.go`
- `internal/help/help.go`

---

## Common Patterns

### Pattern: Cursor Navigation

```go
type MenuModel struct {
    cursor int
    items []string
}

func (m *MenuModel) MoveUp() {
    if m.cursor > 0 {
        m.cursor--
    }
}

func (m *MenuModel) MoveDown() {
    if m.cursor < len(m.items)-1 {
        m.cursor++
    }
}

func (m MenuModel) GetSelected() string {
    return m.items[m.cursor]
}
```

### Pattern: Input Buffer

```go
type EditModel struct {
    inputMode bool
    inputBuffer string
}

func (m *EditModel) StartInput() {
    m.inputMode = true
    m.inputBuffer = ""
}

func (m *EditModel) AppendInput(s string) {
    m.inputBuffer += s
}

func (m *EditModel) Backspace() {
    if len(m.inputBuffer) > 0 {
        m.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
    }
}
```

### Pattern: Multi-Step Workflow

```go
type Step int

const (
    StepOne Step = iota
    StepTwo
    StepThree
)

type WorkflowModel struct {
    currentStep Step
    // Step-specific state
}

func (m *WorkflowModel) NextStep() {
    m.currentStep++
}

func (m *WorkflowModel) PreviousStep() {
    if m.currentStep > 0 {
        m.currentStep--
    }
}
```

### Pattern: Validation

```go
func Validate(value int) error {
    if value < 0 {
        return fmt.Errorf("value cannot be negative: %d", value)
    }
    if value > 100 {
        return fmt.Errorf("value exceeds maximum: %d", value)
    }
    return nil
}

// Usage
if err := Validate(input); err != nil {
    m.errorMsg = err.Error()
    return m, nil
}
```

---

## Extension Exercises

### Beginner

1. **Add a new weapon** to `internal/items/items.go`
2. **Change stat roll formula** in `internal/dice/dice.go`
3. **Add a new menu option** to main menu
4. **Modify theme colors** in `pkg/ui/theme/theme.go`

### Intermediate

5. **Create enemy presets** (load from JSON file)
6. **Add a new spell** with custom effect
7. **Implement combat log scrolling**
8. **Add character export** (different format)

### Advanced

9. **Implement section system** (gamebook navigation)
10. **Add multiplayer** (two characters in combat)
11. **Create plugin system** for custom spells
12. **Build web UI** (keep same business logic)

---

## Troubleshooting Guide

### Build Issues

**Problem**: `package X is not in GOROOT`
**Solution**: Run `go mod tidy` to download dependencies

**Problem**: `cannot find package`
**Solution**: Check import paths match `go.mod` module name

### Runtime Issues

**Problem**: Character doesn't save
**Solution**: Check `SaveDirectory` in config, ensure directory exists

**Problem**: Theme not applying
**Solution**: Verify `theme.Init()` called before rendering

### Testing Issues

**Problem**: Tests fail with random values
**Solution**: Use `MockRoller` or `SeededRoller` for deterministic tests

**Problem**: Test coverage low
**Solution**: Focus on `internal/` packages, skip UI rendering tests

---

## Resources

### Official Documentation

- [Go Documentation](https://go.dev/doc/)
- [Bubble Tea Tutorial](https://github.com/charmbracelet/bubbletea/tree/master/tutorials)
- [Lipgloss Examples](https://github.com/charmbracelet/lipgloss/tree/master/examples)

### Project Documentation

- `README.md` - Project overview
- `saga_demonspawn_ruleset.md` - Game rules
- `PHASE*_COMPLETE.md` - Phase completion summaries
- `learning-phase1.md` - Phase 1 specific guide

### Code Comments

Every exported function, type, and method has documentation comments explaining:
- Purpose
- Parameters
- Return values
- Side effects
- Example usage (where helpful)

---

## Conclusion

This project demonstrates professional Go development practices:

✅ **Clean Architecture** - Clear separation of concerns  
✅ **Test-Driven Development** - Comprehensive test coverage  
✅ **Domain-Driven Design** - Rich domain models  
✅ **The Elm Architecture** - Predictable state management  
✅ **Progressive Enhancement** - Incremental feature development  

By studying this codebase, you'll learn:
- How to structure a Go application
- How to implement The Elm Architecture
- How to write testable code
- How to manage complex state
- How to build terminal UIs

**Next Steps**:
1. Read through the recommended order
2. Run the application and explore features
3. Read the tests to understand behavior
4. Try the extension exercises
5. Build your own features!

Happy learning! 🚀
