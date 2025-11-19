package magic

import (
	"fmt"

	"github.com/benoit/saga-demonspawn/internal/dice"
)

// SpellEffect represents the outcome of applying a spell's effect.
type SpellEffect struct {
	Success        bool
	Message        string
	DamageDealt    int    // For offensive spells
	LPRestored     int    // For healing/restoration
	CombatEnded    bool   // Whether combat should end
	Victory        bool   // Whether combat ended in victory
	EnemyKilled    bool   // Whether enemy was killed
	CharacterDied  bool   // Whether character died (for RESURRECTION)
	RequiresReroll bool   // Whether stats need rerolling (RESURRECTION)
	NavigateTo     string // Section to navigate to (CRYPT, RETRACE)
}

// ApplyARMOUR applies the ARMOUR spell effect.
func ApplyARMOUR() SpellEffect {
	return SpellEffect{
		Success: true,
		Message: "Magical armor of light surrounds you! Incoming damage reduced by 10 points.",
	}
}

// ApplyCRYPT applies the CRYPT spell effect (simplified - just restore POW to max).
func ApplyCRYPT() SpellEffect {
	return SpellEffect{
		Success:    true,
		Message:    "You are transported to the Crypts. Your POWER is fully restored!",
		NavigateTo: "CRYPT",
	}
}

// ApplyFIREBALL applies the FIREBALL spell effect.
func ApplyFIREBALL() SpellEffect {
	return SpellEffect{
		Success:     true,
		Message:     "A ball of flame strikes the enemy!",
		DamageDealt: 50,
	}
}

// ApplyINVISIBILITY applies the INVISIBILITY spell effect.
func ApplyINVISIBILITY(inCombat bool) SpellEffect {
	if inCombat {
		return SpellEffect{
			Success:     true,
			Message:     "You fade from sight. The enemy cannot see you!",
			CombatEnded: true,
			Victory:     true,
		}
	}
	return SpellEffect{
		Success: true,
		Message: "You become invisible, avoiding danger ahead.",
	}
}

// ApplyPARALYSIS applies the PARALYSIS spell effect.
func ApplyPARALYSIS() SpellEffect {
	return SpellEffect{
		Success:     true,
		Message:     "The enemy is paralyzed! You escape combat.",
		CombatEnded: true,
		Victory:     false, // No victory rewards
	}
}

// ApplyPOISONNEEDLE applies the POISON NEEDLE spell effect.
func ApplyPOISONNEEDLE(roller dice.Roller) SpellEffect {
	// Roll 1d6: 1-3 = immune, 4-6 = affected
	roll := roller.Roll1D6()
	
	if roll >= 4 {
		return SpellEffect{
			Success:     true,
			Message:     fmt.Sprintf("The poisoned needle strikes! (rolled %d) The poison is invariably fatal!", roll),
			EnemyKilled: true,
		}
	}
	
	return SpellEffect{
		Success:     true,
		Message:     fmt.Sprintf("The poisoned needle strikes! (rolled %d) But the enemy is immune to the poison.", roll),
		EnemyKilled: false,
	}
}

// ApplyRESURRECTION applies the RESURRECTION spell effect.
func ApplyRESURRECTION() SpellEffect {
	return SpellEffect{
		Success:        true,
		Message:        "Death is not your fate! You are resurrected at the start of this section.",
		RequiresReroll: true,
	}
}

// ApplyRETRACE applies the RETRACE spell effect.
func ApplyRETRACE(sectionName string) SpellEffect {
	return SpellEffect{
		Success:    true,
		Message:    fmt.Sprintf("You trace your steps back to: %s", sectionName),
		NavigateTo: sectionName,
	}
}

// ApplyTIMEWARP applies the TIMEWARP spell effect.
func ApplyTIMEWARP() SpellEffect {
	return SpellEffect{
		Success: true,
		Message: "Time warps around you! You return to the beginning of this section.",
	}
}

// ApplyXENOPHOBIA applies the XENOPHOBIA spell effect.
func ApplyXENOPHOBIA() SpellEffect {
	return SpellEffect{
		Success: true,
		Message: "The enemy is gripped by fear! Their damage is reduced by 5 points.",
	}
}
