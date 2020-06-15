package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

func getSplitArgAndFiles(arguments ...string) ([]string, []string) {

	var options []string
	var files []string

	if arguments[len(arguments)-1] == "all" {
		wd, _ := os.Getwd()
		allFiles := getFilesUnderPwd(wd)
		options = arguments[0 : len(arguments)-1]
		files = getGzFiles(allFiles)

		for _, file := range files {
			fmt.Println(file)
		}

		return options, files
	}

	for _, value := range arguments {
		if strings.Contains(value, ".gz") {
			files = append(files, value)
		} else {
			options = append(options, value)
		}
	}
	return options, files
}

func getGzFiles(files []string) []string {
	var gzFiles []string
	for _, value := range files {
		if strings.Contains(value, ".gz") {
			gzFiles = append(gzFiles, value)
		}
	}
	return gzFiles
}

func getChunks(files []string, numOfCores int) [][]string {

	var chunks [][]string
	chunkSize := (len(files) + numOfCores - 1) / numOfCores

	for i := 0; i < len(files); i += chunkSize {

		end := i + chunkSize

		if end > len(files) {
			end = len(files)
		}

		chunks = append(chunks, files[i:end])
	}
	return chunks
}

func getFilesUnderPwd(currentDir string) []string {

	var files []string

	err := filepath.Walk(currentDir, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	if err != nil {
		panic(err)
	}

	return files
}

func main() {

	var wg sync.WaitGroup
	numOfCores := runtime.NumCPU()

	start := time.Now()

	arguments := os.Args[1:]

	// if os.Args[1] == "--help" {
	// 	fmt.Println("This is a command line utlity created to run zgrep in parallel! Please refere zgrep documentation to use pzgrep.")
	// }

	options, files := getSplitArgAndFiles(arguments...)

	chunks := getChunks(files, numOfCores)

	wg.Add(len(chunks))

	for _, chunk := range chunks {
		args := append(options, chunk...)
		go func() {
			defer wg.Done()
			cmd := exec.Command("zgrep", args...)

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			cmd.Run()
		}()
	}
	wg.Wait()
	duration := time.Since(start)

	fmt.Println(duration)
}
