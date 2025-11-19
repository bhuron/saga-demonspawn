# Technology Stack

## Language & Runtime

- **Go 1.24.0** (toolchain 1.24.10)
- Module: `github.com/benoit/saga-demonspawn`

## Core Dependencies

- **Bubble Tea** (`github.com/charmbracelet/bubbletea`): TUI framework implementing The Elm Architecture (Model-View-Update)
- **Lipgloss** (`github.com/charmbracelet/lipgloss`): Terminal styling and theming
- **Charm libraries**: Color profiles, ANSI handling, terminal utilities

## Build System

### Building

```bash
# Standard build
go build -o saga ./cmd/saga

# Optimized build (smaller binary)
go build -ldflags="-s -w" -o saga ./cmd/saga

# Using build script
./build.sh
```

### Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Test specific package
go test ./internal/character
go test ./internal/combat
go test ./internal/dice
go test ./internal/magic
```

### Running

```bash
./saga
```

## Architecture Pattern

**The Elm Architecture (TEA)** via Bubble Tea:
- **Model**: Immutable state structs
- **Update**: Message-driven state transitions
- **View**: Pure rendering functions

All UI components follow this pattern with a root `Model` that delegates to screen-specific models.
