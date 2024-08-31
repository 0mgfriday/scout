package internal

import (
	"fmt"
	"sync"
)

type MultiScanner struct {
	scan Scanner
}

func NewMultiScanner(scanner Scanner) *MultiScanner {
	newScanner := MultiScanner{
		scan: scanner,
	}

	return &newScanner
}

// func (mutli MultiScanner) Scan(urls []string, outputQueue chan<- Report, maxThreads int) {
// 	for _, u := range urls {
// 		report, err := mutli.scan.Scan(u)
// 		if err == nil {
// 			outputQueue <- *report
// 		}
// 	}

// 	close(outputQueue)
// }

func (multi MultiScanner) Scan(urls []string, outputQueue chan<- Report, maxThreads int) {
	requestQueue := make(chan string, maxThreads)
	wg := &sync.WaitGroup{}
	for i := 0; i < maxThreads; i++ {
		wg.Add(1)
		go multi.scanWorker(wg, requestQueue, outputQueue)
	}

	for _, u := range urls {
		requestQueue <- u
	}
	close(requestQueue)
	wg.Wait()
	close(outputQueue)
}

func (multi MultiScanner) scanWorker(wg *sync.WaitGroup, requestQueue chan string, outputQueue chan<- Report) {
	defer wg.Done()
	for item := range requestQueue {
		result, err := multi.scan.Scan(item)
		if err == nil {
			outputQueue <- *result
		} else {
			fmt.Println(err)
		}
	}
}
