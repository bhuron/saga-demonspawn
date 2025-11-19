package combat

import (
	"testing"

	"github.com/benoit/saga-demonspawn/internal/character"
	"github.com/benoit/saga-demonspawn/internal/items"
)

// MockRoller implements dice.Roller with fixed results for testing.
type MockRoller struct {
	NextRoll int
	Rolls    []int
	Index    int
}

func (m *MockRoller) Roll2D6() int {
	if m.NextRoll > 0 {
		return m.NextRoll
	}
	if m.Index < len(m.Rolls) {
		result := m.Rolls[m.Index]
		m.Index++
		return result
	}
	return 7 // Default
}

func (m *MockRoller) Roll1D6() int {
	if m.NextRoll > 0 {
		return m.NextRoll / 2
	}
	return 3 // Default
}

func (m *MockRoller) RollCharacteristic() int {
	return 64 // Default characteristic
}

func (m *MockRoller) SetSeed(seed int64) {
	// No-op for mock
}

func TestNewEnemy(t *testing.T) {
	tests := []struct {
		name            string
		enemyName       string
		str, spd, sta   int
		crg, lck, skill int
		currentLP, maxLP int
		weaponBonus, armorProtection int
		isDemonspawn    bool
		wantErr         bool
	}{
		{
			name:            "Valid enemy",
			enemyName:       "Goblin",
			str:             40, spd: 35, sta: 30,
			crg:             25, lck: 20, skill: 0,
			currentLP:       150, maxLP: 150,
			weaponBonus:     5, armorProtection: 0,
			isDemonspawn:    false,
			wantErr:         false,
		},
		{
			name:            "Empty name",
			enemyName:       "",
			str:             40, spd: 35, sta: 30,
			crg:             25, lck: 20, skill: 0,
			currentLP:       150, maxLP: 150,
			weaponBonus:     5, armorProtection: 0,
			isDemonspawn:    false,
			wantErr:         true,
		},
		{
			name:            "Negative strength",
			enemyName:       "Goblin",
			str:             -1, spd: 35, sta: 30,
			crg:             25, lck: 20, skill: 0,
			currentLP:       150, maxLP: 150,
			weaponBonus:     5, armorProtection: 0,
			isDemonspawn:    false,
			wantErr:         true,
		},
		{
			name:            "Zero max LP",
			enemyName:       "Goblin",
			str:             40, spd: 35, sta: 30,
			crg:             25, lck: 20, skill: 0,
			currentLP:       0, maxLP: 0,
			weaponBonus:     5, armorProtection: 0,
			isDemonspawn:    false,
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enemy, err := NewEnemy(
				tt.enemyName,
				tt.str, tt.spd, tt.sta,
				tt.crg, tt.lck, tt.skill,
				tt.currentLP, tt.maxLP,
				tt.weaponBonus, tt.armorProtection,
				tt.isDemonspawn,
			)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewEnemy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && enemy == nil {
				t.Error("NewEnemy() returned nil enemy without error")
			}

			if !tt.wantErr && enemy.Name != tt.enemyName {
				t.Errorf("NewEnemy() name = %v, want %v", enemy.Name, tt.enemyName)
			}
		})
	}
}

func TestCalculateInitiative(t *testing.T) {
	player, _ := character.New(64, 56, 72, 48, 80, 40, 56)
	enemy, _ := NewEnemy("Goblin", 40, 35, 30, 25, 20, 0, 150, 150, 5, 0, false)

	tests := []struct {
		name       string
		playerRoll int
		enemyRoll  int
		wantPlayerFirst bool
	}{
		{
			name:       "Player wins initiative",
			playerRoll: 8,
			enemyRoll:  5,
			wantPlayerFirst: true, // 8+56+48+80=192 > 5+35+25+20=85
		},
		{
			name:       "Player wins even with low roll",
			playerRoll: 2,
			enemyRoll:  12,
			wantPlayerFirst: true, // 2+56+48+80=186 > 12+35+25+20=92 (player stats are much higher)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roller := &MockRoller{Rolls: []int{tt.playerRoll, tt.enemyRoll}}
			_, _, playerFirst := CalculateInitiative(player, enemy, roller)

			if playerFirst != tt.wantPlayerFirst {
				t.Errorf("CalculateInitiative() playerFirst = %v, want %v", playerFirst, tt.wantPlayerFirst)
			}
		})
	}
}

