package tui

import (
	"strconv"

	"github.com/goverclock/gopolar/internal/core"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func NewTableModel(tunnelList []core.Tunnel) *table.Model {
	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Name", Width: 16},
		{Title: "Source", Width: 16},
		{Title: "Dest", Width: 20},
		{Title: "Status", Width: 8},
	}
	rows := listToRows(tunnelList)
	tb := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(10),
		table.WithFocused(true),
	)
	tb.KeyMap.HalfPageDown.Unbind() // conflicts with 'd' - delete tunnel, so unbind it
	tb.KeyMap.HalfPageUp.Unbind()
	s := table.DefaultStyles()
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("255")).
		Background(lipgloss.Color("8")).
		Bold(false)
	tb.SetStyles(s)
	return &tb
}

func listToRows(tunnelList []core.Tunnel) []table.Row {
	rows := []table.Row{}
	for _, t := range tunnelList {
		status := "STOPPED"
		if t.Enable {
			status = "RUNNING"
		}
		rows = append(rows, table.Row{
			strconv.FormatUint(t.ID, 10),
			t.Name,
			t.Source,
			t.Dest,
			status,
		})
	}
	return rows
}
