package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
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

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("201"))
	cursorStyle  = focusedStyle
	noStyle      = lipgloss.NewStyle()
	greyStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("242"))
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Tab      key.Binding
	Generate key.Binding
	Quit     key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view.
// It's part of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.Generate, k.Quit}
}

// FullHelp returns keybindings for the expanded help view.
// It's part of the key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return nil
}

var keys = keyMap{
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next field"),
	),
	Generate: key.NewBinding(
		key.WithKeys("ctrl+g"),
		key.WithHelp("ctr+g", "generate files"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

type (
	errMsg error
)

type model struct {
	keys         keyMap
	help         help.Model
	mdfileinput  textinput.Model
	pngfileinput textinput.Model
	textarea     textarea.Model
	err          error
}

func main() {

	// f, err := tea.LogToFile("debug.log", "debug")
	// if err != nil {
	// 	fmt.Println("fatal:", err)
	// 	os.Exit(1)
	// }
	// defer f.Close()

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
	ta.SetHeight(18)
	ta.FocusedStyle.CursorLine = noStyle
	ta.ShowLineNumbers = false
	ta.Placeholder = ""
	ta.Cursor.Style = focusedStyle
	ta.Focus()

	in1 := textinput.New()
	in1.CharLimit = 32
	in1.Placeholder = "output.md"

	in2 := textinput.New()
	in2.CharLimit = 32
	in2.Placeholder = "chart.png"

	return model{
		keys:         keys,
		help:         help.New(),
		mdfileinput:  in1,
		pngfileinput: in2,
		textarea:     ta,
		err:          nil,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Tab):
			if m.textarea.Focused() {
				m.textarea.Blur()
				cmds = append(cmds, m.mdfileinput.Focus())
				m.mdfileinput.PromptStyle = focusedStyle
				m.mdfileinput.TextStyle = focusedStyle
			} else if m.mdfileinput.Focused() {
				m.mdfileinput.Blur()
				m.mdfileinput.PromptStyle = noStyle
				m.mdfileinput.TextStyle = noStyle
				cmds = append(cmds, m.pngfileinput.Focus())
				m.pngfileinput.PromptStyle = focusedStyle
				m.pngfileinput.TextStyle = focusedStyle
			} else if m.pngfileinput.Focused() {
				m.pngfileinput.Blur()
				m.pngfileinput.PromptStyle = noStyle
				m.pngfileinput.TextStyle = noStyle
				cmds = append(cmds, m.textarea.Focus())
			}
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)

	m.mdfileinput, cmd = m.mdfileinput.Update(msg)
	cmds = append(cmds, cmd)

	m.pngfileinput, cmd = m.pngfileinput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	helpView := m.help.View(m.keys)
	intro := "Enter your assignments below\n"
	eg := "eg: 100 x 2 Reading Responses\n" +
		"    250 x 1 Midterm\n" +
		"    1200 Participation\n"
	ta := m.textarea.View()
	fn1 := "Markdown " + m.mdfileinput.View()
	fn2 := "PNG file " + m.pngfileinput.View()

	return fmt.Sprintf("\n%s\n%s\n%s\n\n%s\n%s\n\n%s\n\n",
		intro,
		greyStyle.Render(eg),
		ta,
		fn1,
		fn2,
		helpView,
	)
}
