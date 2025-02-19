package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pawndev/resizeboy/internal/vars"

	"github.com/pawndev/resizeboy/pkg/img"
	"github.com/pawndev/resizeboy/pkg/task"
	"github.com/remeh/sizedwaitgroup"
)

const (
	MaxGoRoutines = 50
)

type App struct {
	Vars   *vars.Vars
	Report func(<-chan *task.Result)
}

func New(vars *vars.Vars) *App {
	return &App{
		Vars: vars,
	}
}

func (a *App) Run() {
	files, err := os.ReadDir(a.Vars.InputDir)
	if err != nil {
		panic(err)
	}

	// Ensure the dist directory exists
	if _, err := os.Stat(a.Vars.OutputDir); os.IsNotExist(err) {
		err := os.MkdirAll(a.Vars.OutputDir, os.ModeDir)
		if err != nil {
			panic(err)
		}
	}

	swg := sizedwaitgroup.New(MaxGoRoutines)
	resChan := make(chan *task.Result, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		swg.Add()

		go func(dirEntry os.DirEntry) {
			defer swg.Done()
			start := time.Now()
			res := &task.Result{
				FileName: dirEntry.Name(),
			}

			defer func(r *task.Result) {
				r.Duration = time.Since(start)
				resChan <- res
			}(res)

			f, err := os.Open(filepath.Join(a.Vars.InputDir, dirEntry.Name()))
			defer func(f *os.File) {
				_ = f.Close()
			}(f)
			if err != nil {
				res.Err = err
				return
			}

			ext := filepath.Ext(dirEntry.Name())
			filename := strings.TrimSuffix(dirEntry.Name(), ext)

			i, err := img.Open(f)
			if err != nil {
				res.Err = err
				return
			}

			w, err := strconv.ParseUint(a.Vars.MaxWidth, 10, 64)
			if err != nil {
				res.Err = err
				return
			}

			i.Resize(uint(w))

			outFilename := filename
			if a.Vars.ShouldAddSuffix {
				outFilename = fmt.Sprintf("%s.%s", filename, a.Vars.FileSuffix)
			}

			out, err := os.Create(filepath.Join(a.Vars.OutputDir, fmt.Sprintf("%s.%s", outFilename, a.Vars.OutFormat)))
			if err != nil {
				res.Err = err
				return
			}
			err = i.Encode(img.PNGFormat, out)
			if err != nil {
				res.Err = err
			}
		}(file)
	}

	go func() {
		swg.Wait()
		close(resChan)
	}()

	a.Report(resChan)
}
