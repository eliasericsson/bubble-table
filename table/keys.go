package table

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines the keybindings for the table when it's focused.
type KeyMap struct {
	RowDown key.Binding
	RowUp   key.Binding

	RowSelectToggle  key.Binding
	CellSelectToggle key.Binding

	PageDown  key.Binding
	PageUp    key.Binding
	PageFirst key.Binding
	PageLast  key.Binding

	// Filter allows the user to start typing and filter the rows.
	Filter key.Binding

	// FilterBlur is the key that stops the user's input from typing into the filter.
	FilterBlur key.Binding

	// FilterClear will clear the filter while it's blurred.
	FilterClear key.Binding

	// ScrollRight will move one column to the right when overflow occurs.
	ScrollRight key.Binding

	// ScrollLeft will move one column to the left when overflow occurs.
	ScrollLeft key.Binding
}

// DefaultKeyMap returns a set of sensible defaults for controlling a focused table with help text.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		// Vertical movement
		RowDown: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "move down"),
		),
		RowUp: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "move up"),
		),

		// Row selection
		RowSelectToggle: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select row"),
		),

		// Cell selection
		CellSelectToggle: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "select cell"),
		),

		// Horizontal movement
		ScrollLeft: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", "move left"),
		),
		ScrollRight: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("→", "move right"),
		),

		// Paging (optional, keep same as bubble-table)
		PageDown: key.NewBinding(
			key.WithKeys("pgdown"),
			key.WithHelp("PgDn", "next page"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup"),
			key.WithHelp("PgUp", "prev page"),
		),
		PageFirst: key.NewBinding(
			key.WithKeys("home"),
			key.WithHelp("Home", "first page"),
		),
		PageLast: key.NewBinding(
			key.WithKeys("end"),
			key.WithHelp("End", "last page"),
		),

		// Filtering
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
		),
		FilterClear: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "clear filter"),
		),
		FilterBlur: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "apply filter"),
		),
	}
}

// FullHelp returns a multi row view of all the helpkeys that are defined. Needed to fullfil the 'help.Model' interface.
// Also appends all user defined extra keys to the help.
func (m Model) FullHelp() [][]key.Binding {
	keyBinds := [][]key.Binding{
		{m.keyMap.RowDown, m.keyMap.RowUp, m.keyMap.RowSelectToggle},
		{m.keyMap.PageDown, m.keyMap.PageUp, m.keyMap.PageFirst, m.keyMap.PageLast},
		{m.keyMap.Filter, m.keyMap.FilterBlur, m.keyMap.FilterClear, m.keyMap.ScrollRight, m.keyMap.ScrollLeft},
	}
	if m.additionalFullHelpKeys != nil {
		keyBinds = append(keyBinds, m.additionalFullHelpKeys())
	}

	return keyBinds
}

// ShortHelp just returns a single row of help views. Needed to fullfil the 'help.Model' interface.
// Also appends all user defined extra keys to the help.
func (m Model) ShortHelp() []key.Binding {
	keyBinds := []key.Binding{
		m.keyMap.RowDown,
		m.keyMap.RowUp,
		m.keyMap.RowSelectToggle,
		m.keyMap.PageDown,
		m.keyMap.PageUp,
		m.keyMap.Filter,
		m.keyMap.FilterBlur,
		m.keyMap.FilterClear,
	}
	if m.additionalShortHelpKeys != nil {
		keyBinds = append(keyBinds, m.additionalShortHelpKeys()...)
	}

	return keyBinds
}
