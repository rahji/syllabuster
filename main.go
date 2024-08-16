package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rahji/syllabuster/config"
	"github.com/rahji/syllabuster/ui"
)

func main() {

	conf, err := config.ReadConfig("syllabuster.yaml")
	if err != nil {
		log.Fatalf("readconfig %s", err)
	}

	p := tea.NewProgram(ui.InitialModel(conf), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

}
