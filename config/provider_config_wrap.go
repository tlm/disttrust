package config

import (
	"github.com/tlmiller/disttrust/provider"
)

type ProviderConfigWrap struct {
	Options map[string]interface{}
	P       provider.Provider
}

func (p *ProviderConfigWrap) Issue(r *provider.Request) (provider.Lease, error) {
	return p.P.Issue(r)
}

func (p *ProviderConfigWrap) Renew(l provider.Lease) (provider.Lease, error) {
	return p.P.Renew(l)
}
