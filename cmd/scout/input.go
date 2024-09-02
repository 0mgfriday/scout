package main

import (
	"bufio"
	"errors"
	"os"
)

func ReadFileLines(filePath string) ([]string, error) {
	if _, err := os.Stat(filePath); err == nil {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		wordListScanner := bufio.NewScanner(file)
		wordListScanner.Split(bufio.ScanLines)

		var lines []string
		for wordListScanner.Scan() {
			lines = append(lines, wordListScanner.Text())
		}

		return lines, nil
	} else {
		return nil, errors.New("File " + filePath + " does not exist")
	}
}
