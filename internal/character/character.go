// Package character manages character state and operations for Fire*Wolf.
package character

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/benoit/saga-demonspawn/internal/items"
)

// Character represents Fire*Wolf with all stats, equipment, and progress.
type Character struct {
	// Core characteristics (rolled at creation)
	Strength   int `json:"strength"`   // STR: Physical power
	Speed      int `json:"speed"`      // SPD: Agility and reaction
	Stamina    int `json:"stamina"`    // STA: Endurance
	Courage    int `json:"courage"`    // CRG: Bravery
	Luck       int `json:"luck"`       // LCK: Fortune
	Charm      int `json:"charm"`      // CHM: Charisma
	Attraction int `json:"attraction"` // ATT: Personal magnetism

	// Derived values
	CurrentLP int `json:"current_lp"` // Current life points
	MaximumLP int `json:"maximum_lp"` // Maximum life points
	Skill     int `json:"skill"`      // SKL: Combat proficiency

	// Magic system (unlocked during adventure)
	CurrentPOW int  `json:"current_pow"` // Current power
	MaximumPOW int  `json:"maximum_pow"` // Maximum power
	MagicUnlocked bool `json:"magic_unlocked"` // Whether magic system is available

	// Equipment
	EquippedWeapon *items.Weapon `json:"equipped_weapon"` // Current weapon
	EquippedArmor  *items.Armor  `json:"equipped_armor"`  // Current armor
	HasShield      bool          `json:"has_shield"`      // Shield equipped

	// Progress tracking
	EnemiesDefeated int       `json:"enemies_defeated"` // Total enemies killed
	CreatedAt       time.Time `json:"created_at"`       // Character creation timestamp
	LastSaved       time.Time `json:"last_saved"`       // Last save timestamp
}

// New creates a new character with the specified characteristics.
// Life points are calculated as the sum of all characteristics.
// Skill and Power start at 0.
func New(str, spd, sta, crg, lck, chm, att int) (*Character, error) {
	// Validate characteristics are within reasonable bounds
	if err := validateCharacteristic("Strength", str); err != nil {
		return nil, err
	}
	if err := validateCharacteristic("Speed", spd); err != nil {
		return nil, err
	}
	if err := validateCharacteristic("Stamina", sta); err != nil {
		return nil, err
	}
	if err := validateCharacteristic("Courage", crg); err != nil {
		return nil, err
	}
	if err := validateCharacteristic("Luck", lck); err != nil {
		return nil, err
	}
	if err := validateCharacteristic("Charm", chm); err != nil {
		return nil, err
	}
	if err := validateCharacteristic("Attraction", att); err != nil {
		return nil, err
	}

	// Calculate maximum LP as sum of all characteristics
	maxLP := str + spd + sta + crg + lck + chm + att

	char := &Character{
		Strength:   str,
		Speed:      spd,
		Stamina:    sta,
		Courage:    crg,
		Luck:       lck,
		Charm:      chm,
		Attraction: att,
		CurrentLP:  maxLP,
		MaximumLP:  maxLP,
		Skill:      0,
		CurrentPOW: 0,
		MaximumPOW: 0,
		MagicUnlocked: false,
		EquippedWeapon: &items.WeaponSword, // Default starting weapon
		EquippedArmor:  &items.ArmorNone,   // No armor by default
		HasShield:      false,
		EnemiesDefeated: 0,
		CreatedAt:      time.Now(),
		LastSaved:      time.Now(),
	}

	return char, nil
}

// validateCharacteristic checks if a characteristic value is reasonable.
// Allows 0-999 but warns about unusual values.
func validateCharacteristic(name string, value int) error {
	if value < 0 {
		return fmt.Errorf("%s cannot be negative: %d", name, value)
	}
	if value > 999 {
		return fmt.Errorf("%s exceeds maximum (999): %d", name, value)
	}
	return nil
}

// ModifyStrength changes the Strength characteristic by the specified amount.
// Returns an error if the result would be negative.
func (c *Character) ModifyStrength(delta int) error {
	newVal := c.Strength + delta
	if newVal < 0 {
		return fmt.Errorf("strength cannot be negative (would be %d)", newVal)
	}
	c.Strength = newVal
	return nil
}

// ModifySpeed changes the Speed characteristic.
func (c *Character) ModifySpeed(delta int) error {
	newVal := c.Speed + delta
	if newVal < 0 {
		return fmt.Errorf("speed cannot be negative (would be %d)", newVal)
	}
	c.Speed = newVal
	return nil
}

// ModifyStamina changes the Stamina characteristic.
func (c *Character) ModifyStamina(delta int) error {
	newVal := c.Stamina + delta
	if newVal < 0 {
		return fmt.Errorf("stamina cannot be negative (would be %d)", newVal)
	}
	c.Stamina = newVal
	return nil
}

// ModifyCourage changes the Courage characteristic.
func (c *Character) ModifyCourage(delta int) error {
	newVal := c.Courage + delta
	if newVal < 0 {
		return fmt.Errorf("courage cannot be negative (would be %d)", newVal)
	}
	c.Courage = newVal
	return nil
}

// ModifyLuck changes the Luck characteristic.
func (c *Character) ModifyLuck(delta int) error {
	newVal := c.Luck + delta
	if newVal < 0 {
		return fmt.Errorf("luck cannot be negative (would be %d)", newVal)
	}
	c.Luck = newVal
	return nil
}

