package config

import (
	"testing"

	"github.com/tlmiller/disttrust/provider"
)

type dummyProvider struct {
	InitialiseCalled bool
	IssueCalled      bool
	NameVal          string
	NameCalled       bool
	RenewCalled      bool
}

func (d *dummyProvider) Initialise() error {
	d.InitialiseCalled = true
	return nil
}

func (d *dummyProvider) Issue(_ *provider.Request) (provider.Lease, error) {
	d.IssueCalled = true
	return nil, nil
}

func (d *dummyProvider) Name() string {
	d.NameCalled = true
	return d.NameVal
}

func (d *dummyProvider) Renew(_ provider.Lease) (provider.Lease, error) {
	d.RenewCalled = true
	return nil, nil
}

func TestProviderConfigWrapPassThrough(t *testing.T) {
	p := &dummyProvider{
		NameVal: "test",
	}

	pcw := &ProviderConfigWrap{
		P: p,
	}

	pcw.Initialise()
	pcw.Issue(nil)
	pcw.Renew(nil)

	if pcw.Name() != "test" {
		t.Error("ProviderConfigWrap return val for Name() was not expected")
	}

	if !p.InitialiseCalled {
		t.Error("ProviderConfigWrap did not pass initialise call through")
	}

	if !p.IssueCalled {
		t.Error("ProviderConfigWrap did not pass issue call through")
	}

	if !p.NameCalled {
		t.Error("ProviderConfigWrap did not pass name call through")
	}

	if !p.RenewCalled {
		t.Error("ProviderConfigWrap did not pass renew call through")
	}
}
