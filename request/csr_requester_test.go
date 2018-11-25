package request

import (
	"crypto/x509"
	"testing"

	"github.com/tlmiller/disttrust/provider"
	"github.com/tlmiller/disttrust/util"
)

func TestCSRFromRequestResult(t *testing.T) {
	req := provider.Request{
		CommonName: "example.com",
		AltNames: provider.AltNames{
			DNSNames:       []string{"www.example.com"},
			EmailAddresses: []string{"test@example.com"},
			IPAddresses:    []string{"fe80:aabb::1"},
		},
	}
	csrr := NewCSRRequester(NewRSAKeyMaker(1024))
	csrData, err := csrr.CSRFromRequest(&req)
	if err != nil {
		t.Fatalf("unexpected error generating csr from request: %v", err)
	}

	csr, err := x509.ParseCertificateRequest(csrData)
	if err != nil {
		t.Fatalf("unexpected error building csr represnetation: %v", err)
	}

	if csr.Subject.CommonName != "example.com" {
		t.Errorf("generated csr common name, expected example.com got %s",
			csr.Subject.CommonName)
	}

	if len(csr.DNSNames) != 1 && !util.StringInSlice("www.example.com", csr.DNSNames) {
		t.Errorf("generated csr dns names, expected [www.example.com] got %v",
			csr.DNSNames)
	}

	if len(csr.EmailAddresses) != 1 && !util.StringInSlice(
		"test@example.com", csr.EmailAddresses) {
		t.Errorf("generated csr email addresses, expected [test@example.com] got %v",
			csr.EmailAddresses)
	}

	if len(csr.IPAddresses) != 1 {
		t.Error("expected 1 ip address in generated csr")
	}
}
