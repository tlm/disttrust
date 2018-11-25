package vault

import (
	"testing"

	"github.com/hashicorp/vault/api"

	"github.com/tlmiller/disttrust/provider"
	"github.com/tlmiller/disttrust/request"
)

func TestCSRVerbatimIssuerRew(t *testing.T) {
	csrRequester := request.NewCSRRequester(request.NewRSAKeyMaker(1024))

	called := false
	w := func(p string, d map[string]interface{}) (*api.Secret, error) {
		called = true

		return &api.Secret{
			Data: map[string]interface{}{
				KeyCertificate:  "certificate",
				KeyIssuingCA:    "issuing-ca",
				KeySerialNumber: "serial",
			},
			LeaseDuration: 120,
		}, nil
	}

	req := provider.Request{
		CommonName: "common-name",
		AltNames: provider.AltNames{
			EmailAddresses: []string{"test@example.com"},
			DNSNames:       []string{"example.com"},
			IPAddresses:    []string{"fe80:ffee::1"},
		},
	}

	issuer := CSRVerbatimIssuer("path", "role", csrRequester)
	_, err := issuer.Issue(&req, WriterFunc(w))
	if err != nil {
		t.Fatalf("unexpected error calling csr verbatim issuer: %v", err)
	}

	if !called {
		t.Error("vault test writer was not called by csr verbatim issuer")
	}
}
