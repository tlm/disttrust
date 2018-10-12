package provider

type DummyProvider struct {
}

func (p *DummyProvider) Issue(_ *Request) (Lease, error) {
	return nil, nil
}

func (p *DummyProvider) Renew(_ Lease) (Lease, error) {
	return nil, nil
}
