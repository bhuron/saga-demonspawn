# Go Study Guide for the Saga Demonspawn Codebase

This guide helps you learn Go development and the Bubble Tea architecture using this compact, educational codebase.

---

## 1) Recommended Reading Order (and why)

1. pkg/ui/model.go
   - Start here to understand the root application state, screen registry, and where the dice roller and character are wired. It orients you to the app’s architecture.

2. pkg/ui/update.go
   - See the Elm-style control flow: how messages (keyboard input, window size) are handled and routed to screen-specific logic.

3. Screen models (top → down)
   - pkg/ui/main_menu.go
   - pkg/ui/game_session.go
   - pkg/ui/load_character.go
   - pkg/ui/character_view.go
   - pkg/ui/character_edit.go
   - pkg/ui/character_creation.go
   - Why: These are focused, self-contained models. Reading them after update.go shows how each screen owns a small piece of state and behavior.

4. Domain logic
   - internal/character/character.go
   - Why: Core game data model and business logic (validation, LP, skill/POW, equipment, save/load).

5. Supporting data and utilities
   - internal/items/items.go
   - internal/dice/dice.go
   - Why: Declarative items and RNG abstraction. You’ll see how the UI and character logic consume these.

6. Project context and rules
   - README.md
   - PHASE1_COMPLETE.md
   - saga_demonspawn_ruleset.md
   - Why: Get the big-picture roadmap, what Phase 1 includes, and rules that future phases will implement.

---

## 2) Codebase Organization and Architectural Decisions

- internal/
  - character/: Character struct and methods (validation, LP/POW/Skill management, equipment, JSON save/load).
  - dice/: RNG abstraction via a Roller interface (seeded and standard), enabling deterministic tests and reproducibility.
  - items/: Typed definitions for weapons, armor, shields, and helper accessors (e.g., StartingWeapons).

- pkg/ui/ (Bubble Tea components following the Elm Architecture)
  - model.go: Root Model holds global state, dice roller, and screen models.
  - update.go: Central message loop; delegates to per-screen handlers.
  - main_menu.go, game_session.go, load_character.go, character_view.go, character_edit.go, character_creation.go: Narrowly scoped models managing their own state and cursors.

- Architectural choices
  - Elm-style Model-Update: Predictable state transitions via message dispatch.
  - Separation of concerns: Domain logic in internal/, UI orchestration in pkg/ui/.
  - Interface-driven design: dice.Roller makes randomness testable and pluggable.
  - Simple persistence: Character JSON save/load with timestamped files for easy inspection.
  - Declarative data: Items defined as structs; UI consumes them directly.

---

## 3) Key Go Development Concepts Demonstrated

- Modules and packages: go.mod defines the module and toolchain; cohesive packages under internal/ and pkg/.
- Interfaces and testability: dice.Roller abstraction; seedable RNG for deterministic behavior.
- Structs with JSON tags: Character fields serialize cleanly to human-readable JSON.
- Pointer receivers: Mutating Character methods use pointer receivers for correctness and efficiency.
- Error handling: Validation returns explicit errors; file/JSON operations wrap errors with context.
- File I/O and time: Save/Load with os, filepath, time for timestamped filenames.
- Composition: Model composes screen models; Character composes equipped items.
- Slice navigation patterns: Cursor-based menu navigation across UI models.
- Elm/Bubble Tea patterns: Controlled state transitions via messages and handler functions.

---

## 4) Learning Objectives by Component

- pkg/ui/model.go
  - Compose root app state and share dependencies (dice, character) across screens.
  - Implement lifecycle operations like LoadCharacter/SaveCharacter.

- pkg/ui/update.go
  - Implement message routing for tea.Msg (keyboard/window events).
  - Structure clean switch-based dispatch and delegate to screen-specific handlers.

- pkg/ui/main_menu.go
  - Cursor-based selection; expose minimal getters for the view and update layers.

- pkg/ui/game_session.go
  - Conditional menu construction (show/hide “Cast Spell” when magic is unlocked).
  - Maintain cursor and selection with simple slice operations.

- pkg/ui/load_character.go
  - Scan filesystem for save files (filepath.Glob).
  - Manage selection and bubble up errors; format basic file metadata.

- pkg/ui/character_view.go
  - Minimal view state holder with clean setter/getter for the loaded character.

- pkg/ui/character_edit.go
  - Implement edit workflows: cursor navigation, input buffer, input mode toggling.
  - Use typed enums (EditField) for safe, indexed field selection.
  - Conditionally hide POW fields until Character.MagicUnlocked is true.

- pkg/ui/character_creation.go
  - Build a multi-step workflow (roll → equipment → review).
  - Use dice.Roller for stat rolls (2d6 × 8 → range 16–96).
  - Compute LP = sum of all characteristics; equip selected items; finalize Character via character.New.

- internal/character/character.go
  - Design a robust domain model: validation, derived LP, progression (Skill, POW), flags, counters.
  - Implement JSON persistence with timestamping.
  - Apply equipment effects (GetWeaponDamageBonus, GetArmorProtection); shield/armor interaction.

- internal/items/items.go
  - Define typed items, constants, and helpers (AllWeapons, StartingWeapons, GetWeaponByName).
  - Keep data declarative and UI-friendly.

- internal/dice/dice.go
  - Build an RNG wrapper with seed control.
  - Implement standard dice ops (Roll1D6, Roll2D6) and domain-specific stat rolling.

---

## 5) How to Navigate the Codebase to Maximize Learning

- Follow message flow top-down
  - Start at pkg/ui/model.go for structure.
  - Read pkg/ui/update.go and trace each handler (handleMainMenuKeys → screen methods).
  - Open the referenced screen model and read its state and cursor methods.

- Trace a feature end-to-end (example: Character Creation)
  - update.go → handleCharacterCreationKeys
  - character_creation.go: rolling, selection, CreateCharacter
  - internal/character/character.go: validation, LP, equipment application
  - Back to model.go: LoadCharacter wiring to view/edit screens

- Make safe, testable tweaks
  - Change RollCharacteristic multiplier in internal/dice/dice.go to see effect on rolled stats.
  - Add a new starting weapon in internal/items/items.go and confirm it appears in selection.
  - Adjust armor or shield values and confirm GetArmorProtection behavior.

- Build and run cleanly
  - Use ./build.sh after changes to avoid stale binaries and ensure a clean rebuild.

- Respect separation of concerns
  - Keep business logic in internal/ and UI logic in pkg/ui/.
  - Use Character.Save/Load for persistence; let the root Model coordinate storage, not screen models.

- Keep rules aligned
  - Use saga_demonspawn_ruleset.md as your source of truth for mechanics—especially when extending combat, inventory, or magic in future phases.

---

## Notes on Game Mechanics Reflected in Code

- Characteristics are rolled 2d6 × 8 (16–96), intentionally excluding 100 to fit the “nobody is perfect” design.
- LP = sum of all characteristics at creation, tracked as CurrentLP and MaximumLP.
- Equipment modifies damage/protection; shields provide reduced protection when worn with armor.
- Phase 1 focuses on creation, viewing, editing, loading, and saving; combat, inventory, and magic are planned for later phases.