package main

import (
	"errors"
	"os"
	"strings"
)

func ReadScopeFile(filePath string) ([]string, error) {
	if _, err := os.Stat(filePath); err == nil {
		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil, errors.New("Unable to read scope file " + filePath)
		}

		lines := strings.Split(string(content), "\n")
		return lines, nil
	} else {
		return nil, errors.New("Scope file " + filePath + " does not exist")
	}
}
