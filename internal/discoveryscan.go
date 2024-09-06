package internal

import (
	"regexp"
	"strings"
)

type DiscoveryScanner struct {
	multiScan *MultiScanner
	scope     []string
	seenUrls  map[string]bool
}

func NewDiscoverScanner(multiScanner *MultiScanner, scope []string) *DiscoveryScanner {
	return &DiscoveryScanner{
		multiScan: multiScanner,
		scope:     scope,
		seenUrls:  make(map[string]bool),
	}
}

func (discovery *DiscoveryScanner) Scan(urls []string, outputQueue chan<- Report, maxThreads int) {
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
			foundDomains = append(foundDomains, getDomainsFromCert(item)...)
			if _, ok := item.Headers["Content-Security-Policy"]; ok {
				foundDomains = append(foundDomains, getDomainsFromCSP(item.Headers["Content-Security-Policy"])...)
			}

			for _, d := range foundDomains {
				if !discovery.seenUrls[d] {
					if discovery.inScope(d) {
						discovery.seenUrls[d] = true
						newTargets = append(newTargets, d)
					}
				}
			}
		}

		targetUrls = &newTargets
	}
	close(outputQueue)
}

func (discover *DiscoveryScanner) inScope(domain string) bool {
	for _, s := range discover.scope {
		if strings.HasSuffix(domain, s) {
			return true
		}
	}

	return false
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

func getDomainsFromCSP(csp string) []string {
	r := regexp.MustCompile(`[\w-]+\.[\w.-]+`)
	return r.FindAllString(csp, -1)
}

func trimProtocol(s string) string {
	d := strings.TrimPrefix(s, "https://")
	d = strings.TrimPrefix(d, "http://")

	return d
}
