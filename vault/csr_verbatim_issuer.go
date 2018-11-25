package vault

import (
	"encoding/pem"
	"fmt"

	"github.com/pkg/errors"

	"github.com/tlmiller/disttrust/provider"
	"github.com/tlmiller/disttrust/request"
)

const (
	KeyCSR = "csr"
)

func CSRVerbatimIssuer(path, role string, requester *request.CSRRequester) Issuer {
	return IssuerFunc(func(r *provider.Request, w Writer) (provider.Lease, error) {
		csrASN1, key, err := requester.CSRFromRequest(r)
		if err != nil {
			return nil, errors.Wrapf(err, "making vault sign csr for %s",
				r.CommonName)
		}

		csrPEMRaw := pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE REQUEST",
			Bytes: csrASN1,
		})
		dest := fmt.Sprintf("%s/sign-verbatim/%s", path, role)
		data := map[string]interface{}{}
		data[KeyCSR] = string(csrPEMRaw)
		data[KeyFormat] = "pem"

		secret, err := w.Write(dest, data)
		if err != nil {
			return nil, errors.Wrapf(err,
				"signing certificate verbatim from vault for %s", r.CommonName)
		}

		keyPKCS8, err := key.PKCS8()
		if err != nil {
			return nil, errors.Wrapf(err, "encoding private key as PKCS8 for %s",
				r.CommonName)
		}

		keyPem := pem.EncodeToMemory(&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: keyPKCS8,
		})
		secret.Data[KeyPrivateKey] = string(keyPem)

		lease, err := LeaseFromSecret(r, secret)
		if err != nil {
			return nil, errors.Wrapf(err,
				"making lease from generated certificate for %s", r.CommonName)
		}
		return lease, nil
	})
}
