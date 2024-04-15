package main

import (
	"encoding/json"
	"fmt"
	"omg/scout/internal"
	"strconv"
	"strings"
)

func printAsJson(obj any) {
	j, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(j))
}

func prettyPrintAsJson(obj any) {
	j, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(j))
}

const LabelColor = "\033[1;36m%s\033[0m"

func printReport(r internal.Report) {
	printValue("URL", r.Url)
	printValue("IPs", strings.Join(r.IPs, ", "))
	printValue("TLS.Subject", r.TLS.SubjectCommonName)
	printValue("TLS.Organization", r.TLS.SubjectOrganization)
	printValue("TLS.Issuer", r.TLS.Issuer)
	printValue("TLS.NotAfter", r.TLS.NotAfter.String())
	printValue("TLS.Domains", strings.Join(r.TLS.Domains, ", "))
	printValue("Status", strconv.Itoa(r.Status))
	printValue("Title", r.Title)
	printValue("Wappalyzer", r.Wappalyzer)

	headersOut, _ := json.MarshalIndent(r.Headers, "", "    ")
	printValue("Headers", string(headersOut))

	jsFileOut, _ := json.MarshalIndent(r.JSFiles, "", "    ")
	printValue("JSFiles", string(jsFileOut))
}

func printValue(label string, value string) {
	fmt.Printf(LabelColor, label)
	fmt.Println(": " + value)
}
