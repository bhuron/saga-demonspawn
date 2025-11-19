# Combat System

<cite>
**Referenced Files in This Document**
- [internal/combat/combat.go](file://internal/combat/combat.go)
- [pkg/ui/combat_view.go](file://pkg/ui/combat_view.go)
- [internal/dice/dice.go](file://internal/dice/dice.go)
- [internal/magic/spells.go](file://internal/magic/spells.go)
- [internal/magic/effects.go](file://internal/magic/effects.go)
- [internal/character/character.go](file://internal/character/character.go)
- [internal/items/items.go](file://internal/items/items.go)
- [pkg/ui/combat_setup.go](file://pkg/ui/combat_setup.go)
- [pkg/ui/game_session.go](file://pkg/ui/game_session.go)
- [internal/help/content/combat.txt](file://internal/help/content/combat.txt)
- [internal/help/content/global.txt](file://internal/help/content/global.txt)
- [saga_demonspawn_ruleset.md](file://saga_demonspawn_ruleset.md)
</cite>

## Table of Contents
1. [Introduction](#introduction)
2. [System Architecture](#system-architecture)
3. [Core Components](#core-components)
4. [Combat Flow and Mechanics](#combat-flow-and-mechanics)
5. [Character Statistics and Derived Values](#character-statistics-and-derived-values)
6. [Damage Calculation System](#damage-calculation-system)
7. [Initiative and Turn Order](#initiative-and-turn-order)
8. [Special Items Integration](#special-items-integration)
9. [Magic System Integration](#magic-system-integration)
10. [UI Implementation](#ui-implementation)
11. [Endurance and Rest System](#endurance-and-rest-system)
12. [Death Saves and Victory Conditions](#death-saves-and-victory-conditions)
13. [Testing and Validation](#testing-and-validation)
14. [Performance Considerations](#performance-considerations)
15. [Troubleshooting Guide](#troubleshooting-guide)

## Introduction

The Combat System in "Sagas of the Demonspawn" is a comprehensive turn-based battle engine that faithfully implements the classic gamebook's combat mechanics. Built with Go and the Bubble Tea framework, it provides automated calculations, intuitive UI controls, and seamless integration with the game's magic and item systems.

The system handles all aspects of combat including initiative determination, attack resolution, damage calculation, endurance tracking, death saves, and special item interactions. It maintains strict adherence to the original rules while adding modern conveniences like automatic calculations and persistent combat logs.

## System Architecture

The combat system follows a modular architecture with clear separation of concerns:

```mermaid
graph TB
subgraph "UI Layer"
CV[CombatViewModel]
CS[CombatSetupModel]
GS[GameSessionModel]
end
subgraph "Business Logic"
CE[CombatEngine]
DM[DamageCalculator]
IM[InitiativeManager]
ES[EnduranceSystem]
end
subgraph "Data Models"
C[Character]
E[Enemy]
CS_STATE[CombatState]
AR[AttackResult]
end
subgraph "Support Systems"
D[DiceRoller]
M[MagicSystem]
I[ItemSystem]
end
CV --> CE
CS --> E
CE --> DM
CE --> IM
CE --> ES
CE --> C
CE --> AR
DM --> D
IM --> D
ES --> C
CE --> M
CE --> I
```

**Diagram sources**
- [internal/combat/combat.go](file://internal/combat/combat.go#L1-L50)
- [pkg/ui/combat_view.go](file://pkg/ui/combat_view.go#L1-L50)
- [internal/dice/dice.go](file://internal/dice/dice.go#L1-L30)

**Section sources**
- [internal/combat/combat.go](file://internal/combat/combat.go#L1-L100)
- [pkg/ui/combat_view.go](file://pkg/ui/combat_view.go#L1-L100)

## Core Components

### Combat State Management

The `CombatState` struct serves as the central hub for tracking combat progress:

```mermaid
classDiagram
class CombatState {
+bool IsActive
+int CurrentRound
+bool PlayerTurn
+bool PlayerFirstStrike
+bool DeathSaveUsed
+int EnduranceLimit
+int RoundsSinceLastRest
+Enemy Enemy
+[]string CombatLog
+int PlayerInitiative
+int EnemyInitiative
+AddLogEntry(message string)
+NewCombatState(enemy, enduranceLimit) CombatState
}
class Enemy {
+string Name
+int Strength
+int Speed
+int Stamina
+int Courage
+int Luck
+int Skill
+int CurrentLP
+int MaximumLP
+int WeaponBonus
+int ArmorProtection
+bool IsDemonspawn
}
class AttackResult {
+int Roll
+int Requirement
+bool Hit
+int DamageBeforeArmor
+int FinalDamage
+int TargetLP
}
CombatState --> Enemy
CombatState --> AttackResult
```

**Diagram sources**
- [internal/combat/combat.go](file://internal/combat/combat.go#L79-L116)
- [internal/combat/combat.go](file://internal/combat/combat.go#L11-L26)

### Dice Rolling System

The dice system provides deterministic randomness for testing while maintaining game unpredictability:

```mermaid
classDiagram
class Roller {
<<interface>>
+Roll2D6() int
+Roll1D6() int
+RollCharacteristic() int
+SetSeed(seed int64)
}
class StandardRoller {
+rng *rand.Rand
+Roll2D6() int
+Roll1D6() int
+RollCharacteristic() int
+SetSeed(seed int64)
-rollDie(sides int) int
}
class RollResult {
+int Value
+string Description
+string Details
}
Roller <|.. StandardRoller
StandardRoller --> RollResult
```

**Diagram sources**
- [internal/dice/dice.go](file://internal/dice/dice.go#L11-L27)

**Section sources**
- [internal/combat/combat.go](file://internal/combat/combat.go#L79-L116)
- [internal/dice/dice.go](file://internal/dice/dice.go#L11-L97)

## Combat Flow and Mechanics

### Turn-Based Resolution

The combat system operates on a strict turn-based model with clear phase transitions:

```mermaid
flowchart TD
Start([Combat Start]) --> InitRoll["Roll Initiative<br/>2d6 + SPD + CRG + LCK"]
InitRoll --> DetermineOrder{"Player First?"}
DetermineOrder --> |Yes| PlayerTurn["Player Turn"]
DetermineOrder --> |No| EnemyTurn["Enemy Turn"]
PlayerTurn --> ActionSelect["Select Action<br/>Attack/Cast/Use/Flee"]
ActionSelect --> ExecuteAction["Execute Action"]
ExecuteAction --> CheckVictory{"Victory?"}
CheckVictory --> |Yes| Victory["Combat Victory"]
CheckVictory --> |No| CheckDefeat{"Defeat?"}
CheckDefeat --> |Yes| DeathSave["Attempt Death Save"]
CheckDefeat --> |No| CheckRest{"Rest Needed?"}
EnemyTurn --> EnemyAction["Enemy Automatic Attack"]
EnemyAction --> CheckVictory
CheckRest --> |Yes| RestPhase["Rest Phase<br/>Free Enemy Attack"]
CheckRest --> |No| NextTurn["Advance Turn"]
RestPhase --> NextTurn
DeathSave --> DeathSaveResult{"Save Success?"}
DeathSaveResult --> |Yes| RestartCombat["Restart Combat"]
DeathSaveResult --> |No| Defeat["Combat Defeat"]
NextTurn --> CheckRoundEnd{"Round End?"}
CheckRoundEnd --> |Yes| NewRound["New Round"]
CheckRoundEnd --> |No| PlayerTurn
RestartCombat --> PlayerTurn
NewRound --> PlayerTurn
Victory --> End([Combat End])
Defeat --> End
```

**Diagram sources**
- [internal/combat/combat.go](file://internal/combat/combat.go#L297-L383)
- [pkg/ui/combat_view.go](file://pkg/ui/combat_view.go#L78-L158)

### Action Types and Implementation

The combat system supports four primary actions with distinct mechanics:

| Action | Description | Combat Effects | Special Rules |
|--------|-------------|----------------|---------------|
| **Attack** | Standard melee attack | Damage calculation, potential blood price | Doombringer curse, The Orb doubling |
| **Cast Spell** | Magic spell casting | Spell effects, POW cost deduction | Natural inclination, FFR check |
| **Use Item** | Special item activation | Healing, equipment changes | Charge consumption |
| **Flee** | Combat escape attempt | Immediate exit | No victory condition |

**Section sources**
- [pkg/ui/combat_view.go](file://pkg/ui/combat_view.go#L161-L385)
- [internal/help/content/combat.txt](file://internal/help/content/combat.txt#L8-L25)

## Character Statistics and Derived Values

### Core Characteristics

The game defines seven core characteristics that form the foundation of combat abilities:

```mermaid
graph LR
subgraph "Core Characteristics"
STR[Strength<br/>Physical Power]
SPD[Speed<br/>Agility & Reaction]
STA[Stamina<br/>Endurance]
CRG[Courage<br/>Bravery]
LCK[Luck<br/>Fortune]
CHM[Charm<br/>Social Skills]
ATT[Attraction<br/>Personal Magnetism]
end
subgraph "Derived Values"
LP[Life Points<br/>Sum of all characteristics]
SKL[Skill<br/>Combat Proficiency]
POW[Power<br/>Magic Resource]
end
STR --> LP
SPD --> LP
STA --> LP
CRG --> LP
LCK --> LP
CHM --> LP
ATT --> LP
LP --> SKL
SKL --> POW
```

**Diagram sources**
- [internal/character/character.go](file://internal/character/character.go#L14-L35)
- [saga_demonspawn_ruleset.md](file://saga_demonspawn_ruleset.md#L5-L16)

### Stat Ranges and Impact

Each characteristic influences combat mechanics differently:

| Characteristic | Range | Combat Impact | Formula Application |
|----------------|-------|---------------|-------------------|
| **STR** | 16-96 | Damage bonus (STR ÷ 10 × 5) | Base damage calculation |
| **SPD** | 16-96 | Initiative modifier | +SPD to initiative rolls |
| **STA** | 16-96 | Endurance limit (STA ÷ 20) | Rest requirement calculation |
| **CRG** | 16-96 | Initiative modifier | +CRG to initiative rolls |
| **LCK** | 16-96 | To-hit modifier (-1 per 16 LCK) | Hit probability adjustment |
| **SKL** | 0-∞ | To-hit modifier (-1 per 10 SKL) | Experience-based hit improvement |

**Section sources**
- [internal/character/character.go](file://internal/character/character.go#L14-L35)
- [saga_demonspawn_ruleset.md](file://saga_demonspawn_ruleset.md#L158-L169)

## Damage Calculation System

### Player Damage Formula

The damage calculation follows the original gamebook mechanics with mathematical precision:

```mermaid
flowchart TD
Roll["Roll 2d6 for Attack"] --> BaseCalc["Base Damage = Roll × 5"]
BaseCalc --> StrBonus["STR Bonus = (STR ÷ 10) × 5"]
StrBonus --> WeaponBonus["Weapon Bonus"]
WeaponBonus --> TotalDamage["Total Damage = Base + STR + Weapon"]
TotalDamage --> ArmorReduction["Armor Reduction = Enemy Armor"]
ArmorReduction --> FinalDamage["Final Damage = Total - Armor"]
FinalDamage --> ApplyDamage["Apply to Enemy LP"]
Roll --> CheckHit{"Roll ≥ Requirement?"}
CheckHit --> |No| Miss["Miss - No Damage"]
CheckHit --> |Yes| TotalDamage
```

**Diagram sources**
- [internal/combat/combat.go](file://internal/combat/combat.go#L157-L174)

### Enemy Damage Calculation

Enemy attacks follow identical mechanics with character protection:

```mermaid
flowchart TD
EnemyRoll["Enemy Rolls 2d6"] --> EnemyBase["Enemy Base Damage"]
EnemyBase --> Xenophobia["Apply XENOPHOBIA Reduction"]
Xenophobia --> PlayerArmor["Calculate Player Armor"]
PlayerArmor --> ShieldCheck{"Has Shield?"}
ShieldCheck --> |Yes| ShieldCalc["Shield + Armor Protection"]
ShieldCheck --> |No| ArmorCalc["Armor Only Protection"]
ShieldCalc --> FinalEnemyDamage["Final Enemy Damage"]
ArmorCalc --> FinalEnemyDamage
FinalEnemyDamage --> ApplyPlayerDamage["Apply to Player LP"]
```

**Diagram sources**
- [internal/combat/combat.go](file://internal/combat/combat.go#L237-L295)

### Special Item Damage Modifications

Certain items modify damage calculations dynamically:

| Item | Effect | Calculation Impact |
|------|--------|-------------------|
| **Doombringer** | Blood price + soul thirst | -10 LP before attack, heal for damage dealt |
| **The Orb** | Double damage to Demonspawn | Damage × 2 when equipped and targeting Demonspawn |
| **ARMOUR Spell** | Damage reduction | +10 damage reduction to player |
| **XENOPHOBIA Spell** | Enemy damage reduction | -5 damage reduction from enemy |

**Section sources**
- [internal/combat/combat.go](file://internal/combat/combat.go#L157-L295)
- [internal/items/items.go](file://internal/items/items.go#L143-L149)

## Initiative and Turn Order

### Initiative Calculation

Initiative determines the first-strike advantage and turn order:

```mermaid
sequenceDiagram
participant Player as Player
participant Dice as Dice Roller
participant Enemy as Enemy
participant Combat as Combat Engine
Combat->>Dice : Roll2D6() for player
Dice-->>Combat : Player roll result
Combat->>Combat : Calculate player initiative
Note over Combat : Initiative = roll + SPD + CRG + LCK
Combat->>Dice : Roll2D6() for enemy
Dice-->>Combat : Enemy roll result
Combat->>Combat : Calculate enemy initiative
Note over Combat : Initiative = roll + SPD + CRG + LCK
Combat->>Combat : Compare initiatives
alt Player initiative > Enemy initiative
Combat->>Combat : Player goes first
Combat->>Combat : Set PlayerFirstStrike = true
else Enemy initiative > Player initiative
Combat->>Combat : Enemy goes first
Combat->>Combat : Set PlayerFirstStrike = false
end
```

**Diagram sources**
- [internal/combat/combat.go](file://internal/combat/combat.go#L123-L133)

### Turn Alternation Logic

The system maintains strict turn alternation with round progression:

```mermaid
stateDiagram-v2
[*] --> PlayerTurn
PlayerTurn --> EnemyTurn : NextTurn()
EnemyTurn --> PlayerTurn : NextTurn()
PlayerTurn --> NewRound : Round End?
EnemyTurn --> NewRound : Round End?
NewRound --> PlayerTurn : Reset turn tracker
note right of NewRound
Round counter increments
Endurance counters reset
Rest requirements checked
end note
```

**Diagram sources**
- [internal/combat/combat.go](file://internal/combat/combat.go#L325-L336)

**Section sources**
- [internal/combat/combat.go](file://internal/combat/combat.go#L123-L133)
- [internal/combat/combat.go](file://internal/combat/combat.go#L325-L336)

## Special Items Integration

### Healing Stone System

The Healing Stone provides emergency healing with charge management:

```mermaid
flowchart TD
UseStone["Use Healing Stone"] --> CheckCharges{"Charges > 0?"}
CheckCharges --> |No| StoneEmpty["Stone Depleted"]
CheckCharges --> |Yes| CheckHP{"Current HP < Max HP?"}
CheckHP --> |No| FullHealth["Already at Full Health"]
CheckHP --> |Yes| RollHeal["Roll 1d6 × 10"]
RollHeal --> ApplyHeal["Apply Healing"]
ApplyHeal --> UpdateCharges["Reduce Charges"]
UpdateCharges --> CheckRemaining{"Charges Remaining?"}
CheckRemaining --> |Yes| MoreActions["Continue Actions"]
CheckRemaining --> |No| RemoveOption["Remove from Action List"]
```

**Diagram sources**
- [pkg/ui/combat_view.go](file://pkg/ui/combat_view.go#L275-L322)

### Doombringer Weapon Mechanics

Doombringer introduces risk-reward dynamics with blood price and soul thirst:

```mermaid
flowchart TD
AttackWithDoombringer["Attack with Doombringer"] --> BloodPrice["Pay Blood Price<br/>-10 LP"]
BloodPrice --> CheckHP{"HP > 0?"}
CheckHP --> |No| DoombringerDeath["Doombringer Death"]
CheckHP --> |Yes| ExecuteAttack["Execute Attack"]
ExecuteAttack --> CheckHit{"Attack Hit?"}
CheckHit --> |No| NoHealing["No Healing"]
CheckHit --> |Yes| CalculateHealing["Healing = Damage Dealt"]
CalculateHealing --> CapHealing["Cap at Max LP"]
CapHealing --> ApplyHealing["Apply Healing"]
```

**Diagram sources**
- [pkg/ui/combat_view.go](file://pkg/ui/combat_view.go#L184-L261)

### The Orb Special Abilities

The Orb provides dual functionality with strategic depth:

| Ability | Condition | Effect | Outcome |
|---------|-----------|--------|---------|
| **Held** | Equipped in left hand | Double damage to Demonspawn | Damage × 2 after all reductions |
| **Thrown** | Combat situation | Demons: instant kill or 200 damage | Instant kill (4+ on 2d6) or 200 damage |
| **Thrown** | Non-Demonspawn | No effect | Orb destroyed, no damage |

**Section sources**
- [pkg/ui/combat_view.go](file://pkg/ui/combat_view.go#L275-L382)
- [pkg/ui/combat_view.go](file://pkg/ui/combat_view.go#L324-L378)
- [saga_demonspawn_ruleset.md](file://saga_demonspawn_ruleset.md#L77-L86)

## Magic System Integration

### Spell Effects in Combat

Magic spells integrate seamlessly with combat mechanics:

```mermaid
classDiagram
class SpellEffect {
+bool Success
+string Message
+int DamageDealt
+int LPRestored
+bool CombatEnded
+bool Victory
+bool EnemyKilled
+bool CharacterDied
+bool RequiresReroll
+string NavigateTo
}
class CombatIntegration {
+ApplyARMOUR() SpellEffect
+ApplyFIREBALL() SpellEffect
+ApplyXENOPHOBIA() SpellEffect
+ApplyPOISON_NEEDLE() SpellEffect
}
CombatIntegration --> SpellEffect
```

**Diagram sources**
- [internal/magic/effects.go](file://internal/magic/effects.go#L1-L47)

### Combat-Specific Spells

Several spells have specialized combat applications:

| Spell | Combat Effect | Integration Point | Duration |
|-------|---------------|-------------------|----------|
| **ARMOUR** | -10 incoming damage | Enemy attack calculation | Until section end |
| **FIREBALL** | 50 LP damage | Player attack resolution | Single use |
| **XENOPHOBIA** | -5 enemy damage | Enemy attack calculation | Until section end |
| **POISON NEEDLE** | 50% instant kill | Enemy attack resolution | Single use |
| **PARALYSIS** | Enemy immobilization | Enemy turn processing | Single use |

### Magic System Coordination

The magic system coordinates with combat through spell effect tracking:

```mermaid
sequenceDiagram
participant Player as Player
participant Magic as Magic System
participant Combat as Combat Engine
participant Character as Character
Player->>Magic : Cast spell
Magic->>Magic : Validate spell availability
Magic->>Magic : Check POW cost
Magic->>Magic : Perform FFR check
Magic->>Combat : Apply spell effect
Combat->>Character : Update spell effects
Character->>Combat : Store effect duration
Combat->>Combat : Integrate effects into damage calculation
```

**Diagram sources**
- [internal/magic/effects.go](file://internal/magic/effects.go#L1-L47)
- [internal/combat/combat.go](file://internal/combat/combat.go#L256-L282)

**Section sources**
- [internal/magic/spells.go](file://internal/magic/spells.go#L1-L137)
- [internal/magic/effects.go](file://internal/magic/effects.go#L1-L47)
- [internal/combat/combat.go](file://internal/combat/combat.go#L256-L282)

## UI Implementation

### Combat View Model Architecture

The UI system uses Bubble Tea's MVU pattern for responsive combat interface:

```mermaid
classDiagram
class CombatViewModel {
+Character player
+CombatState combatState
+Roller roller
+int selectedAction
+bool waitingForInput
+bool victoryState
+bool defeatState
+bool needsRest
+[]string actions
+Update(msg) CombatViewModel
+View() string
+handleAction() CombatViewModel
+processEnemyTurn() CombatViewModel
+checkCombatState() CombatViewModel
}
class ActionMenu {
+actionAttack
+actionCastSpell
+actionFlee
+actionUseHealingStone
+actionThrowOrb
}
CombatViewModel --> ActionMenu
```

**Diagram sources**
- [pkg/ui/combat_view.go](file://pkg/ui/combat_view.go#L13-L75)

### Combat Interface Layout

The combat interface displays comprehensive information in organized sections:

```mermaid
graph TB
subgraph "Combat Interface"
Header["Combat Header<br/>Round Number"]
Stats["Combatant Stats<br/>LP, STR, SPD, Weapon, Armor"]
Endurance["Endurance Status<br/>Rounds remaining"]
Log["Combat Log<br/>Recent events"]
Actions["Action Menu<br/>Attack/Cast/Flee/Use"]
end
Header --> Stats
Stats --> Endurance
Endurance --> Log
Log --> Actions
```

**Diagram sources**
- [pkg/ui/combat_view.go](file://pkg/ui/combat_view.go#L527-L627)

### Input Handling and Navigation

The UI supports keyboard navigation with clear visual feedback:

| Key | Action | Context | Effect |
|-----|--------|---------|--------|
| **↑/↓** | Navigate actions | During player turn | Change selected action |
| **Enter** | Confirm action | During player turn | Execute selected action |
| **Esc** | Back/Quit | Anytime | Return to previous screen |
| **?** | Help | Anytime | Show context-sensitive help |

**Section sources**
- [pkg/ui/combat_view.go](file://pkg/ui/combat_view.go#L78-L158)
- [pkg/ui/combat_view.go](file://pkg/ui/combat_view.go#L527-L627)

## Endurance and Rest System

### Stamina-Based Endurance

The endurance system tracks combat fatigue and rest requirements:

```mermaid
flowchart TD
StartCombat["Start Combat"] --> CalcEndurance["Calculate Endurance Limit<br/>STA ÷ 20"]
CalcEndurance --> TrackRounds["Track Rounds Since Last Rest"]
TrackRounds --> CheckRest{"Rounds ≥ Endurance Limit?"}
CheckRest --> |No| ContinueCombat["Continue Combat"]
CheckRest --> |Yes| RestRequired["Rest Required"]
RestRequired --> FreeAttack["Enemy Gets Free Attack"]
FreeAttack --> ResetCounter["Reset Round Counter"]
ResetCounter --> ContinueCombat
ContinueCombat --> CheckVictory{"Victory/Defeat?"}
CheckVictory --> |No| TrackRounds
CheckVictory --> |Yes| EndCombat["End Combat"]
```

**Diagram sources**
- [internal/combat/combat.go](file://internal/combat/combat.go#L176-L186)

### Rest Mechanics

Rest provides temporary relief from combat fatigue:

| Scenario | Trigger | Effect | Duration |
|----------|---------|--------|----------|
| **Player Rest** | Stamina depletion | Free enemy attack | Single round |
| **Enemy Rest** | Enemy stamina depletion | Player free attack | Single round |
| **Combat End** | Victory/Defeat | Reset all counters | Permanent |

**Section sources**
- [internal/combat/combat.go](file://internal/combat/combat.go#L176-L186)
- [internal/combat/combat.go](file://internal/combat/combat.go#L338-L346)

## Death Saves and Victory Conditions

### Death Save Mechanics

The death save system provides a single chance at revival:

```mermaid
flowchart TD
ZeroLP["LP Drops to 0"] --> CheckDeathSave{"Death Save Available?"}
CheckDeathSave --> |No| Defeat["Combat Defeat"]
CheckDeathSave --> |Yes| ShowPrompt["Show Death Save Prompt"]
ShowPrompt --> PlayerRoll["Roll 2d6 × 10"]
PlayerRoll --> CompareLuck{"Roll ≤ Luck?"}
CompareLuck --> |Yes| Success["Death Save Success"]
CompareLuck --> |No| Failure["Death Save Failure"]
Success --> RestoreHP["Restore to Max LP"]
RestoreHP --> RestartCombat["Restart Combat"]
RestartCombat --> NewInitiative["Re-roll Initiative"]
Failure --> Defeat
```

**Diagram sources**
- [internal/combat/combat.go](file://internal/combat/combat.go#L354-L382)

### Victory and Defeat Conditions

Combat termination occurs under specific conditions:

```mermaid
stateDiagram-v2
[*] --> InProgress
InProgress --> Victory : Enemy LP ≤ 0
InProgress --> Defeat : Player LP ≤ 0
InProgress --> DeathSave : Player LP ≤ 0
DeathSave --> Victory : Save Success
DeathSave --> Defeat : Save Failure
Victory --> [*]
Defeat --> [*]
```

**Diagram sources**
- [internal/combat/combat.go](file://internal/combat/combat.go#L315-L323)

### Combat Resolution

Victory conditions trigger character progression:

| Condition | Outcome | Effects |
|-----------|---------|---------|
| **Enemy Defeated** | Victory achieved | Skill +1, Enemies defeated +1 |
| **Flee Combat** | Partial victory | No stat changes, immediate exit |
| **Death Save Success** | Resurrection | Full LP restored, combat restart |
| **Death Save Failure** | Defeat | Game over, character lost |

**Section sources**
- [internal/combat/combat.go](file://internal/combat/combat.go#L315-L382)
- [internal/combat/combat.go](file://internal/combat/combat.go#L348-L352)

## Testing and Validation

### Unit Testing Strategy

The combat system employs comprehensive testing across all components:

```mermaid
graph TB
subgraph "Test Categories"
CT[Combat Tests]
DT[Damage Tests]
IT[Initiative Tests]
ET[Endurance Tests]
ST[Spell Integration Tests]
end
subgraph "Test Tools"
TR[Test Runner]
MR[Mock Roller]
TF[Test Fixtures]
end
CT --> TR
DT --> TR
IT --> TR
ET --> TR
ST --> TR
TR --> MR
TR --> TF
```

### Key Test Coverage Areas

| Component | Test Focus | Validation Criteria |
|-----------|------------|-------------------|
| **Damage Calculation** | Formula accuracy | Correct mathematical results |
| **Initiative Rolls** | Random distribution | Statistical fairness |
| **Endurance System** | Rest triggers | Proper timing and effects |
| **Death Saves** | Probability mechanics | Luck-based outcomes |
| **Spell Integration** | Effect application | Correct buff/debuff implementation |

**Section sources**
- [internal/combat/combat_test.go](file://internal/combat/combat_test.go)

## Performance Considerations

### Computational Efficiency

The combat system optimizes performance through several strategies:

- **Pre-calculated Constants**: Initiative and damage formulas are pre-computed where possible
- **Minimal Allocations**: String concatenation is minimized in hot paths
- **Efficient Data Structures**: Combat state uses compact, cache-friendly layouts
- **Lazy Evaluation**: Spell effects are only calculated when needed

### Memory Management

Memory usage is optimized through:
- **Object Pooling**: Reuse of frequently allocated objects
- **String Interning**: Shared string constants for UI text
- **Compact Serialization**: JSON serialization minimizes data size

### Scalability Factors

The system scales efficiently with:
- **Linear Complexity**: Most operations scale linearly with combat participants
- **Deterministic Execution**: Predictable performance across different scenarios
- **Modular Design**: Easy extension without performance degradation

## Troubleshooting Guide

### Common Combat Issues

| Problem | Symptoms | Cause | Solution |
|---------|----------|-------|----------|
| **Incorrect Damage** | Unexpected damage amounts | Math formula error | Verify damage calculation logic |
| **Initiative Problems** | Wrong turn order | Initiative calculation bug | Check SPD/CRG/LCK addition |
| **Death Save Failures** | Too many/little deaths | Luck comparison error | Validate 2d6×10 vs Luck logic |
| **Endurance Issues** | Rest not triggering | Counter logic error | Review round counting |
| **Spell Effects** | Spells not applying | Effect tracking bug | Check spell effect storage |

### Debugging Combat Scenarios

For combat debugging, examine these key areas:

1. **Combat State**: Verify all state fields contain expected values
2. **Roll Results**: Check dice roll outcomes match expectations
3. **Calculation Steps**: Trace damage and initiative calculations
4. **Effect Application**: Confirm spell effects are properly applied
5. **Turn Logic**: Validate turn alternation and round progression

### Performance Issues

Common performance bottlenecks and solutions:

- **Excessive Logging**: Reduce combat log verbosity in production
- **String Operations**: Minimize string concatenation in loops
- **Memory Allocation**: Use object pools for frequently allocated objects
- **Algorithm Complexity**: Review nested loops in combat calculations

**Section sources**
- [internal/combat/combat.go](file://internal/combat/combat.go#L1-L383)
- [pkg/ui/combat_view.go](file://pkg/ui/combat_view.go#L1-L653)