// Package dice provides dice rolling functionality for the game.
// It abstracts random number generation to enable testability through
// deterministic seeding.
package dice

import (
	"math/rand"
	"time"
)

// Roller is an interface for dice rolling operations.
// This abstraction allows for mock implementations in tests.
type Roller interface {
	// Roll2D6 simulates rolling two six-sided dice and returns the sum (2-12).
	Roll2D6() int

	// Roll1D6 simulates rolling one six-sided dice and returns the result (1-6).
	Roll1D6() int

	// RollCharacteristic rolls 2d6 and multiplies by 8 for character stat generation (16-96).
	// This represents percentage values where 100% is unattainable - nobody is perfect!
	RollCharacteristic() int

	// SetSeed sets the random number generator seed for deterministic behavior.
	// Primarily used for testing.
	SetSeed(seed int64)
}

// StandardRoller implements the Roller interface using Go's math/rand.
type StandardRoller struct {
	rng *rand.Rand
}

// NewStandardRoller creates a new StandardRoller with a random seed based on current time.
func NewStandardRoller() *StandardRoller {
	return &StandardRoller{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// NewSeededRoller creates a new StandardRoller with a specific seed.
// Useful for testing or for reproducible game sessions.
func NewSeededRoller(seed int64) *StandardRoller {
	return &StandardRoller{
		rng: rand.New(rand.NewSource(seed)),
	}
}

// Roll2D6 rolls two six-sided dice and returns their sum (2-12).
func (r *StandardRoller) Roll2D6() int {
	return r.rollDie(6) + r.rollDie(6)
}

// Roll1D6 rolls one six-sided die and returns the result (1-6).
func (r *StandardRoller) Roll1D6() int {
	return r.rollDie(6)
}

// RollCharacteristic rolls 2d6 and multiplies by 8 for character statistics.
// This produces values in the range 16-96, representing percentage capabilities.
// Nobody is perfect, so 100% is impossible!
func (r *StandardRoller) RollCharacteristic() int {
	return r.Roll2D6() * 8
}

// SetSeed changes the random number generator's seed.
// This is primarily used for testing to create deterministic sequences.
func (r *StandardRoller) SetSeed(seed int64) {
	r.rng = rand.New(rand.NewSource(seed))
}

// rollDie is a helper function that rolls a single die with the specified number of sides.
// Returns a value from 1 to sides (inclusive).
func (r *StandardRoller) rollDie(sides int) int {
	return r.rng.Intn(sides) + 1
}

// RollResult represents the outcome of a dice roll with context.
// This is useful for logging and displaying roll results to the player.
type RollResult struct {
	// Value is the final result of the roll
	Value int
	// Description explains what the roll was for
	Description string
	// Details provides breakdown (e.g., "rolled 3 + 4 = 7")
	Details string
}

// NewRollResult creates a RollResult with the specified values.
func NewRollResult(value int, description, details string) RollResult {
	return RollResult{
		Value:       value,
		Description: description,
		Details:     details,
	}
}
