package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/benoit/saga-demonspawn/internal/character"
	"github.com/benoit/saga-demonspawn/internal/combat"
	"github.com/benoit/saga-demonspawn/internal/dice"
	"github.com/benoit/saga-demonspawn/pkg/ui/theme"
)

// CombatViewModel handles the active combat interface.
type CombatViewModel struct {
	player      *character.Character
	combatState *combat.CombatState
	roller      dice.Roller

	// UI state
	selectedAction int
	waitingForInput bool
	victoryState    bool
	defeatState     bool
	needsRest       bool
	needsEnemyRest  bool
	deathSaveActive bool

	// Action menu
	actions []string
}

const (
	actionAttack = iota
	actionCastSpell
	actionFlee
	actionUseHealingStone
	actionThrowOrb
	actionTotalActions
)

// NewCombatViewModel creates a new combat view model.
func NewCombatViewModel(player *character.Character, combatState *combat.CombatState, roller dice.Roller) CombatViewModel {
	// Build action list based on available items
	actions := []string{"Attack"}
	
	// Add Cast Spell option if magic is unlocked
	if player.MagicUnlocked {
		actions = append(actions, "Cast Spell")
	}
	
	actions = append(actions, "Flee Combat")
	
	// Add Healing Stone option if available
	if player.HealingStoneCharges > 0 {
		actions = append(actions, fmt.Sprintf("Use Healing Stone (%d charges)", player.HealingStoneCharges))
	}
	
	// Add Throw Orb option if possessed and not equipped (can't throw while held)
	if player.OrbPossessed && !player.OrbDestroyed && !player.OrbEquipped {
		actions = append(actions, "Throw The Orb")
	}
	
	return CombatViewModel{
		player:          player,
		combatState:     combatState,
		roller:          roller,
		selectedAction:  actionAttack,
		waitingForInput: true,
		victoryState:    false,
		defeatState:     false,
		needsRest:       false,
		needsEnemyRest:  false,
		deathSaveActive: false,
		actions:         actions,
	}
}

// Update handles combat view input.
func (m CombatViewModel) Update(msg tea.Msg) (CombatViewModel, tea.Cmd) {
	// Handle victory/defeat states
	if m.victoryState || m.defeatState {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "enter" {
				return m, func() tea.Msg {
					return CombatEndMsg{Victory: m.victoryState}
				}
			}
		}
		return m, nil
	}

	// Handle death save
	if m.deathSaveActive {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "enter" {
				roll, success := combat.AttemptDeathSave(m.player, m.combatState, m.roller)
				if success {
					m.combatState.AddLogEntry(fmt.Sprintf("[Death Save] Rolled %d vs Luck %d - SUCCESS!", roll, m.player.Luck))
					m.combatState.AddLogEntry(fmt.Sprintf("[Death Save] Restored to %d LP. Combat restarted!", m.player.CurrentLP))
					m.combatState.AddLogEntry(fmt.Sprintf("[Initiative] Player: %d, Enemy: %d", m.combatState.PlayerInitiative, m.combatState.EnemyInitiative))
					m.deathSaveActive = false
					m.waitingForInput = m.combatState.PlayerTurn
					// If it's enemy turn after death save, trigger enemy turn
					if !m.combatState.PlayerTurn {
						return m, func() tea.Msg {
							return EnemyTurnMsg{}
						}
					}
					return m, nil
				} else {
					m.combatState.AddLogEntry(fmt.Sprintf("[Death Save] Rolled %d vs Luck %d - FAILED!", roll, m.player.Luck))
					m.defeatState = true
					return m, nil
				}
			}
		}
		return m, nil
	}

	// Handle turn processing messages FIRST before auto-trigger
	switch msg.(type) {
	case EnemyTurnMsg:
		return m.processEnemyTurn()
	case PlayerAttackCompleteMsg:
		return m.checkCombatState()
	case EnemyAttackCompleteMsg:
		return m.checkCombatState()
	}

	// Handle player turn input
	if m.waitingForInput && m.combatState.PlayerTurn {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "up", "k":
				if m.selectedAction > 0 {
					m.selectedAction--
				}
			case "down", "j":
				if m.selectedAction < len(m.actions)-1 {
					m.selectedAction++
				}
			case "enter":
				return m.handleAction()
			}
		}
	}

	// If it's enemy turn and we haven't processed a message, trigger enemy turn automatically
	if !m.combatState.PlayerTurn && !m.needsRest {
		return m, func() tea.Msg {
			return EnemyTurnMsg{}
		}
	}

	return m, nil
}

