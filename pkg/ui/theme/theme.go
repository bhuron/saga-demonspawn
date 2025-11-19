// Package theme provides centralized styling and theming for the UI using lipgloss.
package theme

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ColorScheme represents the active color theme.
type ColorScheme string

const (
	// ColorSchemeDark is the default dark theme
	ColorSchemeDark ColorScheme = "dark"
	// ColorSchemeLight is the light theme
	ColorSchemeLight ColorScheme = "light"
)

// Theme holds all the styling configuration.
type Theme struct {
	// Color definitions
	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Success   lipgloss.Color
	Warning   lipgloss.Color
	Danger    lipgloss.Color
	Neutral   lipgloss.Color
	Accent    lipgloss.Color
	Muted     lipgloss.Color

	// Background colors
	BgPrimary   lipgloss.Color
	BgSecondary lipgloss.Color
	BgNeutral   lipgloss.Color

	// Typography styles
	Title      lipgloss.Style
	Heading    lipgloss.Style
	Label      lipgloss.Style
	Value      lipgloss.Style
	Body       lipgloss.Style
	Emphasis   lipgloss.Style
	MutedText  lipgloss.Style
	Error      lipgloss.Style
	SuccessMsg lipgloss.Style
	WarningMsg lipgloss.Style

	// Component styles
	Box         lipgloss.Style
	Header      lipgloss.Style
	MenuItem    lipgloss.Style
	MenuItemSel lipgloss.Style
	Button      lipgloss.Style
	ButtonFocus lipgloss.Style
	Panel       lipgloss.Style
	Border      lipgloss.Style

	// UI Elements
	Cursor         string
	CursorInactive string
	UseUnicode     bool
}

var (
	// Current active theme
	current *Theme
)

// Init initializes the theme system with the specified color scheme.
func Init(scheme ColorScheme, useUnicode bool) {
	switch scheme {
	case ColorSchemeLight:
		current = newLightTheme(useUnicode)
	default:
		current = newDarkTheme(useUnicode)
	}
}

// Current returns the currently active theme.
func Current() *Theme {
	if current == nil {
		Init(ColorSchemeDark, true)
	}
	return current
}

// newDarkTheme creates the dark color theme.
func newDarkTheme(useUnicode bool) *Theme {
	t := &Theme{
		// Color palette
		Primary:   lipgloss.Color("86"),  // Bright Cyan
		Secondary: lipgloss.Color("141"), // Purple/Magenta
		Success:   lipgloss.Color("42"),  // Green
		Warning:   lipgloss.Color("220"), // Yellow/Gold
		Danger:    lipgloss.Color("196"), // Red
		Neutral:   lipgloss.Color("252"), // Light Gray
		Accent:    lipgloss.Color("208"), // Orange
		Muted:     lipgloss.Color("240"), // Dark Gray

		BgPrimary:   lipgloss.Color("235"),
		BgSecondary: lipgloss.Color("237"),
		BgNeutral:   lipgloss.Color("0"),

		UseUnicode:     useUnicode,
		Cursor:         ">",
		CursorInactive: " ",
	}

	// Build typography styles
	t.Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Primary).
		Padding(0, 1)

	t.Heading = lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Secondary)

	t.Label = lipgloss.NewStyle().
		Foreground(t.Neutral)

	t.Value = lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Neutral)

	t.Body = lipgloss.NewStyle().
		Foreground(t.Neutral)

	t.Emphasis = lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Accent)

	t.MutedText = lipgloss.NewStyle().
		Foreground(t.Muted)

	t.Error = lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Danger)

	t.SuccessMsg = lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Success)

	t.WarningMsg = lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Warning)

	// Build component styles
	borderStyle := lipgloss.RoundedBorder()
	if !useUnicode {
		borderStyle = lipgloss.NormalBorder()
	}

	t.Box = lipgloss.NewStyle().
		Border(borderStyle).
		BorderForeground(t.Primary).
		Padding(0, 1)

	t.Header = lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Primary).
		Align(lipgloss.Center).
		Padding(0, 2)

	t.MenuItem = lipgloss.NewStyle().
		Foreground(t.Neutral).
		Padding(0, 1)

	t.MenuItemSel = lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Primary).
		Padding(0, 1)

	t.Button = lipgloss.NewStyle().
		Foreground(t.Neutral).
		Border(borderStyle).
		BorderForeground(t.Muted).
		Padding(0, 1)

	t.ButtonFocus = lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Primary).
		Border(borderStyle).
		BorderForeground(t.Primary).
		Padding(0, 1)

	t.Panel = lipgloss.NewStyle().
		Border(borderStyle).
		BorderForeground(t.Muted).
		Padding(0, 1)

	t.Border = lipgloss.NewStyle().
		Foreground(t.Primary)

	return t
}

