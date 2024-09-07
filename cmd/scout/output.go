package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/0mgfriday/scout/internal"
)

type output interface {
	outputReport(report internal.Report)
	close()
}

type consoleOutput struct {
}

func newConsoleOutput() *consoleOutput {
	return &consoleOutput{}
}

func (consoleOut consoleOutput) outputReport(report internal.Report) {
	printAsJson(report)
}

func (consoleOut consoleOutput) close() {
}

type fileOutput struct {
	file  *os.File
	count int32
}

func newFileOutput(file *os.File) *fileOutput {
	return &fileOutput{
		file:  file,
		count: 0,
	}
}

func (fileOut *fileOutput) outputReport(report internal.Report) {
	j, err := json.Marshal(report)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Fprintln(fileOut.file, string(j))
	fileOut.count++
	fmt.Printf("\r%d scan results collected", fileOut.count)
}

func (fileOut *fileOutput) close() {
	fileOut.file.Close()
}

func getOutput(outputFilePath string) output {
	if outputFilePath != "" {
		outFile, err := os.Create(outputFilePath)
		if err != nil {
			fmt.Println("failed to create file " + outputFilePath)
			os.Exit(0)
		}

		return newFileOutput(outFile)
	} else {
		return newConsoleOutput()
	}
}
