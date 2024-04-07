package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"

	wappalyzer "github.com/projectdiscovery/wappalyzergo"
)

type Report struct {
	Url        string
	IPs        []string
	TLS        TLSInfo
	Status     int
	Title      string
	Wappalyzer string
	Headers    map[string]string
	JSFiles    []string
}

func reportFromResponse(url string, IPAddresses []net.IP, rsp *http.Response) *Report {
	report := Report{
		Url:    url,
		IPs:    IPsToString(IPAddresses),
		Status: rsp.StatusCode,
	}

	if rsp.TLS != nil && len(rsp.TLS.PeerCertificates) > 0 {
		report.TLS = *infoFromCert(rsp.TLS.PeerCertificates[0])
	}

	report.Headers = headersToMap(rsp.Header)

	//body := readAsString(rsp)
	defer rsp.Body.Close()
	if rsp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(rsp.Body)
		if err != nil {
			fmt.Println(err)
		}
		bodyString := string(bodyBytes)

		if len(bodyBytes) > 0 {
			report.Title = getTitle(bodyString)
			report.JSFiles = getJSFiles(bodyString)
			report.Wappalyzer = getWappalyzerResult(rsp.Header, bodyBytes)
		}
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

func getWappalyzerResult(header http.Header, body []byte) string {
	tech := ""
	wappalyzerClient, err := wappalyzer.New()
	if err == nil {
		fingerprints := wappalyzerClient.Fingerprint(header, body)
		for f := range fingerprints {
			tech += f + ", "
		}
	}

	return strings.TrimSuffix(tech, ", ")
}
