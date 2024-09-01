package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/0mgfriday/scout/internal"
)

type Output interface {
	OutputReport(report internal.Report)
}

type ConsoleOutput struct {
}

func NewConsoleOutput() *ConsoleOutput {
	return &ConsoleOutput{}
}

func (consoleOut ConsoleOutput) OutputReport(report internal.Report) {
	printAsJson(report)
}

type FileOutput struct {
	file *os.File
}

func NewFileOutput(file *os.File) *FileOutput {
	return &FileOutput{
		file: file,
	}
}

func (fileOut FileOutput) OutputReport(report internal.Report) {
	j, err := json.Marshal(report)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Fprintln(fileOut.file, string(j))
}

func (fileOut FileOutput) Close() {
	fileOut.file.Close()
}
