package vault

import (
	"testing"

	"github.com/hashicorp/vault/api"

	"github.com/tlmiller/disttrust/provider"
)

func TestGenerateIssuerReq(t *testing.T) {
	req := provider.Request{
		CommonName: "common-name",
		AltNames: provider.AltNames{
			EmailAddresses: []string{"test@example.com"},
			DNSNames:       []string{"example.com"},
			IPAddresses:    []string{"fe80:ffee::1"},
		},
	}

	called := false
	w := func(p string, d map[string]interface{}) (*api.Secret, error) {
		called = true

		return &api.Secret{
			Data: map[string]interface{}{
				KeyCertificate:  "certificate",
				KeyIssuingCA:    "issuing-ca",
				KeyPrivateKey:   "priv-key",
				KeySerialNumber: "serial",
			},
			LeaseDuration: 120,
		}, nil
	}

	issuer := GenerateIssuer("path", "role")
	_, err := issuer.Issue(&req, WriterFunc(w))
	if err != nil {
		t.Fatalf("unexpected error calling generate issuer: %v", err)
	}

	if !called {
		t.Error("vault test writer was not called by generate issuers")
	}
}
