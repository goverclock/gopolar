package main

import (
	"fmt"
	"gopolar"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	noStyle       = lipgloss.NewStyle()
	focusedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Background(lipgloss.Color("8"))
	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", noStyle.Render("Submit"))
)

type sessionState uint // track which model is focused
const (
	tableView sessionState = iota
	createView
	editView
	deleteConfirm
)
const (
	TableHelpMsg string = "c - CREATE, e - EDIT, d - DELETE, r - RUN/STOP"
	EditHelpMsg  string = "enter - CONFIRM, esc - CANCEL"
)

type UIModel struct {
	table   table.Model
	edit    EditModel // multiple textinputs
	helpMsg string

	state sessionState
	// TODO: replace updatedList channel with tickMsg: https://github.com/charmbracelet/bubbletea/blob/master/examples/progress-animated/main.go
	end *CLIEnd
}

func NewUIModel(end *CLIEnd) *UIModel {
	tunnelList, err := end.GetTunnelList()
	if err != nil {
		log.Println("fail to get tunnel list, is gopolar daemon running?")
		os.Exit(1)
	}
	ret := &UIModel{
		table:   *NewTableModel(tunnelList),
		edit:    *NewEditModel(),
		helpMsg: TableHelpMsg,
		state:   tableView,
		end:     end,
	}
	return ret
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m *UIModel) updateListCmd() tea.Msg {
	newTunnels, err := m.end.GetTunnelList()
	if err != nil {
		log.Println("fail to get tunnel list, is gopolar daemon running?")
		os.Exit(1)
	}
	return newTunnels
}

func (m UIModel) Init() tea.Cmd {
	return tickCmd()
}

func (m UIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// tick update
	_, ok := msg.(tickMsg)
	if ok {
		return m,
			tea.Batch(tickCmd(), m.updateListCmd)
	}

	// local update
	msgnt, ok := msg.([]gopolar.Tunnel)
	if ok {
		m.table.SetRows(listToRows(msgnt))
		return m, nil
	}

	msgk, ok := msg.(tea.KeyMsg) // only care about key message
	if !ok {
		return m, nil
	}
	s := msgk.String()

	var cmd tea.Cmd = nil
	// main model
	switch s {
	case "esc":
		m.state = tableView
		m.helpMsg = TableHelpMsg
		return m, nil
	case "ctrl+c":
		return m, tea.Quit
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
			m.edit.SetValues("", "localhost:", "")
			return m, nil
		case "e":
			m.state = editView
			m.helpMsg = EditHelpMsg
			vals := m.table.SelectedRow()
			m.edit.SetValues(vals[1], vals[2], vals[3])
			return m, nil
		case "d":
			m.state = deleteConfirm
			sr := m.table.SelectedRow()
			m.helpMsg = fmt.Sprintf("Delete tunnel %v(%v)?(Y/n)", sr[0], sr[1])
			return m, nil
		case "r":
			id, err := strconv.ParseInt(m.table.SelectedRow()[0], 10, 64)
			if err != nil {
				m.helpMsg = "Fail to parse tunnel ID: " + fmt.Sprint(err)
				break
			}
			// request core
			err = m.end.ToggleTunnel(id)
			strOk := "Stopped"
			strFail := "stop"
			if m.table.SelectedRow()[4] == "STOP" {
				strOk = "Started"
				strFail = "start"
			}
			if err != nil {
				m.helpMsg = "Fail to " + strFail + " tunnel: " + fmt.Sprint(err)
			} else {
				m.helpMsg = strOk + " tunnel " + fmt.Sprint(id) + " successfully"
			}
			return m, m.updateListCmd
		}
		// reset to table help message only when table updates
		if m.state == tableView {
			m.helpMsg = TableHelpMsg
		}
		m.table, cmd = m.table.Update(msg)
	case createView:
		m.edit, cmd = m.edit.Update(msg)
		if cmd == nil {
			break
		}
		if cmd() == "submit" { // submitted
			name, source, dest := m.edit.GetInput()
			// request core
			id, err := m.end.CreateTunnel(name, source, dest)
			if err != nil {
				m.helpMsg = fmt.Sprint(err)
			} else {
				m.helpMsg = "Created tunnel " + fmt.Sprint(id)
			}
			m.state = tableView
			return m, m.updateListCmd
		} else {
			m.helpMsg = cmd().(string)
		}
	case editView:
		m.edit, cmd = m.edit.Update(msg)
		if cmd == nil {
			break
		}
		if cmd() == "submit" { // submitted
			name, source, dest := m.edit.GetInput()
			id, err := strconv.ParseUint(m.table.SelectedRow()[0], 10, 64)
			if err != nil {
				m.helpMsg = "Fail to parse tunnel ID: " + fmt.Sprint(err)
				break
			}
			// request core
			err = m.end.EditTunnel(id, name, source, dest)
			if err != nil {
				m.helpMsg = "Fail to edit tunnel: " + fmt.Sprint(err)
			} else {
				m.helpMsg = "Edited tunnel " + fmt.Sprint(id) + " successfully"
			}
			m.state = tableView
			return m, m.updateListCmd
		} else { // validate fail, got error message
			m.helpMsg = cmd().(string)
		}
	case deleteConfirm:
		switch s {
		case "y", "Y", "enter": // confirm
			id, err := strconv.ParseInt(m.table.SelectedRow()[0], 10, 64)
			if err != nil {
				m.helpMsg = "Fail to parse tunnel ID: " + fmt.Sprint(err)
				break
			}
			err = m.end.DeleteTunnel(id)
			if err != nil {
				m.helpMsg = "Fail to delete tunnel: " + fmt.Sprint(err)
			} else {
				m.helpMsg = "Deleted tunnel " + fmt.Sprint(id) + " successfully"
			}
			m.state = tableView
			return m, m.updateListCmd
		case "n", "N", "esc": // cancel
			m.helpMsg = TableHelpMsg
			m.state = tableView
			return m, nil
		}
	}
	return m, cmd
}

func (m UIModel) View() string {
	ret := m.table.View()
	ret += "\n" + m.helpMsg
	if m.state == createView || m.state == editView {
		ret += "\n" + m.edit.View()
	}
	return ret
}
