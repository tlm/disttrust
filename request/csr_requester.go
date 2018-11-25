package request

import (
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"net"

	"github.com/pkg/errors"

	"github.com/tlmiller/disttrust/provider"
)

type CSRRequester struct {
	KeyMaker KeyMaker
}

func (c *CSRRequester) CSRFromRequest(req *provider.Request) (
	[]byte, Key, error) {

	csrTmpl := x509.CertificateRequest{}
	csrTmpl.Subject.CommonName = req.CommonName
	csrTmpl.Subject.Organization = req.Organization
	csrTmpl.Subject.OrganizationalUnit = req.OrganizationalUnit
	csrTmpl.DNSNames = req.AltNames.DNSNames
	csrTmpl.EmailAddresses = req.AltNames.EmailAddresses

	for _, ipStr := range req.AltNames.IPAddresses {
		ip := net.ParseIP(ipStr)
		if ip == nil {
			return nil, nil, fmt.Errorf("ip address %s is not valid", ipStr)
		}
		csrTmpl.IPAddresses = append(csrTmpl.IPAddresses, ip)
	}

	key, err := c.KeyMaker.MakeKey()
	if err != nil {
		return nil, nil, errors.Wrapf(err, "making csr private key for %s",
			req.CommonName)
	}

	csrData, err := x509.CreateCertificateRequest(rand.Reader, &csrTmpl, key.Raw())
	if err != nil {
		return nil, nil, errors.Wrapf(err, "creating certificate request for %s",
			req.CommonName)
	}

	return csrData, key, nil
}

func NewCSRRequester(keyMaker KeyMaker) *CSRRequester {
	return &CSRRequester{
		KeyMaker: keyMaker,
	}
}
