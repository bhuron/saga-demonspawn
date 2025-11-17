# Character System

<cite>
**Referenced Files in This Document**
- [internal/character/character.go](file://internal/character/character.go)
- [pkg/ui/model.go](file://pkg/ui/model.go)
- [pkg/ui/character_creation.go](file://pkg/ui/character_creation.go)
- [pkg/ui/character_edit.go](file://pkg/ui/character_edit.go)
- [pkg/ui/character_view.go](file://pkg/ui/character_view.go)
- [internal/dice/dice.go](file://internal/dice/dice.go)
- [internal/items/items.go](file://internal/items/items.go)
- [BUGFIX_CHARACTER_EDIT.md](file://BUGFIX_CHARACTER_EDIT.md)
- [README.md](file://README.md)
</cite>

## Table of Contents
1. [Introduction](#introduction)
2. [Character Entity Architecture](#character-entity-architecture)
3. [Core Characteristics System](#core-characteristics-system)
4. [Derived Values and Progression](#derived-values-and-progression)
5. [Equipment Management](#equipment-management)
6. [Character Creation Workflow](#character-creation-workflow)
7. [Character Editing System](#character-editing-system)
8. [Persistence and Serialization](#persistence-and-serialization)
9. [UI Integration Layer](#ui-integration-layer)
10. [Validation and Business Rules](#validation-and-business-rules)
11. [Common Issues and Edge Cases](#common-issues-and-edge-cases)
12. [Practical Implementation Examples](#practical-implementation-examples)

## Introduction

The Character System in the Saga of the Demonspawn application provides comprehensive character management capabilities for the gamebook simulation. Built around a robust data model, it handles everything from initial character creation with dice rolls to ongoing stat modifications and equipment management. The system follows the original gamebook rules while providing modern persistence and UI capabilities.

The character system operates on a foundation of seven core characteristics (STR, SPD, STA, CRG, LCK, CHM, ATT) that determine all aspects of character performance, from combat effectiveness to magical abilities. These characteristics form the basis for derived values like Life Points, Skill levels, and Power reserves, creating a dynamic progression system that evolves throughout gameplay.

## Character Entity Architecture

The Character struct serves as the central data model for all character-related information in the application. It encapsulates both static attributes established during creation and dynamic state that changes throughout gameplay.

```mermaid
classDiagram
class Character {
+int Strength
+int Speed
+int Stamina
+int Courage
+int Luck
+int Charm
+int Attraction
+int CurrentLP
+int MaximumLP
+int Skill
+int CurrentPOW
+int MaximumPOW
+bool MagicUnlocked
+Weapon EquippedWeapon
+Armor EquippedArmor
+bool HasShield
+int EnemiesDefeated
+time.Time CreatedAt
+time.Time LastSaved
+New(str, spd, sta, crg, lck, chm, att) Character
+ModifyStrength(delta) error
+ModifySpeed(delta) error
+ModifyStamina(delta) error
+ModifyCourage(delta) error
+ModifyLuck(delta) error
+ModifyCharm(delta) error
+ModifyAttraction(delta) error
+ModifyLP(delta)
+SetLP(value)
+SetMaxLP(value) error
+ModifySkill(delta) error
+SetSkill(value) error
+UnlockMagic(initialPOW) error
+ModifyPOW(delta)
+SetPOW(value)
+SetMaxPOW(value) error
+EquipWeapon(weapon)
+EquipArmor(armor)
+ToggleShield()
+IncrementEnemiesDefeated()
+IsAlive() bool
+GetArmorProtection() int
+GetWeaponDamageBonus() int
+Save(directory) error
+Load(filepath) Character
}
class Weapon {
+string Name
+int DamageBonus
+string Description
+bool Special
}
class Armor {
+string Name
+int Protection
+string Description
}
Character --> Weapon : "equips"
Character --> Armor : "equips"
```

**Diagram sources**
- [internal/character/character.go](file://internal/character/character.go#L14-L44)
- [internal/items/items.go](file://internal/items/items.go#L20-L41)

**Section sources**
- [internal/character/character.go](file://internal/character/character.go#L14-L44)

## Core Characteristics System

The seven core characteristics represent fundamental aspects of a character's personality and physical abilities. Each characteristic is rolled independently using the gamebook's standard mechanic of rolling 2d6 and multiplying by 8, producing values in the range of 16-96.

### Characteristic Definitions

| Characteristic | Abbreviation | Purpose | Range (0-999) |
|----------------|--------------|---------|----------------|
| Strength | STR | Physical power and melee damage | 16-96 (gamebook), 0-999 (validation) |
| Speed | SPD | Agility, reaction time, initiative | 16-96 (gamebook), 0-999 (validation) |
| Stamina | STA | Endurance and hit points | 16-96 (gamebook), 0-999 (validation) |
| Courage | CRG | Mental fortitude and morale | 16-96 (gamebook), 0-999 (validation) |
| Luck | LCK | Fortune and chance encounters | 16-96 (gamebook), 0-999 (validation) |
| Charm | CHM | Social influence and charisma | 16-96 (gamebook), 0-999 (validation) |
| Attraction | ATT | Personal magnetism and appeal | 16-96 (gamebook), 0-999 (validation) |

### Life Points Calculation

Life Points (LP) are calculated as the mathematical sum of all seven characteristics at character creation. This fundamental mechanic ensures that characters with higher base stats have greater durability and resilience.

```mermaid
flowchart TD
Start([Character Creation]) --> RollStats["Roll 2d6 × 8 for Each Stat<br/>STR, SPD, STA, CRG, LCK, CHM, ATT"]
RollStats --> ValidateRange["Validate Characteristic Range<br/>0-999 (no negatives, cap at 999)"]
ValidateRange --> CalcLP["Calculate LP = Σ(All Characteristics)<br/>Example: 64 + 56 + 72 + 48 + 80 + 40 + 32 = 392"]
CalcLP --> InitValues["Initialize Derived Values:<br/>• CurrentLP = MaximumLP = Calculated LP<br/>• Skill = 0<br/>• POW = 0<br/>• MagicUnlocked = false"]
InitValues --> EquipDefault["Equip Default Equipment:<br/>• Weapon: Sword<br/>• Armor: None<br/>• Shield: False"]
EquipDefault --> SetTimestamps["Set Creation Timestamps"]
SetTimestamps --> Complete([Character Ready])
```

**Diagram sources**
- [internal/character/character.go](file://internal/character/character.go#L47-L96)
- [internal/dice/dice.go](file://internal/dice/dice.go#L60-L63)

**Section sources**
- [internal/character/character.go](file://internal/character/character.go#L47-L96)
- [internal/dice/dice.go](file://internal/dice/dice.go#L60-L63)

## Derived Values and Progression

The character system implements several derived values that evolve throughout gameplay, providing meaningful progression mechanics and strategic depth.

### Life Points (LP) System

Life Points serve as the primary measure of a character's health and survivability. Unlike traditional RPG systems, LP in this implementation can go negative, representing critical injuries or death conditions.

```mermaid
stateDiagram-v2
[*] --> Alive : CurrentLP > 0
Alive --> Injured : CurrentLP ≤ 0
Injured --> NearDeath : CurrentLP < -10
Injured --> Recovering : Heal effect applied
NearDeath --> Dead : No recovery
NearDeath --> Alive : Successful death save
Recovering --> Alive : CurrentLP > 0
Dead --> [*]
note right of Alive : Can fight normally
note right of Injured : Reduced effectiveness
note right of NearDeath : Must make death saves
note right of Dead : Character permanently deceased
```

### Skill Progression System

Skill represents a character's combat proficiency and tactical acumen. It starts at zero and increases by one point for each enemy defeated, affecting hit chances and combat effectiveness.

### Magic System Unlocking

The magic system becomes available during gameplay when the character unlocks magical abilities. This introduces a significant shift in gameplay dynamics, requiring careful resource management of Power (POW) reserves.

```mermaid
sequenceDiagram
participant Player as Player
participant Game as Game Engine
participant Char as Character
participant UI as UI Layer
Player->>Game : Encounter magical event
Game->>Char : Check magic prerequisites
Char->>Char : Evaluate conditions
alt Magic not yet unlocked
Char-->>Game : MagicUnlocked = false
Game-->>Player : Cannot cast magic
else Magic ready to unlock
Game->>Char : UnlockMagic(initialPOW)
Char->>Char : Set MagicUnlocked = true
Char->>Char : Initialize POW values
Char-->>Game : Success
Game-->>Player : Magic system activated
Player->>UI : Enable POW fields in editor
end
```

**Diagram sources**
- [internal/character/character.go](file://internal/character/character.go#L222-L231)

**Section sources**
- [internal/character/character.go](file://internal/character/character.go#L26-L33)
- [internal/character/character.go](file://internal/character/character.go#L222-L231)

## Equipment Management

The equipment system provides comprehensive item management with realistic damage reduction calculations and special item handling.

### Weapon System

Weapons provide damage bonuses and can have special properties that affect combat outcomes. The system supports both standard melee weapons and special items like the cursed Doombringer.

| Weapon | Damage Bonus | Special Properties | Description |
|--------|-------------|-------------------|-------------|
| Sword | 10 | None | Standard melee weapon |
| Dagger | 5 | None | Light, concealable |
| Axe | 15 | None | Standard melee weapon |
| Mace | 14 | None | Heavy melee weapon |
| Doombringer | 20 | Cursed | -10 LP per attack, heals on hit |

### Armor and Protection System

Armor provides damage reduction with sophisticated interaction mechanics. Shields offer additional protection that varies depending on whether armor is worn.

```mermaid
flowchart TD
Attack[Incoming Attack] --> CheckArmor{Armor Equipped?}
CheckArmor --> |Yes| ArmorProtection["Armor Protection: Armor.Protection"]
CheckArmor --> |No| NoArmor["No Armor Protection"]
ArmorProtection --> CheckShield{Shield Equipped?}
NoArmor --> CheckShield
CheckShield --> |Yes| ShieldCalc["Shield Protection:<br/>• With Armor: Shield.ProtectionWithArmor<br/>• Without Armor: Shield.Protection"]
CheckShield --> |No| NoShield["No Shield Protection"]
ShieldCalc --> TotalProtection["Total Protection = Armor + Shield"]
NoShield --> TotalProtection
TotalProtection --> ApplyReduction["Final Damage = Attack Damage - Total Protection"]
ApplyReduction --> DamageTaken[Damage Taken]
```

**Diagram sources**
- [internal/character/character.go](file://internal/character/character.go#L284-L302)
- [internal/items/items.go](file://internal/items/items.go#L183-L192)

### Equipment Interaction Rules

The system implements realistic equipment interaction rules:
- Shields provide full protection when worn alone
- When combined with armor, shield protection is reduced
- Different armor types offer varying protection levels
- Special weapons like Doombringer have unique combat mechanics

**Section sources**
- [internal/character/character.go](file://internal/character/character.go#L258-L277)
- [internal/items/items.go](file://internal/items/items.go#L20-L41)

## Character Creation Workflow

The character creation process follows a structured three-step workflow that mirrors the original gamebook experience while adding digital convenience.

### Creation Step 1: Characteristic Rolling

Players roll characteristics individually or simultaneously using the gamebook's standard mechanic. The system provides both manual and automatic rolling capabilities.

```mermaid
stateDiagram-v2
[*] --> RollCharacteristics
RollCharacteristics --> StrengthRoll : Roll Strength
StrengthRoll --> SpeedRoll : Roll Speed
SpeedRoll --> StaminaRoll : Roll Stamina
StaminaRoll --> CourageRoll : Roll Courage
CourageRoll --> LuckRoll : Roll Luck
LuckRoll --> CharmRoll : Roll Charm
CharmRoll --> AttractionRoll : Roll Attraction
AttractionRoll --> AllRolled : All characteristics rolled
AllRolled --> SelectEquipment : Proceed to equipment selection
SelectEquipment --> ReviewCharacter : Final review
ReviewCharacter --> Complete : Character creation complete
Complete --> [*]
```

### Creation Step 2: Equipment Selection

Players choose starting equipment from predefined options. The system provides cursor-based navigation and immediate preview of equipment effects.

### Creation Step 3: Character Review and Confirmation

The final step allows players to review all character attributes before committing to the creation process.

**Section sources**
- [pkg/ui/character_creation.go](file://pkg/ui/character_creation.go#L12-L20)
- [pkg/ui/character_creation.go](file://pkg/ui/character_creation.go#L229-L257)

## Character Editing System

The character editing system provides comprehensive stat modification capabilities with real-time validation and undo functionality.

### Edit Field Categories

The editing interface organizes character attributes into logical categories for intuitive navigation and modification.

| Category | Fields | Purpose |
|----------|--------|---------|
| Basic Stats | Strength, Speed, Stamina, Courage, Luck, Charm, Attraction | Core physical and mental attributes |
| Health | Current LP, Maximum LP | Life points management |
| Combat | Skill | Combat proficiency and tactics |
| Magic | Current POW, Maximum POW | Magical ability and resources |

### Input Mode Management

The editing system supports both direct value modification and text-based input modes, accommodating different user preferences and accessibility needs.

```mermaid
flowchart TD
EditStart[Start Editing] --> CursorNav["Navigate with Arrow Keys<br/>↑/↓: Move cursor<br/>←/→: Change field"]
CursorNav --> SelectField["Select Field with Enter"]
SelectField --> InputMode{Input Mode?}
InputMode --> |Text Input| TextInput["Type Numeric Values<br/>0-9: Add digit<br/>Backspace: Remove digit<br/>Enter: Confirm"]
InputMode --> |Direct Modification| DirectMod["Use +/- Keys<br/>Increase/Decrease Value"]
TextInput --> ValidateInput["Validate Input<br/>• Numeric only<br/>• Within bounds<br/>• No negative values"]
DirectMod --> ValidateInput
ValidateInput --> ApplyChange["Apply Change to Character"]
ApplyChange --> UpdateDisplay["Update UI Display"]
UpdateDisplay --> CursorNav
TextInput --> CancelInput["ESC: Cancel Input"]
DirectMod --> CancelInput
CancelInput --> CursorNav
```

**Diagram sources**
- [pkg/ui/character_edit.go](file://pkg/ui/character_edit.go#L110-L136)

**Section sources**
- [pkg/ui/character_edit.go](file://pkg/ui/character_edit.go#L8-L21)
- [pkg/ui/character_edit.go](file://pkg/ui/character_edit.go#L23-L56)

## Persistence and Serialization

The character system implements robust persistence using JSON serialization with timestamped filenames for version control and backup capabilities.

### Save File Structure

Character data is serialized to JSON format with comprehensive metadata for tracking creation and modification history.

```mermaid
erDiagram
CHARACTER {
int strength
int speed
int stamina
int courage
int luck
int charm
int attraction
int current_lp
int maximum_lp
int skill
int current_pow
int maximum_pow
bool magic_unlocked
Weapon equipped_weapon
Armor equipped_armor
bool has_shield
int enemies_defeated
datetime created_at
datetime last_saved
}
WEAPON {
string name
int damage_bonus
string description
bool special
}
ARMOR {
string name
int protection
string description
}
CHARACTER ||--|| WEAPON : "equips"
CHARACTER ||--|| ARMOR : "equips"
```

**Diagram sources**
- [internal/character/character.go](file://internal/character/character.go#L14-L44)

### Filename Convention

Save files use timestamped filenames in the format `character_YYYYMMDD-HHMMSS.json`, enabling automatic versioning and easy identification of save slots.

### Load/Save Workflow

The persistence system handles both saving and loading with comprehensive error handling and validation.

```mermaid
sequenceDiagram
participant UI as UI Layer
participant Char as Character
participant FS as File System
participant JSON as JSON Encoder
UI->>Char : Save(directory)
Char->>Char : Update LastSaved timestamp
Char->>Char : Format timestamp for filename
Char->>FS : Create directory if needed
FS-->>Char : Directory ready
Char->>JSON : MarshalIndent(character)
JSON-->>Char : JSON data
Char->>FS : Write to file
FS-->>Char : Save complete
Char-->>UI : Success/error
Note over UI,FS : Load Process
UI->>FS : Read file
FS-->>UI : File content
UI->>JSON : Unmarshal(data)
JSON-->>UI : Character object
UI->>UI : Validate character data
UI-->>UI : Load complete
```

**Diagram sources**
- [internal/character/character.go](file://internal/character/character.go#L312-L339)

**Section sources**
- [internal/character/character.go](file://internal/character/character.go#L312-L355)

## UI Integration Layer

The character system integrates seamlessly with the Bubble Tea UI framework through specialized model components that handle state management and user interaction.

### Model Architecture

The UI layer consists of three primary models that work together to provide comprehensive character management functionality.

```mermaid
classDiagram
class Model {
+Screen CurrentScreen
+Character Character
+Dice Roller
+MainMenuModel MainMenu
+CharacterCreationModel CharCreation
+LoadCharacterModel LoadChar
+GameSessionModel GameSession
+CharacterViewModel CharView
+CharacterEditModel CharEdit
+int Width
+int Height
+error Err
+LoadCharacter(char)
+SaveCharacter() error
}
class CharacterCreationModel {
+Roller dice
+CreationStep step
+int strength, speed, stamina, courage, luck, charm, attraction
+bool allRolled
+int weaponCursor, armorCursor
+Character character
+RollStrength() int
+CreateCharacter() Character
}
class CharacterEditModel {
+Character character
+Character originalChar
+int cursor
+[]string fields
+bool inputMode
+string inputBuffer
+SetCharacter(char)
+StartInput()
+GetCurrentValue() int
}
class CharacterViewModel {
+Character character
+SetCharacter(char)
+GetCharacter() Character
}
Model --> CharacterCreationModel
Model --> CharacterEditModel
Model --> CharacterViewModel
Model --> GameSessionModel
```

**Diagram sources**
- [pkg/ui/model.go](file://pkg/ui/model.go#L33-L50)
- [pkg/ui/character_creation.go](file://pkg/ui/character_creation.go#L21-L44)
- [pkg/ui/character_edit.go](file://pkg/ui/character_edit.go#L23-L31)
- [pkg/ui/character_view.go](file://pkg/ui/character_view.go#L6-L8)

### Screen Navigation

The UI system supports seamless navigation between character management screens with proper state synchronization.

**Section sources**
- [pkg/ui/model.go](file://pkg/ui/model.go#L33-L95)
- [pkg/ui/character_creation.go](file://pkg/ui/character_creation.go#L21-L44)
- [pkg/ui/character_edit.go](file://pkg/ui/character_edit.go#L23-L31)

## Validation and Business Rules

The character system implements comprehensive validation rules to maintain game balance and prevent invalid states.

### Characteristic Validation

All characteristic values must adhere to strict validation rules to ensure game balance and prevent exploitation.

| Rule | Constraint | Error Message |
|------|------------|---------------|
| Minimum Value | ≥ 0 | "{Characteristic} cannot be negative: {value}" |
| Maximum Value | ≤ 999 | "{Characteristic} exceeds maximum (999): {value}" |
| Negative Modifications | Not allowed for characteristics | "Cannot modify below zero" |
| LP Management | Can be negative (death condition) | N/A |

### Modification Constraints

Stat modifications follow specific rules to maintain game balance and prevent abuse of the system.

```mermaid
flowchart TD
ModifyAttempt[Stat Modification Attempt] --> CheckDirection{Positive or Negative?}
CheckDirection --> |Positive| AllowModification["Allow modification<br/>No restrictions"]
CheckDirection --> |Negative| CheckMinimum{Would result in negative?}
CheckMinimum --> |Yes| PreventModification["Prevent modification<br/>Return error"]
CheckMinimum --> |No| AllowModification
AllowModification --> UpdateCharacter["Update character stat"]
UpdateCharacter --> ValidateBounds["Validate new bounds"]
ValidateBounds --> Success["Modification successful"]
PreventModification --> ErrorReturn["Return error message"]
```

### Magic System Constraints

The magic system introduces additional validation requirements for power management and unlocking mechanics.

**Section sources**
- [internal/character/character.go](file://internal/character/character.go#L101-L111)
- [internal/character/character.go](file://internal/character/character.go#L223-L231)

## Common Issues and Edge Cases

The character system addresses several common issues and edge cases that can arise during gameplay and character management.

### Character Edit Screen Bug Fix

A significant UI bug was identified where all character edit fields displayed the same value as the currently selected field. This was caused by incorrect value retrieval logic in the display function.

**Root Cause Analysis:**
The bug occurred because the display function was retrieving the cursor position value for all fields instead of fetching each field's actual value. This resulted in all non-selected fields displaying the same incorrect value.

**Fix Implementation:**
The solution involved moving the value retrieval logic to execute for every field in the loop, ensuring each field displays its correct value regardless of cursor position.

### Equipment Management Edge Cases

Several edge cases exist in equipment management that require careful handling:

- **Shield-Armor Interaction:** Shield protection varies depending on whether armor is worn
- **Special Weapon Effects:** Weapons like Doombringer have unique combat mechanics
- **Equipment Conflicts:** Certain combinations may have unexpected effects
- **Resource Limits:** Power and life point limits must be respected

### Input Validation Issues

The character editing system must handle various input scenarios gracefully:

- **Invalid Characters:** Non-numeric input during text editing mode
- **Boundary Conditions:** Maximum and minimum value limits
- **Empty Input:** Handling of empty input buffers
- **Cancel Operations:** Proper restoration of previous values

**Section sources**
- [BUGFIX_CHARACTER_EDIT.md](file://BUGFIX_CHARACTER_EDIT.md#L1-L122)

## Practical Implementation Examples

The character system provides numerous practical examples of Go programming patterns and game development concepts.

### Dice Rolling Integration

The character creation system demonstrates proper integration between the dice rolling subsystem and character generation, showcasing dependency injection and modular design.

### JSON Serialization Patterns

Character persistence demonstrates idiomatic Go JSON handling with proper error management and file I/O operations.

### State Management Examples

The UI integration showcases effective state management patterns using the Bubble Tea framework's immutable state approach.

### Error Handling Strategies

The system demonstrates comprehensive error handling with meaningful error messages and graceful degradation.

### Equipment Calculation Examples

The armor protection calculation system shows practical implementation of game mechanics with proper mathematical formulas and edge case handling.

**Section sources**
- [internal/character/character.go](file://internal/character/character.go#L47-L96)
- [internal/dice/dice.go](file://internal/dice/dice.go#L60-L63)
- [internal/character/character.go](file://internal/character/character.go#L312-L355)