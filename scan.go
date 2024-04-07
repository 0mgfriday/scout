package main

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"
)

func Scan(u string) (*Report, error) {
	if !strings.HasPrefix(u, "http") {
		u = "https://" + u
	}
	uri, err := url.ParseRequestURI(u)
	if err != nil {
		return nil, errors.New("invalid url")
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Get(uri.String())
	if err != nil {
		panic(err)
	}

	IPs, _ := net.LookupIP(uri.Host)
	return reportFromResponse(uri.String(), IPs, resp), nil
}
