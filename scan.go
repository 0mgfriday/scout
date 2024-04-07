package main

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type scanner struct {
	client http.Client
}

func newScanner(timeout int) *scanner {
	newScanner := scanner{
		client: http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Transport: &http.Transport{
				TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
				ResponseHeaderTimeout: time.Duration(timeout) * time.Second,
				Dial: func(network, addr string) (net.Conn, error) {
					return net.DialTimeout(network, addr, time.Duration(timeout)*time.Second)
				},
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

	rsp, err := scan.getWithRetry(uri.String(), 2, impersonateBrowser)
	if err != nil {
		return nil, errors.New("request failed")
	}

	IPs, _ := net.LookupIP(uri.Host)
	return reportFromResponse(rsp.Request.URL.String(), IPs, rsp), nil
}

func (scan scanner) getWithRetry(url string, attempts int, impersonateBrowser bool) (*http.Response, error) {
	req, err := createReq(url, impersonateBrowser)

	for i := 1; i <= attempts; i++ {
		if err != nil {
			continue
		}

		rsp, err := scan.client.Do(req)

		if err == nil {
			return rsp, nil
		}
	}

	// fallback to ther protocol
	rsp, err := scan.doFallbackRequest(req)
	if err == nil {
		return rsp, nil
	}

	return nil, errors.New("request failed")
}

func (scan scanner) doFallbackRequest(req *http.Request) (*http.Response, error) {
	if req.URL.Scheme == "https" {
		req.URL.Scheme = "http"

		return scan.client.Do(req)
	} else {
		req.URL.Scheme = "https"

		return scan.client.Do(req)
	}
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
