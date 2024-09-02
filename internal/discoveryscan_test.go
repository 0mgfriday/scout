package internal

import "testing"

func TestGetDomainsFromCSP(t *testing.T) {
	domain1 := "test.com"
	domain2 := "test2.com"
	domain3 := "api-dev.test.io"
	csp := "default-src 'none'; base-uri 'self'; block-all-mixed-content; child-src 'self' " + domain1 + " *." + domain2 + " " + domain3
	domains := getDomainsFromCSP(csp)

	if len(domains) != 3 {
		t.Fatalf("3 domains found, want %d", len(domains))
	}

	if domains[0] != domain1 {
		t.Fatalf("domain[0] = %q, want %q", domains[2], domain1)
	}

	if domains[1] != domain2 {
		t.Fatalf("domain[1] = %q, want %q", domains[2], domain2)
	}

	if domains[2] != domain3 {
		t.Fatalf("domain[2] = %q, want %q", domains[2], domain3)
	}
}
