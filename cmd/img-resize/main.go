package main

import (
	"errors"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"os"
	"strconv"
)

var (
	ErrInvalidInput    = errors.New("invalid input")
	ErrDirPathNotExist = errors.New("directory path does not exist")
)

var (
	inputDir        string
	outputDir       string
	fileSuffix      string
	shouldAddSuffix bool
	maxWidth        string
	outFormat       string
)

func main() {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Image directory to convert/resize").Value(&inputDir).Validate(validatePath),
			huh.NewInput().Title("Output directory").Value(&outputDir).Validate(func(s string) error {
				if s == "" {
					return ErrInvalidInput
				}

				return nil
			}),
			huh.NewSelect[string]().Title("Output format").Options(
				huh.NewOption("PNG", "png"),
				huh.NewOption("JPEG", "jpg"),
			).Value(&outFormat),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Add suffix to output files?").
				Value(&shouldAddSuffix).
				Description("MinUI needs a file named like my_game.<rom_extension>.png").
				Affirmative("Yup").Negative("Nop"),
		),
		huh.NewGroup(
			huh.NewInput().Title("File suffix").Value(&fileSuffix),
		).WithHideFunc(func() bool {
			return !shouldAddSuffix
		}),
		huh.NewGroup(
			huh.NewInput().Title("Max width").Value(&maxWidth).Validate(func(str string) error {
				_, err := strconv.ParseUint(str, 10, 64)

				return err
			}),
		),
	)

	if err := form.Run(); err != nil {
		if err != huh.ErrUserAborted {
			panic(err)
		}
	}

	if err := spinner.New().Title("Processing images...").Action(func() {
		App(inputDir, outputDir, fileSuffix, maxWidth, outFormat, shouldAddSuffix)
	}).Run(); err != nil {
		panic(err)
	}
}

func ensureDirExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}

	return nil
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
