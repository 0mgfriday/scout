package main

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
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

	result := reportFromResponse(resp)

	j, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(j))
}

type TLSInfo struct {
	SubjectCommonName   string
	SubjectOrganization string
	Issuer              string
	NotBefore           time.Time
	NotAfter            time.Time
	Domains             []string
}

func infoFromCert(cert *x509.Certificate) *TLSInfo {
	info := TLSInfo{
		SubjectCommonName:   cert.Subject.CommonName,
		SubjectOrganization: strings.Join(cert.Subject.Organization, ", "),
		Issuer:              strings.Join(cert.Issuer.Organization, ", "),
		NotBefore:           cert.NotBefore,
		NotAfter:            cert.NotAfter,
		Domains:             cert.DNSNames,
	}

	return &info
}

type Report struct {
	Url     string
	IP      string
	TLS     TLSInfo
	Status  int
	Title   string
	Tech    string
	Headers map[string]string
	JSFiles []string
}

func reportFromResponse(rsp *http.Response) *Report {
	report := Report{
		Url:    rsp.Request.RequestURI,
		IP:     rsp.Request.RemoteAddr,
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
	r := regexp.MustCompile(`<script src="(?P<Src>.{1,200})"`)
	matches := r.FindAllStringSubmatch(s, -1)
	for _, match := range matches {
		index := r.SubexpIndex("Src")
		if index != -1 {
			files = append(files, match[index])
		}
	}

	return files
}
