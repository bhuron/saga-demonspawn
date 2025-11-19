package magic

import (
	"fmt"

	"github.com/benoit/saga-demonspawn/internal/dice"
)

// CastResult represents the outcome of a spell casting attempt.
type CastResult struct {
	Success           bool
	FFRFailed         bool   // Fundamental Failure Rate failed
	InsufficientPower bool   // Not enough POW
	Message           string // Human-readable result message
	PowerSpent        int    // Amount of POW consumed
	RequiresSacrifice bool   // Whether LP sacrifice is needed
	SacrificeAmount   int    // Amount of LP to sacrifice for POW
}

// NaturalInclinationCheck performs the natural inclination check.
// Returns true if Fire*Wolf overcomes his aversion to magic (roll >= 4).
func NaturalInclinationCheck(roller dice.Roller) (bool, int) {
	roll := roller.Roll2D6()
	return roll >= 4, roll
}

// CanAffordSpell checks if character has enough POW for the spell.
func CanAffordSpell(currentPOW int, spellCost int) bool {
	return currentPOW >= spellCost
}

// CalculateSacrificeNeeded calculates how much LP needs to be sacrificed.
func CalculateSacrificeNeeded(currentPOW int, spellCost int) int {
	if currentPOW >= spellCost {
		return 0
	}
	return spellCost - currentPOW
}

// CanSacrificeLP checks if character can sacrifice enough LP without dying.
func CanSacrificeLP(currentLP int, sacrificeAmount int) bool {
	return currentLP > sacrificeAmount // Must survive with at least 1 LP
}

// FundamentalFailureRate performs the FFR check.
// Returns true if spell succeeds (roll >= 6).
func FundamentalFailureRate(roller dice.Roller) (bool, int) {
	roll := roller.Roll2D6()
	return roll >= 6, roll
}

// ValidateCast checks if a spell can be cast given the current context.
func ValidateCast(spell *Spell, currentPOW int, currentLP int, inCombat bool, isDead bool) CastResult {
	result := CastResult{
		Success:   false,
		FFRFailed: false,
	}

	// Check death-only restriction (RESURRECTION)
	if spell.DeathOnly && !isDead {
		result.Message = "This spell can only be cast when you are dead"
		return result
	}

	// Check combat-only restriction
	if spell.CombatOnly && !inCombat {
		result.Message = "This spell can only be cast during combat"
		return result
	}

	// Check non-combat restriction (CRYPT, RETRACE)
	if !spell.CombatOnly && inCombat && (spell.Name == "CRYPT" || spell.Name == "RETRACE") {
		result.Message = "This spell cannot be cast during combat"
		return result
	}

	// Check power affordability
	if !CanAffordSpell(currentPOW, spell.PowerCost) {
		sacrificeNeeded := CalculateSacrificeNeeded(currentPOW, spell.PowerCost)
		if !CanSacrificeLP(currentLP, sacrificeNeeded) {
			result.InsufficientPower = true
			result.Message = fmt.Sprintf("Insufficient POWER (%d/%d). Cannot sacrifice %d LP without dying.", currentPOW, spell.PowerCost, sacrificeNeeded)
			return result
		}
		// Sacrifice is possible
		result.RequiresSacrifice = true
		result.SacrificeAmount = sacrificeNeeded
		result.Message = fmt.Sprintf("Insufficient POWER (%d/%d). Sacrifice %d LP for %d POW?", currentPOW, spell.PowerCost, sacrificeNeeded, sacrificeNeeded)
		return result
	}

	// All validations passed
	result.Success = true
	return result
}

// PerformCast executes the spell casting after validation.
// Returns the result including FFR check outcome.
func PerformCast(spell *Spell, roller dice.Roller) CastResult {
	result := CastResult{
		PowerSpent: spell.PowerCost,
	}

	// Perform Fundamental Failure Rate check
	ffrSuccess, ffrRoll := FundamentalFailureRate(roller)
	if !ffrSuccess {
		result.Success = false
		result.FFRFailed = true
		result.Message = fmt.Sprintf("The spell fizzles and fails! (rolled %d, needed 6+)", ffrRoll)
		return result
	}

	// Spell succeeded
	result.Success = true
	result.Message = fmt.Sprintf("Spell cast successfully! (FFR roll: %d)", ffrRoll)
	return result
}
