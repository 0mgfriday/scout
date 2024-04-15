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
	timeout := flag.Int("timeout", 5, "Connection and request timeout")
	maxThreads := flag.Int("threads", 1, "Max number of threads to use for requests")
	outputFile := flag.String("o", "", "File to write results to")
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
		if _, err := os.Stat(*targetList); err == nil {
			wordList, err := os.Open(*targetList)
			if err != nil {
				fmt.Println(err)
			}
			defer wordList.Close()

			wordListScanner := bufio.NewScanner(wordList)

			if *outputFile != "" {
				scanToFile(*scanner, *wordListScanner, *outputFile, *maxThreads)
			} else {
				scanToStdOut(*scanner, *wordListScanner, *maxThreads)
			}

		} else if errors.Is(err, os.ErrNotExist) {
			fmt.Println("File " + *targetList + " does not exist")
		}
	} else {
		fmt.Println("Must provide -u or -l parameter. -h for more details")
		os.Exit(0)
	}
}
