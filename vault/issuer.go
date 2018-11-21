package vault

import (
	"github.com/tlmiller/disttrust/provider"
)

type Issuer interface {
	Issue(*provider.Request) (provider.Lease, error)
}

type IssuerFunc func(*provider.Request) (provider.Lease, error)

func (i IssuerFunc) Issue(r *provider.Request) (provider.Lease, error) {
	return i(r)
}
