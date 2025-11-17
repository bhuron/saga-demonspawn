# Sagas of the Demonspawn - Complete Game Rules

## 1. Your Character (Fire*Wolf)

Your capabilities are defined by seven characteristics, each a percentage value you roll at the start:
- **STRENGTH (STR)**
- **SPEED (SPD)**
- **STAMINA (STA)**
- **COURAGE (CRG)**
- **LUCK (LCK)**
- **CHARM (CHM)**
- **ATTRACTION (ATT)**

**Life Points (LP):** Your health. Starting LP = STR + SPD + STA + CRG + LCK + CHM + ATT.
**Skill (SKL):** Starts at 0. You gain 1 SKL point for every enemy you kill. SKL is added to your total LP and helps you hit in combat.

---

## 2. The Flow of Combat

A fight is a series of rounds. In each round, you and your enemy each get one attack.

### Step 1: Determine First Strike
- Roll two dice for yourself, and two for the enemy.
- Add your **SPD + CRG + LCK** to your roll.
- Add the enemy's **SPD + CRG + LCK** to their roll.
- The combatant with the higher total score strikes first in the first round. After that, you alternate attacks.

---

## 3. Making an Attack (The "To-Hit" Roll)

When you attack, roll two dice.
- You need a **7 or higher** to score a hit.

**Modifiers:**
- **Skill Modifier:** For every **10 full points** of SKL, reduce the number you need to hit by 1.  
  *Example: With 25 SKL, you need a 5 or higher (7 - 2 = 5).*
- **Luck Modifier:** If your LUCK is **72 or higher**, reduce the number you need to hit by 1.
- **These modifiers are cumulative.** With 20 SKL and 80 LCK, you would need a 4 or higher (7 - 2 - 1 = 4).
- The minimum score needed to hit is always **2**.

---

## 4. Calculating Damage

If your attack hits, calculate the damage you deal:
- **Base Damage:** Take the number you rolled on the two dice (the "to-hit" roll). Multiply it by **5**.  
  *Example: You roll an 8 to hit. Your base damage is 8 x 5 = 40.*
- **Strength Bonus:** For every **10 full points** of STR, add **5** to the damage.  
  *Example: With 65 STR, you add 30 damage (6 x 5).*
- **Weapon Bonus:** Add the bonus from your weapon (see table below).

**Total Damage = Base Damage + STR Bonus + Weapon Bonus**

Subtract this total from the enemy's current LIFE POINTS.

---

## 5. Enemy Attacks & Damage

The enemy attacks you in the same way. Their stats will be provided in the book (e.g., on page 248). They use the same rules for hitting and damage. If you are wearing armor or using a shield, subtract its protection value from the damage you take.

---

## 6. Special Rules

### Avoiding Death
If your LIFE POINTS drop to 0 or below, you are not automatically dead. You get one last chance. Roll two dice and multiply by **10**. If the result is **less than or equal to** your LUCK score, you survive! Restore your LP to their full starting value and re-run the fight from the beginning. You may only attempt this once per fight.

### Endurance (Stamina)
Your STAMINA determines how long you can fight without rest. Divide your STA by **20** (round down) to find the number of combat rounds you can fight continuously. If a fight lasts longer than this, you must rest for one round (your enemy gets a free attack).

### Healing Stone
If you have this item, you can use it once between combat rounds to restore a number of LIFE POINTS equal to a roll of one die multiplied by 10 (1d6 x 10). It has 50 LP in total and recharges 48 hours after its last use in combat.

### Doombringer (The Cursed Blade)
This legendary black axe is a double-edged weapon of immense power.
- **Damage Bonus:** +20 (10 points above a normal sword).
- **Blood Price:** Each time you attempt to strike with Doombringer, you **immediately lose 10 LP** before the attack is resolved. This cost is paid even if you miss. If this reduces you to 0 LP, you die before the attack lands.
- **Soul Thirst:** If your attack **hits**, you immediately heal **LP equal to the total damage dealt** (after armor/shield reduction). You cannot exceed your maximum starting LP.

### The Orb
An ancient weapon against the Demonspawn with two uses:
- **Held:** When carried in your left hand during combat, any damage you deal to a **Demonspawn** is **doubled** after all other calculations. You cannot wield another weapon or shield in that hand.
- **Thrown:** You may hurl the Orb as a weapon. Roll two dice: a result of **4 or higher** is a hit. A hit **instantly kills** any Demonspawn. A miss still deals **200 damage** to a Demonspawn. In either case, the Orb is destroyed. The Orb has **no effect** on non-Demonspawn creatures.

