package config

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/tlmiller/disttrust/dest"
)

func MakeDest(id string, opts json.RawMessage) (dest.Dest, error) {
	if id != "file" {
		return nil, fmt.Errorf("unknown dest type '%s'", id)
	}
	uopts := map[string]string{}
	err := json.Unmarshal(opts, &uopts)
	if err != nil {
		return nil, errors.Wrap(err, "parsing dest json")
	}
	fdest := dest.File{}
	if cafile, exists := uopts["caFile"]; exists {
		fdest.CAFile = cafile
	}
	if cfile, exists := uopts["certFile"]; exists {
		fdest.CertificateFile = cfile
	}
	if cbfile, exists := uopts["certBundleFile"]; exists {
		fdest.CertificateBundleFile = cbfile
	}
	if pkfile, exists := uopts["privKeyFile"]; exists {
		fdest.PrivateKeyFile = pkfile
	}
	return &fdest, nil
}
