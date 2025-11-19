package ui

import (
	"fmt"
	
	"github.com/benoit/saga-demonspawn/internal/character"
	"github.com/benoit/saga-demonspawn/internal/items"
)

// ItemCategory represents different sections in the inventory
type ItemCategory int

const (
	CategoryWeapons ItemCategory = iota
	CategoryArmor
	CategoryShield
	CategorySpecialItems
)

// InventoryItem represents a selectable item in the inventory
type InventoryItem struct {
	Name        string
	Category    ItemCategory
	IsEquipped  bool
	IsHeader    bool  // True for section headers like "WEAPONS"
	IsNone      bool  // True for "None" options
	Weapon      *items.Weapon
	Armor       *items.Armor
	IsShield    bool
	SpecialItem string // "healing_stone", "doombringer", "orb"
}

// InventoryManagementModel represents the inventory management screen state
type InventoryManagementModel struct {
	character *character.Character
	cursor    int
	items     []InventoryItem
	inCombat  bool
	message   string // For displaying validation messages
}

// NewInventoryManagementModel creates a new inventory management model
func NewInventoryManagementModel(char *character.Character, inCombat bool) InventoryManagementModel {
	model := InventoryManagementModel{
		character: char,
		cursor:    0,
		inCombat:  inCombat,
		message:   "",
	}
	model.rebuildItemList()
	
	// Initialize cursor to first non-header item
	for i, item := range model.items {
		if !item.IsHeader {
			model.cursor = i
			break
		}
	}
	
	return model
}

// rebuildItemList reconstructs the inventory items list based on current state
func (m *InventoryManagementModel) rebuildItemList() {
	m.items = []InventoryItem{}

	// Weapons section
	m.items = append(m.items, InventoryItem{
		Name:     "WEAPONS",
		Category: CategoryWeapons,
		IsHeader: true,
	})

	// None weapon option
	m.items = append(m.items, InventoryItem{
		Name:       "None",
		Category:   CategoryWeapons,
		IsEquipped: m.character.EquippedWeapon == nil,
		IsNone:     true,
	})

	// Standard weapons
	weapons := items.AllWeapons()
	for _, w := range weapons {
		weapon := w // Create a copy for the pointer
		
		// Skip Doombringer if not possessed
		if weapon.Name == items.DoombringerName && !m.character.DoombringerPossessed {
			continue
		}
		
		isEquipped := m.character.EquippedWeapon != nil && m.character.EquippedWeapon.Name == weapon.Name
		m.items = append(m.items, InventoryItem{
			Name:       weapon.Name,
			Category:   CategoryWeapons,
			IsEquipped: isEquipped,
			Weapon:     &weapon,
		})
	}

	// Armor section
	m.items = append(m.items, InventoryItem{
		Name:     "ARMOR",
		Category: CategoryArmor,
		IsHeader: true,
	})

	armors := items.AllArmor()
	for _, a := range armors {
		armor := a // Create a copy for the pointer
		isEquipped := m.character.EquippedArmor != nil && m.character.EquippedArmor.Name == armor.Name
		m.items = append(m.items, InventoryItem{
			Name:       armor.Name,
			Category:   CategoryArmor,
			IsEquipped: isEquipped,
			Armor:      &armor,
		})
	}

	// Shield section
	m.items = append(m.items, InventoryItem{
		Name:     "SHIELD",
		Category: CategoryShield,
		IsHeader: true,
	})

	m.items = append(m.items, InventoryItem{
		Name:       "Shield",
		Category:   CategoryShield,
		IsEquipped: m.character.HasShield,
		IsShield:   true,
	})

	// Special items section
	m.items = append(m.items, InventoryItem{
		Name:     "SPECIAL ITEMS",
		Category: CategorySpecialItems,
		IsHeader: true,
	})

	// Healing Stone - always show, mark as possessed or not
	m.items = append(m.items, InventoryItem{
		Name:        "Healing Stone",
		Category:    CategorySpecialItems,
		SpecialItem: "healing_stone",
	})

	// Doombringer - always show (unless equipped as weapon)
	isDoombringerEquipped := m.character.EquippedWeapon != nil && m.character.EquippedWeapon.Name == items.DoombringerName
	if !isDoombringerEquipped {
		m.items = append(m.items, InventoryItem{
			Name:        items.DoombringerName,
			Category:    CategorySpecialItems,
			SpecialItem: "doombringer",
		})
	}

	// The Orb - always show (unless destroyed)
	if !m.character.OrbDestroyed {
		m.items = append(m.items, InventoryItem{
			Name:        items.TheOrbName,
			Category:    CategorySpecialItems,
			SpecialItem: "orb",
			IsEquipped:  m.character.OrbEquipped,
		})
	}
}

// MoveUp moves the cursor up, skipping headers
func (m *InventoryManagementModel) MoveUp() {
	if m.cursor <= 0 {
		return
	}
	
	// Move up and skip any headers
	m.cursor--
	for m.cursor > 0 && m.items[m.cursor].IsHeader {
		m.cursor--
	}
	
	// If we landed on a header at position 0, find first non-header
	if m.items[m.cursor].IsHeader {
		for i := 0; i < len(m.items); i++ {
			if !m.items[i].IsHeader {
				m.cursor = i
				break
			}
		}
	}
}

