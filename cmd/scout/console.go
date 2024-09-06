package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/0mgfriday/scout/internal"
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

	printJsonValue("Headers", r.Headers)
	printJsonValue("JSFiles", r.JSFiles)
}

func printJsonValue(label string, value any) {
	jsonOut, _ := json.MarshalIndent(value, "", "    ")
	fmt.Printf(LabelColor, label)
	fmt.Println(": " + string(jsonOut))
}

func printValue(label string, value string) {
	fmt.Printf(LabelColor, label)
	fmt.Println(": " + escapeBadCharacters(value))
}

func escapeBadCharacters(s string) string {
	r := fmt.Sprintf("%q", s)
	r = strings.TrimPrefix(r, "\"")
	r = strings.TrimSuffix(r, "\"")

	return strings.ReplaceAll(r, "\\\"", "\"")
}
