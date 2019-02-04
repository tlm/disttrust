package provider

type Id string

type Provider interface {
	Name() string
	Initialise() error
	Issue(*Request) (Lease, error)
	Renew(Lease) (Lease, error)
}
