package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pawndev/minui-image-resizer/pkg/img"
	"github.com/pawndev/minui-image-resizer/pkg/task"
	"github.com/pawndev/minui-image-resizer/pkg/tui"
	"github.com/remeh/sizedwaitgroup"
)

const (
	MaxGoRoutines = 50
)

func App(inputDir, outputDir, fileSuffix, maxWidth, outFormat string, shouldAddSuffix bool) {
	files, err := os.ReadDir(inputDir)
	if err != nil {
		panic(err)
	}

	// Ensure the dist directory exists
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.MkdirAll(outputDir, os.ModeDir)
		if err != nil {
			panic(err)
		}
	}

	swg := sizedwaitgroup.New(MaxGoRoutines)
	resChan := make(chan *task.Result, len(files))
	doneChan := make(chan bool)
	for _, file := range files {
		swg.Add()
		if file.IsDir() {
			// fmt.Println(file.Name(), "is a directory. skipping...")
			continue
		}

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

			f, err := os.Open(filepath.Join(inputDir, dirEntry.Name()))
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

			w, err := strconv.ParseUint(maxWidth, 10, 64)
			if err != nil {
				res.Err = err
				return
			}

			i.Resize(uint(w))

			outFilename := filename
			if shouldAddSuffix {
				outFilename = fmt.Sprintf("%s.%s", filename, fileSuffix)
			}

			out, err := os.Create(filepath.Join(outputDir, fmt.Sprintf("%s.%s", outFilename, outFormat)))
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
		doneChan <- true
	}()

	tui.Report(resChan, doneChan)

	close(resChan)
}
