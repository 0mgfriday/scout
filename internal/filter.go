package internal

import "strings"

var excludeHeaders = map[string]bool{
	"connection":             true,
	"referrer-policy":        true,
	"cf-ray":                 true,
	"cf-cache-status":        true,
	"content-language":       true,
	"content-disposition":    true,
	"expect-ct":              true,
	"etag":                   true,
	"pragma":                 true,
	"server-timing":          true,
	"x-request-id":           true,
	"x-content-type-options": true,
	"x-timer":                true,
}

func isNoiseyHeader(header string) bool {
	_, ok := excludeHeaders[strings.ToLower(header)]
	return ok
}

var excludeJavaScript = [...]string{
	"google_tag.script.js",
	"googletagmanager",
	"gtag.js",
	"optimizely",
	"ruxitagentjs",
}

func isNoiseyJSFile(file string) bool {
	for _, exclusion := range excludeJavaScript {
		if strings.Contains(file, exclusion) {
			return true
		}
	}

	return false
}
