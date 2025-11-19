// Package combat implements the combat system for Sagas of the Demonspawn.
package combat

import (
	"fmt"

	"github.com/benoit/saga-demonspawn/internal/character"
	"github.com/benoit/saga-demonspawn/internal/dice"
)

// Enemy represents an opponent in combat.
type Enemy struct {
	Name             string `json:"name"`              // Enemy identifier
	Strength         int    `json:"strength"`          // STR: Physical power
	Speed            int    `json:"speed"`             // SPD: Agility and reaction
	Stamina          int    `json:"stamina"`           // STA: Endurance
	Courage          int    `json:"courage"`           // CRG: Bravery
	Luck             int    `json:"luck"`              // LCK: Fortune
	Skill            int    `json:"skill"`             // SKL: Combat proficiency
	CurrentLP        int    `json:"current_lp"`        // Current life points
	MaximumLP        int    `json:"maximum_lp"`        // Maximum life points
	WeaponBonus      int    `json:"weapon_bonus"`      // Weapon damage bonus
	ArmorProtection  int    `json:"armor_protection"`  // Armor damage reduction
	IsDemonspawn     bool   `json:"is_demonspawn"`     // For special item interactions
}

// NewEnemy creates a new enemy with validation.
func NewEnemy(name string, str, spd, sta, crg, lck, skill, currentLP, maxLP, weaponBonus, armorProtection int, isDemonspawn bool) (*Enemy, error) {
	if name == "" {
		return nil, fmt.Errorf("enemy name cannot be empty")
	}
	if str < 0 {
		return nil, fmt.Errorf("strength cannot be negative: %d", str)
	}
	if spd < 0 {
		return nil, fmt.Errorf("speed cannot be negative: %d", spd)
	}
	if sta < 0 {
		return nil, fmt.Errorf("stamina cannot be negative: %d", sta)
	}
	if crg < 0 {
		return nil, fmt.Errorf("courage cannot be negative: %d", crg)
	}
	if lck < 0 {
		return nil, fmt.Errorf("luck cannot be negative: %d", lck)
	}
	if skill < 0 {
		return nil, fmt.Errorf("skill cannot be negative: %d", skill)
	}
	if maxLP <= 0 {
		return nil, fmt.Errorf("maximum LP must be positive: %d", maxLP)
	}
	if currentLP < 0 {
		return nil, fmt.Errorf("current LP cannot be negative: %d", currentLP)
	}
	if weaponBonus < 0 {
		return nil, fmt.Errorf("weapon bonus cannot be negative: %d", weaponBonus)
	}
	if armorProtection < 0 {
		return nil, fmt.Errorf("armor protection cannot be negative: %d", armorProtection)
	}

	return &Enemy{
		Name:            name,
		Strength:        str,
		Speed:           spd,
		Stamina:         sta,
		Courage:         crg,
		Luck:            lck,
		Skill:           skill,
		CurrentLP:       currentLP,
		MaximumLP:       maxLP,
		WeaponBonus:     weaponBonus,
		ArmorProtection: armorProtection,
		IsDemonspawn:    isDemonspawn,
	}, nil
}

// CombatState encapsulates the complete state of an active combat encounter.
type CombatState struct {
	IsActive            bool     `json:"is_active"`              // Whether combat is currently ongoing
	CurrentRound        int      `json:"current_round"`          // Round counter (starts at 1)
	PlayerTurn          bool     `json:"player_turn"`            // True if player turn, false if enemy turn
	PlayerFirstStrike   bool     `json:"player_first_strike"`    // Tracks who won initiative
	DeathSaveUsed       bool     `json:"death_save_used"`        // Prevents multiple death saves
	EnduranceLimit      int      `json:"endurance_limit"`        // Calculated max rounds (STA ÷ 20)
	RoundsSinceLastRest int      `json:"rounds_since_last_rest"` // Endurance tracking
	EnemyEnduranceLimit      int      `json:"enemy_endurance_limit"`        // Enemy's max rounds (STA ÷ 20)
	EnemyRoundsSinceLastRest int      `json:"enemy_rounds_since_last_rest"` // Enemy endurance tracking
	Enemy               *Enemy   `json:"enemy"`                  // Complete enemy data
	CombatLog           []string `json:"combat_log"`             // Historical combat messages
	PlayerInitiative    int      `json:"player_initiative"`      // Player's initiative roll result
	EnemyInitiative     int      `json:"enemy_initiative"`       // Enemy's initiative roll result
}

