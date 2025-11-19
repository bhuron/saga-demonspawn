# Phase 4: Magic System Design

## Overview

Phase 4 implements the complete magic system for Sagas of the Demonspawn, enabling Fire*Wolf to cast spells using POWER (POW) as a resource. The system follows the four-step spell casting process defined in the ruleset, including natural inclination checks, power cost management, fundamental failure rates, and ten unique spell effects.

## Design Goals

- Implement all ten spells with accurate rule mechanics
- Enforce the four-step casting process with proper validation
- Manage POWER as a strategic resource with multiple renewal paths
- Integrate spells into combat, inventory, and navigation systems
- Provide clear UI feedback for spell outcomes and power state
- Track spell usage to prevent duplicate casting per section

## System Architecture

### Component Structure

The magic system comprises three main components working together:

| Component | Responsibility | Location |
|-----------|----------------|----------|
| **Magic Engine** | Core spell casting logic, validation, and effect resolution | `internal/magic/` |
| **Character Integration** | POWER tracking, magic unlock state, spell history per section | `internal/character/` |
| **UI Layer** | Spell selection, casting interface, feedback display | `pkg/ui/` |

### Data Flow

The magic system follows this decision flow:

1. Player selects "Cast Spell" from game session menu
2. UI displays available spells with power costs
3. Player selects spell to cast
4. Magic engine validates preconditions and executes four-step process
5. Results propagate back to character state and UI display

## Core Mechanics

### Four-Step Casting Process

Each spell cast must pass through four distinct steps in sequence:

#### Step 1: Natural Inclination Check

Fire*Wolf abhors sorcery and must overcome his reluctance before using magic in a section.

**Behavior:**
- Player initiates check when gamebook requires it (before first spell in a section)
- Roll 2d6; success requires 4 or higher
- If failed, player should refrain from using magic for the section (honor system)
- If successful, player may cast spells according to gamebook rules

**Implementation:**
- Provide "Natural Inclination Check" as a separate action in the spell casting screen
- Display roll result and success/failure message
- No enforcement of check requirement or result
- Trust player to follow gamebook instructions

#### Step 2: Power Cost Payment

Each spell has a fixed POWER cost that must be paid upfront.

**Behavior:**
- Check if character has sufficient CurrentPOW for spell cost
- If insufficient, offer option to sacrifice Life Points for POWER at 1:1 ratio
- Deduct cost from CurrentPOW after confirmation
- Cost is permanently spent regardless of spell success or failure

**Validation:**
- Cannot cast if CurrentPOW < spell cost and player declines LP sacrifice
- LP sacrifice cannot reduce CurrentLP below 1 (character must survive)
- POWER cannot exceed MaximumPOW after restoration

#### Step 3: Fundamental Failure Rate

All spells have an inherent chance of failure even after cost payment.

**Behavior:**
- Roll 2d6; success requires 6 or higher
- Approximately 72% success rate
- If failed, spell has no effect but cost remains spent
- If successful, proceed to spell effect execution

#### Step 4: Spell Effect Application

Each spell applies its unique effect to the game state.

**Effects by Spell:**
| Spell | Effect Type | Impact |
|-------|-------------|--------|
| ARMOUR | Defensive buff | Reduce incoming damage by 10 points for current section |
| CRYPT | Navigation | Return to Crypts for POWER restoration tests |
| FIREBALL | Direct damage | Deal 50 LP damage to target enemy |
| INVISIBILITY | Combat avoidance | Skip combat and proceed as if victorious |
| PARALYSIS | Combat escape | Exit combat immediately, move to next section |
| POISON NEEDLE | Conditional instant kill | Roll 1d6: 4-6 kills target, 1-3 target immune |
| RESURRECTION | Death recovery | Only usable when killed; return to section start, reroll all stats |
| RETRACE | Navigation | Return to any previously visited section |
| TIMEWARP | Combat reset | Reset current section; restore all LP to section starting values |
| XENOPHOBIA | Enemy debuff | Reduce enemy damage by 5 points |

### Spell Catalog

