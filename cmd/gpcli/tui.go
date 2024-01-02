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
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle.Copy()
	// noStyle      = lipgloss.NewStyle()
	// helpStyle           = blurredStyle.Copy()
	// cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	focusedButton       = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton       = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type sessionState uint // track which model is focused
const (
	tableView sessionState = iota
	editView
	deleteConfirm
	pending // waiting for server response, ignore all key messages except for quit
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
		helpMsg: "e - EDIT, d - DELETE, r - RUN/STOP",
		end:     end,
	}
}

func (m UIModel) Init() tea.Cmd {
	return nil
}

func (m UIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	// main model
	switch msgv := msg.(type) {
	case tea.KeyMsg: // ignore all other message types
		switch msgv.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "e":
			m.state = editView
			m.helpMsg = "enter - CONFIRM, esc - CANCEL"
			return m, nil
		case "d":
			m.state = deleteConfirm
			m.helpMsg = "Delete tunnel " + m.table.SelectedRow()[0] + " ?(Y/n)"
			return m, nil
		case "r":
			return m, tea.Batch(
				tea.Println("run/stop ", m.table.SelectedRow()[0]),
			)
		case "esc":
			m.state = tableView
			m.helpMsg = "e - EDIT, d - DELETE, r - RUN/STOP"
			return m, nil
		}
	}

	// else send msg to sub model
	switch m.state {
	case tableView:
		m.table, cmd = m.table.Update(msg)
	case editView:
		m.edit, cmd = m.edit.Update(msg)
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
