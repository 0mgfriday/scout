package internal

import (
	"strings"
)

type DiscoveryScanner struct {
	multiScan *MultiScanner
	scope     map[string]bool
	seenUrls  map[string]bool
}

func NewDiscoverScanner(multiScanner *MultiScanner, scope map[string]bool) *DiscoveryScanner {
	return &DiscoveryScanner{
		multiScan: multiScanner,
		scope:     scope,
		seenUrls:  make(map[string]bool),
	}
}

func (discovery DiscoveryScanner) Scan(urls []string, outputQueue chan<- Report, maxThreads int) {
	for _, u := range urls {
		discovery.seenUrls[trimProtocol(u)] = true
	}
	targetUrls := &urls

	for len(*targetUrls) > 0 {
		multiScanOutputQueue := make(chan Report)
		go discovery.multiScan.Scan(*targetUrls, multiScanOutputQueue, maxThreads)
		var newTargets []string

		for item := range multiScanOutputQueue {
			var foundDomains []string
			outputQueue <- item
			certDomains := getDomainsFromCert(item)
			foundDomains = append(foundDomains, certDomains...)

			for _, d := range foundDomains {
				if !discovery.seenUrls[d] {
					discovery.seenUrls[d] = true
					newTargets = append(newTargets, d)
				}
			}
		}

		targetUrls = &newTargets
	}
	close(outputQueue)
}

func getDomainsFromCert(report Report) []string {
	var domains []string
	if report.TLS.Domains != nil {
		for _, d := range report.TLS.Domains {
			d = strings.TrimPrefix(d, "*")
			d = strings.TrimPrefix(d, ".")
			domains = append(domains, d)
		}
	}

	return domains
}

func trimProtocol(s string) string {
	d := strings.TrimPrefix(s, "https://")
	d = strings.TrimPrefix(d, "http://")

	return d
}
