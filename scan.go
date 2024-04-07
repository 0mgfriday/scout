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

func (scan scanner) Scan(u string) (*Report, error) {
	if !strings.HasPrefix(u, "http") {
		u = "https://" + u
	}
	uri, err := url.ParseRequestURI(u)
	if err != nil {
		return nil, errors.New("invalid url")
	}

	rsp, _ := scan.getWithRetry(uri.String(), 2)
	IPs, _ := net.LookupIP(uri.Host)
	return reportFromResponse(uri.String(), IPs, rsp), nil
}

func (scan scanner) getWithRetry(url string, attempts int) (*http.Response, error) {
	for i := 1; i <= attempts; i++ {
		rsp, err := scan.client.Get(url)
		if err == nil {
			return rsp, nil
		} else if i == attempts {
			return nil, err
		}
	}

	return nil, errors.New("request failed")
}
