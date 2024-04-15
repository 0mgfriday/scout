package internal

import (
	"crypto/x509"
	"strings"
	"time"
)

type TLSInfo struct {
	SubjectCommonName   string
	SubjectOrganization string
	Issuer              string
	NotAfter            time.Time
	Domains             []string
}

func infoFromCert(cert *x509.Certificate) *TLSInfo {
	info := TLSInfo{
		SubjectCommonName:   cert.Subject.CommonName,
		SubjectOrganization: strings.Join(cert.Subject.Organization, ", "),
		Issuer:              strings.Join(cert.Issuer.Organization, ", "),
		NotAfter:            cert.NotAfter,
		Domains:             cert.DNSNames,
	}

	return &info
}
