package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/0mgfriday/scout/internal"
)

func main() {
	targetUrl := flag.String("u", "", "Target URL")
	targetList := flag.String("l", "", "File with list of target URLs")
	scopeList := flag.String("sl", "", "File with list of in-scope URLs for discovery")
	discovery := flag.Bool("d", false, "Discover and scan new in-scope domains")
	impersonate := flag.Bool("i", false, "Impersonate browser when sending requests")
	checkCORS := flag.Bool("cors", false, "Probe for CORS response headers")
	timeout := flag.Int("timeout", 5, "Connection and request timeout in seconds")
	maxThreads := flag.Int("threads", 1, "Max number of threads to use for requests")
	outputFilePath := flag.String("o", "", "File to write results to")
	proxy := flag.String("proxy", "", "Proxy URL")
	jsonOutput := flag.Bool("json", false, "Output as JSON for single URL scan (list always outputs JSON)")

	flag.Parse()

	if *targetUrl == "" && *targetList == "" {
		fmt.Println("Must provide -u or -l parameter. -h for more details")
		os.Exit(0)
	}

	scanner, err := internal.NewScanner(*timeout, *impersonate, *proxy, *checkCORS)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	if *targetUrl != "" {
		SingleTargetScan(*scanner, *targetUrl, *jsonOutput)
	} else if *targetList != "" {
		MultiTargetScan(*scanner, *targetList, *scopeList, *outputFilePath, *discovery, *maxThreads)
	}
}
