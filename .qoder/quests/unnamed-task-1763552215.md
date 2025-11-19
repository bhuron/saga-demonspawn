# Phase 5: Polish & User Experience

## Objective

Enhance the Sagas of the Demonspawn companion application with professional styling, comprehensive help system, and configurable preferences to improve overall user experience and visual appeal.

## Background

Phases 1-4 have established a fully functional application with character management, combat, inventory, and magic systems. Phase 5 focuses on refining the user interface through visual theming, providing contextual help, and enabling user customization through configuration.

The lipgloss library is already available as a dependency but currently unused. All UI rendering uses plain text formatting with basic box-drawing characters.

## Scope

### In Scope

1. Visual theming and styling system using lipgloss
2. Contextual help system accessible throughout the application
3. Configuration management for user preferences
4. Enhanced error messages and validation feedback
5. Quality of life improvements to existing screens

### Out of Scope

- New game mechanics or systems
- Refactoring existing business logic
- Performance optimization beyond UX improvements
- Internationalization or localization
- Network features or multiplayer

## Design Overview

Phase 5 transforms the functional terminal application into a polished, visually appealing experience through three core enhancements:

1. **Theming System**: Unified color schemes and typography
2. **Help System**: Context-sensitive documentation
3. **Configuration**: Persistent user preferences

## Feature Design

### 1. Styling and Theming System

#### Visual Design Language

Establish a consistent visual identity across all screens using lipgloss styling primitives.

**Color Palette**

| Element | Usage | Rationale |
|---------|-------|-----------|
| Primary (Cyan/Bright Blue) | Headers, active selections, focused elements | High visibility, fantasy theme |
| Secondary (Purple/Magenta) | Subheadings, magical elements (POW, spells) | Mystical association |
| Success (Green) | Positive outcomes, successful actions, health | Universal positive indicator |
| Warning (Yellow/Gold) | Warnings, low resources, important notices | Attention without alarm |
| Danger (Red) | Critical states, death, errors, negative effects | Universal danger signal |
| Neutral (Gray/White) | Body text, inactive elements, borders | Readable contrast |
| Accent (Orange) | Combat actions, special items, highlights | Energy and action |

**Typography Hierarchy**

| Style | Application | Visual Treatment |
|-------|-------------|------------------|
| Title | Screen headers, application name | Bold, large, bordered, primary color |
| Heading | Section headers, character name | Bold, secondary color |
| Label | Field labels, stat names | Normal weight, neutral color |
| Value | Stat values, numbers | Bold when important, color-coded by context |
| Body | Help text, descriptions | Normal weight, neutral color |
| Emphasis | Important messages, active state | Bold, color accent |
| Muted | Disabled options, hints | Faint, gray |

**Layout Components**

Define reusable styled components:

| Component | Purpose | Styling Approach |
|-----------|---------|------------------|
| Box | Containers with borders | Rounded or square borders, padding, title support |
| Header | Screen titles | Full-width, centered, decorative borders |
| Menu Item | Selectable list items | Cursor prefix, hover highlight, alignment |
| Stat Display | Character attributes | Label-value pairs, color-coded values |
| Button | Action indicators | Bordered, highlight on focus |
| Message | Feedback and notifications | Type-specific colors (success/warning/error) |
| Panel | Grouped information | Subtle background, borders |

**Screen-Specific Styling**

Apply consistent theming to each screen:

| Screen | Key Styling Elements |
|--------|---------------------|
| Main Menu | Dramatic title with ASCII art border, gradient effect on selections |
| Character Creation | Step indicator, stat value color-coding (high=green, low=yellow/red) |
| Character View | Multi-column layout, grouped stats, equipment highlights |
| Character Edit | Field focus indicator, input mode visual feedback |
| Combat | Health bars with color transitions, action emphasis, log with severity colors |
| Inventory | Item rarity/type indicators, equipped item badges, scrollbar |
| Magic | POW meter, spell cost indicators, effect duration displays |
| Load Character | File metadata formatting, timestamp styling |

#### Implementation Strategy

Create a centralized theming package that provides:

