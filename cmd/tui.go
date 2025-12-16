package cmd

import (
	"fmt"

	"anybakup/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	selectedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#F7DC6F"))
	cursorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#EB984E"))
	unselectedStyle = lipgloss.NewStyle()
)

type model struct {
	cursor   int
	choices  []string
	selected string
	done     bool
}

type tickMsg struct{}

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

func (m model) View() string {
	if m.done {
		return ""
	}

	s := "Choose an option:\n\n"

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
			ss := fmt.Sprintf("%-10s %50s", name, string(c.RepoDir))
			profileNames = append(profileNames, ss)
		}
	}

	if len(profileNames) == 0 {
		return "default", nil
	}

	return ShowOption(profileNames)
}

func ShowOption(options []string) (string, error) {
	if len(options) == 0 {
		return "", fmt.Errorf("no options provided")
	}

	initialModel := model{
		choices: options,
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
