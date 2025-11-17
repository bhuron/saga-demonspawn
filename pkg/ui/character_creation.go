package ui

import (
	"github.com/benoit/saga-demonspawn/internal/character"
	"github.com/benoit/saga-demonspawn/internal/dice"
	"github.com/benoit/saga-demonspawn/internal/items"
)

// CreationStep represents the current step in character creation.
type CreationStep int

const (
	// StepRollCharacteristics is the first step where stats are rolled
	StepRollCharacteristics CreationStep = iota
	// StepSelectEquipment is where starting weapon and armor are chosen
	StepSelectEquipment
	// StepReviewCharacter is the final confirmation screen
	StepReviewCharacter
)

// CharacterCreationModel represents the character creation flow state.
type CharacterCreationModel struct {
	dice dice.Roller
	step CreationStep

	// Rolled characteristics
	strength   int
	speed      int
	stamina    int
	courage    int
	luck       int
	charm      int
	attraction int
	allRolled  bool

	// Equipment selection
	weaponCursor int
	armorCursor  int
	weaponOptions []items.Weapon
	armorOptions  []items.Armor

	// Created character
	character *character.Character
}

// NewCharacterCreationModel creates a new character creation model.
func NewCharacterCreationModel(roller dice.Roller) CharacterCreationModel {
	return CharacterCreationModel{
		dice:          roller,
		step:          StepRollCharacteristics,
		strength:      0,
		speed:         0,
		stamina:       0,
		courage:       0,
		luck:          0,
		charm:         0,
		attraction:    0,
		allRolled:     false,
		weaponCursor:  0,
		armorCursor:   0,
		weaponOptions: items.StartingWeapons(),
		armorOptions:  items.StartingArmor(),
		character:     nil,
	}
}

// GetStep returns the current creation step.
func (m *CharacterCreationModel) GetStep() CreationStep {
	return m.step
}

// RollStrength rolls the Strength characteristic.
func (m *CharacterCreationModel) RollStrength() int {
	m.strength = m.dice.RollCharacteristic()
	m.checkAllRolled()
	return m.strength
}

// RollSpeed rolls the Speed characteristic.
func (m *CharacterCreationModel) RollSpeed() int {
	m.speed = m.dice.RollCharacteristic()
	m.checkAllRolled()
	return m.speed
}

// RollStamina rolls the Stamina characteristic.
func (m *CharacterCreationModel) RollStamina() int {
	m.stamina = m.dice.RollCharacteristic()
	m.checkAllRolled()
	return m.stamina
}

// RollCourage rolls the Courage characteristic.
func (m *CharacterCreationModel) RollCourage() int {
	m.courage = m.dice.RollCharacteristic()
	m.checkAllRolled()
	return m.courage
}

// RollLuck rolls the Luck characteristic.
func (m *CharacterCreationModel) RollLuck() int {
	m.luck = m.dice.RollCharacteristic()
	m.checkAllRolled()
	return m.luck
}

// RollCharm rolls the Charm characteristic.
func (m *CharacterCreationModel) RollCharm() int {
	m.charm = m.dice.RollCharacteristic()
	m.checkAllRolled()
	return m.charm
}

// RollAttraction rolls the Attraction characteristic.
func (m *CharacterCreationModel) RollAttraction() int {
	m.attraction = m.dice.RollCharacteristic()
	m.checkAllRolled()
	return m.attraction
}

// RollAll rolls all characteristics at once.
func (m *CharacterCreationModel) RollAll() {
	m.RollStrength()
	m.RollSpeed()
	m.RollStamina()
	m.RollCourage()
	m.RollLuck()
	m.RollCharm()
	m.RollAttraction()
}

// checkAllRolled checks if all characteristics have been rolled.
func (m *CharacterCreationModel) checkAllRolled() {
	m.allRolled = m.strength > 0 && m.speed > 0 && m.stamina > 0 &&
		m.courage > 0 && m.luck > 0 && m.charm > 0 && m.attraction > 0
}

// AreAllRolled returns true if all characteristics have been rolled.
func (m *CharacterCreationModel) AreAllRolled() bool {
	return m.allRolled
}