func TestCalculateToHitRequirement(t *testing.T) {
	tests := []struct {
		name        string
		skill       int
		luck        int
		wantReq     int
	}{
		{
			name:    "Base requirement - no modifiers",
			skill:   0,
			luck:    50,
			wantReq: 7,
		},
		{
			name:    "Skill 10 - one modifier",
			skill:   10,
			luck:    50,
			wantReq: 6,
		},
		{
			name:    "Skill 25 - two modifiers",
			skill:   25,
			luck:    50,
			wantReq: 5,
		},
		{
			name:    "Luck 72 - one modifier",
			skill:   0,
			luck:    72,
			wantReq: 6,
		},
		{
			name:    "Skill 20 + Luck 80 - three modifiers",
			skill:   20,
			luck:    80,
			wantReq: 4,
		},
		{
			name:    "High skill - capped at 2",
			skill:   60,
			luck:    90,
			wantReq: 2,
		},
		{
			name:    "Exactly at luck threshold",
			skill:   0,
			luck:    72,
			wantReq: 6,
		},
		{
			name:    "Just below luck threshold",
			skill:   0,
			luck:    71,
			wantReq: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateToHitRequirement(tt.skill, tt.luck)
			if got != tt.wantReq {
				t.Errorf("CalculateToHitRequirement(%d, %d) = %d, want %d", tt.skill, tt.luck, got, tt.wantReq)
			}
		})
	}
}

func TestCalculateDamage(t *testing.T) {
	tests := []struct {
		name        string
		roll        int
		strength    int
		weaponBonus int
		wantDamage  int
	}{
		{
			name:        "Roll 7, STR 64, Sword +10",
			roll:        7,
			strength:    64,
			weaponBonus: 10,
			wantDamage:  75, // (7*5) + (6*5) + 10 = 35 + 30 + 10 = 75
		},
		{
			name:        "Roll 9, STR 64, Sword +10",
			roll:        9,
			strength:    64,
			weaponBonus: 10,
			wantDamage:  85, // (9*5) + (6*5) + 10 = 45 + 30 + 10 = 85
		},
		{
			name:        "Roll 12, STR 80, Axe +15",
			roll:        12,
			strength:    80,
			weaponBonus: 15,
			wantDamage:  115, // (12*5) + (8*5) + 15 = 60 + 40 + 15 = 115
		},
		{
			name:        "Roll 2, STR 30, Dagger +5",
			roll:        2,
			strength:    30,
			weaponBonus: 5,
			wantDamage:  30, // (2*5) + (3*5) + 5 = 10 + 15 + 5 = 30
		},
		{
			name:        "Roll 7, STR 5, No weapon",
			roll:        7,
			strength:    5,
			weaponBonus: 0,
			wantDamage:  35, // (7*5) + (0*5) + 0 = 35
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateDamage(tt.roll, tt.strength, tt.weaponBonus)
			if got != tt.wantDamage {
				t.Errorf("CalculateDamage(%d, %d, %d) = %d, want %d", tt.roll, tt.strength, tt.weaponBonus, got, tt.wantDamage)
			}
		})
	}
}

func TestApplyArmorReduction(t *testing.T) {
	tests := []struct {
		name            string
		damage          int
		armorProtection int
		wantFinal       int
	}{
		{
			name:            "Damage exceeds armor",
			damage:          75,
			armorProtection: 8,
			wantFinal:       67,
		},
		{
			name:            "Armor reduces to zero",
			damage:          20,
			armorProtection: 20,
			wantFinal:       0,
		},
		{
			name:            "Armor exceeds damage",
			damage:          5,
			armorProtection: 12,
			wantFinal:       0,
		},
		{
			name:            "No armor",
			damage:          50,
			armorProtection: 0,
			wantFinal:       50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ApplyArmorReduction(tt.damage, tt.armorProtection)
			if got != tt.wantFinal {
				t.Errorf("ApplyArmorReduction(%d, %d) = %d, want %d", tt.damage, tt.armorProtection, got, tt.wantFinal)
			}
		})
	}
}

func TestCheckEndurance(t *testing.T) {
	tests := []struct {
		name                string
		roundsSinceLastRest int
		enduranceLimit      int
		wantRest            bool
	}{
		{
			name:                "Below limit",
			roundsSinceLastRest: 2,
			enduranceLimit:      3,
			wantRest:            false,
		},
		{
			name:                "At limit",
			roundsSinceLastRest: 3,
			enduranceLimit:      3,
			wantRest:            true,
		},
		{
			name:                "Above limit",
			roundsSinceLastRest: 4,
			enduranceLimit:      3,
			wantRest:            true,
		},
		{
			name:                "Zero endurance",
			roundsSinceLastRest: 1,
			enduranceLimit:      0,
			wantRest:            false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckEndurance(tt.roundsSinceLastRest, tt.enduranceLimit)
			if got != tt.wantRest {
				t.Errorf("CheckEndurance(%d, %d) = %v, want %v", tt.roundsSinceLastRest, tt.enduranceLimit, got, tt.wantRest)
			}
		})
	}
}

