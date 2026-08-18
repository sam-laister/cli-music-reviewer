# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project overview

A Go CLI/TUI application for creating and managing music reviews, built with
[Bubble Tea](https://github.com/charmbracelet/bubbletea) (Elm-architecture TUI framework) and
[Lip Gloss](https://github.com/charmbracelet/lipgloss) for styling. Planned functionality (not yet
implemented) includes Spotify integration and linking out to review sites like RateYourMusic and
Album of the Year.

The project was bootstrapped from a copied CLI codebase (a "TikTok clip generator/auto-captioner"),
so some leftover text/branding — notably the splash screen subtitle in
`components/splash_screen.go` — does not describe this app and should be replaced rather than
treated as intentional.

## Commands

- Run the app: `go run .`
- Build: `go build .`
- Tidy/verify deps: `go mod tidy`
- Format: `gofmt -l .` (or `go fmt ./...`)
- Vet: `go vet ./...`
- No test files exist yet. When tests are added, run all with `go test ./...` or a single package
  with `go test ./components/...`.

## Architecture

The app follows the Bubble Tea Elm architecture (Model / Update / View), nested across a small
component tree:

- `main.go` boots a single `tea.Program` running `views.NewHomepage()`.
- `views/homepage.go` (`HomepageModel`) is the root model. It owns a `homepageState` enum
  (`StateSplash`, `StateMenu`, ...) and delegates `Update`/`View` to whichever child component is
  active for that state. `tab` cycles state, `q`/`ctrl+c` quits globally — these are intercepted at
  the homepage level before being handed down to child components.
- `components/` holds child Bubble Tea models, one per screen/widget:
  - `splash_screen.go` — static intro screen (leftover branding, see above).
  - `entry_browser.go` (`EntryBrowserModel`) — a list/menu of `EntryRowModel` children; tracks
    `activeIndex` and exposes `CursorUp`/`CursorDown`, called directly by the parent homepage model
    on arrow-key presses (rather than being routed through `Update`).
  - `entry_row.go` (`EntryRowModel`) — a single selectable row (title + timestamp), rendered
    differently when selected via `SetSelected`.
- `interfaces/component_interface.go` defines a `ComponentInterface` (`Update`/`View`) intended as
  the common shape for components, though current components implement matching methods with
  concrete (not interface) types rather than formally implementing it.
- `styles/styles.go` centralizes all Lip Gloss colors and styles (colors, borders, cursor/selection
  styling) — add new visual styling here rather than inlining `lipgloss.NewStyle()` calls in
  components.
- `config/app_config.go` holds app-wide constants (currently just `AppTitle`).
- `cmd/root.go` is currently an empty package stub — no CLI command wiring (e.g. Cobra) exists yet.

### Adding a new screen/component

Follow the existing pattern: define a `*Model` struct in `components/` with `Update(tea.Msg) (*Model, tea.Cmd)`
and `View() string`, a `NewX(...)` constructor, use `styles/styles.go` for appearance, and wire it
into `HomepageModel`'s state enum and switch statements in `views/homepage.go`.
