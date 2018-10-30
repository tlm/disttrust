package config

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/tlmiller/disttrust/dest"
)

func destTemplateFileMapper(opts json.RawMessage) (dest.Dest, error) {
	uopts := map[string]string{}
	err := json.Unmarshal(opts, &uopts)
	if err != nil {
		return nil, errors.Wrap(err, "parsing dest template json")
	}

	if uopts["source"] == "" {
		return nil, errors.New("dest template does not have a source template path specificed")
	}
	if uopts["out"] == "" {
		return nil, errors.New("dest template does not have an output path specificed")
	}
	outFile, err := destBuilder(uopts["out"], uopts["mode"], uopts["gid"],
		uopts["uid"])
	if err != nil {
		return nil, errors.Wrap(err, "dest template output path")
	}
	return dest.NewTemplateFile(dest.TemplateFileLoader(uopts["source"]), outFile), nil
}

func init() {
	err := MapDest("template", DestMapper(destTemplateFileMapper))
	if err != nil {
		panic(err)
	}
}
