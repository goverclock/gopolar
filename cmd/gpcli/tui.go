package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Background(lipgloss.Color("8"))
	blurredStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	noStyle       = lipgloss.NewStyle()
	cursorStyle   = noStyle
)

type sessionState uint // track which model is focused
const (
	tableView sessionState = iota
	createView
	editView
	deleteConfirm
)
const (
	TableHelpMsg  string = "c - CREATE, e - EDIT, d - DELETE, r - RUN/STOP"
	EditHelpMsg   string = "enter - CONFIRM, esc - CANCEL"
	DeleteHelpMsg string = "Delete selected tunnel?(Y/n)"
)

type UIModel struct {
	table   table.Model
	edit    EditModel // multiple textinputs
	helpMsg string

	state sessionState
	end   *CLIEnd
}

func NewUIModel(end *CLIEnd) *UIModel {
	tunnelList, err := end.GetTunnelsList()
	if err != nil {
		log.Println("fail to get tunnel list, is gopolar daemon running?")
		os.Exit(1)
	}
	return &UIModel{
		table:   *NewTableModel(tunnelList),
		edit:    *NewEditModel(),
		helpMsg: TableHelpMsg,
		end:     end,
	}
}

func (m UIModel) Init() tea.Cmd {
	return nil
}

func (m UIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	msgv, ok := msg.(tea.KeyMsg) // only care about key message
	if !ok {
		return m, nil
	}
	s := msgv.String()

	var cmd tea.Cmd
	// main model
	switch s {
	case "esc":
		m.state = tableView
		m.helpMsg = TableHelpMsg
		return m, nil
	}

	// else send msg to sub model
	switch m.state {
	case tableView:
		switch s {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "c":
			m.state = createView
			m.helpMsg = EditHelpMsg
			return m, nil
		case "e":
			m.state = editView
			m.helpMsg = EditHelpMsg
			vals := m.table.SelectedRow()
			m.edit.SetValue(vals[1], vals[2], vals[3])
			return m, nil
		case "d":
			m.state = deleteConfirm
			m.helpMsg = DeleteHelpMsg
			return m, nil
		case "r":
			return m, tea.Batch(
				tea.Println("run/stop ", m.table.SelectedRow()[0]),
			)
		}
		m.table, cmd = m.table.Update(msg)
	case editView:
		m.edit, cmd = m.edit.Update(msg)
		if cmd != nil { // submitted
			name, source, dest := m.edit.GetInput()
			m.end.EditTunnel(123, name, source, dest)

		}
	}

	return m, cmd
}

func (m UIModel) View() string {
	ret := m.table.View()
	ret += "\n" + m.helpMsg
	if m.state == editView {
		ret += "\n" + m.edit.View()
	}
	return ret
}
