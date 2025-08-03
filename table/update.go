package table

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) moveHighlightUp() {
	m.rowCursorIndex--

	if m.rowCursorIndex < 0 {
		m.rowCursorIndex = len(m.GetVisibleRows()) - 1
	}

	m.currentPage = m.expectedPageForRowIndex(m.rowCursorIndex)
}

func (m *Model) moveHighlightDown() {
	m.rowCursorIndex++

	if m.rowCursorIndex >= len(m.GetVisibleRows()) {
		m.rowCursorIndex = 0
	}

	m.currentPage = m.expectedPageForRowIndex(m.rowCursorIndex)
}

func (m *Model) moveHighlightLeft() {
	m.cursorColIndex--
	if m.cursorColIndex < 0 {
		m.cursorColIndex = len(m.columns) - 1
	}
}

func (m *Model) moveHighlightRight() {
	m.cursorColIndex++
	if m.cursorColIndex >= len(m.columns) {
		m.cursorColIndex = 0
	}
}

func (m *Model) toggleCellSelect() {
	if len(m.GetVisibleRows()) == 0 || len(m.columns) == 0 {
		return
	}

	row := m.GetVisibleRows()[m.rowCursorIndex]
	col := m.columns[m.cursorColIndex]

	key := CellKey{RowID: row.id, ColKey: col.key}

	if m.selectedCells[key] {
		delete(m.selectedCells, key)
	} else {
		m.selectedCells[key] = true
	}

	m.appendUserEvent(UserEventCellSelectToggled{
		RowIndex:    m.rowCursorIndex,
		ColumnIndex: m.cursorColIndex,
		IsSelected:  m.selectedCells[key],
	})
}

func (m *Model) toggleSelect() {
	if !m.selectableRows || len(m.GetVisibleRows()) == 0 {
		return
	}

	rows := m.GetVisibleRows()

	rowID := rows[m.rowCursorIndex].id

	currentSelectedState := false

	for i := range m.rows {
		if m.rows[i].id == rowID {
			currentSelectedState = m.rows[i].selected
			m.rows[i].selected = !m.rows[i].selected
		}
	}

	m.visibleRowCacheUpdated = false

	m.appendUserEvent(UserEventRowSelectToggled{
		RowIndex:   m.rowCursorIndex,
		IsSelected: !currentSelectedState,
	})
}

func (m Model) updateFilterTextInput(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keyMap.FilterBlur) {
			m.filterTextInput.Blur()
		}
	}
	m.filterTextInput, cmd = m.filterTextInput.Update(msg)
	m.pageFirst()
	m.visibleRowCacheUpdated = false

	return m, cmd
}

// This is a series of Matches tests with minimal logic
//
//nolint:cyclop
func (m *Model) handleKeypress(msg tea.KeyMsg) {
	previousRowIndex := m.rowCursorIndex

	if key.Matches(msg, m.keyMap.RowDown) {
		m.moveHighlightDown()
	}

	if key.Matches(msg, m.keyMap.RowUp) {
		m.moveHighlightUp()
	}

	if key.Matches(msg, m.keyMap.RowSelectToggle) {
		m.toggleSelect()
	}
	// Horizontal movement
	if key.Matches(msg, m.keyMap.ScrollLeft) {
		m.moveHighlightLeft()
	}

	if key.Matches(msg, m.keyMap.ScrollRight) {
		m.moveHighlightRight()
	}

	// Cell select toggle
	if key.Matches(msg, m.keyMap.CellSelectToggle) {
		m.toggleCellSelect()
	}

	if key.Matches(msg, m.keyMap.PageDown) {
		m.pageDown()
	}

	if key.Matches(msg, m.keyMap.PageUp) {
		m.pageUp()
	}

	if key.Matches(msg, m.keyMap.PageFirst) {
		m.pageFirst()
	}

	if key.Matches(msg, m.keyMap.PageLast) {
		m.pageLast()
	}

	if key.Matches(msg, m.keyMap.Filter) {
		m.filterTextInput.Focus()
		m.appendUserEvent(UserEventFilterInputFocused{})
	}

	if key.Matches(msg, m.keyMap.FilterClear) {
		m.visibleRowCacheUpdated = false
		m.filterTextInput.Reset()
	}

	if key.Matches(msg, m.keyMap.ScrollRight) {
		m.scrollRight()
	}

	if key.Matches(msg, m.keyMap.ScrollLeft) {
		m.scrollLeft()
	}

	if m.rowCursorIndex != previousRowIndex {
		m.appendUserEvent(UserEventHighlightedIndexChanged{
			PreviousRowIndex: previousRowIndex,
			SelectedRowIndex: m.rowCursorIndex,
		})
	}
}

// Update responds to input from the user or other messages from Bubble Tea.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	m.clearUserEvents()

	if !m.focused {
		return m, nil
	}

	if m.filterTextInput.Focused() {
		var cmd tea.Cmd
		m, cmd = m.updateFilterTextInput(msg)

		if !m.filterTextInput.Focused() {
			m.appendUserEvent(UserEventFilterInputUnfocused{})
		}

		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.handleKeypress(msg)
	}

	return m, nil
}
