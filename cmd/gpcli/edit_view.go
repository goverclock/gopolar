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
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32
		switch i {
		case 0:
			t.Focus()
			t.Placeholder = "Nickname"
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Email"
			t.CharLimit = 64
		case 2:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}
		m.inputs[i] = t
	}

	return &m
}

func (m EditModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m EditModel) Update(msg tea.Msg) (EditModel, tea.Cmd) {
	// switch msgv := msg.(type) {
	// case tea.KeyMsg:

	// }
	// log.Println("Edit Model rec", msg.(tea.KeyMsg).String())
	return m, nil
}

func (m EditModel) View() string {
	var b strings.Builder
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i != len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}
	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	return b.String()
}
