package provider

type AltNames struct {
	DNSNames       []string
	EmailAddresses []string
	IPAddresses    []string
}

type Request struct {
	AltNames           AltNames
	CommonName         string
	Organization       []string
	OrganizationalUnit []string
}
