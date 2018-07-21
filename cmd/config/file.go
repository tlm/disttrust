package config

import (
	"io/ioutil"

	"github.com/pkg/errors"
)

func FromFiles(files ...string) (*Config, error) {
	conf := DefaultConfig()
	for _, file := range files {
		raw, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, errors.Wrap(err, "reading config file")
		}

		c, err := New(raw)
		if err != nil {
			return nil, errors.Wrap(err, "parsing config file")
		}
		conf = mergeConfigs(conf, c)
	}
	return conf, nil
}
