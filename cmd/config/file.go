package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func FromFiles(files ...string) (*Config, error) {
	conf := DefaultConfig()
	for _, file := range files {
		finfo, err := os.Stat(file)
		if err != nil {
			return nil, errors.Wrap(err, "geting config file info")
		}

		if finfo.Mode().IsDir() {
			conf, err = FromDirectory(conf, file)
			if err != nil {
				return nil, err
			}
			continue
		}

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

func FromDirectory(conf *Config, dir string) (*Config, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return conf, errors.Wrap(err, "reading config files from dir")
	}

	for _, f := range files {
		if f.IsDir() || filepath.Ext(f.Name()) != ".json" {
			continue
		}

		raw, err := ioutil.ReadFile(filepath.Join(dir, f.Name()))
		if err != nil {
			return conf, errors.Wrap(err, "reading config file")
		}
		c, err := New(raw)
		if err != nil {
			return conf, errors.Wrap(err, "parsing config file")
		}
		conf = mergeConfigs(conf, c)
	}

	return conf, nil
}
