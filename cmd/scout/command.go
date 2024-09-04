package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/0mgfriday/scout/internal"
)

func singleTargetScan(scanner internal.Scanner, targetUrl string, jsonOutput bool) {
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

func multiTargetScan(scanner internal.Scanner, targetList string, scopeList string, outputFilePath string, discovery bool, maxThreads int) {
	multiScan := getMultiScanner(scanner, discovery, scopeList)
	targets, err := readFileLines(targetList)
	if err == nil {

		outputQueue := make(chan internal.Report)
		go multiScan.Scan(targets, outputQueue, maxThreads)

		output := getOutput(outputFilePath)

		for item := range outputQueue {
			output.outputReport(item)
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
			scope, err = readFileLines(scopeListFilePath)
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
		}

		return internal.NewDiscoverScanner(multiScanner, scope)
	}

	return multiScanner
}
