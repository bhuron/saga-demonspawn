// Package items defines the game's item system including weapons, armor, and special items.
package items

// ItemType represents the category of an item.
type ItemType string

const (
	// ItemTypeWeapon represents weapons that can be equipped and used in combat.
	ItemTypeWeapon ItemType = "weapon"
	// ItemTypeArmor represents armor that provides damage reduction.
	ItemTypeArmor ItemType = "armor"
	// ItemTypeShield represents shields that provide additional protection.
	ItemTypeShield ItemType = "shield"
	// ItemTypeSpecial represents unique items with custom mechanics.
	ItemTypeSpecial ItemType = "special"
	// ItemTypeConsumable represents items that can be used and consumed.
	ItemTypeConsumable ItemType = "consumable"
)

// Weapon represents a weapon item that can deal damage.
type Weapon struct {
	// Name is the display name of the weapon
	Name string
	// DamageBonus is added to damage calculations
	DamageBonus int
	// Description provides flavor text and usage notes
	Description string
	// Special indicates if this weapon has special rules (e.g., Doombringer)
	Special bool
}

// Armor represents armor that provides damage reduction.
type Armor struct {
	// Name is the display name of the armor
	Name string
	// Protection is the damage reduction value
	Protection int
	// Description provides flavor text
	Description string
}

// Shield represents a shield for additional protection.
type Shield struct {
	// Name is the display name
	Name string
	// Protection is the damage reduction (affected by armor combination)
	Protection int
	// ProtectionWithArmor is reduced protection when worn with armor
	ProtectionWithArmor int
	// Description provides flavor text
	Description string
}

// Special item names as constants for easy reference
const (
	HealingStoneName = "Healing Stone"
	DoombringerName  = "Doombringer"
	TheOrbName       = "The Orb"
)

// Predefined weapons following the game rules
var (
	// WeaponArrow is a ranged single-use weapon
	WeaponArrow = Weapon{
		Name:        "Arrow",
		DamageBonus: 10,
		Description: "Ranged, single use",
		Special:     false,
	}

	// WeaponAxe is a standard melee weapon
	WeaponAxe = Weapon{
		Name:        "Axe",
		DamageBonus: 15,
		Description: "Standard melee weapon",
		Special:     false,
	}

	// WeaponClub is a basic melee weapon
	WeaponClub = Weapon{
		Name:        "Club",
		DamageBonus: 8,
		Description: "Basic melee weapon",
		Special:     false,
	}

	// WeaponDagger is a light concealable weapon
	WeaponDagger = Weapon{
		Name:        "Dagger",
		DamageBonus: 5,
		Description: "Light, concealable",
		Special:     false,
	}

	// WeaponFlail is a standard melee weapon
	WeaponFlail = Weapon{
		Name:        "Flail",
		DamageBonus: 7,
		Description: "Standard melee weapon",
		Special:     false,
	}

	// WeaponHalberd is a two-handed weapon
	WeaponHalberd = Weapon{
		Name:        "Halberd",
		DamageBonus: 12,
		Description: "Two-handed weapon",
		Special:     false,
	}

	// WeaponLance is used mounted or for charges
	WeaponLance = Weapon{
		Name:        "Lance",
		DamageBonus: 12,
		Description: "Mounted/charge weapon",
		Special:     false,
	}

	// WeaponMace is a heavy melee weapon
	WeaponMace = Weapon{
		Name:        "Mace",
		DamageBonus: 14,
		Description: "Heavy melee weapon",
		Special:     false,
	}

	// WeaponSpear can be thrown
	WeaponSpear = Weapon{
		Name:        "Spear",
		DamageBonus: 12,
		Description: "Can be thrown",
		Special:     false,
	}

	// WeaponSword is the standard adventurer's weapon
	WeaponSword = Weapon{
		Name:        "Sword",
		DamageBonus: 10,
		Description: "Standard melee weapon",
		Special:     false,
	}

	// WeaponDoombringer is the cursed legendary axe
	WeaponDoombringer = Weapon{
		Name:        DoombringerName,
		DamageBonus: 20,
		Description: "Cursed blade: -10 LP per attack, heal on hit",
		Special:     true,
	}
)

// Predefined armor following the game rules
var (
	// ArmorNone represents no armor equipped
	ArmorNone = Armor{
		Name:        "None",
		Protection:  0,
		Description: "No armor equipped",
	}

	// ArmorLeather provides basic protection
	ArmorLeather = Armor{
		Name:        "Leather Armor",
		Protection:  5,
		Description: "Light armor, no movement penalty",
	}

	// ArmorChain provides moderate protection
	ArmorChain = Armor{
		Name:        "Chain Mail",
		Protection:  8,
		Description: "Medium armor, no movement penalty",
	}

	// ArmorPlate provides heavy protection
	ArmorPlate = Armor{
		Name:        "Plate Mail",
		Protection:  12,
		Description: "Heavy armor, no movement penalty",
	}
)

// Predefined shield
var (
	// ShieldStandard is the standard shield
	ShieldStandard = Shield{
		Name:                "Shield",
		Protection:          7,
		ProtectionWithArmor: 5,
		Description:         "Protection varies if worn with armor",
	}
)

// AllWeapons returns a slice of all available weapons for selection.
func AllWeapons() []Weapon {
	return []Weapon{
		WeaponArrow,
		WeaponAxe,
		WeaponClub,
		WeaponDagger,
		WeaponFlail,
		WeaponHalberd,
		WeaponLance,
		WeaponMace,
		WeaponSpear,
		WeaponSword,
		WeaponDoombringer,
	}
}

// StartingWeapons returns weapons available at character creation.
func StartingWeapons() []Weapon {
	return []Weapon{
		WeaponSword,
		WeaponDagger,
		WeaponClub,
	}
}

// AllArmor returns a slice of all available armor for selection.
func AllArmor() []Armor {
	return []Armor{
		ArmorNone,
		ArmorLeather,
		ArmorChain,
		ArmorPlate,
	}
}

// StartingArmor returns armor available at character creation.
func StartingArmor() []Armor {
	return []Armor{
		ArmorNone,
		ArmorLeather,
	}
}

// GetWeaponByName finds a weapon by name, returns nil if not found.
func GetWeaponByName(name string) *Weapon {
	for _, w := range AllWeapons() {
		if w.Name == name {
			return &w
		}
	}
	return nil
}

// GetArmorByName finds armor by name, returns nil if not found.
func GetArmorByName(name string) *Armor {
	for _, a := range AllArmor() {
		if a.Name == name {
			return &a
		}
	}
	return nil
}
