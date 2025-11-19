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
	CurrentPOW        int            `json:"current_pow"`         // Current power
	MaximumPOW        int            `json:"maximum_pow"`         // Maximum power
	MagicUnlocked     bool           `json:"magic_unlocked"`      // Whether magic system is available
	ActiveSpellEffects map[string]int `json:"active_spell_effects"` // Active spell buffs/debuffs

	// Equipment
	EquippedWeapon *items.Weapon `json:"equipped_weapon"` // Current weapon
	EquippedArmor  *items.Armor  `json:"equipped_armor"`  // Current armor
	HasShield      bool          `json:"has_shield"`      // Shield equipped

	// Special items (Phase 3)
	HealingStoneCharges  int  `json:"healing_stone_charges"`  // Current Healing Stone charges (max 50)
	DoombringerPossessed bool `json:"doombringer_possessed"` // Whether Doombringer is possessed
	OrbPossessed         bool `json:"orb_possessed"`         // Whether The Orb is possessed
	OrbEquipped          bool `json:"orb_equipped"`          // Whether The Orb is held in left hand
	OrbDestroyed         bool `json:"orb_destroyed"`         // Whether The Orb has been thrown

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
		ActiveSpellEffects: make(map[string]int),
		EquippedWeapon: &items.WeaponSword, // Default starting weapon
		EquippedArmor:  &items.ArmorNone,   // No armor by default
		HasShield:      false,
		// Special items start not possessed (acquired during adventure)
		HealingStoneCharges:  0,     // Will be set to 50 when acquired
		DoombringerPossessed: false,
		OrbPossessed:         false,
		OrbEquipped:          false,
		OrbDestroyed:         false,
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

// AcquireHealingStone gives the character the Healing Stone with full charges.
func (c *Character) AcquireHealingStone() {
	c.HealingStoneCharges = 50
}

// RechargeHealingStone recharges the Healing Stone to full capacity.
func (c *Character) RechargeHealingStone() error {
	if c.HealingStoneCharges >= 50 {
		return fmt.Errorf("healing stone is already fully charged")
	}
	c.HealingStoneCharges = 50
	return nil
}

// UseHealingStone heals the character and depletes the stone.
// Returns the amount healed and any error.
func (c *Character) UseHealingStone(healAmount int) (int, error) {
	if c.HealingStoneCharges <= 0 {
		return 0, fmt.Errorf("healing stone is depleted")
	}
	if c.CurrentLP >= c.MaximumLP {
		return 0, fmt.Errorf("already at full health")
	}

	// Calculate actual healing (capped at max LP and available charges)
	actualHeal := healAmount
	if actualHeal > c.HealingStoneCharges {
		actualHeal = c.HealingStoneCharges
	}
	if c.CurrentLP+actualHeal > c.MaximumLP {
		actualHeal = c.MaximumLP - c.CurrentLP
	}

	// Apply healing and deplete charges
	c.CurrentLP += actualHeal
	c.HealingStoneCharges -= healAmount // Always deplete by roll amount
	if c.HealingStoneCharges < 0 {
		c.HealingStoneCharges = 0
	}

	return actualHeal, nil
}

// AcquireDoombringer gives the character Doombringer.
func (c *Character) AcquireDoombringer() {
	c.DoombringerPossessed = true
}

// AcquireOrb gives the character The Orb.
func (c *Character) AcquireOrb() {
	c.OrbPossessed = true
	c.OrbDestroyed = false
	c.OrbEquipped = false
}

// DestroyOrb marks The Orb as destroyed (after throwing).
func (c *Character) DestroyOrb() {
	c.OrbDestroyed = true
	c.OrbEquipped = false
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

// validateSpecialItems validates special item state consistency.
func validateSpecialItems(c *Character) error {
	// Healing Stone charges must be 0-50
	if c.HealingStoneCharges < 0 {
		return fmt.Errorf("healing stone charges cannot be negative: %d", c.HealingStoneCharges)
	}
	if c.HealingStoneCharges > 50 {
		return fmt.Errorf("healing stone charges exceed maximum: %d", c.HealingStoneCharges)
	}
	
	// The Orb cannot be both equipped and destroyed
	if c.OrbEquipped && c.OrbDestroyed {
		return fmt.Errorf("the orb cannot be both equipped and destroyed")
	}
	
	// Cannot equip The Orb if not possessed
	if c.OrbEquipped && !c.OrbPossessed {
		return fmt.Errorf("cannot equip orb that is not possessed")
	}
	
	// Cannot destroy The Orb if not possessed
	if c.OrbDestroyed && !c.OrbPossessed {
		return fmt.Errorf("cannot destroy orb that is not possessed")
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
	
	// Initialize map if nil (backward compatibility)
	if char.ActiveSpellEffects == nil {
		char.ActiveSpellEffects = make(map[string]int)
	}
	
	// Validate special item state
	if err := validateSpecialItems(&char); err != nil {
		return nil, fmt.Errorf("invalid special item state: %w", err)
	}

	return &char, nil
}

// AddSpellEffect adds or updates an active spell effect.
func (c *Character) AddSpellEffect(effectName string, value int) {
	if c.ActiveSpellEffects == nil {
		c.ActiveSpellEffects = make(map[string]int)
	}
	c.ActiveSpellEffects[effectName] = value
}

// RemoveSpellEffect removes an active spell effect.
func (c *Character) RemoveSpellEffect(effectName string) {
	delete(c.ActiveSpellEffects, effectName)
}

// GetSpellEffect returns the value of an active spell effect, or 0 if not present.
func (c *Character) GetSpellEffect(effectName string) int {
	return c.ActiveSpellEffects[effectName]
}

// HasSpellEffect checks if a spell effect is active.
func (c *Character) HasSpellEffect(effectName string) bool {
	_, exists := c.ActiveSpellEffects[effectName]
	return exists
}

// ClearAllSpellEffects removes all active spell effects.
func (c *Character) ClearAllSpellEffects() {
	c.ActiveSpellEffects = make(map[string]int)
}
