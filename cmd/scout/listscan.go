package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"omg/scout/internal"
	"os"
	"sync"
)

func scanToStdOut(scanner internal.Scanner, wordListScanner bufio.Scanner, maxThreads int) {
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

func scanToStdOutWorker(wg *sync.WaitGroup, requestQueue chan string, scanner internal.Scanner) {
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

func scanToFile(scanner internal.Scanner, wordListScanner bufio.Scanner, outfile string, maxThreads int) error {
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

func scanToFileWorker(wg *sync.WaitGroup, requestQueue chan string, scanner internal.Scanner, file *os.File, completed *int) {
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
