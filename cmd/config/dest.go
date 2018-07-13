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
	if cfile, exists := uopts["certificateFile"]; exists {
		fdest.CertificateFile = cfile
	}
	if pkfile, exists := uopts["privateKeyFile"]; exists {
		fdest.PrivateKeyFile = pkfile
	}
	return &fdest, nil
}
