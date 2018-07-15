package dest

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"

	"github.com/tlmiller/disttrust/provider"
)

type File struct {
	CAFile                string
	CertificateFile       string
	CertificateBundleFile string
	PrivateKeyFile        string
}

func (f *File) Send(res *provider.Response) error {

	if res.CA != "" && f.CAFile != "" {
		err := ioutil.WriteFile(f.CAFile, []byte(res.CA), os.FileMode(0644))
		if err != nil {
			return errors.Wrap(err, "writing ca file")
		}
	}

	if res.Certificate != "" && f.CertificateFile != "" {
		err := ioutil.WriteFile(f.CertificateFile, []byte(res.Certificate), os.FileMode(0644))
		if err != nil {
			return errors.Wrap(err, "writing certificate file")
		}
	}

	if res.CABundle != "" && f.CertificateBundleFile != "" {
		f, err := os.OpenFile(f.CertificateBundleFile,
			os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.FileMode(0644))
		defer f.Close()
		if err != nil {
			return errors.Wrap(err, "writing certificate bundle file")
		}
		_, err = f.WriteString(res.CABundle + "\n" + res.Certificate)
		if err != nil {
			return errors.Wrap(err, "writing certificate bundle file")
		}
	}

	if res.PrivateKey != "" && f.PrivateKeyFile != "" {
		err := ioutil.WriteFile(f.PrivateKeyFile, []byte(res.PrivateKey), os.FileMode(0600))
		if err != nil {
			return errors.Wrap(err, "writing private key file")
		}
	}
	return nil
}
