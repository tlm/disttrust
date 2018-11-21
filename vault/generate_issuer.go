package vault

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/tlmiller/disttrust/provider"
)

const (
	KeyAltNames   = "alt_names"
	KeyCommonName = "common_name"
	KeyFormat     = "format"
	KeyIPSans     = "ip_sans"
)

func GenerateIssuer(path, role string, writer Writer) Issuer {
	return IssuerFunc(func(r *provider.Request) (provider.Lease, error) {
		dest := fmt.Sprintf("%s/issue/%s", path, role)

		data := map[string]interface{}{}
		data[KeyAltNames] = strings.Join(append(r.AltNames.DNSNames,
			r.AltNames.EmailAddresses...), ",")
		data[KeyIPSans] = strings.Join(r.AltNames.IPAddresses, ",")
		data[KeyCommonName] = r.CommonName
		data[KeyFormat] = "pem"

		secret, err := writer.Write(dest, data)
		if err != nil {
			return nil, errors.Wrap(err, "generating certificate from vault")
		}

		lease, err := LeaseFromSecret(r, secret)
		if err != nil {
			return nil, errors.Wrap(err, "making lease from generated certificate")
		}
		return lease, nil
	})
}
