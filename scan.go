package main

import (
	"errors"
	"net"
	"net/http"
	"net/url"
)

func Scan(u string) (*Report, error) {
	uri, err := url.ParseRequestURI(u)
	if err != nil {
		return nil, errors.New("Invalid url")
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(uri.String())
	if err != nil {
		panic(err)
	}

	IPs, _ := net.LookupIP(uri.Host)
	return reportFromResponse(uri.String(), IPs, resp), nil
}