All ten spells from the ruleset must be implemented with exact costs and effects:

| Spell Name | POWER Cost | Category | Usage Context |
|------------|------------|----------|---------------|
| ARMOUR | 25 | Defensive | Pre-combat or during combat |
| CRYPT | 150 | Strategic | Outside combat |
| FIREBALL | 15 | Offensive | During combat |
| INVISIBILITY | 30 | Tactical | Pre-combat or during combat |
| PARALYSIS | 30 | Escape | During combat |
| POISON NEEDLE | 25 | Offensive | During combat |
| RESURRECTION | 50 | Recovery | Only when CurrentLP ≤ 0 |
| RETRACE | 20 | Navigation | Outside combat |
| TIMEWARP | 10 | Combat reset | During combat |
| XENOPHOBIA | 15 | Debuff | During combat |

### POWER Management

POWER serves as the magic system's resource economy with multiple renewal paths:

#### Restoration Methods

| Method | Mechanism | Constraints |
|--------|-----------|-------------|
| **Exploration** | Gain 1 POW when entering new section | Cannot exceed MaximumPOW |
| **Sacrifice** | Trade LP for POW at 1:1 ratio | Cannot reduce CurrentLP below 1 |
| **CRYPT Spell** | Return to Crypts for restoration tests | Costs 150 POW; tests may increase MaximumPOW |

#### Strategic Considerations

POWER management creates meaningful decisions:

- **Conservative play**: Save POW for critical moments, rely on exploration gains
- **Aggressive play**: Sacrifice LP for immediate POW, use spells liberally
- **CRYPT gamble**: Spend 150 POW for chance to increase maximum capacity

### Spell Usage Restrictions

The magic system enforces minimal restrictions, trusting the player to follow gamebook rules:

| Rule | Enforcement |
|------|-------------|
| **RESURRECTION Exclusivity** | Only available when CurrentLP ≤ 0 |
| **Combat Context** | Some spells only castable during combat |
| **Non-Combat Context** | Some spells only castable outside combat |

**Player Responsibility:**
- Players manage their own spell usage per the gamebook's "once per section" rule
- No technical enforcement of duplicate casting within sections
- Natural inclination checks are player-managed (honor system)
- POWER restoration from exploration (+1 per section) is manually applied via Character Edit

## Magic System Unlock Mechanism

### How Players Unlock Magic

The magic system must be manually unlocked during gameplay, as Fire*Wolf discovers or is forced to use sorcery during the adventure. The unlock mechanism provides flexibility for both manual player control and gamebook-directed activation.

### Unlock Entry Points

#### Option 1: Character Edit Screen (Recommended for Phase 4)

**Location:** Character Edit screen (accessible via Game Session → Edit Character Stats)

**Behavior:**
- Add a special action at the bottom of the character edit field list
- When magic is locked: Display "[U] Unlock Magic System" option
- When pressed, prompt: "Unlock magic? Enter initial POW value:"
- Accept numeric input for initial POW (e.g., 50, 100)
- Call `Character.UnlockMagic(initialPOW)` to activate
- Display confirmation: "Magic system unlocked with X POW"
- Update Game Session menu to show "Cast Spell" option

**User Flow:**
1. Player reads gamebook instruction: "You gain 50 POWER and can now use magic"
2. Navigate to Edit Character Stats
3. Press 'U' to unlock magic
4. Enter 50 as initial POW
5. Confirm unlock
6. Return to Game Session, see "Cast Spell" now available

**Advantages:**
- Consistent with other manual stat modifications
- Clear, discoverable location
- Allows gamebook to specify initial POW amount
- Single, focused responsibility

#### Option 2: Inventory Acquisition Pattern

**Location:** Inventory Management screen

**Behavior:**
- Add "Magic Powers" as a special pseudo-item in inventory list
- Appears as "[Not Acquired] Magic Powers" when locked
- Press 'A' to acquire, prompts for initial POW
- Changes to "[Possessed] Magic Powers (X POW)" when unlocked

