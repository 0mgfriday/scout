package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
)

type Report struct {
	Url     string
	IPs     []string
	TLS     TLSInfo
	Status  int
	Title   string
	Tech    string
	Headers map[string]string
	JSFiles []string
}

func reportFromResponse(url string, IPAddresses []net.IP, rsp *http.Response) *Report {
	report := Report{
		Url:    url,
		IPs:    IPsToString(IPAddresses),
		Status: rsp.StatusCode,
	}

	report.TLS = *infoFromCert(rsp.TLS.PeerCertificates[0])
	report.Headers = headersToMap(rsp.Header)

	body := readAsString(rsp)
	if body != "" {
		report.Title = getTitle(body)
		report.JSFiles = getJSFiles(body)
	}

	return &report
}

func IPsToString(IPs []net.IP) []string {
	result := []string{}
	for _, IP := range IPs {
		if ipv4 := IP.To4(); ipv4 != nil {
			result = append(result, ipv4.String())
		}
	}

	return result
}

func headersToMap(header http.Header) map[string]string {
	headers := make(map[string]string)
	for name, values := range header {
		for _, val := range values {
			headers[name] = val // just replace if multiple for now
		}
	}

	return headers
}

func readAsString(rsp *http.Response) string {
	defer rsp.Body.Close()
	bodyString := ""

	if rsp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(rsp.Body)
		if err != nil {
			fmt.Println(err)
		}
		bodyString = string(bodyBytes)
	}

	return bodyString
}

func getTitle(s string) string {
	r := regexp.MustCompile(`<(title|Title|TITLE)>(?P<Title>.{1,150})</(title|Title|TITLE)>`)
	match := r.FindStringSubmatch(s)
	index := r.SubexpIndex("Title")
	if index == -1 {
		return ""
	}

	return match[index]
}

func getJSFiles(s string) []string {
	files := []string{}
	r := regexp.MustCompile(`<script src="(?P<Src>[\w\/\-.]{1,200})["?]`)
	matches := r.FindAllStringSubmatch(s, -1)
	for _, match := range matches {
		index := r.SubexpIndex("Src")
		if index != -1 {
			files = append(files, match[index])
		}
	}

	return files
}