// newLightTheme creates the light color theme.
func newLightTheme(useUnicode bool) *Theme {
	t := &Theme{
		// Color palette (adjusted for light backgrounds)
		Primary:   lipgloss.Color("33"),  // Blue
		Secondary: lipgloss.Color("129"), // Purple
		Success:   lipgloss.Color("28"),  // Dark Green
		Warning:   lipgloss.Color("136"), // Dark Yellow
		Danger:    lipgloss.Color("160"), // Dark Red
		Neutral:   lipgloss.Color("235"), // Dark Gray (text)
		Accent:    lipgloss.Color("166"), // Dark Orange
		Muted:     lipgloss.Color("244"), // Medium Gray

		BgPrimary:   lipgloss.Color("255"),
		BgSecondary: lipgloss.Color("252"),
		BgNeutral:   lipgloss.Color("231"),

		UseUnicode:     useUnicode,
		Cursor:         ">",
		CursorInactive: " ",
	}

	// Build typography styles (similar structure to dark theme)
	borderStyle := lipgloss.RoundedBorder()
	if !useUnicode {
		borderStyle = lipgloss.NormalBorder()
	}

	t.Title = lipgloss.NewStyle().Bold(true).Foreground(t.Primary).Padding(0, 1)
	t.Heading = lipgloss.NewStyle().Bold(true).Foreground(t.Secondary)
	t.Label = lipgloss.NewStyle().Foreground(t.Neutral)
	t.Value = lipgloss.NewStyle().Bold(true).Foreground(t.Neutral)
	t.Body = lipgloss.NewStyle().Foreground(t.Neutral)
	t.Emphasis = lipgloss.NewStyle().Bold(true).Foreground(t.Accent)
	t.MutedText = lipgloss.NewStyle().Foreground(t.Muted)
	t.Error = lipgloss.NewStyle().Bold(true).Foreground(t.Danger)
	t.SuccessMsg = lipgloss.NewStyle().Bold(true).Foreground(t.Success)
	t.WarningMsg = lipgloss.NewStyle().Bold(true).Foreground(t.Warning)

	t.Box = lipgloss.NewStyle().Border(borderStyle).BorderForeground(t.Primary).Padding(0, 1)
	t.Header = lipgloss.NewStyle().Bold(true).Foreground(t.Primary).Align(lipgloss.Center).Padding(0, 2)
	t.MenuItem = lipgloss.NewStyle().Foreground(t.Neutral).Padding(0, 1)
	t.MenuItemSel = lipgloss.NewStyle().Bold(true).Foreground(t.Primary).Padding(0, 1)
	t.Button = lipgloss.NewStyle().Foreground(t.Neutral).Border(borderStyle).BorderForeground(t.Muted).Padding(0, 1)
	t.ButtonFocus = lipgloss.NewStyle().Bold(true).Foreground(t.Primary).Border(borderStyle).BorderForeground(t.Primary).Padding(0, 1)
	t.Panel = lipgloss.NewStyle().Border(borderStyle).BorderForeground(t.Muted).Padding(0, 1)
	t.Border = lipgloss.NewStyle().Foreground(t.Primary)

	return t
}

// RenderTitle renders a styled screen title with decorative border.
func RenderTitle(title string) string {
	t := Current()
	width := len(title) + 4 // Account for padding and borders

	var top, bottom string
	if t.UseUnicode {
		top = "╔" + strings.Repeat("═", width) + "╗"
		bottom = "╚" + strings.Repeat("═", width) + "╝"
	} else {
		top = "+" + strings.Repeat("-", width) + "+"
		bottom = "+" + strings.Repeat("-", width) + "+"
	}

	titleLine := t.Title.Render(title)
	
	// Build the title box
	var b strings.Builder
	b.WriteString(t.Border.Render(top) + "\n")
	if t.UseUnicode {
		b.WriteString(t.Border.Render("║") + " " + titleLine + " " + t.Border.Render("║") + "\n")
	} else {
		b.WriteString(t.Border.Render("|") + " " + titleLine + " " + t.Border.Render("|") + "\n")
	}
	b.WriteString(t.Border.Render(bottom))

	return b.String()
}

// RenderMenuItem renders a menu item with appropriate styling based on selection.
func RenderMenuItem(text string, selected bool) string {
	t := Current()
	cursor := t.CursorInactive
	style := t.MenuItem
	
	if selected {
		cursor = t.Cursor
		style = t.MenuItemSel
	}

	return cursor + " " + style.Render(text)
}

