package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
)

type model struct {
	table table.Model
}

func newModel() model {
	// Define columns
	cols := []table.Column{
		table.NewColumn("id", "ID", 5),
		table.NewColumn("name", "Name", 12),
		table.NewColumn("score", "Score", 6),
	}

	// Define rows
	rows := []table.Row{
		table.NewRow(table.RowData{"id": 1, "name": "Alice", "score": 42}),
		table.NewRow(table.RowData{"id": 2, "name": "Bob", "score": 88}),
		table.NewRow(table.RowData{"id": 3, "name": "Charlie", "score": 73}),
	}

	// Create table model
	tbl := table.New(cols).
		WithRows(rows).
		Focused(true).
		SelectableRows(true)

	return model{table: tbl}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			fmt.Println("\nSelected cells:")
			for _, c := range m.table.SelectedCells() {
				fmt.Printf("RowID=%d Col=%s\n", c.RowID, c.ColKey)
			}
			return m, tea.Quit
		}
	}

	return m, cmd
}

func (m model) View() string {
	return m.table.View() + "\n\nMove with arrows (↑↓←→), select cell with [space], quit with [q]."
}

func main() {
	if err := tea.NewProgram(newModel()).Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