- Pre-configured lipgloss styles for all component types
- Color scheme definition with fallback for limited terminal support
- Utility functions for common formatting patterns (stat values, health bars, borders)
- Consistent spacing and padding constants

All existing view functions will be updated to use the theme system rather than raw string formatting.

### 2. Help System

#### Architecture

**Multi-Layer Help Approach**

| Layer | Trigger | Content | Location |
|-------|---------|---------|----------|
| Inline Help | Always visible | Brief navigation hints | Bottom of each screen |
| Context Help | '?' key press | Screen-specific detailed help | Overlay modal |
| Global Help | 'H' key from main menu | Complete application guide | Full-screen documentation |
| Tooltips | Hover/focus (future) | Field-level explanations | Adjacent to element |

#### Help Content Structure

**Global Help Topics**

The comprehensive help system covers:

1. Getting Started
   - Application overview
   - Navigation basics
   - Saving and loading

2. Character Management
   - Creating characters
   - Understanding characteristics
   - Derived values (LP, Skill, POW)
   - Editing stats
   
3. Combat System
   - Initiative mechanics
   - To-hit calculations
   - Damage and armor
   - Death saves

4. Inventory Management
   - Equipment types
   - Special items
   - Item restrictions
   - Combat equipment locks

5. Magic System
   - Unlocking magic
   - Spell casting process
   - POW management
   - Natural inclination
   - Fundamental Failure Rate
   - Spell descriptions

6. Game Rules Reference
   - Quick reference for formulas
   - Skill progression
   - Equipment bonuses
   - Special mechanics

**Screen-Specific Help**

Each screen provides focused help relevant to current context:

| Screen | Help Content |
|--------|--------------|
| Main Menu | Navigation keys, creating vs loading characters |
| Character Creation | Rolling mechanics, characteristic meanings, equipment choices |
| Character View | Reading the character sheet, stat explanations |
| Character Edit | Field navigation, unlocking magic, value constraints |
| Combat Setup | Enemy stats, starting combat |
| Combat | Combat flow, action selection, spell casting in combat, death saves |
| Inventory | Acquiring items, equipping/unequipping, special item rules |
| Magic | Available spells, POW costs, LP sacrifice, casting restrictions |
| Load Character | File selection, character metadata |

#### Help Display Design

**Help Modal Specification**

- Overlay on top of current screen with semi-transparent background effect
- Scrollable content area for long help text
- Clearly marked sections with styled headers
- Key bindings highlighted in accent color
- Footer navigation: "↑/↓ Scroll, Esc/? to close"
- Preserve underlying screen state

**Help Content Format**

Help text uses structured formatting:

- Section headers with visual separators
- Bullet lists for navigation keys
- Example boxes for calculations
- Tables for spell/item reference
- Callout boxes for important notes

### 3. Configuration Management

#### Configuration Architecture

**Configuration File Structure**

Store user preferences in a simple, human-readable format:

| Setting Category | Settings | Default Values |
|------------------|----------|----------------|
| Appearance | Color scheme (dark/light/custom), Use ASCII vs Unicode, Show animations | dark, unicode, true |
| Gameplay | Confirm dangerous actions, Auto-save on exit, Show roll details | true, true, true |
| Accessibility | High contrast mode, Reduced motion, Larger text | false, false, false |
| Paths | Save directory, Character backup enabled | ~/.saga-demonspawn, true |

**Configuration Location**

- Default location: `~/.saga-demonspawn/config.json`
- Fallback to current directory if home directory unavailable
- Create default config on first run

**Configuration Schema**

Define configuration structure with validation rules:

| Field | Type | Validation | Purpose |
|-------|------|------------|---------|
| theme | string | Enum: dark, light, custom | Color scheme selection |
| use_unicode | boolean | - | Enable Unicode characters vs ASCII |
| show_animations | boolean | - | Enable visual transitions |
| confirm_actions | boolean | - | Require confirmation for risky actions |
| auto_save | boolean | - | Save character on application exit |
| show_roll_details | boolean | - | Display dice roll breakdowns |
| high_contrast | boolean | - | Accessibility: enhanced contrast |
| save_directory | string | Valid directory path | Character save location |
| character_backup | boolean | - | Create timestamped backups |