// RenderStatValue renders a stat value with color coding based on value.
func RenderStatValue(value, min, max int) string {
	t := Current()
	
	// Calculate percentage of range
	rangeSize := max - min
	if rangeSize <= 0 {
		return t.Value.Render(fmt.Sprintf("%d", value))
	}
	
	percentage := float64(value-min) / float64(rangeSize)
	
	var style lipgloss.Style
	switch {
	case percentage >= 0.75:
		style = lipgloss.NewStyle().Bold(true).Foreground(t.Success)
	case percentage >= 0.50:
		style = lipgloss.NewStyle().Bold(true).Foreground(t.Neutral)
	case percentage >= 0.25:
		style = lipgloss.NewStyle().Bold(true).Foreground(t.Warning)
	default:
		style = lipgloss.NewStyle().Bold(true).Foreground(t.Danger)
	}
	
	return style.Render(fmt.Sprintf("%d", value))
}

// RenderHealthBar renders a health/LP bar with color gradient.
func RenderHealthBar(current, max, width int) string {
	t := Current()
	
	if max <= 0 {
		return ""
	}
	
	percentage := float64(current) / float64(max)
	if percentage < 0 {
		percentage = 0
	}
	if percentage > 1 {
		percentage = 1
	}
	
	filled := int(float64(width) * percentage)
	empty := width - filled
	
	var barChar, emptyChar string
	if t.UseUnicode {
		barChar = "█"
		emptyChar = "░"
	} else {
		barChar = "#"
		emptyChar = "-"
	}
	
	var color lipgloss.Color
	switch {
	case percentage >= 0.75:
		color = t.Success
	case percentage >= 0.50:
		color = t.Warning
	case percentage >= 0.25:
		color = t.Accent
	default:
		color = t.Danger
	}
	
	bar := strings.Repeat(barChar, filled) + strings.Repeat(emptyChar, empty)
	styled := lipgloss.NewStyle().Foreground(color).Render(bar)
	
	return fmt.Sprintf("[%s] %d/%d", styled, current, max)
}

// RenderPOWMeter renders a POWER gauge similar to health bar.
func RenderPOWMeter(current, max, width int) string {
	t := Current()
	
	if max <= 0 {
		return ""
	}
	
	percentage := float64(current) / float64(max)
	if percentage < 0 {
		percentage = 0
	}
	if percentage > 1 {
		percentage = 1
	}
	
	filled := int(float64(width) * percentage)
	empty := width - filled
	
	var barChar, emptyChar string
	if t.UseUnicode {
		barChar = "▓"
		emptyChar = "░"
	} else {
		barChar = "*"
		emptyChar = "-"
	}
	
	bar := strings.Repeat(barChar, filled) + strings.Repeat(emptyChar, empty)
	styled := lipgloss.NewStyle().Foreground(t.Secondary).Render(bar)
	
	return fmt.Sprintf("[%s] %d/%d POW", styled, current, max)
}

// RenderError formats an error message with styling.
func RenderError(title, description, suggestion string) string {
	t := Current()
	
	var b strings.Builder
	
	icon := "✗"
	if !t.UseUnicode {
		icon = "X"
	}
	
	b.WriteString(t.Error.Render(icon + " " + title) + "\n")
	if description != "" {
		b.WriteString(t.Body.Render(description) + "\n")
	}
	if suggestion != "" {
		b.WriteString(t.MutedText.Render("→ " + suggestion))
	}
	
	return b.String()
}

// RenderWarning formats a warning message with styling.
func RenderWarning(title, description string) string {
	t := Current()
	
	var b strings.Builder
	
	icon := "⚠"
	if !t.UseUnicode {
		icon = "!"
	}
	
	b.WriteString(t.WarningMsg.Render(icon + " " + title) + "\n")
	if description != "" {
		b.WriteString(t.Body.Render(description))
	}
	
	return b.String()
}

// RenderSuccess formats a success message with styling.
func RenderSuccess(message string) string {
	t := Current()
	
	icon := "✓"
	if !t.UseUnicode {
		icon = "+"
	}
	
	return t.SuccessMsg.Render(icon + " " + message)
}

// RenderBox renders content in a styled box with optional title.
func RenderBox(content string, title string) string {
	t := Current()
	
	if title != "" {
		// Render box with title manually since BorderTitle may not be available
		var b strings.Builder
		b.WriteString(t.Heading.Render(title) + "\n")
		b.WriteString(t.Box.Render(content))
		return b.String()
	}
	
	return t.Box.Render(content)
}

// RenderKeyHelp renders keyboard shortcuts with consistent styling.
func RenderKeyHelp(keys ...string) string {
	t := Current()
	
	var parts []string
	for _, key := range keys {
		parts = append(parts, t.MutedText.Render(key))
	}
	
	return strings.Join(parts, t.MutedText.Render("  |  "))
}

// RenderLabel renders a label-value pair with proper styling.
func RenderLabel(label, value string) string {
	t := Current()
	return t.Label.Render(label+": ") + t.Value.Render(value)
}

// RenderSeparator renders a visual separator line.
func RenderSeparator(width int) string {
	t := Current()
	
	var char string
	if t.UseUnicode {
		char = "─"
	} else {
		char = "-"
	}
	
	return t.MutedText.Render(strings.Repeat(char, width))
}
