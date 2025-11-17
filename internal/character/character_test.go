package character

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/benoit/saga-demonspawn/internal/items"
)

// TestNew verifies character creation with valid stats.
func TestNew(t *testing.T) {
	char, err := New(50, 60, 55, 45, 70, 40, 35)
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	// Check characteristics are set correctly
	if char.Strength != 50 {
		t.Errorf("Strength = %d; want 50", char.Strength)
	}
	if char.Speed != 60 {
		t.Errorf("Speed = %d; want 60", char.Speed)
	}

	// Check LP calculation (sum of all characteristics)
	expectedLP := 50 + 60 + 55 + 45 + 70 + 40 + 35
	if char.MaximumLP != expectedLP {
		t.Errorf("MaximumLP = %d; want %d", char.MaximumLP, expectedLP)
	}
	if char.CurrentLP != expectedLP {
		t.Errorf("CurrentLP = %d; want %d (should equal MaximumLP)", char.CurrentLP, expectedLP)
	}

	// Check initial values
	if char.Skill != 0 {
		t.Errorf("Skill = %d; want 0", char.Skill)
	}
	if char.CurrentPOW != 0 {
		t.Errorf("CurrentPOW = %d; want 0", char.CurrentPOW)
	}
	if char.MagicUnlocked {
		t.Error("MagicUnlocked = true; want false")
	}

	// Check default equipment
	if char.EquippedWeapon == nil || char.EquippedWeapon.Name != "Sword" {
		t.Error("Expected default weapon to be Sword")
	}
	if char.EquippedArmor == nil || char.EquippedArmor.Name != "None" {
		t.Error("Expected default armor to be None")
	}
	if char.HasShield {
		t.Error("HasShield = true; want false")
	}
}