#### Configuration UI

**Settings Screen**

Add "Settings" option to Main Menu:

- Navigate settings with arrow keys
- Toggle boolean values with Enter/Space
- Edit text values with inline input
- Preview theme changes in real-time
- Save/Cancel/Reset to defaults actions

**Settings Display Layout**

```
╔═══════════════════════════════════════╗
║            SETTINGS                   ║
╚═══════════════════════════════════════╝

  Appearance
  > Color Scheme: [Dark] Light Custom
    Use Unicode: [✓] Yes  [ ] No
    Show Animations: [✓] Enabled

  Gameplay
    Confirm Actions: [✓] Enabled
    Auto-save on Exit: [✓] Enabled
    Show Roll Details: [✓] Enabled

  Accessibility
    High Contrast: [ ] Disabled
    Reduced Motion: [ ] Disabled

  Files
    Save Directory: ~/.saga-demonspawn

  [Save] [Cancel] [Reset to Defaults]

  Navigate: ↑/↓  Toggle: Enter  Edit: E  Esc: Cancel
```

#### Configuration Loading

Configuration loading process:

1. Application startup checks for config file
2. Load existing config or create default
3. Validate all settings
4. Apply theme and preferences to UI components
5. Invalid settings fall back to defaults with warning

### 4. Enhanced Error Messages and Validation Feedback

#### Error Message Improvements

**Current State**: Simple error text displayed without context or visual emphasis.

**Enhanced Approach**

| Error Type | Visual Treatment | Content Enhancement |
|------------|------------------|---------------------|
| Validation Errors | Yellow border box, warning icon | Specific field mentioned, valid range shown, example provided |
| File Errors | Red border box, error icon | Friendly explanation, suggested action, path displayed |
| Combat Errors | Orange highlight in log | Clear cause, impact on combat state |
| Magic Errors | Purple/magenta box | POW cost shown, alternative actions suggested |

**Error Message Structure**

All error messages follow pattern:

1. Error type/title (bold, colored)
2. What happened (clear description)
3. Why it happened (context)
4. How to fix it (actionable guidance)

Example transformations:

| Before | After |
|--------|-------|
| "Invalid value" | "⚠ Invalid Value\nSTR must be between 16-96.\nYou entered: 150\nPress 'r' to re-roll characteristics." |
| "Cannot cast spell" | "✗ Insufficient POWER\nFIREBALL costs 15 POW.\nYou have: 10 POW (Need: 5 more)\nSacrifice 5 LP to cast? (Y/N)" |
| "File not found" | "⚠ Save File Not Found\nCannot load: FireWolf_2025-01-15.json\nThe file may have been moved or deleted.\nPress Esc to return to menu." |

#### Validation Feedback

**Real-Time Validation**

Provide immediate feedback during input:

| Input Context | Validation Behavior |
|---------------|---------------------|
| Character Edit | Show valid range next to field, color-code input (red=invalid, yellow=warning, green=valid) |
| Combat Setup | Validate enemy stats as typed, show warnings for unusual values |
| Equipment Selection | Highlight incompatible combinations (shield + orb) before selection |
| Spell Casting | Show POW availability, indicate if insufficient before cast attempt |

**Success Feedback**

Balance error messages with positive confirmation:

- Successful saves: "✓ Character saved successfully"
- Successful spell cast: "✓ FIREBALL cast! Enemy takes 50 damage"
- Equipment equipped: "✓ Doombringer equipped (+15 damage)"
- Magic unlocked: "✓ Magic unlocked! POWER: 50/50"

### 5. Quality of Life Improvements

#### Navigation Enhancements

| Enhancement | Benefit |
|-------------|---------|
| Breadcrumb trail | Show navigation path (Main Menu > Game Session > Combat) |
| Quick jump keys | Press number keys to jump to menu items (1-9) |
| Recent characters | Show last 3 loaded characters on main menu for quick access |
| Confirm quit | Prevent accidental exits with confirmation |

#### Display Improvements

