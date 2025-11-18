package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/benoit/saga-demonspawn/internal/character"
	"github.com/benoit/saga-demonspawn/internal/combat"
	"github.com/benoit/saga-demonspawn/internal/dice"
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
	deathSaveActive bool

	// Action menu
	actions []string
}

const (
	actionAttack = iota
	actionFlee
	actionTotalActions
)

// NewCombatViewModel creates a new combat view model.
func NewCombatViewModel(player *character.Character, combatState *combat.CombatState, roller dice.Roller) CombatViewModel {
	return CombatViewModel{
		player:          player,
		combatState:     combatState,
		roller:          roller,
		selectedAction:  actionAttack,
		waitingForInput: true,
		victoryState:    false,
		defeatState:     false,
		needsRest:       false,
		deathSaveActive: false,
		actions:         []string{"Attack", "Flee Combat"},
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

		// Execute player attack
		result := combat.ExecutePlayerAttack(m.combatState, m.player, m.roller)
		
		if result.Hit {
			m.combatState.AddLogEntry(fmt.Sprintf("[R%d] You rolled %d (need %d+) - HIT!", m.combatState.CurrentRound, result.Roll, result.Requirement))
			m.combatState.AddLogEntry(fmt.Sprintf("[R%d] Damage: (%d×5) + STR + Weapon - Armor = %d", m.combatState.CurrentRound, result.Roll, result.FinalDamage))
			m.combatState.AddLogEntry(fmt.Sprintf("[R%d] Enemy takes %d damage (%d LP remaining)", m.combatState.CurrentRound, result.FinalDamage, result.TargetLP))
		} else {
			m.combatState.AddLogEntry(fmt.Sprintf("[R%d] You rolled %d (need %d+) - MISS!", m.combatState.CurrentRound, result.Roll, result.Requirement))
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
	}

	return m, nil
}

func (m CombatViewModel) processEnemyTurn() (CombatViewModel, tea.Cmd) {
	// If resting, enemy gets free attack
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

	s.WriteString("═══════════════════════════════════════════════════════════════\n")
	s.WriteString(fmt.Sprintf("                COMBAT: Round %d\n", m.combatState.CurrentRound))
	s.WriteString("═══════════════════════════════════════════════════════════════\n\n")

	// Combatant stats
	s.WriteString("FIRE*WOLF                          ENEMY: " + m.combatState.Enemy.Name + "\n")
	s.WriteString(fmt.Sprintf("LP: %d/%d (%d%%)              LP: %d/%d (%d%%)\n",
		m.player.CurrentLP, m.player.MaximumLP, (m.player.CurrentLP*100)/max(m.player.MaximumLP, 1),
		m.combatState.Enemy.CurrentLP, m.combatState.Enemy.MaximumLP, (m.combatState.Enemy.CurrentLP*100)/max(m.combatState.Enemy.MaximumLP, 1)))
	
	s.WriteString(fmt.Sprintf("STR: %d  SPD: %d  STA: %d       STR: %d  SPD: %d  STA: %d\n",
		m.player.Strength, m.player.Speed, m.player.Stamina,
		m.combatState.Enemy.Strength, m.combatState.Enemy.Speed, m.combatState.Enemy.Stamina))
	
	s.WriteString(fmt.Sprintf("Skill: %d                         Skill: %d\n",
		m.player.Skill, m.combatState.Enemy.Skill))

	weaponName := "None"
	weaponBonus := 0
	if m.player.EquippedWeapon != nil {
		weaponName = m.player.EquippedWeapon.Name
		weaponBonus = m.player.EquippedWeapon.DamageBonus
	}
	s.WriteString(fmt.Sprintf("Weapon: %s (+%d)           Weapon Bonus: +%d\n", weaponName, weaponBonus, m.combatState.Enemy.WeaponBonus))

	armorProtection := m.player.GetArmorProtection()
	s.WriteString(fmt.Sprintf("Armor Protection: %d               Armor Protection: %d\n", armorProtection, m.combatState.Enemy.ArmorProtection))

	// Endurance status
	if m.combatState.EnduranceLimit > 0 {
		remaining := m.combatState.EnduranceLimit - m.combatState.RoundsSinceLastRest
		if remaining <= 0 {
			s.WriteString("\n⚠️  ENDURANCE DEPLETED - Rest required!\n")
		} else if remaining == 1 {
			s.WriteString(fmt.Sprintf("\n⚠️  Endurance: %d round remaining before rest\n", remaining))
		} else {
			s.WriteString(fmt.Sprintf("\nEndurance: %d rounds remaining\n", remaining))
		}
	}

	s.WriteString("\n───────────────────────────────────────────────────────────────\n")
	s.WriteString("COMBAT LOG (Recent)\n")
	s.WriteString("───────────────────────────────────────────────────────────────\n")

	// Show last 6 log entries to keep display compact
	logStart := len(m.combatState.CombatLog) - 6
	if logStart < 0 {
		logStart = 0
	}
	for i := logStart; i < len(m.combatState.CombatLog); i++ {
		s.WriteString(m.combatState.CombatLog[i] + "\n")
	}

	s.WriteString("\n═══════════════════════════════════════════════════════════════\n")

	// Victory/Defeat messages
	if m.victoryState {
		s.WriteString("\n🎉 VICTORY! 🎉\n\n")
		s.WriteString("Press Enter to return to game menu\n")
		return s.String()
	}

	if m.defeatState {
		s.WriteString("\n💀 DEFEAT 💀\n\n")
		s.WriteString("Press Enter to return to game menu\n")
		return s.String()
	}

	// Death save prompt
	if m.deathSaveActive {
		s.WriteString("\n⚠️  DEATH SAVE ⚠️\n\n")
		s.WriteString("Press Enter to attempt death save\n")
		return s.String()
	}

	// Turn indicator and actions
	if m.combatState.PlayerTurn {
		s.WriteString("\nYour Turn\n")
		s.WriteString("═══════════════════════════════════════════════════════════════\n")
		
		for i, action := range m.actions {
			cursor := "  "
			if i == m.selectedAction {
				cursor = "> "
			}
			s.WriteString(fmt.Sprintf("%s%s\n", cursor, action))
		}
		
		s.WriteString("\n[↑/↓: Select | Enter: Confirm | Esc: Menu]\n")
	} else {
		s.WriteString("\nEnemy Turn...\n")
		s.WriteString("═══════════════════════════════════════════════════════════════\n")
		s.WriteString("\n[Processing enemy action...]\n")
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

// EnemyTurnMsg signals that the enemy should take a turn.
type EnemyTurnMsg struct{}

// PlayerAttackCompleteMsg signals that the player's attack is complete.
type PlayerAttackCompleteMsg struct{}

// EnemyAttackCompleteMsg signals that the enemy's attack is complete.
type EnemyAttackCompleteMsg struct{}
