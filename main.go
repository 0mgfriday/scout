package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"sync"
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

	scanner, err := newScanner(*timeout, *impersonate, *proxy)
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

func scanToStdOut(scanner scanner, wordListScanner bufio.Scanner, maxThreads int) {
	requestQueue := make(chan string, maxThreads)
	wg := &sync.WaitGroup{}
	for i := 0; i < maxThreads; i++ {
		wg.Add(1)
		go scanToStdOutWorker(wg, requestQueue, scanner)
	}

	wordListScanner.Split(bufio.ScanLines)

	for wordListScanner.Scan() {
		requestQueue <- wordListScanner.Text()
	}
	close(requestQueue)
	wg.Wait()
}

func scanToStdOutWorker(wg *sync.WaitGroup, requestQueue chan string, scanner scanner) {
	defer wg.Done()
	for item := range requestQueue {
		result, err := scanner.Scan(item)
		if err == nil {
			printAsJson(result)
		} else {
			fmt.Println(err)
		}
	}
}

func scanToFile(scanner scanner, wordListScanner bufio.Scanner, outfile string, maxThreads int) error {
	f, err := os.Create(outfile)
	if err != nil {
		return errors.New("failed to create file " + outfile)
	}
	defer f.Close()

	completed := 0
	requestQueue := make(chan string, maxThreads)
	wg := &sync.WaitGroup{}
	for i := 0; i < maxThreads; i++ {
		wg.Add(1)
		go scanToFileWorker(wg, requestQueue, scanner, f, &completed)
	}

	for wordListScanner.Scan() {
		requestQueue <- wordListScanner.Text()
	}
	close(requestQueue)
	wg.Wait()

	return nil
}

func scanToFileWorker(wg *sync.WaitGroup, requestQueue chan string, scanner scanner, file *os.File, completed *int) {
	defer wg.Done()
	for item := range requestQueue {
		result, err := scanner.Scan(item)
		if err == nil {
			j, err := json.Marshal(result)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Fprintln(file, string(j))
		}

		*completed++
		fmt.Printf("\r%d URLs scanned", *completed)
	}
}
