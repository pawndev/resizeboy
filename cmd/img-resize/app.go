package main

import (
	"fmt"
	"github.com/pawndev/minui-image-resizer/pkg/img"
	"github.com/remeh/sizedwaitgroup"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func App(inputDir, outputDir, fileSuffix, maxWidth, outFormat string, shouldAddSuffix bool) {
	files, err := os.ReadDir(inputDir)
	if err != nil {
		panic(err)
	}

	// Ensure the dist directory exists
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	swg := sizedwaitgroup.New(50)
	resChan := make(chan *Result, len(files))
	doneChan := make(chan bool)
	for _, file := range files {
		swg.Add()
		if file.IsDir() {
			fmt.Println(file.Name(), "is a directory. skipping...")
			continue
		}

		go func(dirEntry os.DirEntry) {
			defer swg.Done()
			start := time.Now()
			res := &Result{
				FileName: dirEntry.Name(),
			}

			defer func(result *Result) {
				res.Duration = time.Since(start)
				resChan <- res
			}(res)

			f, err := os.Open(filepath.Join(inputDir, dirEntry.Name()))
			defer f.Close()
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
				fmt.Println("Error encoding image: ", err)
			}
		}(file)
	}

	go func() {
		swg.Wait()
		doneChan <- true
	}()

	Report(resChan, doneChan)

	close(resChan)

	//for res := range resChan {
	//	if res.Err != nil {
	//		fmt.Println("Error processing file: ", res.FileName, " - ", res.Err)
	//	} else {
	//		fmt.Println("Processed file: ", res.FileName)
	//	}
	//}
}
