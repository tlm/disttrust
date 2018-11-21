package config

import (
	"testing"

	"github.com/tlmiller/disttrust/provider"
)

type dummyProvider struct {
	IssueCalled bool
	RenewCalled bool
}

func (d *dummyProvider) Issue(_ *provider.Request) (provider.Lease, error) {
	d.IssueCalled = true
	return nil, nil
}

func (d *dummyProvider) Renew(_ provider.Lease) (provider.Lease, error) {
	d.RenewCalled = true
	return nil, nil
}

func TestProviderConfigWrapPassThrough(t *testing.T) {
	p := &dummyProvider{
		IssueCalled: false,
		RenewCalled: false,
	}

	pcw := &ProviderConfigWrap{
		P: p,
	}

	pcw.Issue(nil)
	pcw.Renew(nil)

	if !p.IssueCalled {
		t.Fatal("ProviderConfigWrap did not pass issue call through")
	}

	if !p.RenewCalled {
		t.Fatal("ProviderConfigWrap did not pass renew call through")
	}
}
