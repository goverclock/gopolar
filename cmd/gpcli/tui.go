package main

import (
	"context"
	"net"
	"net/http"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type UIModel struct {
	tunnelsTable table.Model
	end          *CLIEnd
}

func NewUIModel(end *CLIEnd) *UIModel {
	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Name", Width: 10},
		{Title: "Source", Width: 8},
		{Title: "Dest", Width: 8},
	}
	rows := []table.Row{
		{"1", "http", "8080", "9999"},
		{"2", "random", "2345", "2222"},
		{"3", "db", "3306", "7777"},
	}
	tb := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)
	s := table.DefaultStyles()
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("255")).
		Background(lipgloss.Color("8")).
		Bold(false)
	tb.SetStyles(s)

	return &UIModel{
		tunnelsTable: tb,
		end:          end,
	}
}

func (m UIModel) Init() tea.Cmd {
	return nil
}

func (m UIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msgv := msg.(type) {
	case tea.KeyMsg: // ignore all other message types
		switch msgv.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			httpc := http.Client{
				Transport: &http.Transport{
					DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
						return net.Dial("unix", "/tmp/gopolar.sock")
					},
				},
			}

			response, err := httpc.Get("http://unix" + "/hello")
			check(err)
			buf := make([]byte, 10)
			response.Body.Read(buf)

			return m, tea.Batch(
				tea.Println("selected ", m.tunnelsTable.Cursor()),
				tea.Println(string(buf)),
			)
		}
	}
	m.tunnelsTable, cmd = m.tunnelsTable.Update(msg)
	return m, cmd
}

func (m UIModel) View() string {
	return m.tunnelsTable.View()
}
