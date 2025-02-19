package tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pawndev/resizeboy/pkg/task"
)

var (
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Margin(1, 0)
	appStyle     = lipgloss.NewStyle().Margin(1, 2, 0, 2) //nolint:mnd // because of the margin values
)

type model struct {
	spinner  spinner.Model
	results  []*task.Result
	quitting bool
}

func newModel() model {
	s := spinner.New()
	s.Style = spinnerStyle
	return model{
		spinner:  s,
		results:  []*task.Result{},
		quitting: false,
	}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) View() string {
	var s string

	if m.quitting {
		s += "That’s all for today!"
	} else {
		s += m.spinner.View() + " Converting and resizing your beautiful images..."
	}

	s += "\n\n"

	for _, res := range m.results {
		s += res.String() + "\n"
	}

	if !m.quitting {
		s += helpStyle.Render("Press any key to exit")
	}

	if m.quitting {
		s += "\n"
	}

	return appStyle.Render(s)
}

//nolint:ireturn // I do it the way bubbletea wants me to do
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.quitting = true
		return m, tea.Quit
	case *task.Result:
		m.results = append(m.results, msg)
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func Report(results <-chan *task.Result) {
	p := tea.NewProgram(newModel())
	go func() {
		for r := range results {
			p.Send(r)
		}
		p.Send(tea.KeyMsg{Type: tea.KeyEnter})
	}()

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %e\n", err)
		os.Exit(1)
	}
}
