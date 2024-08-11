package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

/* koanf config stuff */

type Scale struct {
	Letter string  `yaml:"Letter"`
	Min    float64 `yaml:"Min"`
}

type Config struct {
	Scale       []Scale  `yaml:"scale"`
	Assignments []string `yaml:"assignments"`
}

// Global koanf instance. Use "." as the key path delimiter. This can be "/" or any character.
var (
	k      = koanf.New(".")
	parser = yaml.Parser()
)

/* bubbletea */

type (
	errMsg error
)

type model struct {
	textarea textarea.Model
	err      error
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	if err := k.Load(file.Provider("mock.yaml"), parser); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	var config Config
	if err := k.Unmarshal("", &config); err != nil {
		log.Fatalf("error unmarshaling config: %v", err)
	}

	// fmt.Printf("%+v\n", config.Assignments)

}

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = ""
	ta.Focus()

	return model{
		textarea: ta,
		err:      nil,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)

}

func (m model) View() string {
	text := fmt.Sprintf(
		"Enter your assignments below\n\n"+
			"eg: 100 x 2 Reading Responses\n"+
			"    250 x 1 Midterm\n"+
			"    1200 Participation\n"+
			"\n%s\n\n%s",
		m.textarea.View(),
		"(ctrl+c to quit)",
	) + "\n"

	style := lipgloss.NewStyle().
		SetString(text)
	return style.Render()
}
