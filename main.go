package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: scout https://example.com")
		os.Exit(0)
	}

	u, err := url.ParseRequestURI(os.Args[1])
	if err != nil {
		fmt.Println("Invalid url")
		os.Exit(0)
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(u.String())
	if err != nil {
		panic(err)
	}

	IPs, _ := net.LookupIP(u.Host)
	result := reportFromResponse(u.String(), IPs, resp)

	j, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(j))
}