// MoveDown moves the cursor down, skipping headers
func (m *InventoryManagementModel) MoveDown() {
	if m.cursor >= len(m.items)-1 {
		return
	}
	
	// Move down and skip any headers
	m.cursor++
	for m.cursor < len(m.items) && m.items[m.cursor].IsHeader {
		m.cursor++
	}
	
	// If we went past the end, find the last non-header item
	if m.cursor >= len(m.items) {
		for i := len(m.items) - 1; i >= 0; i-- {
			if !m.items[i].IsHeader {
				m.cursor = i
				break
			}
		}
	}
}

// HandleEnter processes the Enter key press for equipping items
func (m *InventoryManagementModel) HandleEnter() {
	if m.inCombat {
		m.message = "[LOCKED IN COMBAT] Cannot change equipment during combat"
		return
	}

	if m.cursor >= len(m.items) {
		return
	}

	item := m.items[m.cursor]

	switch item.Category {
	case CategoryWeapons:
		if item.IsNone {
			m.character.EquipWeapon(nil)
		} else if item.Weapon != nil {
			m.character.EquipWeapon(item.Weapon)
		}
		m.message = ""

	case CategoryArmor:
		if item.Armor != nil {
			m.character.EquipArmor(item.Armor)
		}
		m.message = ""

	case CategoryShield:
		m.character.ToggleShield()
		m.message = ""

	case CategorySpecialItems:
		if item.SpecialItem == "orb" && !m.character.OrbDestroyed {
			m.character.OrbEquipped = !m.character.OrbEquipped
			if m.character.OrbEquipped && m.character.HasShield {
				// Cannot equip shield while Orb is held
				m.character.HasShield = false
				m.message = "Shield unequipped - cannot use shield while holding The Orb"
			} else {
				m.message = ""
			}
		}
	}

	m.rebuildItemList()
}

// HandleUse processes the 'u' key for using items
func (m *InventoryManagementModel) HandleUse() {
	if m.cursor >= len(m.items) {
		return
	}

	item := m.items[m.cursor]

	if item.SpecialItem == "healing_stone" {
		if m.inCombat {
			m.message = "Cannot use Healing Stone outside of combat turns"
			return
		}

		if m.character.HealingStoneCharges <= 0 {
			m.message = "The Healing Stone is depleted"
			return
		}

		if m.character.CurrentLP >= m.character.MaximumLP {
			m.message = "You are already at full health"
			return
		}

		// This is just a preview - actual healing happens in combat
		m.message = "Healing Stone can be used during combat (see combat actions)"
	} else {
		m.message = "This item cannot be used here"
	}
}

// HandleRecharge processes the 'r' key for recharging Healing Stone
func (m *InventoryManagementModel) HandleRecharge() bool {
	if m.cursor >= len(m.items) {
		return false
	}

	item := m.items[m.cursor]

	if item.SpecialItem == "healing_stone" {
		if m.character.HealingStoneCharges >= 50 {
			m.message = "The Healing Stone is already fully charged"
			return false
		}
		// Return true to request confirmation
		return true
	}

	m.message = "This item cannot be recharged"
	return false
}

// HandleAcquire processes the 'a' key for acquiring special items
func (m *InventoryManagementModel) HandleAcquire() {
	if m.cursor >= len(m.items) {
		return
	}

	item := m.items[m.cursor]

	if item.SpecialItem == "healing_stone" {
		if m.character.HealingStoneCharges > 0 {
			m.message = "You already possess the Healing Stone"
		} else {
			m.character.AcquireHealingStone()
			m.message = "Acquired Healing Stone! Restore LP during combat."
			m.rebuildItemList()
		}
	} else if item.SpecialItem == "doombringer" {
		if !m.character.DoombringerPossessed {
			m.character.AcquireDoombringer()
			m.message = "Acquired Doombringer! Beware its cursed power..."
			m.rebuildItemList()
		} else {
			m.message = "You already possess Doombringer"
		}
	} else if item.SpecialItem == "orb" {
		if !m.character.OrbPossessed {
			m.character.AcquireOrb()
			m.message = "Acquired The Orb! A powerful weapon against Demonspawn."
			m.rebuildItemList()
		} else {
			m.message = "You already possess The Orb"
		}
	} else {
		m.message = "This item cannot be acquired (use this when finding items in the adventure)"
	}
}

// ConfirmRecharge recharges the Healing Stone
func (m *InventoryManagementModel) ConfirmRecharge() {
	err := m.character.RechargeHealingStone()
	if err != nil {
		m.message = fmt.Sprintf("Cannot recharge: %v", err)
	} else {
		m.message = "Healing Stone recharged to 50 LP"
	}
	m.rebuildItemList()
}

// GetCurrentItem returns the currently selected item
func (m *InventoryManagementModel) GetCurrentItem() *InventoryItem {
	if m.cursor >= 0 && m.cursor < len(m.items) {
		return &m.items[m.cursor]
	}
	return nil
}

// GetItems returns all inventory items
func (m *InventoryManagementModel) GetItems() []InventoryItem {
	return m.items
}

// GetCursor returns the current cursor position
func (m *InventoryManagementModel) GetCursor() int {
	return m.cursor
}

// GetMessage returns the current message
func (m *InventoryManagementModel) GetMessage() string {
	return m.message
}

// ClearMessage clears the current message
func (m *InventoryManagementModel) ClearMessage() {
	m.message = ""
}
