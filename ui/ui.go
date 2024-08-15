package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rahji/syllabuster/assignment"
	"github.com/rahji/syllabuster/config"
	"github.com/rahji/syllabuster/pie"
	"github.com/rahji/syllabuster/scale"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("201"))
	cursorStyle  = focusedStyle
	noStyle      = lipgloss.NewStyle()
	greyStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("242"))
)

type (
	errMsg error
)

type model struct {
	keys         keyMap
	help         help.Model
	mdfileinput  textinput.Model
	pngfileinput textinput.Model
	textarea     textarea.Model
	conf         config.Config
	err          error
}

func InitialModel(cfg config.Config) model {

	ta := textarea.New()
	ta.FocusedStyle.CursorLine = noStyle
	ta.ShowLineNumbers = false
	ta.SetHeight(18)
	ta.SetWidth(40)
	ta.Placeholder = ""
	ta.Cursor.Style = focusedStyle
	ta.Focus()

	in1 := textinput.New()
	in1.CharLimit = 32
	in1.SetValue("output.md")

	in2 := textinput.New()
	in2.CharLimit = 32
	in2.SetValue("chart.png")

	return model{
		keys:         keys,
		help:         help.New(),
		mdfileinput:  in1,
		pngfileinput: in2,
		textarea:     ta,
		conf:         cfg,
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
		case key.Matches(msg, m.keys.Generate):
			al := assignment.NewAssignmentList(m.textarea.Value())

			letterscale := scale.Rescale(m.conf.Scale, al.SemesterPoints)
			outputbytes := []byte(letterscale + al.Markdown())
			err := os.WriteFile(
				m.mdfileinput.Value(),
				outputbytes,
				0666,
			)
			if err != nil {
				m.err = err
				return m, nil
			}

			err = pie.Draw(m.pngfileinput.Value(), al.ChartVals())
			if err == nil {
				m.err = err
				return m, nil
			}
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
	eg := "eg: 100 x 2 Reading Responses (Readings)\n" +
		"    250 x 1 Midterm Exam (Midterm)\n" +
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