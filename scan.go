package main

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"
)

type scanner struct {
	client http.Client
}

func newScanner() *scanner {
	newScanner := scanner{
		client: http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}

	return &newScanner
}

func (scan scanner) Scan(u string, impersonateBrowser bool) (*Report, error) {
	if !strings.HasPrefix(u, "http") {
		u = "https://" + u
	}
	uri, err := url.ParseRequestURI(u)
	if err != nil {
		return nil, errors.New("invalid url")
	}

	rsp, _ := scan.getWithRetry(uri.String(), 2, impersonateBrowser)
	IPs, _ := net.LookupIP(uri.Host)
	return reportFromResponse(uri.String(), IPs, rsp), nil
}

func (scan scanner) getWithRetry(url string, attempts int, impersonateBrowser bool) (*http.Response, error) {
	for i := 1; i <= attempts; i++ {
		req, err := createReq(url, impersonateBrowser)
		if err != nil {
			continue
		}

		rsp, err := scan.client.Do(req)

		if err == nil {
			return rsp, nil
		} else if i == attempts {
			return nil, err
		}
	}

	return nil, errors.New("request failed")
}

func createReq(url string, impersonateBrowser bool) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	if impersonateBrowser {
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:124.0) Gecko/20100101 Firefox/124.0")
		req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
		req.Header.Add("Accept-Language", "en-US,en;q=0.5")
		req.Header.Add("Accept-Encoding", "gzip, deflate, br")
		req.Header.Add("Upgrade-Insecure-Requests", "1")
		req.Header.Add("Sec-Fetch-Dest", "document")
		req.Header.Add("Sec-Fetch-Mode", "navigate")
		req.Header.Add("Sec-Fetch-Site", "none")
		req.Header.Add("Sec-Fetch-User", "?1")
	} else {
		req.Header.Add("Accept", "*/*")
	}

	return req, nil
}

func createBrowserReq(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return req, nil
}