func TestExecuteDeathSave(t *testing.T) {
	tests := []struct {
		name    string
		luck    int
		roll    int
		wantSuccess bool
	}{
		{
			name:    "Success - roll below luck",
			luck:    80,
			roll:    6, // Will be 60 after *10
			wantSuccess: true,
		},
		{
			name:    "Success - roll equals luck",
			luck:    80,
			roll:    8, // Will be 80 after *10
			wantSuccess: true,
		},
		{
			name:    "Failure - roll above luck",
			luck:    80,
			roll:    10, // Will be 100 after *10
			wantSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roller := &MockRoller{NextRoll: tt.roll}
			rollResult, success := ExecuteDeathSave(tt.luck, roller)

			if success != tt.wantSuccess {
				t.Errorf("ExecuteDeathSave(%d) success = %v, want %v", tt.luck, success, tt.wantSuccess)
			}

			expectedRoll := tt.roll * 10
			if rollResult != expectedRoll {
				t.Errorf("ExecuteDeathSave(%d) roll = %d, want %d", tt.luck, rollResult, expectedRoll)
			}
		})
	}
}

func TestExecutePlayerAttack(t *testing.T) {
	player, _ := character.New(64, 56, 72, 48, 80, 40, 56)
	player.EquipWeapon(&items.WeaponSword)
	
	enemy, _ := NewEnemy("Goblin", 40, 35, 30, 25, 20, 0, 150, 150, 5, 8, false)
	cs := NewCombatState(enemy, 3)

	tests := []struct {
		name        string
		roll        int
		wantHit     bool
		wantDamage  int
	}{
		{
			name:    "Hit and damage",
			roll:    9,
			wantHit: true,
			wantDamage: 77, // (9*5) + (6*5) + 10 - 8 = 45 + 30 + 10 - 8 = 77
		},
		{
			name:    "Miss",
			roll:    5,
			wantHit: false,
			wantDamage: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset enemy HP
			enemy.CurrentLP = 150
			
			roller := &MockRoller{NextRoll: tt.roll}
			result := ExecutePlayerAttack(cs, player, roller)

			if result.Hit != tt.wantHit {
				t.Errorf("ExecutePlayerAttack() hit = %v, want %v", result.Hit, tt.wantHit)
			}

			if result.Hit && result.FinalDamage != tt.wantDamage {
				t.Errorf("ExecutePlayerAttack() damage = %d, want %d", result.FinalDamage, tt.wantDamage)
			}

			if result.Hit {
				expectedLP := 150 - tt.wantDamage
				if enemy.CurrentLP != expectedLP {
					t.Errorf("Enemy LP = %d, want %d", enemy.CurrentLP, expectedLP)
				}
			}
		})
	}
}

func TestExecuteEnemyAttack(t *testing.T) {
	player, _ := character.New(64, 56, 72, 48, 80, 40, 56)
	player.EquipArmor(&items.ArmorChain)
	player.HasShield = true
	
	enemy, _ := NewEnemy("Goblin", 40, 35, 30, 25, 20, 0, 150, 150, 5, 0, false)
	cs := NewCombatState(enemy, 3)

	roller := &MockRoller{NextRoll: 8}
	initialLP := player.CurrentLP
	
	result := ExecuteEnemyAttack(cs, player, roller)

	if !result.Hit {
		t.Error("ExecuteEnemyAttack() should have hit")
	}

	// Damage: (8*5) + (4*5) + 5 = 40 + 20 + 5 = 65
	// Armor: ChainMail (8) + Shield with armor (5) = 13
	// Final: 65 - 13 = 52
	expectedDamage := 52
	if result.FinalDamage != expectedDamage {
		t.Errorf("ExecuteEnemyAttack() damage = %d, want %d", result.FinalDamage, expectedDamage)
	}

	expectedLP := initialLP - expectedDamage
	if player.CurrentLP != expectedLP {
		t.Errorf("Player LP = %d, want %d", player.CurrentLP, expectedLP)
	}
}

func TestStartCombat(t *testing.T) {
	player, _ := character.New(64, 56, 72, 48, 80, 40, 56)
	enemy, _ := NewEnemy("Goblin", 40, 35, 30, 25, 20, 0, 150, 150, 5, 0, false)
	
	roller := &MockRoller{Rolls: []int{8, 5}}
	cs := StartCombat(player, enemy, roller)

	if !cs.IsActive {
		t.Error("StartCombat() should create active combat")
	}

	if cs.CurrentRound != 1 {
		t.Errorf("StartCombat() round = %d, want 1", cs.CurrentRound)
	}

	// Player should win initiative: 8+56+48+80 > 5+35+25+20
	if !cs.PlayerFirstStrike {
		t.Error("StartCombat() player should have won initiative")
	}

	if !cs.PlayerTurn {
		t.Error("StartCombat() should be player's turn")
	}

	expectedEndurance := player.Stamina / 10
	if cs.EnduranceLimit != expectedEndurance {
		t.Errorf("StartCombat() endurance = %d, want %d", cs.EnduranceLimit, expectedEndurance)
	}
}

