package tui

import "github.com/charmbracelet/huh/spinner"

type Loader struct {
	spinner *spinner.Spinner
}

func NewLoader(title string) *Loader {
	return &Loader{
		spinner: spinner.New().Title(title),
	}
}

func (s *Loader) Run(action func()) error {
	return s.spinner.Action(action).Run()
}
