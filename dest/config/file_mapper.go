package config

import (
	"github.com/pkg/errors"

	"github.com/tlmiller/disttrust/dest"
)

type File struct {
	CAFile             string
	CAFileMode         string
	CAFileGid          string
	CAFileUid          string
	CABundleFile       string
	CABundleFileMode   string
	CABundleFileGid    string
	CABundleFileUid    string
	CertFile           string
	CertFileMode       string
	CertFileGid        string
	CertFileUid        string
	CertBundleFile     string
	CertBundleFileMode string
	CertBundleFileGid  string
	CertBundleFileUid  string
	PrivKeyFile        string
	PrivKeyFileMode    string
	PrivKeyFileGid     string
	PrivKeyFileUid     string
}

func FileMapper(v interface{}) (dest.Dest, error) {
	fileDests := []dest.Dest{}
	conf, ok := v.(*File)
	if !ok {
		return nil, errors.New("parsing file dest config unknown config type")
	}

	if conf.CAFile != "" {
		caFile, err := DestFileBuilder(conf.CAFile, conf.CAFileMode,
			conf.CAFileGid, conf.CAFileUid)
		if err != nil {
			return nil, errors.Wrap(err, "caFile")
		}
		fileDests = append(fileDests, dest.NewTemplateFile(dest.CAFile, caFile))
	}
	if conf.CABundleFile != "" {
		caBundleFile, err := DestFileBuilder(conf.CABundleFile,
			conf.CABundleFileMode, conf.CABundleFileGid, conf.CABundleFileUid)
		if err != nil {
			return nil, errors.Wrap(err, "caBundleFile")
		}
		fileDests = append(fileDests, dest.NewTemplateFile(dest.CABundleFile,
			caBundleFile))
	}
	if conf.CertFile != "" {
		certFile, err := DestFileBuilder(conf.CertFile, conf.CertFileMode,
			conf.CertFileGid, conf.CertFileUid)
		if err != nil {
			return nil, errors.Wrap(err, "certFile")
		}
		fileDests = append(fileDests, dest.NewTemplateFile(dest.CertificateFile,
			certFile))
	}
	if conf.CertBundleFile != "" {
		certBundleFile, err := DestFileBuilder(conf.CertBundleFile,
			conf.CertBundleFileMode, conf.CertBundleFileGid,
			conf.CertBundleFileUid)
		if err != nil {
			return nil, errors.Wrap(err, "certBundleFile")
		}
		fileDests = append(fileDests,
			dest.NewTemplateFile(dest.CertificateBundleFile, certBundleFile))
	}
	if conf.PrivKeyFile != "" {
		privKeyFile, err := DestFileBuilder(conf.PrivKeyFile,
			conf.PrivKeyFileMode, conf.PrivKeyFileGid,
			conf.PrivKeyFileUid)
		if err != nil {
			return nil, errors.Wrap(err, "privKeyFile")
		}
		fileDests = append(fileDests, dest.NewTemplateFile(dest.PrivateKeyFile,
			privKeyFile))
	}
	return dest.NewAggregate(fileDests...), nil
}

func NewFile() interface{} {
	return &File{}
}
