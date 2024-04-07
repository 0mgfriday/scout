package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
)

func main() {
	u := flag.String("u", "", "Target URL")
	l := flag.String("l", "", "File with list of target URLs")
	i := flag.Bool("i", false, "Impersonate browser when sending requests")
	timeout := flag.Int("timeout", 5, "Connection and request timeout")
	outputFile := flag.String("o", "", "File to write results to")

	flag.Parse()
	scanner := newScanner(*timeout)

	if *u != "" {
		result, err := scanner.Scan(*u, *i)
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
				scanToFile(*scanner, *wordListScanner, *outputFile, *i)
			} else {
				scanToStdOut(*scanner, *wordListScanner, *i)
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

func scanToStdOut(scanner scanner, wordListScanner bufio.Scanner, impersonateBrowser bool) {
	wordListScanner.Split(bufio.ScanLines)

	for wordListScanner.Scan() {
		result, err := scanner.Scan(wordListScanner.Text(), impersonateBrowser)
		if err == nil {
			printAsJson(result)
		}
	}
}

func scanToFile(scanner scanner, wordListScanner bufio.Scanner, outfile string, impersonateBrowser bool) error {
	f, err := os.Create(outfile)
	if err != nil {
		return errors.New("failed to create file " + outfile)
	}

	num := 1
	for wordListScanner.Scan() {
		result, err := scanner.Scan(wordListScanner.Text(), impersonateBrowser)
		if err == nil {
			j, err := json.Marshal(result)
			if err != nil {
				fmt.Println(err)
				return errors.New("failed to serialize result for " + result.Url)
			}

			fmt.Fprintln(f, string(j))
		}

		fmt.Printf("\r%d URLs scanned", num)
		num++
	}

	f.Close()

	return nil
}
