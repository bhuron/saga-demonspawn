package ui

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/benoit/saga-demonspawn/pkg/ui/theme"
)

// CombatSetupModel handles manual enemy entry before combat.
type CombatSetupModel struct {
	// Enemy data fields
	name            string
	strength        string
	speed           string
	stamina         string
	courage         string
	luck            string
	skill           string
	currentLP       string
	maximumLP       string
	weaponBonus     string
	armorProtection string
	isDemonspawn    bool

	// UI state
	focusedField int
	inputMode    bool
	errorMsg     string
	fields       []string
}

const (
	fieldName = iota
	fieldStrength
	fieldSpeed
	fieldStamina
	fieldCourage
	fieldLuck
	fieldSkill
	fieldCurrentLP
	fieldMaximumLP
	fieldWeaponBonus
	fieldArmorProtection
	fieldIsDemonspawn
	fieldStartCombat
	fieldTotalFields
)

// NewCombatSetupModel creates a new combat setup model.
func NewCombatSetupModel() CombatSetupModel {
	return CombatSetupModel{
		name:            "",
		strength:        "",
		speed:           "",
		stamina:         "",
		courage:         "",
		luck:            "",
		skill:           "",
		currentLP:       "",
		maximumLP:       "",
		weaponBonus:     "",
		armorProtection: "",
		isDemonspawn:    false,
		focusedField:    fieldName,
		inputMode:       false,
		errorMsg:        "",
		fields: []string{
			"Name", "Strength", "Speed", "Stamina", "Courage",
			"Luck", "Skill", "Current LP", "Maximum LP",
			"Weapon Bonus", "Armor Protection", "Demonspawn?",
			"Start Combat",
		},
	}
}

// Reset clears all fields.
func (m *CombatSetupModel) Reset() {
	*m = NewCombatSetupModel()
}

// GetFieldValue returns the current value of the focused field.
func (m *CombatSetupModel) GetFieldValue(field int) string {
	switch field {
	case fieldName:
		return m.name
	case fieldStrength:
		return m.strength
	case fieldSpeed:
		return m.speed
	case fieldStamina:
		return m.stamina
	case fieldCourage:
		return m.courage
	case fieldLuck:
		return m.luck
	case fieldSkill:
		return m.skill
	case fieldCurrentLP:
		return m.currentLP
	case fieldMaximumLP:
		return m.maximumLP
	case fieldWeaponBonus:
		return m.weaponBonus
	case fieldArmorProtection:
		return m.armorProtection
	case fieldIsDemonspawn:
		if m.isDemonspawn {
			return "Yes"
		}
		return "No"
	default:
		return ""
	}
}

// SetFieldValue sets the value of the specified field.
func (m *CombatSetupModel) SetFieldValue(field int, value string) {
	switch field {
	case fieldName:
		m.name = value
	case fieldStrength:
		m.strength = value
	case fieldSpeed:
		m.speed = value
	case fieldStamina:
		m.stamina = value
	case fieldCourage:
		m.courage = value
	case fieldLuck:
		m.luck = value
	case fieldSkill:
		m.skill = value
	case fieldCurrentLP:
		m.currentLP = value
	case fieldMaximumLP:
		m.maximumLP = value
	case fieldWeaponBonus:
		m.weaponBonus = value
	case fieldArmorProtection:
		m.armorProtection = value
	}
}

// ValidateAndPrepare validates all fields and returns error if invalid.
func (m *CombatSetupModel) ValidateAndPrepare() error {
	if strings.TrimSpace(m.name) == "" {
		return fmt.Errorf("enemy name is required")
	}

	// Validate numeric fields
	if _, err := strconv.Atoi(m.strength); err != nil || m.strength == "" {
		return fmt.Errorf("strength must be a valid number")
	}
	if _, err := strconv.Atoi(m.speed); err != nil || m.speed == "" {
		return fmt.Errorf("speed must be a valid number")
	}
	if _, err := strconv.Atoi(m.stamina); err != nil || m.stamina == "" {
		return fmt.Errorf("stamina must be a valid number")
	}
	if _, err := strconv.Atoi(m.courage); err != nil || m.courage == "" {
		return fmt.Errorf("courage must be a valid number")
	}
	if _, err := strconv.Atoi(m.luck); err != nil || m.luck == "" {
		return fmt.Errorf("luck must be a valid number")
	}
	if _, err := strconv.Atoi(m.skill); err != nil || m.skill == "" {
		return fmt.Errorf("skill must be a valid number")
	}
	if _, err := strconv.Atoi(m.currentLP); err != nil || m.currentLP == "" {
		return fmt.Errorf("current LP must be a valid number")
	}
	if _, err := strconv.Atoi(m.maximumLP); err != nil || m.maximumLP == "" {
		return fmt.Errorf("maximum LP must be a valid number")
	}
	if _, err := strconv.Atoi(m.weaponBonus); err != nil || m.weaponBonus == "" {
		return fmt.Errorf("weapon bonus must be a valid number")
	}
	if _, err := strconv.Atoi(m.armorProtection); err != nil || m.armorProtection == "" {
		return fmt.Errorf("armor protection must be a valid number")
	}

	return nil
}

