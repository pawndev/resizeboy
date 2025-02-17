package appform

import (
	"errors"
	"github.com/pawndev/resizeboy/internal/vars"
	"os"
	"strconv"

	"github.com/charmbracelet/huh"
)

var (
	ErrInvalidInput    = errors.New("invalid input")
	ErrDirPathNotExist = errors.New("directory path does not exist")
	ErrUserAborted     = huh.ErrUserAborted
)

type Form struct {
	form *huh.Form
	Vars *vars.Vars
}

func New() *Form {
	mainForm := &Form{
		Vars: &vars.Vars{},
	}
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Image directory to convert/resize").Value(&mainForm.Vars.InputDir).Validate(validatePath),
			huh.NewInput().Title("Output directory").Value(&mainForm.Vars.OutputDir).Validate(func(s string) error {
				if s == "" {
					return ErrInvalidInput
				}

				return nil
			}),
			huh.NewSelect[string]().Title("Output format").Options(
				huh.NewOption("PNG", "png"),
				huh.NewOption("JPEG", "jpg"),
			).Value(&mainForm.Vars.OutFormat),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Add suffix to output files?").
				Value(&mainForm.Vars.ShouldAddSuffix).
				Description("MinUI needs a file named like my_game.<rom_extension>.png").
				Affirmative("Yup").Negative("Nop"),
		),
		huh.NewGroup(
			huh.NewInput().Title("File suffix").Value(&mainForm.Vars.FileSuffix),
		).WithHideFunc(func() bool {
			return !mainForm.Vars.ShouldAddSuffix
		}),
		huh.NewGroup(
			huh.NewInput().Title("Max width").Value(&mainForm.Vars.MaxWidth).Validate(func(str string) error {
				_, err := strconv.ParseUint(str, 10, 64)

				return err
			}),
		),
	)

	mainForm.form = form

	return mainForm
}

func (f *Form) Run() error {
	return f.form.Run()
}

func validatePath(str string) error {
	if str == "" {
		return ErrInvalidInput
	}

	if _, err := os.Stat(str); err != nil {
		return ErrDirPathNotExist
	}

	return nil
}