// TestNewInvalidCharacteristics verifies error handling for invalid stats.
func TestNewInvalidCharacteristics(t *testing.T) {
	tests := []struct {
		name        string
		str, spd, sta, crg, lck, chm, att int
		wantErr     bool
	}{
		{"negative strength", -10, 50, 50, 50, 50, 50, 50, true},
		{"negative speed", 50, -10, 50, 50, 50, 50, 50, true},
		{"exceeds max", 1000, 50, 50, 50, 50, 50, 50, true},
		{"all zeros valid", 0, 0, 0, 0, 0, 0, 0, false},
		{"all max valid", 999, 999, 999, 999, 999, 999, 999, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.str, tt.spd, tt.sta, tt.crg, tt.lck, tt.chm, tt.att)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestModifyCharacteristics verifies characteristic modification methods.
func TestModifyCharacteristics(t *testing.T) {
	char, _ := New(50, 50, 50, 50, 50, 50, 50)

	// Test positive modification
	if err := char.ModifyStrength(10); err != nil {
		t.Errorf("ModifyStrength(10) unexpected error: %v", err)
	}
	if char.Strength != 60 {
		t.Errorf("Strength = %d; want 60", char.Strength)
	}

	// Test negative modification
	if err := char.ModifyStrength(-30); err != nil {
		t.Errorf("ModifyStrength(-30) unexpected error: %v", err)
	}
	if char.Strength != 30 {
		t.Errorf("Strength = %d; want 30", char.Strength)
	}

	// Test error on negative result
	if err := char.ModifyStrength(-40); err == nil {
		t.Error("ModifyStrength(-40) expected error for negative result")
	}
}

// TestModifyLP verifies life point modification.
func TestModifyLP(t *testing.T) {
	char, _ := New(50, 50, 50, 50, 50, 50, 50)
	initialLP := char.CurrentLP

	// Modify down
	char.ModifyLP(-20)
	if char.CurrentLP != initialLP-20 {
		t.Errorf("CurrentLP = %d; want %d", char.CurrentLP, initialLP-20)
	}

	// Can go negative (death)
	char.ModifyLP(-1000)
	if char.CurrentLP >= 0 {
		t.Errorf("CurrentLP = %d; want negative value", char.CurrentLP)
	}

	// Modify back up
	char.SetLP(50)
	if char.CurrentLP != 50 {
		t.Errorf("CurrentLP = %d; want 50", char.CurrentLP)
	}
}

// TestSetMaxLP verifies maximum LP modification.
func TestSetMaxLP(t *testing.T) {
	char, _ := New(50, 50, 50, 50, 50, 50, 50)

	if err := char.SetMaxLP(500); err != nil {
		t.Errorf("SetMaxLP(500) unexpected error: %v", err)
	}
	if char.MaximumLP != 500 {
		t.Errorf("MaximumLP = %d; want 500", char.MaximumLP)
	}

	// Test error on negative
	if err := char.SetMaxLP(-10); err == nil {
		t.Error("SetMaxLP(-10) expected error")
	}
}

// TestSkillModification verifies skill changes.
func TestSkillModification(t *testing.T) {
	char, _ := New(50, 50, 50, 50, 50, 50, 50)

	if err := char.ModifySkill(5); err != nil {
		t.Errorf("ModifySkill(5) unexpected error: %v", err)
	}
	if char.Skill != 5 {
		t.Errorf("Skill = %d; want 5", char.Skill)
	}

	if err := char.SetSkill(25); err != nil {
		t.Errorf("SetSkill(25) unexpected error: %v", err)
	}
	if char.Skill != 25 {
		t.Errorf("Skill = %d; want 25", char.Skill)
	}

	// Test error on negative
	if err := char.SetSkill(-1); err == nil {
		t.Error("SetSkill(-1) expected error")
	}
}

// TestUnlockMagic verifies magic system activation.
func TestUnlockMagic(t *testing.T) {
	char, _ := New(50, 50, 50, 50, 50, 50, 50)

	// Initially locked
	if char.MagicUnlocked {
		t.Error("MagicUnlocked should be false initially")
	}

	// Unlock with initial POW
	if err := char.UnlockMagic(100); err != nil {
		t.Errorf("UnlockMagic(100) unexpected error: %v", err)
	}

	if !char.MagicUnlocked {
		t.Error("MagicUnlocked should be true after unlock")
	}
	if char.CurrentPOW != 100 {
		t.Errorf("CurrentPOW = %d; want 100", char.CurrentPOW)
	}
	if char.MaximumPOW != 100 {
		t.Errorf("MaximumPOW = %d; want 100", char.MaximumPOW)
	}

	// Test error on negative POW
	char2, _ := New(50, 50, 50, 50, 50, 50, 50)
	if err := char2.UnlockMagic(-10); err == nil {
		t.Error("UnlockMagic(-10) expected error")
	}
}

// TestPOWModification verifies power manipulation.
func TestPOWModification(t *testing.T) {
	char, _ := New(50, 50, 50, 50, 50, 50, 50)
	char.UnlockMagic(100)

	// Modify down
	char.ModifyPOW(-25)
	if char.CurrentPOW != 75 {
		t.Errorf("CurrentPOW = %d; want 75", char.CurrentPOW)
	}

	// Cannot go below 0
	char.ModifyPOW(-200)
	if char.CurrentPOW != 0 {
		t.Errorf("CurrentPOW = %d; want 0 (should be clamped)", char.CurrentPOW)
	}

	// Set directly
	char.SetPOW(50)
	if char.CurrentPOW != 50 {
		t.Errorf("CurrentPOW = %d; want 50", char.CurrentPOW)
	}

	// Set max
	if err := char.SetMaxPOW(150); err != nil {
		t.Errorf("SetMaxPOW(150) unexpected error: %v", err)
	}
	if char.MaximumPOW != 150 {
		t.Errorf("MaximumPOW = %d; want 150", char.MaximumPOW)
	}
}

// TestEquipment verifies equipment changes.
func TestEquipment(t *testing.T) {
	char, _ := New(50, 50, 50, 50, 50, 50, 50)

	// Equip different weapon
	char.EquipWeapon(&items.WeaponAxe)
	if char.EquippedWeapon.Name != "Axe" {
		t.Errorf("EquippedWeapon = %s; want Axe", char.EquippedWeapon.Name)
	}

	// Equip armor
	char.EquipArmor(&items.ArmorChain)
	if char.EquippedArmor.Name != "Chain Mail" {
		t.Errorf("EquippedArmor = %s; want Chain Mail", char.EquippedArmor.Name)
	}

	// Toggle shield
	char.ToggleShield()
	if !char.HasShield {
		t.Error("HasShield should be true after toggle")
	}
	char.ToggleShield()
	if char.HasShield {
		t.Error("HasShield should be false after second toggle")
	}
}

// TestGetArmorProtection verifies armor protection calculation.
func TestGetArmorProtection(t *testing.T) {
	tests := []struct {
		name           string
		armor          *items.Armor
		hasShield      bool
		wantProtection int
	}{
		{"no armor no shield", &items.ArmorNone, false, 0},
		{"leather armor only", &items.ArmorLeather, false, 5},
		{"shield only", &items.ArmorNone, true, 7},
		{"leather + shield", &items.ArmorLeather, true, 10}, // 5 + 5 (reduced shield)
		{"chain + shield", &items.ArmorChain, true, 13},     // 8 + 5 (reduced shield)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			char, _ := New(50, 50, 50, 50, 50, 50, 50)
			char.EquipArmor(tt.armor)
			if tt.hasShield {
				char.ToggleShield()
			}

			got := char.GetArmorProtection()
			if got != tt.wantProtection {
				t.Errorf("GetArmorProtection() = %d; want %d", got, tt.wantProtection)
			}
		})
	}
}