func TestNextTurn(t *testing.T) {
	enemy, _ := NewEnemy("Goblin", 40, 35, 30, 25, 20, 0, 150, 150, 5, 0, false)
	cs := NewCombatState(enemy, 3)
	cs.PlayerFirstStrike = true
	cs.PlayerTurn = true
	cs.CurrentRound = 1

	// First turn switch (player to enemy)
	NextTurn(cs)
	if cs.PlayerTurn {
		t.Error("NextTurn() should switch to enemy turn")
	}
	if cs.CurrentRound != 1 {
		t.Error("NextTurn() should not increment round on first switch")
	}

	// Second turn switch (enemy to player) - should increment round
	NextTurn(cs)
	if !cs.PlayerTurn {
		t.Error("NextTurn() should switch to player turn")
	}
	if cs.CurrentRound != 2 {
		t.Error("NextTurn() should increment round when returning to first striker")
	}
	if cs.RoundsSinceLastRest != 1 {
		t.Error("NextTurn() should increment rounds since rest")
	}
}

func TestAttemptDeathSave(t *testing.T) {
	player, _ := character.New(64, 56, 72, 48, 80, 40, 56)
	player.ModifyLP(-player.CurrentLP) // Reduce to 0
	
	enemy, _ := NewEnemy("Goblin", 40, 35, 30, 25, 20, 0, 100, 150, 5, 0, false)
	cs := NewCombatState(enemy, 3)
	cs.PlayerTurn = false
	cs.PlayerFirstStrike = true

	// Successful death save
	roller := &MockRoller{Rolls: []int{7, 8, 5}} // First roll for death save, next two for initiative
	roll, success := AttemptDeathSave(player, cs, roller)

	if !success {
		t.Error("AttemptDeathSave() should succeed with roll 70 vs luck 80")
	}

	if roll != 70 {
		t.Errorf("AttemptDeathSave() roll = %d, want 70", roll)
	}

	if player.CurrentLP != player.MaximumLP {
		t.Errorf("AttemptDeathSave() player LP = %d, want max %d", player.CurrentLP, player.MaximumLP)
	}

	if !cs.DeathSaveUsed {
		t.Error("AttemptDeathSave() should mark death save as used")
	}

	if cs.CurrentRound != 1 {
		t.Error("AttemptDeathSave() should reset round to 1")
	}

	// Second death save should fail (already used)
	player.ModifyLP(-player.CurrentLP)
	roll2, success2 := AttemptDeathSave(player, cs, roller)

	if success2 {
		t.Error("AttemptDeathSave() should fail when already used")
	}

	if roll2 != 0 {
		t.Errorf("AttemptDeathSave() second attempt roll should be 0, got %d", roll2)
	}
}

func TestCheckVictory(t *testing.T) {
	enemy, _ := NewEnemy("Goblin", 40, 35, 30, 25, 20, 0, 0, 150, 5, 0, false)
	cs := NewCombatState(enemy, 3)

	if !CheckVictory(cs) {
		t.Error("CheckVictory() should return true when enemy LP is 0")
	}

	cs.Enemy.CurrentLP = 1
	if CheckVictory(cs) {
		t.Error("CheckVictory() should return false when enemy LP is positive")
	}
}

func TestCheckDefeat(t *testing.T) {
	player, _ := character.New(64, 56, 72, 48, 80, 40, 56)
	enemy, _ := NewEnemy("Goblin", 40, 35, 30, 25, 20, 0, 150, 150, 5, 0, false)
	cs := NewCombatState(enemy, 3)

	player.ModifyLP(-player.CurrentLP) // Reduce to 0

	if !CheckDefeat(player, cs) {
		t.Error("CheckDefeat() should return true when player LP is 0")
	}

	player.SetLP(50)
	if CheckDefeat(player, cs) {
		t.Error("CheckDefeat() should return false when player LP is positive")
	}
}

func TestResolveCombatVictory(t *testing.T) {
	player, _ := character.New(64, 56, 72, 48, 80, 40, 56)
	initialDefeated := player.EnemiesDefeated
	initialSkill := player.Skill

	ResolveCombatVictory(player)

	if player.EnemiesDefeated != initialDefeated+1 {
		t.Errorf("ResolveCombatVictory() enemies defeated = %d, want %d", player.EnemiesDefeated, initialDefeated+1)
	}

	if player.Skill != initialSkill+1 {
		t.Errorf("ResolveCombatVictory() skill = %d, want %d", player.Skill, initialSkill+1)
	}
}