| Improvement | Description |
|-------------|-------------|
| Health bars | Visual HP/LP bars with color gradients (green→yellow→red) |
| POW meter | Visual POWER gauge in magic screen |
| Stat change indicators | Show +/- when stats change (temporary effects, equipment) |
| Combat log scroll | Display last N combat messages with scroll controls |
| Item tooltips | Brief item descriptions in inventory |
| Active effects summary | Clear display of active spell effects and durations |

#### Input Improvements

| Enhancement | Description |
|-------------|-------------|
| Input history | Arrow up/down to recall previous inputs in edit fields |
| Autocomplete | Suggest saved character names when loading |
| Copy character | Duplicate existing character as template |
| Bulk actions | Select multiple items in inventory (future enhancement) |

#### Character Management

| Feature | Description |
|---------|-------------|
| Character backup | Automatic backup before editing or combat |
| Character notes | Add free-form notes to character (journal entries) |
| Character portrait | ASCII art portrait selection (optional) |
| Stats history | Track stat changes over time |

## Technical Considerations

### Theming Implementation

**Lipgloss Integration Points**

- Create `pkg/ui/theme/theme.go` package
- Define style constants and builder functions
- Export pre-configured styles for all component types
- Ensure terminal capability detection (fallback for limited color support)
- Maintain separation: theme package provides styles, view functions apply them

**Color Compatibility**

- Detect terminal color capability (8-color, 16-color, 256-color, true color)
- Provide degraded color schemes for limited terminals
- Test on common terminals: xterm, gnome-terminal, iTerm2, Windows Terminal
- Fallback to monochrome with bold/underline emphasis if needed

### Help System Implementation

**Help Content Storage**

Options for storing help content:

| Approach | Pros | Cons | Recommendation |
|----------|------|------|----------------|
| Embedded strings in code | No external dependencies | Hard to edit, clutters code | Not recommended |
| Embedded text files | Easy to edit, clean code | Requires embed directive | Recommended |
| External files | Most flexible | Deployment complexity | Optional for advanced users |

Recommended: Use Go's `embed` package to embed markdown or plain text help files at compile time.

**Help Modal State Management**

- Help modal is a special overlay state, not a screen
- Preserve underlying screen state while help is displayed
- Help scroll position independent of main content scroll
- Handle help dismissal cleanly without disrupting current screen

### Configuration Implementation

**Configuration Package Structure**

Create `internal/config/config.go`:

- Define Config struct with all settings
- Provide Load, Save, Validate, and Reset functions
- Export singleton instance or pass via Model
- Handle migration for config schema changes in future versions

**Configuration Validation**

Validate configuration at load time:

- Check directory paths exist or can be created
- Validate enum values against allowed options
- Ensure boolean and numeric values are correct types
- Provide clear error messages for invalid config
- Fall back to defaults for invalid individual settings

**Configuration Application**

Configuration affects:

- Theme package: Apply selected color scheme
- View rendering: Use Unicode or ASCII characters
- Update logic: Enable/disable confirmations
- Application shutdown: Trigger auto-save if enabled

### Error Handling Architecture

**Centralized Error Formatting**

Create error formatting utilities:

- Define error severity levels (Info, Warning, Error, Critical)
- Provide formatting functions that apply theme styles
- Support structured error data (field, value, constraint)
- Generate actionable error messages with suggestions

**Error State Management**

- Store last error in Model.Err
- Clear errors on successful actions
- Display errors prominently without blocking interaction
- Allow dismissing errors explicitly

## User Interaction Flows

### Theming Experience

**User Flow: Changing Theme**

1. User navigates to Settings from Main Menu
2. Selects "Color Scheme" option
3. Cycles through Dark → Light → Custom using arrow keys or Enter
4. Preview shows sample UI with selected theme in real-time
5. User presses 'S' to save or Esc to cancel
6. Theme applies immediately to all screens

### Help System Usage

**User Flow: Accessing Context Help**

1. User is viewing Character Edit screen
2. Presses '?' key
3. Help modal overlays current screen with semi-transparent background
4. Help content shows Character Edit instructions, field meanings, and magic unlock process
5. User scrolls help content with arrow keys if needed
6. Presses '?' or Esc to dismiss help
7. Returns to Character Edit screen in exact same state

**User Flow: Accessing Global Help**

