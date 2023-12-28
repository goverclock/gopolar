package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func setup() {
	log.SetPrefix("[core]")
	log.SetFlags(0)
}

func main() {
	setup()
	end := NewCLIEnd()
	defer end.Quit()
	m := NewUIModel(end)
	if _, err := tea.NewProgram(m).Run(); err != nil {
		log.Println("fail to setup UI model", err)
		os.Exit(1)
	}
}