// ModifyCharm changes the Charm characteristic.
func (c *Character) ModifyCharm(delta int) error {
	newVal := c.Charm + delta
	if newVal < 0 {
		return fmt.Errorf("charm cannot be negative (would be %d)", newVal)
	}
	c.Charm = newVal
	return nil
}

// ModifyAttraction changes the Attraction characteristic.
func (c *Character) ModifyAttraction(delta int) error {
	newVal := c.Attraction + delta
	if newVal < 0 {
		return fmt.Errorf("attraction cannot be negative (would be %d)", newVal)
	}
	c.Attraction = newVal
	return nil
}

// ModifyLP changes current life points. Can go negative (death).
func (c *Character) ModifyLP(delta int) {
	c.CurrentLP += delta
}

// SetLP sets current life points to a specific value.
func (c *Character) SetLP(value int) {
	c.CurrentLP = value
}

// SetMaxLP sets maximum life points to a specific value.
func (c *Character) SetMaxLP(value int) error {
	if value < 0 {
		return fmt.Errorf("maximum LP cannot be negative: %d", value)
	}
	c.MaximumLP = value
	return nil
}

// ModifySkill changes skill level.
func (c *Character) ModifySkill(delta int) error {
	newVal := c.Skill + delta
	if newVal < 0 {
		return fmt.Errorf("skill cannot be negative (would be %d)", newVal)
	}
	c.Skill = newVal
	return nil
}

// SetSkill sets skill to a specific value.
func (c *Character) SetSkill(value int) error {
	if value < 0 {
		return fmt.Errorf("skill cannot be negative: %d", value)
	}
	c.Skill = value
	return nil
}

// UnlockMagic activates the magic system and sets initial POW.
func (c *Character) UnlockMagic(initialPOW int) error {
	if initialPOW < 0 {
		return fmt.Errorf("initial POW cannot be negative: %d", initialPOW)
	}
	c.MagicUnlocked = true
	c.CurrentPOW = initialPOW
	c.MaximumPOW = initialPOW
	return nil
}

// ModifyPOW changes current power.
func (c *Character) ModifyPOW(delta int) {
	c.CurrentPOW += delta
	if c.CurrentPOW < 0 {
		c.CurrentPOW = 0
	}
}

// SetPOW sets current power to a specific value.
func (c *Character) SetPOW(value int) {
	if value < 0 {
		value = 0
	}
	c.CurrentPOW = value
}

// SetMaxPOW sets maximum power.
func (c *Character) SetMaxPOW(value int) error {
	if value < 0 {
		return fmt.Errorf("maximum POW cannot be negative: %d", value)
	}
	c.MaximumPOW = value
	return nil
}

// EquipWeapon changes the equipped weapon.
func (c *Character) EquipWeapon(weapon *items.Weapon) {
	c.EquippedWeapon = weapon
}

// EquipArmor changes the equipped armor.
func (c *Character) EquipArmor(armor *items.Armor) {
	c.EquippedArmor = armor
}

// ToggleShield equips or unequips the shield.
func (c *Character) ToggleShield() {
	c.HasShield = !c.HasShield
}

// IncrementEnemiesDefeated adds one to the enemies defeated counter.
// This is typically called after winning combat.
func (c *Character) IncrementEnemiesDefeated() {
	c.EnemiesDefeated++
}

// IsAlive returns true if current LP is greater than 0.
func (c *Character) IsAlive() bool {
	return c.CurrentLP > 0
}

// GetArmorProtection returns the total damage reduction from armor and shield.
func (c *Character) GetArmorProtection() int {
	protection := 0

	if c.EquippedArmor != nil {
		protection += c.EquippedArmor.Protection
	}

	if c.HasShield {
		// Shield protection is reduced when worn with armor
		if c.EquippedArmor != nil && c.EquippedArmor.Name != "None" {
			protection += items.ShieldStandard.ProtectionWithArmor
		} else {
			protection += items.ShieldStandard.Protection
		}
	}

	return protection
}

// GetWeaponDamageBonus returns the damage bonus from the equipped weapon.
func (c *Character) GetWeaponDamageBonus() int {
	if c.EquippedWeapon != nil {
		return c.EquippedWeapon.DamageBonus
	}
	return 0
}

// Save saves the character to a JSON file in the specified directory.
// The filename includes a timestamp for versioning.
func (c *Character) Save(directory string) error {
	c.LastSaved = time.Now()

	// Create filename with timestamp
	timestamp := c.LastSaved.Format("20060102-150405")
	filename := fmt.Sprintf("character_%s.json", timestamp)
	filepath := filepath.Join(directory, filename)

	// Ensure directory exists
	if err := os.MkdirAll(directory, 0755); err != nil {
		return fmt.Errorf("failed to create save directory: %w", err)
	}

	// Marshal character to JSON with indentation for readability
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal character: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write save file: %w", err)
	}

	return nil
}

// Load loads a character from a JSON file.
func Load(filepath string) (*Character, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read save file: %w", err)
	}

	var char Character
	if err := json.Unmarshal(data, &char); err != nil {
		return nil, fmt.Errorf("failed to unmarshal character: %w", err)
	}

	return &char, nil
}