// GetEnemyData returns the parsed enemy data as integers.
func (m *CombatSetupModel) GetEnemyData() (name string, str, spd, sta, crg, lck, skill, currentLP, maxLP, weaponBonus, armorProtection int, isDemonspawn bool) {
	name = strings.TrimSpace(m.name)
	str, _ = strconv.Atoi(m.strength)
	spd, _ = strconv.Atoi(m.speed)
	sta, _ = strconv.Atoi(m.stamina)
	crg, _ = strconv.Atoi(m.courage)
	lck, _ = strconv.Atoi(m.luck)
	skill, _ = strconv.Atoi(m.skill)
	currentLP, _ = strconv.Atoi(m.currentLP)
	maxLP, _ = strconv.Atoi(m.maximumLP)
	weaponBonus, _ = strconv.Atoi(m.weaponBonus)
	armorProtection, _ = strconv.Atoi(m.armorProtection)
	isDemonspawn = m.isDemonspawn
	return
}

// Update handles combat setup input.
func (m CombatSetupModel) Update(msg tea.Msg) (CombatSetupModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.inputMode {
			return m.handleInputMode(msg)
		}
		return m.handleNavigationMode(msg)
	}
	return m, nil
}

func (m CombatSetupModel) handleNavigationMode(msg tea.KeyMsg) (CombatSetupModel, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.focusedField > 0 {
			m.focusedField--
		}
	case "down", "j":
		if m.focusedField < fieldTotalFields-1 {
			m.focusedField++
		}
	case "enter":
		// Handle toggle for Demonspawn field
		if m.focusedField == fieldIsDemonspawn {
			m.isDemonspawn = !m.isDemonspawn
			return m, nil
		}
		// Enter input mode for editable fields
		if m.focusedField < fieldIsDemonspawn {
			m.inputMode = true
			return m, nil
		}
		// Start combat if on that field
		if m.focusedField == fieldStartCombat {
			if err := m.ValidateAndPrepare(); err != nil {
				m.errorMsg = err.Error()
				return m, nil
			}
			// Signal to parent to start combat
			return m, func() tea.Msg {
				return CombatStartMsg{}
			}
		}
	}
	return m, nil
}

func (m CombatSetupModel) handleInputMode(msg tea.KeyMsg) (CombatSetupModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.inputMode = false
		return m, nil
	case "enter":
		m.inputMode = false
		m.errorMsg = "" // Clear error on successful edit
		return m, nil
	case "backspace":
		currentValue := m.GetFieldValue(m.focusedField)
		if len(currentValue) > 0 {
			m.SetFieldValue(m.focusedField, currentValue[:len(currentValue)-1])
		}
	default:
		// Only allow printable characters
		if len(msg.String()) == 1 {
			currentValue := m.GetFieldValue(m.focusedField)
			m.SetFieldValue(m.focusedField, currentValue+msg.String())
		}
	}
	return m, nil
}

// View renders the combat setup screen.
func (m CombatSetupModel) View() string {
	var s strings.Builder
	t := theme.Current()

	s.WriteString("\n")
	s.WriteString(theme.RenderTitle("COMBAT SETUP - Enter Enemy"))
	s.WriteString("\n\n")

	// Render fields
	for i := 0; i < fieldStartCombat; i++ {
		focused := i == m.focusedField
		editing := focused && m.inputMode

		fieldName := m.fields[i]
		value := m.GetFieldValue(i)

		if value == "" {
			value = t.MutedText.Render("___")
		} else if editing {
			value = t.Emphasis.Render(value + "_")
		} else {
			value = t.Value.Render(value)
		}

		if focused && !editing {
			s.WriteString("  " + theme.RenderMenuItem(fmt.Sprintf("%-20s: %s", fieldName, value), true) + "\n")
		} else if editing {
			s.WriteString("  " + t.Emphasis.Render("* ") + t.Label.Render(fmt.Sprintf("%-20s: ", fieldName)) + value + "\n")
		} else {
			s.WriteString("  " + t.Label.Render(fmt.Sprintf("%-20s: ", fieldName)) + value + "\n")
		}
	}

	s.WriteString("\n")

	// Start Combat button
	focused := m.focusedField == fieldStartCombat
	if focused {
		s.WriteString("  " + t.ButtonFocus.Render(" Start Combat ") + "\n")
	} else {
		s.WriteString("  " + t.Button.Render(" Start Combat ") + "\n")
	}

	s.WriteString("\n")

	// Instructions
	if m.inputMode {
		s.WriteString(theme.RenderKeyHelp("Type to edit", "Enter Confirm", "Esc Cancel") + "\n")
	} else {
		s.WriteString(theme.RenderKeyHelp("↑/↓ Navigate", "Enter Edit/Confirm", "Esc Back", "? Help") + "\n")
	}

	// Error message
	if m.errorMsg != "" {
		s.WriteString("\n" + theme.RenderError("Input Error", m.errorMsg, "Check your values and try again") + "\n")
	}

	return s.String()
}

// CombatStartMsg signals that combat should begin.
type CombatStartMsg struct{}
