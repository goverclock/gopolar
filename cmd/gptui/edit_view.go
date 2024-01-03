package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type EditModel struct {
	focusIndex int
	inputs     []textinput.Model
}

func NewEditModel() *EditModel {
	m := EditModel{
		inputs: make([]textinput.Model, 3),
	}
	for i := range m.inputs {
		t := textinput.New()
		t.CharLimit = len("111.111.111.111:65535")
		switch i {
		case 0:
			t.Focus()
			t.Placeholder = "new name"
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "new source"
		case 2:
			t.Placeholder = "new dest"
		}
		m.inputs[i] = t
	}

	return &m
}

// reset edit view to default, then set values
func (m *EditModel) SetValues(name, source, dest string) {
	m.Reset()
	m.inputs[0].SetValue(name)
	m.inputs[1].SetValue(source)
	m.inputs[2].SetValue(dest)
}

func (m EditModel) GetInput() (name string, source string, dest string) {
	return m.inputs[0].Value(), m.inputs[1].Value(), m.inputs[2].Value()
}

func (m *EditModel) Reset() {
	m.focusIndex = 0
	m.updateFocusStyle()
	for i := range m.inputs {
		m.inputs[i].SetValue("")
		m.inputs[i].SetCursor(0)
	}
}

func (m EditModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m EditModel) Update(msg tea.Msg) (EditModel, tea.Cmd) {
	msgv, ok := msg.(tea.KeyMsg) // only care about key message
	if !ok {
		return m, nil
	}
	s := msgv.String()

	switch s {
	case "tab", "shift+tab", "enter", "up", "down":
		// submit
		if s == "enter" && m.focusIndex == len(m.inputs) {
			return m, func() tea.Msg {
				return tea.KeyEnter
			}
		}
		// cycle indexes
		if s == "up" || s == "shift+tab" {
			if m.focusIndex > 0 {
				m.focusIndex--
			}
		} else if m.focusIndex < len(m.inputs) {
			m.focusIndex++
		}
		m.updateFocusStyle()
		return m, nil
	}

	_ = m.updateInputs(msg)
	return m, nil
}

func (m EditModel) View() string {
	var b strings.Builder
	prompts := []string{"Name  ", "Source", "Dest  "}
	for i := range m.inputs {
		b.WriteString(prompts[i])
		// padding for background color
		limit := m.inputs[i].CharLimit
		old := m.inputs[i].Value()
		padded := strings.Clone(old)
		for len(padded) < limit {
			padded += " "
		}
		m.inputs[i].SetValue(padded)
		if len(old) == 0 {
			m.inputs[i].SetCursor(0)
		}
		b.WriteString(m.inputs[i].View())
		m.inputs[i].SetValue(old)
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}
	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\t%s\n\n", *button)

	return b.String()
}

func (m *EditModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m *EditModel) updateFocusStyle() {
	for i := 0; i < len(m.inputs); i++ {
		if i == m.focusIndex {
			m.inputs[i].Focus()
			m.inputs[i].PromptStyle = focusedStyle
			m.inputs[i].TextStyle = focusedStyle
		} else { // remove focus
			m.inputs[i].Blur()
			m.inputs[i].PromptStyle = noStyle
			m.inputs[i].TextStyle = noStyle
		}
	}

}
