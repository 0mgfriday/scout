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

	flag.Parse()

	if *u != "" {
		result, err := Scan(*u)
		if err == nil {
			prettyPrintAsJson(result)
		} else {
			fmt.Println(err)
		}
	} else if *l != "" {
		if _, err := os.Stat(*l); err == nil {
			readFile, err := os.Open(*l)

			if err != nil {
				fmt.Println(err)
			}
			fileScanner := bufio.NewScanner(readFile)

			fileScanner.Split(bufio.ScanLines)

			for fileScanner.Scan() {
				result, err := Scan(fileScanner.Text())
				if err == nil {
					printAsJson(result)
				}
			}

			readFile.Close()
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
