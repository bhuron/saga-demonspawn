package magic

// SpellCategory represents the type of spell.
type SpellCategory string

const (
	CategoryOffensive SpellCategory = "Offensive"
	CategoryDefensive SpellCategory = "Defensive"
	CategoryNavigation SpellCategory = "Navigation"
	CategoryTactical  SpellCategory = "Tactical"
	CategoryRecovery  SpellCategory = "Recovery"
)

// Spell represents a magic spell with its properties and constraints.
type Spell struct {
	Name        string
	PowerCost   int
	Description string
	Category    SpellCategory
	CombatOnly  bool // Whether spell can only be cast during combat
	DeathOnly   bool // Whether spell can only be cast when LP <= 0 (RESURRECTION)
}

// AllSpells contains all 10 spells from the ruleset.
var AllSpells = []Spell{
	{
		Name:        "ARMOUR",
		PowerCost:   25,
		Description: "Creates magical armor. Reduces incoming damage by 10 points.",
		Category:    CategoryDefensive,
		CombatOnly:  false,
		DeathOnly:   false,
	},
	{
		Name:        "CRYPT",
		PowerCost:   150,
		Description: "Returns you to the Crypts for POWER restoration.",
		Category:    CategoryNavigation,
		CombatOnly:  false,
		DeathOnly:   false,
	},
	{
		Name:        "FIREBALL",
		PowerCost:   15,
		Description: "Hurls a ball of flame. Deals 50 LP damage to enemy.",
		Category:    CategoryOffensive,
		CombatOnly:  true,
		DeathOnly:   false,
	},
	{
		Name:        "INVISIBILITY",
		PowerCost:   30,
		Description: "Renders you invisible. Avoid combat and proceed as if victorious.",
		Category:    CategoryTactical,
		CombatOnly:  false,
		DeathOnly:   false,
	},
	{
		Name:        "PARALYSIS",
		PowerCost:   30,
		Description: "Paralyzes enemy. Escape combat immediately without victory.",
		Category:    CategoryTactical,
		CombatOnly:  true,
		DeathOnly:   false,
	},
	{
		Name:        "POISON NEEDLE",
		PowerCost:   25,
		Description: "Shoots poisoned needle. Roll 1d6: 4-6 kills enemy, 1-3 immune.",
		Category:    CategoryOffensive,
		CombatOnly:  true,
		DeathOnly:   false,
	},
	{
		Name:        "RESURRECTION",
		PowerCost:   50,
		Description: "Returns to section start when killed. Reroll all stats.",
		Category:    CategoryRecovery,
		CombatOnly:  false,
		DeathOnly:   true,
	},
	{
		Name:        "RETRACE",
		PowerCost:   20,
		Description: "Returns to any previously visited section.",
		Category:    CategoryNavigation,
		CombatOnly:  false,
		DeathOnly:   false,
	},
	{
		Name:        "TIMEWARP",
		PowerCost:   10,
		Description: "Resets section to starting state. Restores all LP.",
		Category:    CategoryNavigation,
		CombatOnly:  false,
		DeathOnly:   false,
	},
	{
		Name:        "XENOPHOBIA",
		PowerCost:   15,
		Description: "Causes enemy to fear you. Reduces their damage by 5 points.",
		Category:    CategoryOffensive,
		CombatOnly:  true,
		DeathOnly:   false,
	},
}

// GetSpellByName returns a spell by its name, or nil if not found.
func GetSpellByName(name string) *Spell {
	for i := range AllSpells {
		if AllSpells[i].Name == name {
			return &AllSpells[i]
		}
	}
	return nil
}

// GetAvailableSpells returns spells that are available based on context.
func GetAvailableSpells(inCombat bool, isDead bool) []Spell {
	available := []Spell{}
	for _, spell := range AllSpells {
		// RESURRECTION only available when dead
		if spell.DeathOnly && !isDead {
			continue
		}
		if !spell.DeathOnly && isDead {
			continue
		}
		// Combat-only spells not available outside combat
		if spell.CombatOnly && !inCombat {
			continue
		}
		available = append(available, spell)
	}
	return available
}
