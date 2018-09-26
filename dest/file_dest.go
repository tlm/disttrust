package dest

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"

	"github.com/tlmiller/disttrust/file"
	"github.com/tlmiller/disttrust/provider"
)

type File struct {
	CA                file.File
	Certificate       file.File
	CertificateBundle file.File
	PrivateKey        file.File
}

func (f *File) Send(res *provider.Response) error {

	if res.CA != "" && f.CA.HasPath() {
		err := ioutil.WriteFile(f.CA.Path, []byte(res.CA), f.CA.Mode)
		if err != nil {
			return errors.Wrap(err, "writing ca file")
		}
		err = f.CA.Chown()
		if err != nil {
			return errors.Wrap(err, "chown ca file")
		}
	}

	if res.Certificate != "" && f.Certificate.HasPath() {
		err := ioutil.WriteFile(f.Certificate.Path, []byte(res.Certificate),
			f.Certificate.Mode)
		if err != nil {
			return errors.Wrap(err, "writing certificate file")
		}
		err = f.Certificate.Chown()
		if err != nil {
			return errors.Wrap(err, "chown certificate file")
		}
	}

	if res.CABundle != "" && f.CertificateBundle.HasPath() {
		s, err := os.OpenFile(f.CertificateBundle.Path,
			os.O_WRONLY|os.O_TRUNC|os.O_CREATE, f.CertificateBundle.Mode)
		defer s.Close()
		if err != nil {
			return errors.Wrap(err, "writing certificate bundle file")
		}
		_, err = s.WriteString(res.CABundle + "\n" + res.Certificate)
		if err != nil {
			return errors.Wrap(err, "writing certificate bundle file")
		}
		err = f.CertificateBundle.Chown()
		if err != nil {
			return errors.Wrap(err, "chown certificate bundle file")
		}
	}

	if res.PrivateKey != "" && f.PrivateKey.HasPath() {
		err := ioutil.WriteFile(f.PrivateKey.Path, []byte(res.PrivateKey),
			f.PrivateKey.Mode)
		if err != nil {
			return errors.Wrap(err, "writing private key file")
		}
		err = f.PrivateKey.Chown()
		if err != nil {
			return errors.Wrap(err, "chown private key file")
		}
	}
	return nil
}
