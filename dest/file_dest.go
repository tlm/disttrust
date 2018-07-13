package dest

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"

	"github.com/tlmiller/disttrust/provider"
)

type File struct {
	CertificateFile string
	PrivateKeyFile  string
}

func (f *File) Send(res *provider.Response) error {

	if res.Certificate != "" && f.CertificateFile != "" {
		err := ioutil.WriteFile(f.CertificateFile, []byte(res.Certificate), os.FileMode(0644))
		if err != nil {
			return errors.Wrap(err, "writing certificate file")
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
