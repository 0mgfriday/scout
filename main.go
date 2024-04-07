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

var printLock sync.Mutex
var fileLock sync.Mutex

func main() {
	u := flag.String("u", "", "Target URL")
	l := flag.String("l", "", "File with list of target URLs")
	i := flag.Bool("i", false, "Impersonate browser when sending requests")
	timeout := flag.Int("timeout", 5, "Connection and request timeout")
	maxThreads := flag.Int("threads", 1, "Max number of threads to use for requests")
	outputFile := flag.String("o", "", "File to write results to")

	flag.Parse()
	scanner := newScanner(*timeout, *i)

	if *u != "" {
		result, err := scanner.Scan(*u)
		if err == nil {
			prettyPrintAsJson(result)
		} else {
			fmt.Println(err)
		}
	} else if *l != "" {
		if _, err := os.Stat(*l); err == nil {
			wordList, err := os.Open(*l)

			if err != nil {
				fmt.Println(err)
			}
			wordListScanner := bufio.NewScanner(wordList)

			if *outputFile != "" {
				scanToFile(*scanner, *wordListScanner, *outputFile, *maxThreads)
			} else {
				scanToStdOut(*scanner, *wordListScanner, *maxThreads)
			}

			wordList.Close()
		} else if errors.Is(err, os.ErrNotExist) {
			fmt.Println("File " + *l + " does not exist")
		}
	} else {
		fmt.Println("Must provide -u or -l parameter. -h for more details")
		os.Exit(0)
	}
}

func printAsJson(obj any) {
	j, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(j))
}

func prettyPrintAsJson(obj any) {
	j, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(j))
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
			printLock.Lock()
			printAsJson(result)
			printLock.Unlock()
		}
	}
}

func scanToFile(scanner scanner, wordListScanner bufio.Scanner, outfile string, maxThreads int) error {
	f, err := os.Create(outfile)
	if err != nil {
		return errors.New("failed to create file " + outfile)
	}

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
	f.Close()

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

			fileLock.Lock()
			fmt.Fprintln(file, string(j))
			fileLock.Unlock()
		}

		*completed++
		fmt.Printf("\r%d URLs scanned", *completed)
	}
}