func (m CombatViewModel) handleAction() (CombatViewModel, tea.Cmd) {
	actionName := m.actions[m.selectedAction]
	
	// Check for Cast Spell action
	if actionName == "Cast Spell" {
		// Signal to switch to magic screen
		return m, func() tea.Msg {
			return CastSpellMsg{}
		}
	}
	
	switch m.selectedAction {
	case actionAttack:
		// Check if rest is needed before attack
		if combat.CheckEndurance(m.combatState.RoundsSinceLastRest, m.combatState.EnduranceLimit) {
			m.combatState.AddLogEntry(fmt.Sprintf("[Round %d] Endurance depleted! Must rest.", m.combatState.CurrentRound))
			m.needsRest = true
			m.waitingForInput = false
			return m, func() tea.Msg {
				return EnemyTurnMsg{}
			}
		}

		// Check if Doombringer is equipped - apply blood price BEFORE attack
		isDoombringerEquipped := m.player.EquippedWeapon != nil && m.player.EquippedWeapon.Name == "Doombringer"
		if isDoombringerEquipped {
			// Blood price: 10 LP before attack
			m.player.ModifyLP(-10)
			m.combatState.AddLogEntry(fmt.Sprintf("[R%d] Doombringer thirsts for blood... -10 LP", m.combatState.CurrentRound))
			m.combatState.AddLogEntry(fmt.Sprintf("[R%d] Current LP: %d/%d", m.combatState.CurrentRound, m.player.CurrentLP, m.player.MaximumLP))
			
			// Check if blood price killed the player
			if m.player.CurrentLP <= 0 {
				m.combatState.AddLogEntry("[Defeat] Doombringer has drained your life!")
				m.defeatState = true
				return m, nil
			}
		}

		// Execute player attack
		result := combat.ExecutePlayerAttack(m.combatState, m.player, m.roller)
		
		if result.Hit {
			m.combatState.AddLogEntry(fmt.Sprintf("[R%d] You rolled %d (need %d+) - HIT!", m.combatState.CurrentRound, result.Roll, result.Requirement))
			
			// Check if The Orb is equipped and enemy is Demonspawn - double damage
			if m.player.OrbEquipped && m.combatState.Enemy.IsDemonspawn {
				// Calculate how much damage was actually dealt
				originalDamage := result.FinalDamage
				doubledDamage := originalDamage * 2
				
				// Apply additional damage to enemy
				m.combatState.Enemy.CurrentLP -= originalDamage // Subtract the original amount again
				if m.combatState.Enemy.CurrentLP < 0 {
					m.combatState.Enemy.CurrentLP = 0
				}
				
				m.combatState.AddLogEntry(fmt.Sprintf("[R%d] The Orb pulses with power! Damage doubled: %d → %d", m.combatState.CurrentRound, originalDamage, doubledDamage))
				m.combatState.AddLogEntry(fmt.Sprintf("[R%d] Damage: (%d×5) + STR + Weapon - Armor = %d", m.combatState.CurrentRound, result.Roll, doubledDamage))
				m.combatState.AddLogEntry(fmt.Sprintf("[R%d] Enemy takes %d damage (%d LP remaining)", m.combatState.CurrentRound, doubledDamage, m.combatState.Enemy.CurrentLP))
				
				// Update result for Doombringer healing calculation
				result.FinalDamage = doubledDamage
				result.TargetLP = m.combatState.Enemy.CurrentLP
			} else {
				m.combatState.AddLogEntry(fmt.Sprintf("[R%d] Damage: (%d×5) + STR + Weapon - Armor = %d", m.combatState.CurrentRound, result.Roll, result.FinalDamage))
				m.combatState.AddLogEntry(fmt.Sprintf("[R%d] Enemy takes %d damage (%d LP remaining)", m.combatState.CurrentRound, result.FinalDamage, result.TargetLP))
			}
			
			// Doombringer soul thirst: heal LP equal to damage dealt (capped at enemy's current LP and MaximumLP)
			if isDoombringerEquipped && result.FinalDamage > 0 {
				oldLP := m.player.CurrentLP
				healAmount := result.FinalDamage
				
				// If enemy died from this hit, cap healing at enemy's LP before the hit
				if result.TargetLP <= 0 {
					enemyLPBeforeHit := result.TargetLP + result.FinalDamage
					if enemyLPBeforeHit < healAmount {
						healAmount = enemyLPBeforeHit
					}
				}
				
				// Cap healing at MaximumLP
				if m.player.CurrentLP + healAmount > m.player.MaximumLP {
					healAmount = m.player.MaximumLP - m.player.CurrentLP
				}
				
				if healAmount > 0 {
					m.player.ModifyLP(healAmount)
					actualHeal := m.player.CurrentLP - oldLP
					m.combatState.AddLogEntry(fmt.Sprintf("[R%d] Doombringer feeds on pain... +%d LP healed!", m.combatState.CurrentRound, actualHeal))
					m.combatState.AddLogEntry(fmt.Sprintf("[R%d] Current LP: %d/%d", m.combatState.CurrentRound, m.player.CurrentLP, m.player.MaximumLP))
				} else {
					m.combatState.AddLogEntry(fmt.Sprintf("[R%d] Doombringer feeds on pain... (already at maximum LP)", m.combatState.CurrentRound))
				}
			}
		} else {
			m.combatState.AddLogEntry(fmt.Sprintf("[R%d] You rolled %d (need %d+) - MISS!", m.combatState.CurrentRound, result.Roll, result.Requirement))
			if isDoombringerEquipped {
				m.combatState.AddLogEntry(fmt.Sprintf("[R%d] No healing from Doombringer on miss", m.combatState.CurrentRound))
			}
		}

		m.waitingForInput = false
		return m, func() tea.Msg {
			return PlayerAttackCompleteMsg{}
		}

	case actionFlee:
		m.combatState.AddLogEntry("[Fled] You fled from combat!")
		return m, func() tea.Msg {
			return CombatEndMsg{Victory: false}
		}
	
	case actionUseHealingStone:
		// Check if Healing Stone is available
		if m.player.HealingStoneCharges <= 0 {
			m.combatState.AddLogEntry("[Healing Stone] The stone is depleted!")
			m.waitingForInput = false
			return m, func() tea.Msg {
				return PlayerAttackCompleteMsg{}
			}
		}
		
		// Check if already at max LP
		if m.player.CurrentLP >= m.player.MaximumLP {
			m.combatState.AddLogEntry("[Healing Stone] You are already at full health!")
			m.waitingForInput = false
			return m, func() tea.Msg {
				return PlayerAttackCompleteMsg{}
			}
		}
		
		// Roll healing amount (1d6 × 10)
		roll := m.roller.Roll1D6()
		healAmount := roll * 10
		
		// Use the healing stone
		actualHeal, err := m.player.UseHealingStone(healAmount)
		if err != nil {
			m.combatState.AddLogEntry(fmt.Sprintf("[Healing Stone] Error: %v", err))
			m.waitingForInput = false
			return m, func() tea.Msg {
				return PlayerAttackCompleteMsg{}
			}
		}
		
		// Log healing
		m.combatState.AddLogEntry(fmt.Sprintf("[R%d] You invoke the Healing Stone... (rolled %d)", m.combatState.CurrentRound, roll))
		m.combatState.AddLogEntry(fmt.Sprintf("[R%d] +%d LP restored! (Charges: %d/50)", m.combatState.CurrentRound, actualHeal, m.player.HealingStoneCharges))
		m.combatState.AddLogEntry(fmt.Sprintf("[R%d] Current LP: %d/%d", m.combatState.CurrentRound, m.player.CurrentLP, m.player.MaximumLP))
		
		// Update action list if charges depleted
		if m.player.HealingStoneCharges <= 0 {
			m.actions = []string{"Attack", "Flee Combat"}
		}
		
		// Keep waiting for input - it's still player's turn
		m.waitingForInput = true
		return m, func() tea.Msg {
			return PlayerAttackCompleteMsg{}
		}
	
	case actionThrowOrb:
		// Check if Orb is still available
		if m.player.OrbDestroyed || !m.player.OrbPossessed {
			m.combatState.AddLogEntry("[The Orb] The Orb has been destroyed!")
			m.waitingForInput = false
			return m, func() tea.Msg {
				return PlayerAttackCompleteMsg{}
			}
		}
		
		// Check if Orb is equipped (can't throw while held)
		if m.player.OrbEquipped {
			m.combatState.AddLogEntry("[The Orb] Unequip The Orb before throwing!")
			m.waitingForInput = false
			return m, func() tea.Msg {
				return PlayerAttackCompleteMsg{}
			}
		}
		
		// Roll 2d6 for throw (hit on 4+)
		roll := m.roller.Roll2D6()
		hit := roll >= 4
		
		m.combatState.AddLogEntry(fmt.Sprintf("[R%d] You hurl The Orb at the enemy!", m.combatState.CurrentRound))
		m.combatState.AddLogEntry(fmt.Sprintf("[R%d] Rolled %d (need 4+)", m.combatState.CurrentRound, roll))
		
		if m.combatState.Enemy.IsDemonspawn {
			if hit {
				// Instant kill
				m.combatState.Enemy.CurrentLP = 0
				m.combatState.AddLogEntry("[The Orb] The Orb strikes true! The Demonspawn is annihilated in brilliant light!")
			} else {
				// Deal 200 damage
				m.combatState.Enemy.CurrentLP -= 200
				if m.combatState.Enemy.CurrentLP < 0 {
					m.combatState.Enemy.CurrentLP = 0
				}
				m.combatState.AddLogEntry(fmt.Sprintf("[The Orb] The Orb's light sears the Demonspawn! 200 damage dealt! (%d LP remaining)", m.combatState.Enemy.CurrentLP))
			}
		} else {
			// No effect on non-Demonspawn
			m.combatState.AddLogEntry("[The Orb] The Orb has no effect on this creature!")
		}
		
		// Destroy the Orb
		m.player.DestroyOrb()
		m.combatState.AddLogEntry("[The Orb] The Orb explodes and is destroyed!")
		
		// Update actions list to remove Throw Orb option
		m.actions = []string{"Attack", "Flee Combat"}
		if m.player.HealingStoneCharges > 0 {
			m.actions = append(m.actions, fmt.Sprintf("Use Healing Stone (%d charges)", m.player.HealingStoneCharges))
		}
		
		m.waitingForInput = false
		return m, func() tea.Msg {
			return PlayerAttackCompleteMsg{}
		}
	}

	return m, nil
}