// NewCombatState creates a new combat state with the given enemy.
func NewCombatState(enemy *Enemy, enduranceLimit int) *CombatState {
	// Calculate enemy endurance limit
	enemyEnduranceLimit := enemy.Stamina / 10
	
	return &CombatState{
		IsActive:            true,
		CurrentRound:        1,
		PlayerTurn:          false, // Will be set by initiative
		PlayerFirstStrike:   false, // Will be set by initiative
		DeathSaveUsed:       false,
		EnduranceLimit:      enduranceLimit,
		RoundsSinceLastRest: 0,
		EnemyEnduranceLimit:      enemyEnduranceLimit,
		EnemyRoundsSinceLastRest: 0,
		Enemy:               enemy,
		CombatLog:           make([]string, 0),
		PlayerInitiative:    0,
		EnemyInitiative:     0,
	}
}

// AddLogEntry appends a message to the combat log.
func (cs *CombatState) AddLogEntry(message string) {
	cs.CombatLog = append(cs.CombatLog, message)
}

// CalculateInitiative determines who strikes first in combat.
// Returns the player's initiative score, enemy's initiative score, and whether player goes first.
func CalculateInitiative(player *character.Character, enemy *Enemy, roller dice.Roller) (int, int, bool) {
	playerRoll := roller.Roll2D6()
	enemyRoll := roller.Roll2D6()

	playerInitiative := playerRoll + player.Speed + player.Courage + player.Luck
	enemyInitiative := enemyRoll + enemy.Speed + enemy.Courage + enemy.Luck

	return playerInitiative, enemyInitiative, playerInitiative > enemyInitiative
}

// CalculateToHitRequirement determines the number needed on 2d6 to hit.
// Base requirement is 7, reduced by skill (1 per 10 points) and luck (1 if >= 72).
// Minimum requirement is always 2.
func CalculateToHitRequirement(skill, luck int) int {
	requirement := 7

	// Skill modifier: -1 per 10 full points
	requirement -= skill / 10

	// Luck modifier: -1 if luck >= 72
	if luck >= 72 {
		requirement--
	}

	// Minimum requirement is always 2
	if requirement < 2 {
		requirement = 2
	}

	return requirement
}

// CalculateDamage computes the total damage before armor reduction.
// Formula: (roll × 5) + (STR ÷ 10 × 5) + weapon bonus
func CalculateDamage(rollResult, strength, weaponBonus int) int {
	baseDamage := rollResult * 5
	strengthBonus := (strength / 10) * 5
	totalDamage := baseDamage + strengthBonus + weaponBonus
	return totalDamage
}

// ApplyArmorReduction subtracts armor protection from damage.
// Returns the final damage (minimum 0).
func ApplyArmorReduction(damage, armorProtection int) int {
	finalDamage := damage - armorProtection
	if finalDamage < 0 {
		finalDamage = 0
	}
	return finalDamage
}

// CheckEndurance determines if rest is required based on rounds fought.
func CheckEndurance(roundsSinceLastRest, enduranceLimit int) bool {
	return roundsSinceLastRest >= enduranceLimit && enduranceLimit > 0
}

// ExecuteDeathSave performs a death save roll.
// Returns true if successful (result <= luck), false otherwise.
func ExecuteDeathSave(luck int, roller dice.Roller) (int, bool) {
	roll := roller.Roll2D6() * 10
	return roll, roll <= luck
}

// AttackResult contains the outcome of an attack.
type AttackResult struct {
	Roll            int    // The 2d6 roll result
	Requirement     int    // Required roll to hit
	Hit             bool   // Whether the attack hit
	DamageBeforeArmor int  // Damage before armor reduction
	FinalDamage     int    // Damage after armor reduction
	TargetLP        int    // Target's LP after damage
}

// ExecutePlayerAttack performs a player attack and updates combat state.
func ExecutePlayerAttack(cs *CombatState, player *character.Character, roller dice.Roller) AttackResult {
	// Calculate to-hit requirement
	requirement := CalculateToHitRequirement(player.Skill, player.Luck)
	
	// Roll to hit
	roll := roller.Roll2D6()
	hit := roll >= requirement
	
	result := AttackResult{
		Roll:        roll,
		Requirement: requirement,
		Hit:         hit,
	}
	
	if hit {
		// Calculate damage
		weaponBonus := 0
		if player.EquippedWeapon != nil {
			weaponBonus = player.EquippedWeapon.DamageBonus
		}
		
		damageBeforeArmor := CalculateDamage(roll, player.Strength, weaponBonus)
		finalDamage := ApplyArmorReduction(damageBeforeArmor, cs.Enemy.ArmorProtection)
		
		// Apply damage
		cs.Enemy.CurrentLP -= finalDamage
		if cs.Enemy.CurrentLP < 0 {
			cs.Enemy.CurrentLP = 0
		}
		
		result.DamageBeforeArmor = damageBeforeArmor
		result.FinalDamage = finalDamage
		result.TargetLP = cs.Enemy.CurrentLP
	}
	
	return result
}