**Advantages:**
- Consistent with Healing Stone/Doombringer/Orb acquisition
- Treats magic as a "special ability" like special items

**Disadvantages:**
- Less intuitive (magic is not an item)
- Inventory screen becomes overloaded with non-inventory concerns

#### Option 3: Dedicated Unlock Dialog from Game Session

**Location:** New "Unlock Magic" option in Game Session menu

**Behavior:**
- Show "Unlock Magic" option in Game Session menu when magic is locked
- Replace with "Cast Spell" once unlocked
- Selecting "Unlock Magic" shows dialog with POW input
- Cannot be un-unlocked once activated

**Advantages:**
- Very discoverable and prominent
- Clearly separates "unlock" from "cast"

**Disadvantages:**
- Menu item only used once, then disappears
- Adds complexity to Game Session menu state

### Recommended Approach for Phase 4

**Use Option 1: Character Edit Screen unlock**

This approach best fits the existing architecture because:

1. **Consistency**: Character Edit already handles manual stat modifications per gamebook instructions
2. **Simplicity**: Single location for all character state changes
3. **Discoverability**: Players already use Edit Character for gamebook-directed changes
4. **Minimal UI Changes**: Only affects one existing screen
5. **Persistence**: Once unlocked, state persists naturally through save/load

### Implementation Details

**Character Edit Screen Changes:**

| Element | Before Magic Unlock | After Magic Unlock |
|---------|---------------------|--------------------|
| Field List | Ends with "Attraction" | Ends with "Maximum POW" |
| Special Action | "[U] Unlock Magic System" at bottom | Hidden (no longer needed) |
| POW Fields | Hidden/greyed out | Fully editable |
| Help Text | "Press U to unlock magic when gamebook allows" | Standard edit help |

**Unlock Dialog Flow:**

1. User presses 'U' on Character Edit screen
2. Show modal prompt: 
   ```
   Unlock Magic System?
   
   Enter initial POWER value: [____]
   
   (Enter to confirm, Esc to cancel)
   ```
3. Accept numeric input (1-999)
4. Validate: POW > 0
5. Call `Character.UnlockMagic(enteredPOW)`
6. Show confirmation: "Magic unlocked! You now have X POW."
7. Return to Character Edit, POW fields now visible and editable
8. Game Session menu automatically shows "Cast Spell" on next view

**Edge Cases:**

| Scenario | Handling |
|----------|----------|
| User presses 'U' when already unlocked | Ignore keypress or show "Magic already unlocked" |
| User enters 0 or negative POW | Show error: "Initial POW must be positive" |
| User enters non-numeric input | Show error: "Please enter a number" |
| User loads saved character with magic unlocked | No unlock needed, "Cast Spell" automatically visible |
| User unlocks, then saves/loads | Magic stays unlocked with saved POW values |

### Alternative: Auto-Unlock Based on POW Value

For simpler implementation, magic could auto-unlock when POW is edited to a positive value:

**Simplified Behavior:**
- In Character Edit, POW fields are always visible but greyed out when magic locked
- When user edits "Maximum POW" or "Current POW" to value > 0, automatically call `UnlockMagic()`
- No explicit "unlock" action required

**Trade-offs:**
- Simpler UI (no special 'U' action needed)
- More intuitive (editing POW = having magic)
- Less explicit (users may not realize they unlocked magic)
- Risk of accidental unlock by editing wrong field

**Recommendation:** Use explicit unlock action for Phase 4 to maintain clarity and intentionality.

## UI Integration

### Spell Casting Screen

A new UI component displays spell selection and casting flow:

**Screen Elements:**
- List of available spells with names, costs, and descriptions
- Current POWER display (CurrentPOW / MaximumPOW)
- Current section identifier
- Spell status indicators (available, already cast, insufficient POW)
- Casting result feedback area

**Navigation:**
- Arrow keys to select spell
- Enter to cast selected spell
- Esc/q to return to game session menu
- Display confirmation prompt for LP sacrifice option

### Spell Feedback

The UI provides clear feedback for each casting step:

| Step | Success Feedback | Failure Feedback |
|------|------------------|------------------|
| Natural Inclination | "Fire*Wolf overcomes his aversion to magic" | "Fire*Wolf refuses to use sorcery" |
| Power Cost | "Spent X POW (Remaining: Y)" | "Insufficient POWER (Need X, Have Y)" |
| LP Sacrifice Offer | "Sacrifice X LP for X POW? (Y/N)" | N/A |
| Fundamental Failure | N/A (proceeds to effect) | "The spell fizzles and fails" |
| Spell Effect | Effect-specific message | N/A |

### Game Session Menu Integration

The existing "Cast Spell" option in the game session menu becomes active when magic is unlocked:

**Visibility Logic:**
- Show "Cast Spell" option only when Character.MagicUnlocked is true
- Greyed out or hidden when magic not yet unlocked
- Position between "Combat" and "Manage Inventory" options

### Combat Integration

Several spells directly affect combat state:

**Combat Spells:**
- FIREBALL: Immediately deal damage to current enemy
- POISON NEEDLE: Potentially instant-kill current enemy
- TIMEWARP: Reset combat to initial state
- XENOPHOBIA: Apply damage reduction buff for combat duration
- INVISIBILITY: Exit combat immediately with victory state
- PARALYSIS: Exit combat immediately without victory

**Combat UI Changes:**
- Add "Cast Spell" action during player's turn
- Display active spell effects (ARMOUR, XENOPHOBIA) in status area
- Show spell outcomes in combat log

## Data Model

### Spell Definition Structure

Each spell is represented with the following attributes:

| Attribute | Type | Description |
|-----------|------|-------------|
| Name | string | Spell identifier (e.g., "FIREBALL") |
| PowerCost | int | POWER required to cast |
| Description | string | Human-readable effect description |
| Category | enum | Combat, Navigation, Defensive, Offensive |
| CombatOnly | bool | Whether spell requires active combat |
| DeathOnly | bool | Whether spell requires CurrentLP ≤ 0 (RESURRECTION) |

### Character Magic State

The Character struct already includes magic fields; additional state needed:

| Field | Type | Purpose |
|-------|------|---------|
| ActiveSpellEffects | map[string]int | Active buff/debuff effects with durations (e.g., ARMOUR: 10) |

### Spell Effect State

Certain spells create temporary effects that persist:

| Effect | Duration | Storage |
|--------|----------|---------|
| ARMOUR | Current section | ActiveSpellEffects["ARMOUR"] = 10 (damage reduction) |
| XENOPHOBIA | Current combat | Combat state field |
| CRYPT | Immediate navigation | No persistent state |

## Spell Effect Implementations

### Offensive Spells

#### FIREBALL (Cost: 15)

Deals fixed damage to enemy.

**Behavior:**
- Must be cast during active combat
- Immediately deduct 50 LP from enemy's CurrentLP
- Display damage in combat log
- No effect outside combat (prevent casting or show error)

**Validation:**
- Requires active combat encounter
- Enemy must exist and be alive

#### POISON NEEDLE (Cost: 25)

Conditional instant-kill with immunity check.

**Behavior:**
- Must be cast during active combat
- Roll 1d6 for immunity check
- If 1-3: Enemy immune, no effect, display message
- If 4-6: Enemy dies instantly, set enemy CurrentLP to 0
- Display outcome in combat log

**Validation:**
- Requires active combat encounter
- Enemy must exist and be alive

### Defensive/Buff Spells

#### ARMOUR (Cost: 25)

Creates magical damage reduction for section duration.

**Behavior:**
- Add entry to ActiveSpellEffects: "ARMOUR" → 10
- During damage calculation, subtract 10 from incoming damage
- Effect persists until section changes
- Display active status in character view

**Validation:**
- Can be cast anytime (combat or non-combat)
- Does not stack with multiple casts (single-use per section prevents this)

#### XENOPHOBIA (Cost: 15)

Reduces enemy damage output for current combat.

