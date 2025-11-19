# Product Overview

Sagas of the Demonspawn is a command-line companion application for the "Sagas of the Demonspawn" gamebook. Built with Go and Bubble Tea, it serves as a rules engine and character management tool.

## Core Features

- **Character Management**: Create, edit, save/load characters with 7 characteristics (STR, SPD, STA, CRG, LCK, CHM, ATT)
- **Combat System**: Turn-based combat with initiative, to-hit rolls, damage calculation, and death saves
- **Magic System**: 10 spells with POW cost, natural inclination checks, and fundamental failure rates
- **Inventory Management**: Equipment system with weapons, armor, shields, and special items (Healing Stone, Doombringer, The Orb)
- **Configuration**: User preferences for themes, auto-save, Unicode display
- **Help System**: Context-sensitive help accessible via `?` key

## Game Mechanics

The application implements the complete ruleset from the gamebook, including:
- Dice-based stat rolling (2d6 × 8 for characteristics)
- Life Points (LP) = sum of all characteristics
- Skill progression (+1 per enemy defeated)
- Combat formulas with modifiers for skill, luck, strength
- Stamina-based endurance system
- Special item mechanics and mutual exclusions