// ExecuteEnemyAttack performs an enemy attack and updates combat state.
func ExecuteEnemyAttack(cs *CombatState, player *character.Character, roller dice.Roller) AttackResult {
	// Calculate to-hit requirement
	requirement := CalculateToHitRequirement(cs.Enemy.Skill, cs.Enemy.Luck)
	
	// Roll to hit
	roll := roller.Roll2D6()
	hit := roll >= requirement
	
	result := AttackResult{
		Roll:        roll,
		Requirement: requirement,
		Hit:         hit,
	}
	
	if hit {
		// Calculate damage
		damageBeforeArmor := CalculateDamage(roll, cs.Enemy.Strength, cs.Enemy.WeaponBonus)
		
		// Calculate player armor protection
		armorProtection := 0
		if player.EquippedArmor != nil {
			armorProtection = player.EquippedArmor.Protection
		}
		if player.HasShield {
			if player.EquippedArmor != nil {
				armorProtection += 5 // Shield with armor is -5
			} else {
				armorProtection += 7 // Shield alone is -7
			}
		}
		
		finalDamage := ApplyArmorReduction(damageBeforeArmor, armorProtection)
		
		// Apply damage
		player.ModifyLP(-finalDamage)
		
		result.DamageBeforeArmor = damageBeforeArmor
		result.FinalDamage = finalDamage
		result.TargetLP = player.CurrentLP
	}
	
	return result
}

// StartCombat initializes combat with initiative roll.
func StartCombat(player *character.Character, enemy *Enemy, roller dice.Roller) *CombatState {
	// Calculate endurance limit
	enduranceLimit := player.Stamina / 10
	
	// Create combat state
	cs := NewCombatState(enemy, enduranceLimit)
	
	// Roll initiative
	playerInit, enemyInit, playerFirst := CalculateInitiative(player, enemy, roller)
	cs.PlayerInitiative = playerInit
	cs.EnemyInitiative = enemyInit
	cs.PlayerFirstStrike = playerFirst
	cs.PlayerTurn = playerFirst
	
	return cs
}

// CheckVictory returns true if the enemy is defeated.
func CheckVictory(cs *CombatState) bool {
	return cs.Enemy.CurrentLP <= 0
}

// CheckDefeat returns true if the player is defeated (LP <= 0 and no death save available).
func CheckDefeat(player *character.Character, cs *CombatState) bool {
	return player.CurrentLP <= 0
}

// NextTurn advances combat to the next turn, handling turn alternation and round progression.
func NextTurn(cs *CombatState) {
	// Toggle turn
	cs.PlayerTurn = !cs.PlayerTurn
	
	// If we're back to the first striker, increment round and endurance trackers
	if cs.PlayerTurn == cs.PlayerFirstStrike {
		cs.CurrentRound++
		cs.RoundsSinceLastRest++
		cs.EnemyRoundsSinceLastRest++
	}
}

// ProcessRest handles the rest mechanic when endurance is depleted.
func ProcessRest(cs *CombatState) {
	cs.RoundsSinceLastRest = 0
}

// ProcessEnemyRest handles the rest mechanic when enemy endurance is depleted.
func ProcessEnemyRest(cs *CombatState) {
	cs.EnemyRoundsSinceLastRest = 0
}

// ResolveCombatVictory updates player stats after winning combat.
func ResolveCombatVictory(player *character.Character) {
	player.IncrementEnemiesDefeated()
	player.ModifySkill(1)
}

// AttemptDeathSave performs a death save and restores player if successful.
// Returns the roll result and whether the save was successful.
func AttemptDeathSave(player *character.Character, cs *CombatState, roller dice.Roller) (int, bool) {
	if cs.DeathSaveUsed {
		return 0, false
	}
	
	roll, success := ExecuteDeathSave(player.Luck, roller)
	cs.DeathSaveUsed = true
	
	if success {
		// Restore player to max LP
		player.SetLP(player.MaximumLP)
		
		// Reset combat to beginning (but enemy keeps current LP)
		cs.CurrentRound = 1
		cs.RoundsSinceLastRest = 0
		cs.EnemyRoundsSinceLastRest = 0
		
		// Re-roll initiative
		playerInit, enemyInit, playerFirst := CalculateInitiative(player, cs.Enemy, roller)
		cs.PlayerInitiative = playerInit
		cs.EnemyInitiative = enemyInit
		cs.PlayerFirstStrike = playerFirst
		cs.PlayerTurn = playerFirst
	}
	
	return roll, success
}
