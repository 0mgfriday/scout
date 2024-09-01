package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/0mgfriday/scout/internal"
)

func main() {
	targetUrl := flag.String("u", "", "Target URL")
	targetList := flag.String("l", "", "File with list of target URLs")
	impersonate := flag.Bool("i", false, "Impersonate browser when sending requests")
	timeout := flag.Int("timeout", 5, "Connection and request timeout in seconds")
	maxThreads := flag.Int("threads", 1, "Max number of threads to use for requests")
	outputFilePath := flag.String("o", "", "File to write results to")
	proxy := flag.String("proxy", "", "Proxy URL")
	jsonOutput := flag.Bool("json", false, "Output as JSON for single URL scan (list always outputs JSON)")

	flag.Parse()

	scanner, err := internal.NewScanner(*timeout, *impersonate, *proxy)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	if *targetUrl != "" {
		result, err := scanner.Scan(*targetUrl)
		if err == nil {
			if *jsonOutput {
				prettyPrintAsJson(result)
			} else {
				printReport(*result)
			}
		} else {
			fmt.Println(err)
		}
	} else if *targetList != "" {
		multiScanner := internal.NewMultiScanner(*scanner)
		if _, err := os.Stat(*targetList); err == nil {
			wordList, err := os.Open(*targetList)
			if err != nil {
				fmt.Println(err)
			}
			defer wordList.Close()

			wordListScanner := bufio.NewScanner(wordList)
			wordListScanner.Split(bufio.ScanLines)

			var lines []string
			for wordListScanner.Scan() {
				lines = append(lines, wordListScanner.Text())
			}

			outputQueue := make(chan internal.Report)
			go multiScanner.Scan(lines, outputQueue, *maxThreads)

			var output Output
			if *outputFilePath != "" {
				outFile, err := os.Create(*outputFilePath)
				if err != nil {
					fmt.Println("failed to create file " + *outputFilePath)
					os.Exit(0)
				}
				defer outFile.Close()

				output = NewFileOutput(outFile)
			} else {
				output = NewConsoleOutput()
			}

			for item := range outputQueue {
				output.OutputReport(item)
			}

		} else if errors.Is(err, os.ErrNotExist) {
			fmt.Println("File " + *targetList + " does not exist")
		}
	} else {
		fmt.Println("Must provide -u or -l parameter. -h for more details")
		os.Exit(0)
	}
}
