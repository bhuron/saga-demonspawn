package magic

import (
	"testing"

	"github.com/benoit/saga-demonspawn/internal/dice"
)

// TestNaturalInclinationCheck verifies the natural inclination check mechanics.
func TestNaturalInclinationCheck(t *testing.T) {
	tests := []struct {
		name    string
		seed    int64
		wantMin int
		wantMax int
	}{
		{"valid roll range", 42, 2, 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roller := dice.NewSeededRoller(tt.seed)
			success, roll := NaturalInclinationCheck(roller)
			
			if roll < tt.wantMin || roll > tt.wantMax {
				t.Errorf("NaturalInclinationCheck() roll = %d, want between %d and %d", roll, tt.wantMin, tt.wantMax)
			}
			
			expectedSuccess := roll >= 4
			if success != expectedSuccess {
				t.Errorf("NaturalInclinationCheck() success = %v, want %v (roll was %d)", success, expectedSuccess, roll)
			}
		})
	}
}

// TestCanAffordSpell verifies power cost checking.
func TestCanAffordSpell(t *testing.T) {
	tests := []struct {
		name       string
		currentPOW int
		spellCost  int
		want       bool
	}{
		{"sufficient power", 50, 25, true},
		{"exact power", 25, 25, true},
		{"insufficient power", 10, 25, false},
		{"zero power", 0, 15, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CanAffordSpell(tt.currentPOW, tt.spellCost)
			if got != tt.want {
				t.Errorf("CanAffordSpell(%d, %d) = %v, want %v", tt.currentPOW, tt.spellCost, got, tt.want)
			}
		})
	}
}

// TestCalculateSacrificeNeeded verifies LP sacrifice calculation.
func TestCalculateSacrificeNeeded(t *testing.T) {
	tests := []struct {
		name       string
		currentPOW int
		spellCost  int
		want       int
	}{
		{"no sacrifice needed", 50, 25, 0},
		{"partial sacrifice", 10, 25, 15},
		{"full sacrifice", 0, 50, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateSacrificeNeeded(tt.currentPOW, tt.spellCost)
			if got != tt.want {
				t.Errorf("CalculateSacrificeNeeded(%d, %d) = %d, want %d", tt.currentPOW, tt.spellCost, got, tt.want)
			}
		})
	}
}

// TestCanSacrificeLP verifies LP sacrifice validation.
func TestCanSacrificeLP(t *testing.T) {
	tests := []struct {
		name            string
		currentLP       int
		sacrificeAmount int
		want            bool
	}{
		{"safe sacrifice", 100, 50, true},
		{"borderline sacrifice", 51, 50, true},
		{"fatal sacrifice", 50, 50, false},
		{"clearly fatal", 30, 50, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CanSacrificeLP(tt.currentLP, tt.sacrificeAmount)
			if got != tt.want {
				t.Errorf("CanSacrificeLP(%d, %d) = %v, want %v", tt.currentLP, tt.sacrificeAmount, got, tt.want)
			}
		})
	}
}

// TestFundamentalFailureRate verifies the FFR check mechanics.
func TestFundamentalFailureRate(t *testing.T) {
	tests := []struct {
		name    string
		seed    int64
		wantMin int
		wantMax int
	}{
		{"valid roll range", 123, 2, 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roller := dice.NewSeededRoller(tt.seed)
			success, roll := FundamentalFailureRate(roller)
			
			if roll < tt.wantMin || roll > tt.wantMax {
				t.Errorf("FundamentalFailureRate() roll = %d, want between %d and %d", roll, tt.wantMin, tt.wantMax)
			}
			
			expectedSuccess := roll >= 6
			if success != expectedSuccess {
				t.Errorf("FundamentalFailureRate() success = %v, want %v (roll was %d)", success, expectedSuccess, roll)
			}
		})
	}
}

// TestValidateCast verifies spell casting validation.
func TestValidateCast(t *testing.T) {
	tests := []struct {
		name          string
		spell         *Spell
		currentPOW    int
		currentLP     int
		inCombat      bool
		isDead        bool
		wantSuccess   bool
		wantSacrifice bool
	}{
		{
			name:          "valid cast with sufficient power",
			spell:         &Spell{Name: "FIREBALL", PowerCost: 15, CombatOnly: true},
			currentPOW:    50,
			currentLP:     100,
			inCombat:      true,
			isDead:        false,
			wantSuccess:   true,
			wantSacrifice: false,
		},
		{
			name:          "combat spell outside combat",
			spell:         &Spell{Name: "FIREBALL", PowerCost: 15, CombatOnly: true},
			currentPOW:    50,
			currentLP:     100,
			inCombat:      false,
			isDead:        false,
			wantSuccess:   false,
			wantSacrifice: false,
		},
		{
			name:          "insufficient power but can sacrifice",
			spell:         &Spell{Name: "FIREBALL", PowerCost: 20, CombatOnly: true},
			currentPOW:    10,
			currentLP:     50,
			inCombat:      true,
			isDead:        false,
			wantSuccess:   false,
			wantSacrifice: true,
		},
		{
			name:          "resurrection when alive",
			spell:         &Spell{Name: "RESURRECTION", PowerCost: 50, DeathOnly: true},
			currentPOW:    60,
			currentLP:     50,
			inCombat:      false,
			isDead:        false,
			wantSuccess:   false,
			wantSacrifice: false,
		},
		{
			name:          "resurrection when dead",
			spell:         &Spell{Name: "RESURRECTION", PowerCost: 50, DeathOnly: true},
			currentPOW:    60,
			currentLP:     0,
			inCombat:      false,
			isDead:        true,
			wantSuccess:   true,
			wantSacrifice: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateCast(tt.spell, tt.currentPOW, tt.currentLP, tt.inCombat, tt.isDead)
			
			if result.Success != tt.wantSuccess {
				t.Errorf("ValidateCast() Success = %v, want %v", result.Success, tt.wantSuccess)
			}
			
			if result.RequiresSacrifice != tt.wantSacrifice {
				t.Errorf("ValidateCast() RequiresSacrifice = %v, want %v", result.RequiresSacrifice, tt.wantSacrifice)
			}
		})
	}
}

// TestPerformCast verifies the casting execution with FFR check.
func TestPerformCast(t *testing.T) {
	spell := &Spell{Name: "FIREBALL", PowerCost: 15}
	
	// Test with successful FFR (roll >= 6)
	t.Run("successful cast", func(t *testing.T) {
		roller := dice.NewSeededRoller(100) // Use a seed that gives high roll
		result := PerformCast(spell, roller)
		
		if result.PowerSpent != spell.PowerCost {
			t.Errorf("PerformCast() PowerSpent = %d, want %d", result.PowerSpent, spell.PowerCost)
		}
		
		// Success depends on the roll, so we check the logic
		// If FFRFailed is false, Success should be true
		if result.FFRFailed {
			if result.Success {
				t.Error("PerformCast() Success should be false when FFRFailed is true")
			}
		} else {
			if !result.Success {
				t.Error("PerformCast() Success should be true when FFRFailed is false")
			}
		}
	})
}