**Behavior:**
- Apply debuff to current enemy: reduce damage by 5
- Effect persists for duration of current combat only
- Clear effect when combat ends
- Display active status in combat UI

**Validation:**
- Requires active combat encounter
- Effect applies to current enemy only

### Navigation Spells

#### CRYPT (Cost: 150)

Returns player to Crypts for POWER restoration tests.

**Behavior:**
- Trigger navigation to special "Crypt" section
- Present test challenges to restore or increase MaximumPOW
- Tests are interactive UI flows with dice rolls
- Player may gain additional POW capacity based on test outcomes

**Validation:**
- Cannot be cast during combat
- Requires 150 POW (very expensive)

**Implementation Note:**
- Crypt tests may be simplified for Phase 4 MVP
- Could auto-restore to MaximumPOW without complex test system
- Full test implementation can be future enhancement

#### RETRACE (Cost: 20)

Returns to any previously visited section.

**Behavior:**
- Display list of previously visited sections
- Player selects destination section
- Navigate to selected section
- CurrentLP and CurrentPOW remain unchanged

**Validation:**
- Cannot be cast during combat
- Requires at least one previously visited section

**Section Tracking:**
- Maintain visited section history as stack or list
- Store section identifiers with human-readable names

#### TIMEWARP (Cost: 10)

Resets current section to starting state.

**Behavior:**
- Restore player's CurrentLP to value at section entry
- Restore enemy's CurrentLP to initial value (if in combat)
- Clear spell cast history for section (allow re-casting)
- Reset natural inclination check requirement
- Display reset confirmation message

**Validation:**
- Can be cast during or outside combat
- Most useful during losing combat scenarios

### Escape/Tactical Spells

#### INVISIBILITY (Cost: 30)

Avoid combat entirely.

**Behavior:**
- If cast before combat: Prevent combat initiation, proceed to next section
- If cast during combat: Exit combat immediately with "victory" flag set
- Player cannot attack while invisible
- Effect lasts for remainder of current section

**Validation:**
- Can be cast before or during combat
- Clear invisibility status when section changes

**Combat Integration:**
- Set combat outcome to "avoided" or "victory via magic"
- Trigger same post-combat flow as normal victory

#### PARALYSIS (Cost: 30)

Escape from combat without victory.

**Behavior:**
- Exit combat immediately
- Do not grant victory rewards (no skill increase, no loot)
- Navigate to next section
- Display escape message

**Validation:**
- Requires active combat encounter
- Different from INVISIBILITY (no victory flag)

### Special Condition Spells

#### RESURRECTION (Cost: 50)

Recover from death.

**Behavior:**
- Only available when CurrentLP ≤ 0
- Return player to start of current section
- Enemy retains CurrentLP from moment of player death
- Reroll all player characteristics (STR, SPD, STA, CRG, LCK, CHM, ATT)
- Reroll MaximumLP and CurrentPOW based on new characteristics
- Recalculate MaximumPOW (may change)
- Display stat changes to player

**Validation:**
- Only castable when CurrentLP ≤ 0
- Hidden from spell list when CurrentLP > 0
- Automatically offered as option upon death (if sufficient POW)

**UI Flow:**
- Upon death, check if CurrentPOW ≥ 50
- Prompt: "Use RESURRECTION spell? (50 POW)"
- If accepted, trigger resurrection sequence
- If declined or insufficient POW, proceed to normal death handling

## Integration Points

### Combat System Integration

The combat engine must check for active spell effects:

**Damage Calculation Modification:**
- Before applying damage to player, check for ARMOUR effect: subtract 10
- Before applying damage from enemy, check for XENOPHOBIA effect: subtract 5
- Update combat log to show spell-modified damage

**Combat Action Options:**
- Add "Cast Spell" to player action menu during combat turn
- Filter spell list to combat-appropriate spells
- Execute spell immediately within turn flow

**Combat Exit Conditions:**
- INVISIBILITY and PARALYSIS trigger immediate combat exit
- Set appropriate victory/escape flags for post-combat processing

### Character Edit Integration

The character edit screen already handles POW editing:

**Existing Support:**
- CurrentPOW and MaximumPOW are editable fields
- Validation prevents negative values
- No additional changes required for Phase 4

**Additional Consideration:**
- May add "Unlock Magic" button to manually activate magic system during adventure
- Input field for initial POW value when unlocking

### Inventory System Integration

No direct integration required, but spells may interact with items:

**Potential Future Integration:**
- Spell-specific items that reduce casting costs
- Items that provide POW regeneration
- Magic-blocking items that prevent spell casting

**Phase 4 Scope:**
- No item-magic interaction for MVP
- Focus on core spell mechanics

### Save/Load System Integration

Magic state must persist across sessions:

**Serialization Fields:**
All magic-related fields in Character struct are already JSON-tagged:
- CurrentPOW
- MaximumPOW
- MagicUnlocked

**Additional State to Persist:**
- ActiveSpellEffects (serialize as JSON map for ongoing buffs like ARMOUR)

## Edge Cases and Error Handling

### Insufficient POWER Scenarios

**Case:** Player attempts to cast spell without sufficient POW

**Handling:**
1. Display error message: "Insufficient POWER (Need X, Have Y)"
2. Offer LP sacrifice option: "Sacrifice X LP for X POW?"
3. If accepted, validate LP sacrifice legality (would not kill player)
4. If declined or impossible, abort casting, return to spell menu

### Natural Inclination Failure

**Case:** Player fails natural inclination check (roll < 4)

**Handling:**
1. Display message: "Fire*Wolf's hatred of sorcery overwhelms him. Magic is unavailable this section."
2. Disable "Cast Spell" menu option for current section
3. Re-enable when section changes

### Player-Managed Restrictions

**Case:** Gamebook rules limit spell usage (once per section, natural inclination, etc.)

**Handling:**
1. Display all spells as available (no restriction enforcement)
2. Trust player to follow gamebook rules
3. Provide natural inclination check as optional action
4. Player manages their own spell usage tracking

### RESURRECTION Edge Cases

**Case:** Player dies with insufficient POW for RESURRECTION

**Handling:**
1. Skip resurrection prompt
2. Proceed to normal death handling (game over or load save)

**Case:** Player dies with sufficient POW, declines RESURRECTION

**Handling:**
1. Respect choice, proceed to game over
2. May offer confirmation: "Are you sure? This character will be lost."

**Case:** Stat reroll results in significantly different character

**Handling:**
1. Display before/after comparison
2. Update all derived values (MaximumLP, etc.)
3. Continue with new stats

### Combat-Only Spell Outside Combat

**Case:** Player attempts to cast combat-only spell (FIREBALL, POISON NEEDLE, etc.) outside combat

**Handling:**
1. Mark spell as unavailable in spell list with indicator: "(Combat only)"
2. If somehow selected, display error: "This spell can only be cast during combat"
3. Return to spell menu

### Non-Combat Spell During Combat

**Case:** Player attempts to cast non-combat spell (CRYPT, RETRACE) during combat

**Handling:**
1. Mark spell as unavailable with indicator: "(Cannot cast in combat)"
2. If selected, display error: "You cannot cast this spell during combat"
3. Return to spell menu

### Active Effect Persistence

**Case:** Player has ARMOUR active, then section changes

**Handling:**
1. Clear ActiveSpellEffects map on section transition
2. Display message: "Your magical protections fade"
3. ARMOUR no longer applies damage reduction

**Case:** Player has XENOPHOBIA active, combat ends

**Handling:**
1. Clear combat-specific effects when combat ends
2. XENOPHOBIA does not persist to next combat

## Testing Strategy

### Unit Testing

Each component requires comprehensive unit tests:

**Magic Engine Tests:**
- Natural inclination check success/failure
- Power cost validation and deduction
- LP sacrifice calculation
- Fundamental failure rate mechanics
- Each spell effect in isolation
- Spell restriction enforcement (combat-only, death-only, etc.)

**Character Integration Tests:**
- UnlockMagic function
- POW modification functions
- Section transition logic
- Spell cast history tracking
- Active effect management

