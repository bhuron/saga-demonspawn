// Package main is the entry point for the Saga of the Demonspawn companion application.
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/benoit/saga-demonspawn/pkg/ui"
)

func main() {
	// Create the root model
	model := ui.NewModel()

	// Create the Bubble Tea program
	p := tea.NewProgram(model, tea.WithAltScreen())

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running application: %v\n", err)
		os.Exit(1)
	}
}
