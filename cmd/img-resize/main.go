package main

import (
	"errors"

	"github.com/pawndev/minui-image-resizer/internal/appform"
	"github.com/pawndev/minui-image-resizer/pkg/tui"
)

func main() {
	form := appform.New()
	loader := tui.NewLoader("Processing images...")

	if err := form.Run(); err != nil {
		if !errors.Is(err, appform.ErrUserAborted) {
			panic(err)
		}
	}

	if err := loader.Run(func() {
		App(form.Vars.InputDir, form.Vars.OutputDir, form.Vars.FileSuffix, form.Vars.MaxWidth, form.Vars.OutFormat, form.Vars.ShouldAddSuffix)
	}); err != nil {
		panic(err)
	}
}
