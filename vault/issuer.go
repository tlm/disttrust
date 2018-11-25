package vault

import (
	"github.com/tlmiller/disttrust/provider"
)

type Issuer interface {
	Issue(*provider.Request, Writer) (provider.Lease, error)
}

type IssuerFunc func(*provider.Request, Writer) (provider.Lease, error)

func (i IssuerFunc) Issue(r *provider.Request, w Writer) (provider.Lease, error) {
	return i(r, w)
}
