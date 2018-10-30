package config

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/tlmiller/disttrust/dest"
)

func destFileMapper(opts json.RawMessage) (dest.Dest, error) {
	fileDests := []dest.Dest{}
	uopts := map[string]string{}
	err := json.Unmarshal(opts, &uopts)
	if err != nil {
		return nil, errors.Wrap(err, "parsing dest file json")
	}

	if uopts["caFile"] != "" {
		caFile, err := destBuilder(uopts["caFile"], uopts["caFileMode"],
			uopts["caFileGid"], uopts["caFileUid"])
		if err != nil {
			return nil, errors.Wrap(err, "caFile")
		}
		fileDests = append(fileDests, dest.NewTemplateFile(dest.CAFile, caFile))
	}
	if uopts["certFile"] != "" {
		cFile, err := destBuilder(uopts["certFile"], uopts["certFileMode"],
			uopts["certFileGid"], uopts["certFileUid"])
		if err != nil {
			return nil, errors.Wrap(err, "certFile")
		}
		fileDests = append(fileDests, dest.NewTemplateFile(dest.CertificateFile, cFile))
	}
	if uopts["certBundleFile"] != "" {
		cbFile, err := destBuilder(uopts["certBundleFile"], uopts["certBundleFileMode"],
			uopts["certBundleFileGid"], uopts["certBundleFileUid"])
		if err != nil {
			return nil, errors.Wrap(err, "certBundleFile")
		}
		fileDests = append(fileDests, dest.NewTemplateFile(dest.CertificateBundleFile, cbFile))
	}
	if uopts["privKeyFile"] != "" {
		pkFile, err := destBuilder(uopts["privKeyFile"], uopts["privKeyFileMode"],
			uopts["privKeyFileGid"], uopts["privKeyFileUid"])
		if err != nil {
			return nil, errors.Wrap(err, "privKeyFile")
		}
		fileDests = append(fileDests, dest.NewTemplateFile(dest.PrivateKeyFile, pkFile))
	}
	return dest.NewAggregate(fileDests...), nil
}

func init() {
	err := MapDest("file", DestMapper(destFileMapper))
	if err != nil {
		panic(err)
	}
}
