package dice

import (
	"testing"
)

// TestRoll2D6Range verifies that Roll2D6 produces values in the expected range (2-12).
func TestRoll2D6Range(t *testing.T) {
	roller := NewStandardRoller()

	// Run multiple rolls to ensure we're within range
	for i := 0; i < 100; i++ {
		result := roller.Roll2D6()
		if result < 2 || result > 12 {
			t.Errorf("Roll2D6() = %d; want value between 2 and 12", result)
		}
	}
}

// TestRoll1D6Range verifies that Roll1D6 produces values in the expected range (1-6).
func TestRoll1D6Range(t *testing.T) {
	roller := NewStandardRoller()

	for i := 0; i < 100; i++ {
		result := roller.Roll1D6()
		if result < 1 || result > 6 {
			t.Errorf("Roll1D6() = %d; want value between 1 and 6", result)
		}
	}
}

// TestRollCharacteristicRange verifies characteristic rolls are in range (16-96).
func TestRollCharacteristicRange(t *testing.T) {
	roller := NewStandardRoller()

	for i := 0; i < 100; i++ {
		result := roller.RollCharacteristic()
		if result < 16 || result > 96 {
			t.Errorf("RollCharacteristic() = %d; want value between 16 and 96", result)
		}
		// Should be multiples of 8 (since it's 2d6 × 8)
		if result%8 != 0 {
			t.Errorf("RollCharacteristic() = %d; want value divisible by 8", result)
		}
	}
}

// TestSeededRollerDeterminism verifies that seeded rollers produce deterministic results.
func TestSeededRollerDeterminism(t *testing.T) {
	seed := int64(12345)

	roller1 := NewSeededRoller(seed)
	roller2 := NewSeededRoller(seed)

	// Both rollers should produce identical sequences
	for i := 0; i < 10; i++ {
		result1 := roller1.Roll2D6()
		result2 := roller2.Roll2D6()

		if result1 != result2 {
			t.Errorf("Seeded rollers diverged at roll %d: roller1=%d, roller2=%d", i, result1, result2)
		}
	}
}

// TestSetSeed verifies that SetSeed changes the random sequence.
func TestSetSeed(t *testing.T) {
	roller := NewStandardRoller()
	seed := int64(99999)

	// Record initial sequence
	roller.SetSeed(seed)
	first1 := roller.Roll2D6()
	first2 := roller.Roll2D6()

	// Reset seed and verify we get the same sequence
	roller.SetSeed(seed)
	second1 := roller.Roll2D6()
	second2 := roller.Roll2D6()

	if first1 != second1 || first2 != second2 {
		t.Errorf("SetSeed() did not reset sequence: first=(%d,%d), second=(%d,%d)",
			first1, first2, second1, second2)
	}
}

// TestRollResult verifies RollResult construction.
func TestRollResult(t *testing.T) {
	result := NewRollResult(7, "Initiative", "rolled 3 + 4")

	if result.Value != 7 {
		t.Errorf("RollResult.Value = %d; want 7", result.Value)
	}
	if result.Description != "Initiative" {
		t.Errorf("RollResult.Description = %q; want %q", result.Description, "Initiative")
	}
	if result.Details != "rolled 3 + 4" {
		t.Errorf("RollResult.Details = %q; want %q", result.Details, "rolled 3 + 4")
	}
}

// TestRollDistribution is a table-driven test to verify reasonable distribution.
// While we can't guarantee perfect distribution in a small sample, we check
// that all possible values appear over many rolls.
func TestRollDistribution(t *testing.T) {
	tests := []struct {
		name     string
		rollFunc func(Roller) int
		minVal   int
		maxVal   int
		samples  int
	}{
		{"Roll2D6", func(r Roller) int { return r.Roll2D6() }, 2, 12, 1000},
		{"Roll1D6", func(r Roller) int { return r.Roll1D6() }, 1, 6, 600},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roller := NewStandardRoller()
			seen := make(map[int]bool)

			for i := 0; i < tt.samples; i++ {
				result := tt.rollFunc(roller)
				seen[result] = true
			}

			// Check that we've seen all possible values
			for val := tt.minVal; val <= tt.maxVal; val++ {
				if !seen[val] {
					t.Errorf("%s did not produce value %d in %d samples", tt.name, val, tt.samples)
				}
			}
		})
	}
}

// BenchmarkRoll2D6 benchmarks the Roll2D6 function.
func BenchmarkRoll2D6(b *testing.B) {
	roller := NewStandardRoller()
	for i := 0; i < b.N; i++ {
		roller.Roll2D6()
	}
}

// BenchmarkRollCharacteristic benchmarks the RollCharacteristic function.
func BenchmarkRollCharacteristic(b *testing.B) {
	roller := NewStandardRoller()
	for i := 0; i < b.N; i++ {
		roller.RollCharacteristic()
	}
}
