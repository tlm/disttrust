package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"

	"github.com/tlmiller/disttrust/dest"
	"github.com/tlmiller/disttrust/file"
)

func ToDest(id string, opts json.RawMessage) (dest.Dest, error) {
	if id != "file" {
		return nil, fmt.Errorf("unknown dest type '%s'", id)
	}
	uopts := map[string]string{}
	err := json.Unmarshal(opts, &uopts)
	if err != nil {
		return nil, errors.Wrap(err, "parsing dest json")
	}
	fdest := dest.File{}

	caFile, err := destBuilder(uopts["caFile"], uopts["caFileMode"],
		uopts["caFileGid"], uopts["caFileUid"])
	if err != nil {
		return nil, errors.Wrap(err, "caFile")
	}
	fdest.CA = caFile

	cFile, err := destBuilder(uopts["certFile"], uopts["certFileMode"],
		uopts["certFileGid"], uopts["certFileUid"])
	if err != nil {
		return nil, errors.Wrap(err, "certFile")
	}
	fdest.Certificate = cFile

	cbFile, err := destBuilder(uopts["certBundleFile"], uopts["certBundleFileMode"],
		uopts["certBundleFileGid"], uopts["certBundleFileUid"])
	if err != nil {
		return nil, errors.Wrap(err, "certBundleFile")
	}
	fdest.CertificateBundle = cbFile

	pkfile, err := destBuilder(uopts["privKeyFile"], uopts["privKeyFileMode"],
		uopts["privKeyFileGid"], uopts["privKeyFileUid"])
	if err != nil {
		return nil, errors.Wrap(err, "privKeyFile")
	}
	fdest.PrivateKey = pkfile

	return &fdest, nil
}

func destBuilder(path, mode, gid, uid string) (file.File, error) {
	builder := file.New(path)
	if mode != "" {
		conv, err := strconv.ParseUint(mode, 8, 32)
		if err != nil {
			return file.File{}, errors.Wrap(err, "invalid mode uint")
		}
		builder.Mode = os.FileMode(conv)
	}
	builder.Gid = gid
	builder.Uid = uid
	return builder, nil
}
