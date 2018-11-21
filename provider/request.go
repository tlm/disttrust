package provider

type AltNames struct {
	EmailAddresses []string
	DNSNames       []string
	IPAddresses    []string
}

type Request struct {
	CommonName string
	AltNames   AltNames
}