func (m CombatViewModel) processEnemyTurn() (CombatViewModel, tea.Cmd) {
	// If player resting, enemy gets free attack
	if m.needsRest {
		m.combatState.AddLogEntry(fmt.Sprintf("[R%d] Enemy attacks while you rest...", m.combatState.CurrentRound))
		result := combat.ExecuteEnemyAttack(m.combatState, m.player, m.roller)
		
		if result.Hit {
			m.combatState.AddLogEntry(fmt.Sprintf("[R%d] Enemy rolled %d (need %d+) - HIT!", m.combatState.CurrentRound, result.Roll, result.Requirement))
			m.combatState.AddLogEntry(fmt.Sprintf("[R%d] Enemy deals %d damage (%d LP remaining)", m.combatState.CurrentRound, result.FinalDamage, result.TargetLP))
		} else {
			m.combatState.AddLogEntry(fmt.Sprintf("[R%d] Enemy rolled %d (need %d+) - MISS!", m.combatState.CurrentRound, result.Roll, result.Requirement))
		}

		combat.ProcessRest(m.combatState)
		m.combatState.AddLogEntry(fmt.Sprintf("[R%d] Rested! Endurance restored.", m.combatState.CurrentRound))
		m.needsRest = false
		m.waitingForInput = true
		
		return m, func() tea.Msg {
			return EnemyAttackCompleteMsg{}
		}
	}

	// Check if enemy needs to rest
	if combat.CheckEndurance(m.combatState.EnemyRoundsSinceLastRest, m.combatState.EnemyEnduranceLimit) {
		m.combatState.AddLogEntry(fmt.Sprintf("[R%d] Enemy endurance depleted! Enemy must rest.", m.combatState.CurrentRound))
		m.needsEnemyRest = true
		
		// Player gets free attack while enemy rests
		m.combatState.AddLogEntry(fmt.Sprintf("[R%d] You attack while the enemy rests...", m.combatState.CurrentRound))
		
		// Check if Doombringer is equipped - apply blood price BEFORE attack
		isDoombringerEquipped := m.player.EquippedWeapon != nil && m.player.EquippedWeapon.Name == "Doombringer"
		if isDoombringerEquipped {
			m.player.ModifyLP(-10)
			m.combatState.AddLogEntry(fmt.Sprintf("[R%d] Doombringer thirsts for blood... -10 LP", m.combatState.CurrentRound))
			
			if m.player.CurrentLP <= 0 {
				m.combatState.AddLogEntry("[Defeat] Doombringer has drained your life!")
				m.defeatState = true
				return m, nil
			}
		}
		
		result := combat.ExecutePlayerAttack(m.combatState, m.player, m.roller)
		
		if result.Hit {
			m.combatState.AddLogEntry(fmt.Sprintf("[R%d] You rolled %d (need %d+) - HIT!", m.combatState.CurrentRound, result.Roll, result.Requirement))
			m.combatState.AddLogEntry(fmt.Sprintf("[R%d] Enemy takes %d damage (%d LP remaining)", m.combatState.CurrentRound, result.FinalDamage, result.TargetLP))
			
			// Doombringer soul thirst during enemy rest
			if isDoombringerEquipped && result.FinalDamage > 0 {
				oldLP := m.player.CurrentLP
				healAmount := result.FinalDamage
				
				// If enemy died from this hit, cap healing at enemy's LP before the hit
				if result.TargetLP <= 0 {
					enemyLPBeforeHit := result.TargetLP + result.FinalDamage
					if enemyLPBeforeHit < healAmount {
						healAmount = enemyLPBeforeHit
					}
				}
				
				if m.player.CurrentLP + healAmount > m.player.MaximumLP {
					healAmount = m.player.MaximumLP - m.player.CurrentLP
				}
				
				if healAmount > 0 {
					m.player.ModifyLP(healAmount)
					actualHeal := m.player.CurrentLP - oldLP
					m.combatState.AddLogEntry(fmt.Sprintf("[R%d] Doombringer feeds on pain... +%d LP healed!", m.combatState.CurrentRound, actualHeal))
				}
			}
		} else {
			m.combatState.AddLogEntry(fmt.Sprintf("[R%d] You rolled %d (need %d+) - MISS!", m.combatState.CurrentRound, result.Roll, result.Requirement))
		}
		
		combat.ProcessEnemyRest(m.combatState)
		m.combatState.AddLogEntry(fmt.Sprintf("[R%d] Enemy rested! Enemy endurance restored.", m.combatState.CurrentRound))
		m.needsEnemyRest = false
		
		return m, func() tea.Msg {
			return EnemyAttackCompleteMsg{}
		}
	}

	// Normal enemy turn
	result := combat.ExecuteEnemyAttack(m.combatState, m.player, m.roller)
	
	if result.Hit {
		m.combatState.AddLogEntry(fmt.Sprintf("[R%d] Enemy rolled %d (need %d+) - HIT!", m.combatState.CurrentRound, result.Roll, result.Requirement))
		m.combatState.AddLogEntry(fmt.Sprintf("[R%d] Enemy deals %d damage (%d LP remaining)", m.combatState.CurrentRound, result.FinalDamage, result.TargetLP))
	} else {
		m.combatState.AddLogEntry(fmt.Sprintf("[R%d] Enemy rolled %d (need %d+) - MISS!", m.combatState.CurrentRound, result.Roll, result.Requirement))
	}

	return m, func() tea.Msg {
		return EnemyAttackCompleteMsg{}
	}
}

