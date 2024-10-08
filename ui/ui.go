package ui

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/rahji/syllabuster/assignment"
	"github.com/rahji/syllabuster/config"
	"github.com/rahji/syllabuster/pie"
	"github.com/rahji/syllabuster/scale"
)

var (
	focusedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("201"))
	cursorStyle   = focusedStyle
	noStyle       = lipgloss.NewStyle()
	greyStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("242"))
	leftHalfStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("25")).
			Padding(1).Width(50).Height(24)
	rightHalfStyle = lipgloss.NewStyle().Padding(4)
)

type (
	errMsg    error
	statusMsg string
)

type model struct {
	keys         keyMap
	help         help.Model
	mdfileinput  textinput.Model
	pngfileinput textinput.Model
	textarea     textarea.Model
	conf         config.Config
	timer        timer.Model
	status       string
	preview      string
	err          error
}

const timeout = time.Second * 2

func InitialModel(cfg config.Config) model {

	ta := textarea.New()
	ta.FocusedStyle.CursorLine = noStyle
	ta.ShowLineNumbers = false
	ta.SetHeight(18)
	ta.SetWidth(45)
	ta.Placeholder = ""
	ta.Cursor.Style = focusedStyle
	ta.Focus()

	in1 := textinput.New()
	in1.CharLimit = 32
	in1.SetValue("output.md")

	in2 := textinput.New()
	in2.CharLimit = 32
	in2.SetValue("chart.png")

	m := model{
		keys:         keys,
		help:         help.New(),
		mdfileinput:  in1,
		pngfileinput: in2,
		textarea:     ta,
		conf:         cfg,
		err:          nil,
	}
	m.timer = timer.NewWithInterval(0, time.Microsecond*50)
	return m
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case timer.TimeoutMsg:
		m.timer, cmd = m.timer.Update(msg)
		cmds = append(cmds, cmd)
	case timer.TickMsg:
		m.timer, cmd = m.timer.Update(msg)
		cmds = append(cmds, cmd)
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Generate):
			preview, err := m.generateFiles()
			if err != nil {
				m.status = fmt.Sprintf("❌ %s", err)
			} else {
				m.status = "✅ Saved files"
				m.preview = preview
			}
			m.timer = timer.NewWithInterval(timeout, time.Millisecond*50)
			return m, m.timer.Init()
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

// generates the markdown and PNG files.
// returns a preview string and an error
func (m model) generateFiles() (string, error) {
	al := assignment.NewAssignmentList(m.textarea.Value())

	md := fmt.Sprintf(`
### Grade Distribution

There are %.0f possible points to be earned throughout the semester. 
*Late submissions earn no points.*

%s

### Letter Grade Key

%s
`,
		al.SemesterPoints,
		al.Markdown(),
		scale.Rescale(m.conf.Scale, al.SemesterPoints),
	)
	outputbytes := []byte(md)
	err := os.WriteFile(
		m.mdfileinput.Value(),
		outputbytes,
		0666,
	)
	if err != nil {
		return "", err
	}

	chartvals := al.ChartVals()
	err = pie.Draw(m.pngfileinput.Value(), chartvals)
	if err != nil {
		return "", err
	}

	return md, nil
}

func (m model) View() string {
	helpView := m.help.View(m.keys)
	intro := "Enter your assignments below\n"
	eg := "eg: 100 x 2 Reading Responses (Readings)\n" +
		"    250 x 1 Midterm Exam (Midterm)\n" +
		"    1200 Participation\n"
	ta := m.textarea.View()
	fn1 := "\nMarkdown " + m.mdfileinput.View()
	fn2 := "PNG file " + m.pngfileinput.View() + "\n"

	var statusLine string
	if !m.timer.Timedout() {
		statusLine = m.status
		m.status = ""
	}

	leftSide := lipgloss.JoinVertical(
		lipgloss.Left,
		intro,
		greyStyle.Render(eg),
		ta,
		fn1,
		fn2,
		helpView,
		statusLine,
	)
	glamoured, _ := glamour.Render(m.preview, "dark")
	bothSides := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftHalfStyle.Render(leftSide),
		glamoured,
	)
	return bothSides
}
