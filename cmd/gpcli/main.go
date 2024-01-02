package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func setup() {
	log.SetPrefix("[gpcli]")
	log.SetFlags(0)
}

func main() {
	setup()
	end := NewCLIEnd()
	m := NewUIModel(end)
	if _, err := tea.NewProgram(m).Run(); err != nil {
		log.Println("fail to setup UI model", err)
		os.Exit(1)
	}
}
