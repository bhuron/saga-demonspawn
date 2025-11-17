# Character Edit Screen Bug Fix

## Problem Description
When viewing the character edit screen, all fields displayed the same value as the currently selected (cursor) field. This was a pure display bug - the actual character data was correct.

## Root Cause
In `pkg/ui/view.go`, the `viewCharacterEdit()` function was incorrectly getting the field value:

```go
// BUGGY CODE (BEFORE):
for i, field := range fields {
    value := m.CharEdit.GetCurrentValue()  // Gets value for cursor position only!
    
    if i == cursor && m.CharEdit.IsInputMode() {
        // Show input buffer
    } else {
        if i == cursor {  // Switch only ran when i == cursor
            switch EditField(i) {
                case EditFieldStrength:
                    value = m.Character.Strength
                // ... etc
            }
        }
        // All fields printed the same 'value' (cursor's value)
        b.WriteString(fmt.Sprintf("%s%-15s: %d\n", prefix, field, value))
    }
}
```

The problem: `value` was initialized with `m.CharEdit.GetCurrentValue()` which only returns the value for the current cursor position. The switch statement inside `if i == cursor` only updated `value` for the selected field, so all other fields displayed that same value.

## Fix Applied
Moved the value retrieval to execute for EVERY field in the loop, not just the cursor:

```go
// FIXED CODE (AFTER):
for i, field := range fields {
    if i == cursor && m.CharEdit.IsInputMode() {
        // Show input buffer when editing
        b.WriteString(fmt.Sprintf("%s%-15s: [%s_]\n", prefix, field, m.CharEdit.GetInputBuffer()))
    } else {
        // Get the actual value for EACH field (not just cursor)
        var value int
        switch EditField(i) {  // This now runs for every i
            case EditFieldStrength:
                value = m.Character.Strength
            case EditFieldSpeed:
                value = m.Character.Speed
            case EditFieldStamina:
                value = m.Character.Stamina
            case EditFieldCourage:
                value = m.Character.Courage
            case EditFieldLuck:
                value = m.Character.Luck
            case EditFieldCharm:
                value = m.Character.Charm
            case EditFieldAttraction:
                value = m.Character.Attraction
            case EditFieldCurrentLP:
                value = m.Character.CurrentLP
            case EditFieldMaxLP:
                value = m.Character.MaximumLP
            case EditFieldSkill:
                value = m.Character.Skill
            case EditFieldCurrentPOW:
                value = m.Character.CurrentPOW
            case EditFieldMaxPOW:
                value = m.Character.MaximumPOW
        }
        b.WriteString(fmt.Sprintf("%s%-15s: %d\n", prefix, field, value))
    }
}
```

## Result
Each field now correctly displays its own value:

**Before (buggy):**
```
  Strength       : 48
  Speed          : 48
  Stamina        : 48
> Courage        : 48    ← selected
  Luck           : 48    ← all showing cursor value!
  Charm          : 48
  ...
```

**After (fixed):**
```
  Strength       : 64
  Speed          : 56
  Stamina        : 72
> Courage        : 48    ← selected
  Luck           : 80    ← each shows correct value!
  Charm          : 40
  ...
```

## To Apply the Fix

The code has been updated in `pkg/ui/view.go`. To rebuild:

```bash
# Clean and rebuild
rm -f saga
go build -o saga ./cmd/saga

# Or use the build script
chmod +x build.sh
./build.sh
```

## File Modified
- `pkg/ui/view.go` - lines 314-354 (viewCharacterEdit function)

## Testing
The fix has been verified:
1. Code compiles without errors
2. Logic is correct - each field gets its own value from the character
3. EditField enum values (0-11) match loop indices (0-11) perfectly
