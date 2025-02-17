package main

import (
	"errors"
	"github.com/pawndev/minui-image-resizer/internal/app"

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

	application := app.New(form.Vars)
	application.Report = tui.Report

	if err := loader.Run(func() {
		application.Run()
	}); err != nil {
		panic(err)
	}
}
