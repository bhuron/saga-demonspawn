package ui

import (
	"fmt"
	"strings"

	"github.com/benoit/saga-demonspawn/internal/items"
)

// renderInventoryView renders the inventory management screen
func renderInventoryView(m Model) string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString("╔══════════════════════════════════════════════════════════╗\n")
	b.WriteString("║                   INVENTORY MANAGEMENT                    ║\n")
	b.WriteString("╠══════════════════════════════════════════════════════════╣\n")
	b.WriteString("║                                                           ║\n")

	// Currently Equipped section
	b.WriteString("║  CURRENTLY EQUIPPED                                       ║\n")
	b.WriteString("║  ┌─────────────────────────────────────────────────────┐ ║\n")

	// Weapon
	weaponName := "None"
	weaponBonus := 0
	if m.Character.EquippedWeapon != nil {
		weaponName = m.Character.EquippedWeapon.Name
		weaponBonus = m.Character.EquippedWeapon.DamageBonus
	}
	weaponLine := fmt.Sprintf("║  │ Weapon:  %-18s (+%d damage)%-11s│ ║\n", weaponName, weaponBonus, "")
	b.WriteString(weaponLine)

	// Armor
	armorName := "None"
	armorProtection := 0
	if m.Character.EquippedArmor != nil {
		armorName = m.Character.EquippedArmor.Name
		armorProtection = m.Character.EquippedArmor.Protection
	}
	armorLine := fmt.Sprintf("║  │ Armor:   %-18s (%d protection)%-9s│ ║\n", armorName, armorProtection, "")
	b.WriteString(armorLine)

	// Shield
	shieldStatus := "Not Equipped"
	shieldProtection := 0
	if m.Character.HasShield {
		shieldStatus = "Equipped"
		if m.Character.EquippedArmor != nil && m.Character.EquippedArmor.Name != "None" {
			shieldProtection = items.ShieldStandard.ProtectionWithArmor
		} else {
			shieldProtection = items.ShieldStandard.Protection
		}
	}
	shieldLine := fmt.Sprintf("║  │ Shield:  %-18s (%d protection)%-9s│ ║\n", shieldStatus, shieldProtection, "")
	b.WriteString(shieldLine)

	b.WriteString("║  │                                                      │ ║\n")

	totalProtection := m.Character.GetArmorProtection()
	protectionLine := fmt.Sprintf("║  │ Total Protection: %-35d│ ║\n", totalProtection)
	b.WriteString(protectionLine)

	b.WriteString("║  └─────────────────────────────────────────────────────┘ ║\n")
	b.WriteString("║                                                           ║\n")

	// Available Equipment and Special Items
	invItems := m.Inventory.GetItems()
	cursor := m.Inventory.GetCursor()

	// Implement scrolling viewport to prevent screen overflow
	// Maximum items to show at once (to fit on screen)
	maxVisibleItems := 12
	startIdx := 0
	endIdx := len(invItems)

	if len(invItems) > maxVisibleItems {
		// Calculate viewport window around cursor
		// Try to keep cursor in middle of viewport
		viewportMid := maxVisibleItems / 2
		startIdx = cursor - viewportMid
		endIdx = cursor + viewportMid

		// Adjust if at start of list
		if startIdx < 0 {
			startIdx = 0
			endIdx = maxVisibleItems
		}

		// Adjust if at end of list
		if endIdx > len(invItems) {
			endIdx = len(invItems)
			startIdx = endIdx - maxVisibleItems
			if startIdx < 0 {
				startIdx = 0
			}
		}
	}

	// Show scroll indicators if needed
	if startIdx > 0 {
		b.WriteString("║  ↑ More items above...                                    ║\n")
	}

	// Render visible items only
	for i := startIdx; i < endIdx; i++ {
		item := invItems[i]
		prefix := "  "
		if i == cursor {
			prefix = "> "
		}

		if item.IsHeader {
			// Section header - width is 59 chars total (2 for prefix + 57 for content)
			headerLine := fmt.Sprintf("║  %-59s║\n", prefix+item.Name)
			b.WriteString(headerLine)
		} else {
			// Regular item
			equipped := ""
			if item.IsEquipped {
				equipped = "[EQUIPPED]"
			}

			// Calculate description
			desc := ""
			if item.Category == CategoryWeapons {
				if item.IsNone {
					desc = "No weapon"
				} else if item.Weapon != nil {
					desc = fmt.Sprintf("+%d damage", item.Weapon.DamageBonus)
				}
			} else if item.Category == CategoryArmor {
				if item.Armor != nil {
					desc = fmt.Sprintf("%d protection", item.Armor.Protection)
				}
			} else if item.Category == CategoryShield {
				desc = fmt.Sprintf("%d/-%d protection", items.ShieldStandard.Protection, items.ShieldStandard.ProtectionWithArmor)
			} else if item.Category == CategorySpecialItems {
				if item.SpecialItem == "healing_stone" {
					desc = fmt.Sprintf("Charges: %d/50", m.Character.HealingStoneCharges)
					if m.Character.HealingStoneCharges > 0 {
						equipped = "[AVAILABLE]"
					} else {
						equipped = "[DEPLETED]"
					}
				} else if item.SpecialItem == "doombringer" {
					desc = "+20 damage, cursed"
					if !m.Character.DoombringerPossessed {
						equipped = "[NOT POSSESSED]"
					}
				} else if item.SpecialItem == "orb" {
					desc = "Anti-Demonspawn"
					if m.Character.OrbDestroyed {
						equipped = "[DESTROYED]"
					} else if m.Character.OrbEquipped {
						equipped = "[EQUIPPED]"
					} else if m.Character.OrbPossessed {
						equipped = "[POSSESSED]"
					} else {
						equipped = "[NOT POSSESSED]"
					}
				}
			}

			// Format item line with proper spacing
			// Total content width is 59 chars between ║ symbols
			// Build the full line then ensure it's exactly 59 chars
			itemText := fmt.Sprintf("%s%-18s %-22s %s", prefix, item.Name, desc, equipped)
			// Pad or trim to exactly 59 characters
			if len(itemText) > 59 {
				itemText = itemText[:59]
			} else if len(itemText) < 59 {
				itemText = fmt.Sprintf("%-59s", itemText)
			}
			itemLine := fmt.Sprintf("║  %s║\n", itemText)
			b.WriteString(itemLine)
		}
	}

	// Show scroll indicator if there are more items below
	if endIdx < len(invItems) {
		b.WriteString("║  ↓ More items below...                                    ║\n")
	}

	b.WriteString("║                                                           ║\n")

	// Actions
	b.WriteString("║  ACTIONS                                                  ║\n")
	b.WriteString("║  [Enter] Equip  [U] Use  [R] Recharge  [A] Acquire      ║\n")
	b.WriteString("║  [I] Info  [Q/Esc] Back to Game Session                 ║\n")

	// Message display
	if m.Inventory.GetMessage() != "" {
		b.WriteString("║                                                           ║\n")
		msg := m.Inventory.GetMessage()
		if len(msg) > 55 {
			msg = msg[:55]
		}
		msgLine := fmt.Sprintf("║  %-57s║\n", msg)
		b.WriteString(msgLine)
	}

	b.WriteString("╚══════════════════════════════════════════════════════════╝\n")

	return b.String()
}
