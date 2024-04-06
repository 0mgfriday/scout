package main

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const LabelColor = "\033[1;34m%s\033[0m"

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: certdump https://example.com")
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

	tlsInfo := infoFromCert(resp.TLS.PeerCertificates[0])

	j, err := json.MarshalIndent(tlsInfo, "", "    ")
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
	Port    int
	TLS     TLSInfo
	Title   string
	Status  int
	Tech    string
	Headers map[string]string
	JSFiles []string
}
