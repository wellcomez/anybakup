package cmd

import (
	"fmt"

	"anybakup/util"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	selectedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#F7DC6F"))
	cursorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#EB984E"))
	unselectedStyle = lipgloss.NewStyle()
	titleStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#2ECC71")).Bold(true)
)

type model struct {
	cursor   int
	choices  []string
	selected string
	done     bool
	title    string
}

type inputModel struct {
	textInput textinput.Model
	title     string
	done      bool
	value     string
	err       error
}

// type tickMsg struct{}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.done {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			m.selected = m.choices[m.cursor]
			m.done = true
			return m, tea.Quit

		case tea.KeyUp, tea.KeyShiftTab:
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.choices) - 1
			}

		case tea.KeyDown, tea.KeyTab:
			m.cursor++
			if m.cursor >= len(m.choices) {
				m.cursor = 0
			}
		}
	}

	return m, nil
}

func (m inputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.done = true
			m.value = m.textInput.Value()
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case error:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m inputModel) View() string {
	if m.done {
		return ""
	}

	title := "Enter value:"
	if m.title != "" {
		title = m.title
	}

	titleText := titleStyle.Render(title)
	inputField := m.textInput.View()

	return fmt.Sprintf("%s\n\n%s\n\nPress Enter to confirm, Esc/Ctrl+C to quit", titleText, inputField)
}

func (m model) View() string {
	if m.done {
		return ""
	}

	title := "Choose an option:"
	if m.title != "" {
		title = m.title
	}
	s := titleStyle.Render(title) + "\n\n"

	for i, choice := range m.choices {
		cursor := " "
		style := unselectedStyle
		if m.cursor == i {
			cursor = ">"
			style = selectedStyle
		}

		cursor = cursorStyle.Render(cursor)
		choice = style.Render(choice)
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\nPress Enter to select, Esc to quit"

	return s
}

func GetTagOption(c *util.Config) (string, error) {
	tags, _ := GetAllTags(c)
	for _, v := range tags {
		fmt.Println(v)
	}
	return ShowInput("Enter tag name:", "tag name")
}

func ShowProfileOption() (string, error) {
	config := util.NewConfig()

	// Extract profile names from the map
	var profileNames []string
	if config.Profile != nil {
		for name := range config.Profile {
			c := config.GetProfile(name)
			if c == nil {
				continue
			}
			ss := fmt.Sprintf("%-10s %-150s", name, string(c.RepoDir))
			profileNames = append(profileNames, ss)
		}
	}

	if len(profileNames) == 0 {
		return "default", nil
	}

	return ShowOption(profileNames, "Select a profile:")
}

func ShowOption(options []string, title ...string) (string, error) {
	if len(options) == 0 {
		return "", fmt.Errorf("no options provided")
	}

	initialModel := model{
		choices: options,
		title:   "Choose an option:",
	}

	if len(title) > 0 && title[0] != "" {
		initialModel.title = title[0]
	}

	p := tea.NewProgram(initialModel)
	m, err := p.Run()
	if err != nil {
		return "", err
	}

	modelInstance, ok := m.(model)
	if !ok {
		return "", fmt.Errorf("unexpected model type")
	}

	if modelInstance.selected == "" {
		return "", fmt.Errorf("no option selected")
	}

	return modelInstance.selected, nil
}

func ShowInput(title string, placeholder string) (string, error) {
	input := textinput.New()
	input.Placeholder = placeholder
	input.Focus()
	input.Prompt = "> "

	initialModel := inputModel{
		textInput: input,
		title:     title,
	}

	p := tea.NewProgram(initialModel)
	m, err := p.Run()
	if err != nil {
		return "", err
	}

	modelInstance, ok := m.(inputModel)
	if !ok {
		return "", fmt.Errorf("unexpected model type")
	}

	return modelInstance.value, nil
}