func (m CombatViewModel) checkCombatState() (CombatViewModel, tea.Cmd) {
	// Check victory
	if combat.CheckVictory(m.combatState) {
		m.combatState.AddLogEntry(fmt.Sprintf("[Victory] %s defeated!", m.combatState.Enemy.Name))
		combat.ResolveCombatVictory(m.player)
		m.combatState.AddLogEntry(fmt.Sprintf("[Victory] Skill increased to %d. Enemies defeated: %d", m.player.Skill, m.player.EnemiesDefeated))
		m.victoryState = true
		return m, nil
	}

	// Check defeat
	if combat.CheckDefeat(m.player, m.combatState) {
		// Check if death save is available
		if !m.combatState.DeathSaveUsed {
			m.combatState.AddLogEntry(fmt.Sprintf("[Critical] Your LP dropped to %d!", m.player.CurrentLP))
			m.combatState.AddLogEntry("[Death Save] Press Enter to roll death save (2d6×10 vs Luck)...")
			m.deathSaveActive = true
			return m, nil
		} else {
			m.combatState.AddLogEntry("[Defeat] You have been defeated!")
			m.defeatState = true
			return m, nil
		}
	}

	// Continue combat - advance turn
	combat.NextTurn(m.combatState)
	m.waitingForInput = m.combatState.PlayerTurn

	// If it's now enemy turn, trigger enemy turn processing
	if !m.combatState.PlayerTurn {
		return m, func() tea.Msg {
			return EnemyTurnMsg{}
		}
	}

	return m, nil
}