// TestGetWeaponDamageBonus verifies weapon damage bonus retrieval.
func TestGetWeaponDamageBonus(t *testing.T) {
	char, _ := New(50, 50, 50, 50, 50, 50, 50)

	// Default sword
	if bonus := char.GetWeaponDamageBonus(); bonus != 10 {
		t.Errorf("GetWeaponDamageBonus() = %d; want 10 (sword)", bonus)
	}

	// Equip Doombringer
	char.EquipWeapon(&items.WeaponDoombringer)
	if bonus := char.GetWeaponDamageBonus(); bonus != 20 {
		t.Errorf("GetWeaponDamageBonus() = %d; want 20 (Doombringer)", bonus)
	}
}

// TestIsAlive verifies alive status check.
func TestIsAlive(t *testing.T) {
	char, _ := New(50, 50, 50, 50, 50, 50, 50)

	if !char.IsAlive() {
		t.Error("IsAlive() = false; want true")
	}

	char.SetLP(0)
	if char.IsAlive() {
		t.Error("IsAlive() = true; want false when LP = 0")
	}

	char.SetLP(-10)
	if char.IsAlive() {
		t.Error("IsAlive() = true; want false when LP < 0")
	}

	char.SetLP(1)
	if !char.IsAlive() {
		t.Error("IsAlive() = false; want true when LP > 0")
	}
}

// TestIncrementEnemiesDefeated verifies enemy counter.
func TestIncrementEnemiesDefeated(t *testing.T) {
	char, _ := New(50, 50, 50, 50, 50, 50, 50)

	if char.EnemiesDefeated != 0 {
		t.Errorf("EnemiesDefeated = %d; want 0", char.EnemiesDefeated)
	}

	char.IncrementEnemiesDefeated()
	if char.EnemiesDefeated != 1 {
		t.Errorf("EnemiesDefeated = %d; want 1", char.EnemiesDefeated)
	}

	char.IncrementEnemiesDefeated()
	char.IncrementEnemiesDefeated()
	if char.EnemiesDefeated != 3 {
		t.Errorf("EnemiesDefeated = %d; want 3", char.EnemiesDefeated)
	}
}

// TestSaveAndLoad verifies character persistence.
func TestSaveAndLoad(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create and configure a character
	original, _ := New(65, 55, 60, 50, 75, 45, 40)
	original.ModifyLP(-20)
	original.ModifySkill(5)
	original.UnlockMagic(80)
	original.EquipWeapon(&items.WeaponAxe)
	original.EquipArmor(&items.ArmorChain)
	original.ToggleShield()

	// Save character
	if err := original.Save(tempDir); err != nil {
		t.Fatalf("Save() unexpected error: %v", err)
	}

	// Find the saved file
	files, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to read temp dir: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("Expected 1 save file, got %d", len(files))
	}

	// Load character
	savePath := filepath.Join(tempDir, files[0].Name())
	loaded, err := Load(savePath)
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	// Verify all fields match
	if loaded.Strength != original.Strength {
		t.Errorf("Loaded Strength = %d; want %d", loaded.Strength, original.Strength)
	}
	if loaded.CurrentLP != original.CurrentLP {
		t.Errorf("Loaded CurrentLP = %d; want %d", loaded.CurrentLP, original.CurrentLP)
	}
	if loaded.Skill != original.Skill {
		t.Errorf("Loaded Skill = %d; want %d", loaded.Skill, original.Skill)
	}
	if loaded.MagicUnlocked != original.MagicUnlocked {
		t.Errorf("Loaded MagicUnlocked = %v; want %v", loaded.MagicUnlocked, original.MagicUnlocked)
	}
	if loaded.CurrentPOW != original.CurrentPOW {
		t.Errorf("Loaded CurrentPOW = %d; want %d", loaded.CurrentPOW, original.CurrentPOW)
	}
	if loaded.EquippedWeapon.Name != original.EquippedWeapon.Name {
		t.Errorf("Loaded weapon = %s; want %s", loaded.EquippedWeapon.Name, original.EquippedWeapon.Name)
	}
	if loaded.EquippedArmor.Name != original.EquippedArmor.Name {
		t.Errorf("Loaded armor = %s; want %s", loaded.EquippedArmor.Name, original.EquippedArmor.Name)
	}
	if loaded.HasShield != original.HasShield {
		t.Errorf("Loaded HasShield = %v; want %v", loaded.HasShield, original.HasShield)
	}
}

// TestLoadNonexistentFile verifies error handling for missing files.
func TestLoadNonexistentFile(t *testing.T) {
	_, err := Load("/nonexistent/path/character.json")
	if err == nil {
		t.Error("Load() expected error for nonexistent file")
	}
}
