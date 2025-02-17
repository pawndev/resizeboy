package task

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	errStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	successStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))
	durationStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	dotStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type Result struct {
	FileName string
	Duration time.Duration
	Err      error
}

func (r *Result) String() string {
	if r.Duration == 0 {
		//nolint:mnd // why would I want to extract this?
		return dotStyle.Render(strings.Repeat(".", 30))
	}
	if r.Err != nil {
		return errStyle.Render(fmt.Sprintf("Error processing file: %s - %s %s", r.FileName, r.Err, durationStyle.Render(r.Duration.String())))
		//return fmt.Sprintf("Error processing file: %s - %s", r.FileName, r.Err)
	}

	return successStyle.Render(fmt.Sprintf("Processed file: %s %s", r.FileName, durationStyle.Render(r.Duration.String())))
	//return fmt.Sprintf("Processed file: %s in %s", r.FileName, r.Duration)
}