---

## 7. Magic & Sorcery

Fire*Wolf abhors sorcery, but survival may force him to use it. Magic is powered by **POWER (POW)**, a resource separate from LIFE POINTS.

### Casting a Spell
Follow these steps in order:

1. **Natural Inclination Check:** Before using *any* magic in a section, roll two dice. If you score **4 or better**, Fire*Wolf overcomes his aversion and may cast spells. If you fail, he refuses to use magic for the entire section, regardless of danger.

2. **Pay Power Cost:** Deduct the spell's POWER cost from your current total. You cannot cast a spell if you lack sufficient POWER. If you wish, you may trade LIFE POINTS for POWER at a **1:1 ratio** to meet the cost.

3. **Fundamental Failure Rate (FFR):** Roll two dice. You must score **6 or better** for the spell to succeed. The POWER cost is spent regardless of success or failure.

4. **Spell Effect:** If successful, apply the spell's effect immediately.

### Additional Restrictions
- You may **never cast the same spell twice** in a single section.
- POWER spent is gone until restored.

### Power Renewal
You can regain POWER in three ways:
- **Exploration:** You automatically gain **1 POWER** when entering a new section. This cannot exceed your original starting total.
- **Sacrifice:** Trade **LIFE POINTS for POWER** at a 1:1 ratio (see step 2 above).
- **Crypt Spell:** Use the CRYPT spell to fully restore or even increase your POWER.

---

## Spell Table

| Spell | Power Cost | Effect |
|-------|------------|--------|
| **ARMOUR** | 25 | Creates magical armor of light for the section. Subtract **10 points** from any damage scored against you in combat. |
| **CRYPT** | 150 | Returns you to the Crypts where you may take tests to restore or **increase** your POWER. |
| **FIREBALL** | 15 | Hurls a ball of flame. If successful, causes **50 LP damage** to an enemy. |
| **INVISIBILITY** | 30 | Renders you invisible for the remainder of the section. You cannot attack, but may avoid combat and proceed as if victorious. |
| **PARALYSIS** | 30 | Paralyzes a single enemy long enough for you to escape to the next section. |
| **POISON NEEDLE** | 25 | Shoots a poisoned needle at a single enemy within combat range. If successful and the enemy is not immune (roll 1d6: 1-3 = immune, 4-6 = affected), the poison is **invariably fatal**. |
| **RESURRECTION** | 50 | **Only usable when killed.** Returns you to the start of the current section. The enemy retains only the LP they had when you died, but you must **reroll all of Fire*Wolf's stats** (including LP and POWER). |
| **RETRACE** | 20 | Returns you to any previously visited section to proceed from there. LIFE POINTS and POWER are **not restored**. |
| **TIMEWARP** | 10 | Resets time to the beginning of the current section. Your LIFE POINTS and the enemy's LIFE POINTS are both restored to their starting values for the section. |
| **XENOPHOBIA** | 15 | Causes a single opponent to fear you. **Subtract 5 points** from any damage they score against you. |

---

## Weapons & Armour Table

### Weapons
| Weapon | Damage Bonus |
|--------|--------------|
| Arrow | +10 |
| Axe | +15 |
| Club | +8 |
| Dagger | +5 |
| Flail | +7 |
| Halberd | +12 |
| Lance | +12 |
| Mace | +14 |
| Spear | +12 |
| Sword | +10 |
| **Doombringer*** | **+20** |

### Armour & Shields
| Armour/Shield | Protection (Damage Reduction) |
|---------------|--------------------------------|
| Leather Armour | -5 |
| Chain Mail | -8 |
| Plate Mail | -12 |
| Shield | -7 |
| Shield (with armour) | -5 (slower usage) |

*See Special Rules for Doombringer's life-draining curse.

---

## Why This is Better

- **Clarity:** Every term is defined, eliminating the original's ambiguous language.
- **Speed:** Damage calculations use multiplication by 5 instead of 10 for easier math.
- **Balance:** High-risk, high-reward items like Doombringer and the Orb now have precise mechanics that create meaningful choices.
- **Completeness:** **All** systems—combat, special items, and sorcery—are integrated into one coherent ruleset with consistent formatting.
- **Authenticity:** Magic rules preserve the original's exact mechanics: Fire*Wolf's moral reluctance, the punishing failure rate, and the strategic POWER economy exactly as written.