**Edge Case Tests:**
- Insufficient POW with and without LP sacrifice
- Duplicate spell casting attempts
- Invalid section transitions
- Boundary conditions (CurrentPOW = 0, CurrentLP = 1, etc.)

### Integration Testing

Test interactions between systems:

**Combat-Magic Integration:**
- Cast FIREBALL during combat, verify damage application
- Cast ARMOUR before combat, verify damage reduction
- Cast INVISIBILITY during combat, verify combat exit
- Cast XENOPHOBIA, verify enemy damage reduction
- Cast same spell multiple times, verify no technical restriction

**Navigation-Magic Integration:**
- Cast RETRACE, verify section navigation
- Cast TIMEWARP, verify state reset
- Cast CRYPT, verify Crypt section entry

**Character State Integration:**
- Save character with magic state, load, verify persistence
- Edit character POW values, verify spell casting respects changes
- Level up skill through combat, verify unaffected by magic use

### Manual/Playthrough Testing

Validate complete user experience:

**Test Scenarios:**
1. Unlock magic mid-game, cast first spell
2. Perform natural inclination check, see result message
3. Run out of POW, sacrifice LP, continue casting
4. Die in combat, use RESURRECTION, verify stat reroll
5. Cast all 10 spells in various contexts
6. Cast same spell multiple times (verify no restriction)
7. Use TIMEWARP during losing combat, verify reset
8. Save with active ARMOUR, load, verify effect persists or clears as designed

## Non-Goals for Phase 4

To maintain focused scope, the following are explicitly excluded:

- **Advanced Crypt Tests:** CRYPT spell will perform simple restoration, not complex test sequences
- **Spell Combos:** No special interactions between multiple spells
- **Magic Items:** No items that modify spell behavior or costs
- **Spell Upgrading:** All spells have fixed costs and effects
- **AI Enemy Spellcasting:** Enemies do not cast spells in Phase 4
- **Animated Spell Effects:** Text-based feedback only, no graphical effects

These features may be considered for Phase 5 or future enhancements.

## Success Criteria

Phase 4 is complete when:

1. ✅ All 10 spells are implemented with correct costs and effects
2. ✅ Casting process implements cost payment, FFR check, and effect application
3. ✅ POWER management supports restoration via Character Edit
4. ✅ Natural inclination check available as optional action
5. ✅ Combat spells integrate seamlessly with combat system
6. ✅ Navigation spells correctly modify game flow
7. ✅ UI provides clear spell selection and feedback
8. ✅ Magic state persists across save/load cycles
9. ✅ Unit tests cover all spell mechanics and edge cases
10. ✅ Manual testing confirms complete spell casting workflow

## Implementation Order

Recommended sequence for development:

### Step 1: Magic Engine Foundation
- Create spell definition catalog
- Implement natural inclination check as optional action
- Implement power cost validation and deduction
- Implement fundamental failure rate check

### Step 2: Simple Spells
- Implement ARMOUR (defensive buff)
- Implement FIREBALL (direct damage)
- Implement XENOPHOBIA (enemy debuff)
- Test with existing combat system

### Step 3: Navigation Spells
- Implement TIMEWARP (section reset)
- Implement RETRACE (section navigation)
- Implement simplified CRYPT (POW restoration)

### Step 4: Tactical Spells
- Implement INVISIBILITY (combat avoidance)
- Implement PARALYSIS (combat escape)
- Implement POISON NEEDLE (conditional kill)

### Step 5: Special Spell
- Implement RESURRECTION (death recovery)
- Create death-triggered UI flow

### Step 6: UI Integration
- Create spell casting screen
- Integrate "Cast Spell" into game session menu
- Add spell feedback display
- Update combat UI for spell actions

### Step 7: Testing and Polish
- Complete unit test coverage
- Integration testing with combat and navigation
- Manual playthrough testing
- Bug fixes and refinement

### Step 8: Documentation
- Create PHASE4_COMPLETE.md summary
- Update README.md with magic system description
- Document spell mechanics for players
