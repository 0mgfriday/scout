package internal

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Scanner struct {
	client             http.Client
	impersonateBrowser bool
	checkCORS          bool
}

func NewScanner(timeout int, impersonateBrowser bool, proxy string, checkCORS bool) (*Scanner, error) {
	newScanner := Scanner{
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
		impersonateBrowser: impersonateBrowser,
		checkCORS:          checkCORS,
	}

	if proxy != "" {
		proxyUrl, err := url.Parse(proxy)
		if err == nil && proxyUrl.Scheme != "" {
			newScanner.client.Transport.(*http.Transport).Proxy = http.ProxyURL(proxyUrl)
		} else {
			return nil, errors.New("Invalid proxy URL: " + proxy)
		}
	}

	return &newScanner, nil
}

func (scan Scanner) Scan(u string) (*Report, error) {
	if !strings.HasPrefix(u, "http") {
		u = "https://" + u
	}
	uri, err := url.ParseRequestURI(u)
	if err != nil {
		return nil, errors.New("invalid url " + u)
	}

	rsp, err := scan.getWithRetry(uri.String(), 2)
	if err != nil {
		return nil, err
	}

	IPs, _ := net.LookupIP(uri.Host)
	return reportFromResponse(rsp.Request.URL.String(), IPs, rsp), nil
}

func (scan Scanner) getWithRetry(url string, attempts int) (*http.Response, error) {
	req, err := scan.createReq(url)

	for i := 1; i <= attempts; i++ {
		if err != nil {
			continue
		}

		rsp, err := scan.client.Do(req)

		if err == nil {
			return rsp, nil
		}
	}

	// fallback to other protocol
	rsp, err := scan.doFallbackRequest(req)
	if err == nil {
		return rsp, nil
	}

	return nil, errors.New("request failed: " + err.Error())
}

func (scan Scanner) doFallbackRequest(req *http.Request) (*http.Response, error) {
	if req.URL.Scheme == "https" {
		req.URL.Scheme = "http"

		return scan.client.Do(req)
	} else {
		req.URL.Scheme = "https"

		return scan.client.Do(req)
	}
}

func (scan Scanner) createReq(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if scan.impersonateBrowser {
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")
		req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		req.Header.Add("Accept-Language", "en-US,en;q=0.9")
		req.Header.Add("Upgrade-Insecure-Requests", "1")
		req.Header.Add("Sec-Fetch-Dest", "document")
		req.Header.Add("Sec-Fetch-Mode", "navigate")
		req.Header.Add("Sec-Fetch-Site", "none")
		req.Header.Add("Sec-Fetch-User", "?1")
	} else {
		req.Header.Add("Accept", "*/*")
	}

	if scan.checkCORS {
		req.Header.Add("Origin", "https://example.com")
	}

	return req, nil
}