// GetCharacteristics returns all rolled characteristics.
func (m *CharacterCreationModel) GetCharacteristics() (int, int, int, int, int, int, int) {
	return m.strength, m.speed, m.stamina, m.courage, m.luck, m.charm, m.attraction
}

// GetCalculatedLP returns the sum of all characteristics.
func (m *CharacterCreationModel) GetCalculatedLP() int {
	return m.strength + m.speed + m.stamina + m.courage + m.luck + m.charm + m.attraction
}

// NextStep advances to the next creation step.
func (m *CharacterCreationModel) NextStep() {
	m.step++
}

// PreviousStep goes back to the previous step.
func (m *CharacterCreationModel) PreviousStep() {
	if m.step > StepRollCharacteristics {
		m.step--
	}
}

// MoveWeaponCursorUp moves the weapon selection cursor up.
func (m *CharacterCreationModel) MoveWeaponCursorUp() {
	if m.weaponCursor > 0 {
		m.weaponCursor--
	}
}

// MoveWeaponCursorDown moves the weapon selection cursor down.
func (m *CharacterCreationModel) MoveWeaponCursorDown() {
	if m.weaponCursor < len(m.weaponOptions)-1 {
		m.weaponCursor++
	}
}

// MoveArmorCursorUp moves the armor selection cursor up.
func (m *CharacterCreationModel) MoveArmorCursorUp() {
	if m.armorCursor > 0 {
		m.armorCursor--
	}
}

// MoveArmorCursorDown moves the armor selection cursor down.
func (m *CharacterCreationModel) MoveArmorCursorDown() {
	if m.armorCursor < len(m.armorOptions)-1 {
		m.armorCursor++
	}
}

// GetSelectedWeapon returns the currently selected weapon.
func (m *CharacterCreationModel) GetSelectedWeapon() *items.Weapon {
	if m.weaponCursor < len(m.weaponOptions) {
		return &m.weaponOptions[m.weaponCursor]
	}
	return nil
}

// GetSelectedArmor returns the currently selected armor.
func (m *CharacterCreationModel) GetSelectedArmor() *items.Armor {
	if m.armorCursor < len(m.armorOptions) {
		return &m.armorOptions[m.armorCursor]
	}
	return nil
}

// GetWeaponOptions returns all available weapons.
func (m *CharacterCreationModel) GetWeaponOptions() []items.Weapon {
	return m.weaponOptions
}

// GetArmorOptions returns all available armor.
func (m *CharacterCreationModel) GetArmorOptions() []items.Armor {
	return m.armorOptions
}

// GetWeaponCursor returns the weapon cursor position.
func (m *CharacterCreationModel) GetWeaponCursor() int {
	return m.weaponCursor
}

// GetArmorCursor returns the armor cursor position.
func (m *CharacterCreationModel) GetArmorCursor() int {
	return m.armorCursor
}

// CreateCharacter finalizes the character creation with selected equipment.
func (m *CharacterCreationModel) CreateCharacter() (*character.Character, error) {
	char, err := character.New(
		m.strength,
		m.speed,
		m.stamina,
		m.courage,
		m.luck,
		m.charm,
		m.attraction,
	)
	if err != nil {
		return nil, err
	}

	// Equip selected items
	selectedWeapon := m.GetSelectedWeapon()
	if selectedWeapon != nil {
		char.EquipWeapon(selectedWeapon)
	}

	selectedArmor := m.GetSelectedArmor()
	if selectedArmor != nil {
		char.EquipArmor(selectedArmor)
	}

	m.character = char
	return char, nil
}

// GetCharacter returns the created character.
func (m *CharacterCreationModel) GetCharacter() *character.Character {
	return m.character
}

// Reset resets the character creation to start over.
func (m *CharacterCreationModel) Reset() {
	m.step = StepRollCharacteristics
	m.strength = 0
	m.speed = 0
	m.stamina = 0
	m.courage = 0
	m.luck = 0
	m.charm = 0
	m.attraction = 0
	m.allRolled = false
	m.weaponCursor = 0
	m.armorCursor = 0
	m.character = nil
}
