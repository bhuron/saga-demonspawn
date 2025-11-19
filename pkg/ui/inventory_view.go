package ui

import (
	"fmt"
	"strings"

	"github.com/benoit/saga-demonspawn/internal/items"
	"github.com/benoit/saga-demonspawn/pkg/ui/theme"
)

// renderInventoryView renders the inventory management screen
func renderInventoryView(m Model) string {
	var b strings.Builder
	t := theme.Current()

	b.WriteString("\n")
	b.WriteString(theme.RenderTitle("INVENTORY MANAGEMENT"))
	b.WriteString("\n\n")

	// Currently Equipped section
	b.WriteString(t.Heading.Render("  Currently Equipped") + "\n")
	b.WriteString(theme.RenderSeparator(60) + "\n")

	// Weapon
	weaponName := "None"
	weaponBonus := 0
	if m.Character.EquippedWeapon != nil {
		weaponName = m.Character.EquippedWeapon.Name
		weaponBonus = m.Character.EquippedWeapon.DamageBonus
	}
	b.WriteString("  " + theme.RenderLabel("Weapon", fmt.Sprintf("%s (+%d damage)", weaponName, weaponBonus)) + "\n")

	// Armor
	armorName := "None"
	armorProtection := 0
	if m.Character.EquippedArmor != nil {
		armorName = m.Character.EquippedArmor.Name
		armorProtection = m.Character.EquippedArmor.Protection
	}
	b.WriteString("  " + theme.RenderLabel("Armor", fmt.Sprintf("%s (-%d damage)", armorName, armorProtection)) + "\n")

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
	b.WriteString("  " + theme.RenderLabel("Shield", fmt.Sprintf("%s (-%d damage)", shieldStatus, shieldProtection)) + "\n")

	totalProtection := m.Character.GetArmorProtection()
	b.WriteString("\n  " + t.Emphasis.Render(fmt.Sprintf("Total Protection: -%d damage", totalProtection)) + "\n\n")

	b.WriteString(theme.RenderSeparator(60) + "\n")

	// Available Equipment and Special Items
	invItems := m.Inventory.GetItems()
	cursor := m.Inventory.GetCursor()

	// Implement scrolling viewport to prevent screen overflow
	maxVisibleItems := 12
	startIdx := 0
	endIdx := len(invItems)

	if len(invItems) > maxVisibleItems {
		viewportMid := maxVisibleItems / 2
		startIdx = cursor - viewportMid
		endIdx = cursor + viewportMid

		if startIdx < 0 {
			startIdx = 0
			endIdx = maxVisibleItems
		}

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
		b.WriteString(t.MutedText.Render("  ↑ More items above...") + "\n\n")
	}

	// Render visible items only
	for i := startIdx; i < endIdx; i++ {
		item := invItems[i]
		selected := i == cursor

		if item.IsHeader {
			b.WriteString("\n" + t.Heading.Render("  "+item.Name) + "\n")
		} else {
			// Regular item
			equipped := ""
			if item.IsEquipped {
				equipped = t.SuccessMsg.Render("[EQUIPPED]")
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
					desc = fmt.Sprintf("-%d protection", item.Armor.Protection)
				}
			} else if item.Category == CategoryShield {
				desc = fmt.Sprintf("-%d/-%d protection", items.ShieldStandard.Protection, items.ShieldStandard.ProtectionWithArmor)
			} else if item.Category == CategorySpecialItems {
				if item.SpecialItem == "healing_stone" {
					desc = fmt.Sprintf("Charges: %d/50", m.Character.HealingStoneCharges)
					if m.Character.HealingStoneCharges > 0 {
						equipped = t.SuccessMsg.Render("[AVAILABLE]")
					} else {
						equipped = t.MutedText.Render("[DEPLETED]")
					}
				} else if item.SpecialItem == "doombringer" {
					desc = "+20 damage, cursed"
					if !m.Character.DoombringerPossessed {
						equipped = t.MutedText.Render("[NOT POSSESSED]")
					}
				} else if item.SpecialItem == "orb" {
					desc = "Anti-Demonspawn"
					if m.Character.OrbDestroyed {
						equipped = t.Error.Render("[DESTROYED]")
					} else if m.Character.OrbEquipped {
						equipped = t.SuccessMsg.Render("[EQUIPPED]")
					} else if m.Character.OrbPossessed {
						equipped = t.Emphasis.Render("[POSSESSED]")
					} else {
						equipped = t.MutedText.Render("[NOT POSSESSED]")
					}
				}
			}

			itemName := fmt.Sprintf("%-20s", item.Name)
			itemDesc := fmt.Sprintf("%-25s", desc)
			
			if selected {
				b.WriteString("  " + theme.RenderMenuItem(itemName+" "+itemDesc+" "+equipped, true) + "\n")
			} else {
				b.WriteString("  " + t.MenuItem.Render(itemName) + " " + t.MutedText.Render(itemDesc) + " " + equipped + "\n")
			}
		}
	}

	// Show scroll indicator if there are more items below
	if endIdx < len(invItems) {
		b.WriteString("\n" + t.MutedText.Render("  ↓ More items below...") + "\n")
	}

	b.WriteString("\n" + theme.RenderSeparator(60) + "\n")

	// Actions
	b.WriteString(t.Heading.Render("  Actions") + "\n")
	b.WriteString(theme.RenderKeyHelp("Enter Equip", "U Use", "R Recharge", "A Acquire", "I Info", "Q/Esc Back") + "\n")

	// Message display
	if m.Inventory.GetMessage() != "" {
		b.WriteString("\n" + t.Emphasis.Render("  "+m.Inventory.GetMessage()) + "\n")
	}

	return b.String()
}