1. User is at Main Menu
2. Presses 'H' key for comprehensive help
3. Full-screen help displays with section navigation
4. User navigates sections with arrow keys or page up/down
5. Help organized: Getting Started, Character, Combat, Inventory, Magic, Rules
6. User presses Esc or 'H' to return to Main Menu

### Configuration Workflow

**User Flow: First-Time Configuration**

1. Application starts for first time
2. No config file exists
3. Application creates default configuration in `~/.saga-demonspawn/config.json`
4. Application runs with default settings (dark theme, Unicode, confirmations enabled)
5. User can access Settings later to customize

**User Flow: Modifying Configuration**

1. User selects Settings from Main Menu
2. Settings screen displays all categories
3. User navigates to "Auto-save on Exit"
4. Presses Enter to toggle from Enabled → Disabled
5. Visual feedback shows change immediately
6. User navigates to bottom and selects "Save"
7. Configuration persists to disk
8. Settings apply immediately to current session

## Validation and Constraints

### Visual Constraints

- All colors must have sufficient contrast ratio for accessibility (WCAG AA minimum)
- Text must remain readable on monochrome terminals
- Borders and boxes must not exceed terminal width
- Unicode characters must have ASCII fallbacks

### Configuration Constraints

- Save directory must be writable
- Config file must be valid JSON
- Invalid settings fall back to defaults without crashing
- Config changes must not require application restart

### Help System Constraints

- Help content must fit terminal height or provide scrolling
- Help modal must not interfere with background screen state
- Help text must be concise yet comprehensive
- All features must have help documentation

## Testing Strategy

### Visual Testing

- Manual testing on multiple terminal emulators
- Verify color schemes in 8-color, 16-color, and 256-color modes
- Test on small terminal sizes (80x24 minimum)
- Verify Unicode fallback to ASCII

### Configuration Testing

**Test Cases**

| Test Case | Expected Behavior |
|-----------|------------------|
| No config file | Create default config, application runs normally |
| Invalid JSON | Show error, use defaults, attempt to create valid config |
| Invalid theme value | Use default theme, warn user |
| Non-existent save directory | Create directory or use current directory |
| Read-only config file | Warn user, run with existing config, cannot save changes |
| Config with future version | Gracefully handle unknown fields, use defaults for missing fields |

### Help System Testing

- Verify help content accuracy
- Test help modal overlay rendering
- Verify scrolling in long help content
- Test help dismissal preserves screen state
- Verify context help available on all screens

### Error Message Testing

- Trigger all error conditions
- Verify error messages are clear and actionable
- Test error display across different screen sizes
- Verify error dismissal workflow

## Documentation Requirements

### User Documentation

Update README.md with:

- Configuration file location and structure
- Help system usage
- Theming options
- Quality of life feature descriptions

### Developer Documentation

Create developer guide covering:

- Theme system architecture
- Adding new themed components
- Extending help content
- Adding new configuration options
- Style guide for consistency

## Success Criteria

Phase 5 is complete when:

- ✓ Lipgloss theming applied to all screens with consistent visual style
- ✓ At least 2 color schemes implemented (dark, light)
- ✓ Context help ('?') available on all screens
- ✓ Global help system accessible with comprehensive content
- ✓ Configuration system implemented with Settings screen
- ✓ Configuration persists and loads correctly
- ✓ Error messages enhanced with visual styling and actionable guidance
- ✓ At least 5 quality of life improvements implemented
- ✓ Application maintains responsiveness with new features
- ✓ All features tested on multiple terminal emulators
- ✓ Documentation updated with Phase 5 features
- ✓ User experience is noticeably improved from Phase 4

## Dependencies

- lipgloss (already in go.mod) - styling and theming
- Go's embed package (standard library) - embedding help content
- Go's encoding/json (standard library) - configuration persistence
- Go's os and path/filepath (standard library) - file system operations

## Future Enhancements

Potential additions beyond Phase 5:

- Custom color scheme editor
- Export character as formatted text or PDF
- Game session logging and playback
- Spell effect animations
- Item rarity visual indicators
- Achievement system
- Multiple configuration profiles
- Plugin system for custom themes