// View renders the combat screen.
func (m CombatViewModel) View() string {
	var s strings.Builder
	t := theme.Current()

	s.WriteString("\n")
	s.WriteString(theme.RenderTitle(fmt.Sprintf("COMBAT - Round %d", m.combatState.CurrentRound)))
	s.WriteString("\n\n")

	// Combatant stats - two columns
	s.WriteString(t.Heading.Render("  Fire*Wolf") + strings.Repeat(" ", 30) + t.Heading.Render("Enemy: "+m.combatState.Enemy.Name) + "\n")
	
	// Health bars
	playerHP := theme.RenderHealthBar(m.player.CurrentLP, m.player.MaximumLP, 20)
	enemyHP := theme.RenderHealthBar(m.combatState.Enemy.CurrentLP, m.combatState.Enemy.MaximumLP, 20)
	s.WriteString("  " + playerHP + "    " + enemyHP + "\n\n")
	
	// Stats comparison
	s.WriteString(t.Label.Render(fmt.Sprintf("  STR: %d  SPD: %d  STA: %d", m.player.Strength, m.player.Speed, m.player.Stamina)))
	s.WriteString(strings.Repeat(" ", 10))
	s.WriteString(t.Label.Render(fmt.Sprintf("STR: %d  SPD: %d  STA: %d", m.combatState.Enemy.Strength, m.combatState.Enemy.Speed, m.combatState.Enemy.Stamina)) + "\n")
	
	s.WriteString(t.Label.Render(fmt.Sprintf("  Skill: %d", m.player.Skill)))
	s.WriteString(strings.Repeat(" ", 35))
	s.WriteString(t.Label.Render(fmt.Sprintf("Skill: %d", m.combatState.Enemy.Skill)) + "\n")

	weaponName := "None"
	weaponBonus := 0
	if m.player.EquippedWeapon != nil {
		weaponName = m.player.EquippedWeapon.Name
		weaponBonus = m.player.EquippedWeapon.DamageBonus
	}
	s.WriteString(t.Label.Render(fmt.Sprintf("  Weapon: %s (+%d)", weaponName, weaponBonus)))
	s.WriteString(strings.Repeat(" ", 25-len(weaponName)))
	s.WriteString(t.Label.Render(fmt.Sprintf("Weapon Bonus: +%d", m.combatState.Enemy.WeaponBonus)) + "\n")

	armorProtection := m.player.GetArmorProtection()
	s.WriteString(t.Label.Render(fmt.Sprintf("  Armor: -%d", armorProtection)))
	s.WriteString(strings.Repeat(" ", 35))
	s.WriteString(t.Label.Render(fmt.Sprintf("Armor: -%d", m.combatState.Enemy.ArmorProtection)) + "\n")

	// Endurance status
	if m.combatState.EnduranceLimit > 0 {
		remaining := m.combatState.EnduranceLimit - m.combatState.RoundsSinceLastRest
		if remaining <= 0 {
			s.WriteString("\n" + theme.RenderWarning("ENDURANCE DEPLETED", "Rest required!") + "\n")
		} else if remaining == 1 {
			s.WriteString("\n" + t.WarningMsg.Render(fmt.Sprintf("  ⚠ Endurance: %d round remaining", remaining)) + "\n")
		} else {
			s.WriteString("\n" + t.MutedText.Render(fmt.Sprintf("  Endurance: %d rounds remaining", remaining)) + "\n")
		}
	}

	s.WriteString("\n" + theme.RenderSeparator(60) + "\n")
	s.WriteString(t.Heading.Render("  Combat Log") + "\n")
	s.WriteString(theme.RenderSeparator(60) + "\n")

	// Show last 4 log entries to keep display compact and prevent scrolling
	logStart := len(m.combatState.CombatLog) - 4
	if logStart < 0 {
		logStart = 0
	}
	for i := logStart; i < len(m.combatState.CombatLog); i++ {
		s.WriteString("  " + t.Body.Render(m.combatState.CombatLog[i]) + "\n")
	}

	s.WriteString("\n" + theme.RenderSeparator(60) + "\n")

	// Victory/Defeat messages
	if m.victoryState {
		s.WriteString("\n" + theme.RenderSuccess("VICTORY!") + "\n\n")
		s.WriteString(t.Body.Render("  Press Enter to return to game menu") + "\n")
		return s.String()
	}

	if m.defeatState {
		s.WriteString("\n" + theme.RenderError("DEFEAT", "You have been defeated", "") + "\n\n")
		s.WriteString(t.Body.Render("  Press Enter to return to game menu") + "\n")
		return s.String()
	}

	// Death save prompt
	if m.deathSaveActive {
		s.WriteString("\n" + theme.RenderWarning("DEATH SAVE", "Roll 2d6×10 vs Luck to survive") + "\n\n")
		s.WriteString(t.Emphasis.Render("  Press Enter to attempt death save") + "\n")
		return s.String()
	}

	// Turn indicator and actions
	if m.combatState.PlayerTurn {
		s.WriteString("\n" + t.Heading.Render("  Your Turn") + "\n")
		s.WriteString(theme.RenderSeparator(60) + "\n\n")
		
		for i, action := range m.actions {
			s.WriteString("  " + theme.RenderMenuItem(action, i == m.selectedAction) + "\n")
		}
		
		s.WriteString("\n" + theme.RenderKeyHelp("↑/↓ Select", "Enter Confirm", "Esc Menu") + "\n")
	} else {
		s.WriteString("\n" + t.Heading.Render("  Enemy Turn...") + "\n")
		s.WriteString(theme.RenderSeparator(60) + "\n\n")
		s.WriteString(t.MutedText.Render("  Processing enemy action...") + "\n")
	}

	return s.String()
}

// Helper function for max
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// CombatEndMsg signals that combat has ended.
type CombatEndMsg struct {
	Victory bool
}

// CastSpellMsg signals to switch to spell casting screen during combat.
type CastSpellMsg struct{}

// EnemyTurnMsg signals that the enemy should take a turn.
type EnemyTurnMsg struct{}

// PlayerAttackCompleteMsg signals that the player's attack is complete.
type PlayerAttackCompleteMsg struct{}

// EnemyAttackCompleteMsg signals that the enemy's attack is complete.
type EnemyAttackCompleteMsg struct{}
