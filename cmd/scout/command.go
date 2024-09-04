package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/0mgfriday/scout/internal"
)

func SingleTargetScan(scanner internal.Scanner, targetUrl string, jsonOutput bool) {
	result, err := scanner.Scan(targetUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	if jsonOutput {
		prettyPrintAsJson(result)
	} else {
		printReport(*result)
	}
}

func MultiTargetScan(scanner internal.Scanner, targetList string, scopeList string, outputFilePath string, discovery bool, maxThreads int) {
	multiScan := getMultiScanner(scanner, discovery, scopeList)
	targets, err := ReadFileLines(targetList)
	if err == nil {

		outputQueue := make(chan internal.Report)
		go multiScan.Scan(targets, outputQueue, maxThreads)

		output := GetOutput(outputFilePath)

		for item := range outputQueue {
			output.OutputReport(item)
		}

	} else if errors.Is(err, os.ErrNotExist) {
		fmt.Println("File " + targetList + " does not exist")
	}
}

func getMultiScanner(scanner internal.Scanner, discovery bool, scopeListFilePath string) internal.MultiScan {
	multiScanner := internal.NewMultiScanner(scanner)

	if discovery {
		var scope []string
		var err error
		if scopeListFilePath != "" {
			scope, err = ReadFileLines(scopeListFilePath)
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
		}

		return internal.NewDiscoverScanner(multiScanner, scope)
	}

	return multiScanner
}
