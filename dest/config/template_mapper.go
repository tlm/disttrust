package config

import (
	"github.com/pkg/errors"

	"github.com/tlmiller/disttrust/dest"
)

type Template struct {
	Gid    string
	Source string
	Mode   string
	Out    string
	Uid    string
}

func TemplateMapper(v interface{}) (dest.Dest, error) {
	conf, ok := v.(*Template)
	if !ok {
		return nil, errors.New("parsing template dest config unknown config type")
	}

	if conf.Source == "" {
		return nil, errors.New("dest template does not have a source template path specified")
	}
	if conf.Out == "" {
		return nil, errors.New("dest template does not have an output path specified")
	}
	outFile, err := DestFileBuilder(conf.Out, conf.Mode, conf.Gid, conf.Uid)
	if err != nil {
		return nil, errors.Wrap(err, "building template output file meta")
	}
	return dest.NewTemplateFile(dest.TemplateFileLoader(conf.Source), outFile), nil
}

func NewTemplate() interface{} {
	return &Template{}
}
